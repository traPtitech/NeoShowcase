import { ApplicationOutput } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { LogContainer } from '/@/components/Log'
import { client, handleAPIError } from '/@/libs/api'
import { toWithAnsi } from '/@/libs/buffers'
import { isScrolledToBottom } from '/@/libs/scroll'
import { addTimestamp, lessTimestamp, minTimestamp } from '/@/libs/timestamp'
import { vars } from '/@/theme'
import { Code, ConnectError } from '@bufbuild/connect'
import { Timestamp } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, For, Ref, Show, createEffect, createMemo, createResource, createSignal, onCleanup } from 'solid-js'

const LoadMoreContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    marginBottom: '6px',
    fontSize: '16px',
  },
})

const LoadMoreButton = styled('div', {
  base: {
    background: vars.bg.black1,
    border: `1px solid ${vars.bg.white2}`,
    borderRadius: '2px',
    padding: '2px',
    color: vars.text.white1,
    selectors: {
      '&:hover': {
        background: vars.text.black3,
      },
    },
  },
})

const loadLimitSeconds = 7 * 86400
const loadDuration = 86400n

const loadLogChunk = async (appID: string, before: Timestamp): Promise<ApplicationOutput[]> => {
  const res = await client.getOutput({ applicationId: appID, before: before })
  return res.outputs
}

const oldestTimestamp = (ts: ApplicationOutput[]): Timestamp =>
  ts.reduce((acc, t) => minTimestamp(acc, t.time), Timestamp.now())
const sortByTimestamp = (ts: ApplicationOutput[]) => ts.sort((a, b) => (lessTimestamp(a.time, b.time) ? -1 : 1))

export interface ContainerLogProps {
  appID: string
  showTimestamp: boolean
}

export const ContainerLog: Component<ContainerLogProps> = (props) => {
  const [loadedUntil, setLoadedUntil] = createSignal(Timestamp.now())
  const [olderLogs, setOlderLogs] = createSignal<ApplicationOutput[]>([])

  const loadDisabled = () => Timestamp.now().seconds - loadedUntil().seconds >= loadLimitSeconds
  const [loading, setLoading] = createSignal(false)
  const load = async () => {
    setLoading(true)
    try {
      const loadedOlderLogs = await loadLogChunk(props.appID, loadedUntil())
      if (loadedOlderLogs.length === 0) {
        setLoadedUntil(addTimestamp(loadedUntil(), -loadDuration))
      } else {
        setLoadedUntil(oldestTimestamp(loadedOlderLogs))
      }
      sortByTimestamp(loadedOlderLogs)
      setOlderLogs((prev) => loadedOlderLogs.concat(prev))
    } catch (e) {
      handleAPIError(e, 'ログの読み込み中にエラーが発生しました')
    }
    setLoading(false)
  }

  const logStreamAbort = new AbortController()
  const [logStream] = createResource(
    () => props.appID,
    (id) => client.getOutputStream({ id }, { signal: logStreamAbort.signal }),
  )
  const [streamedLog, setStreamedLog] = createSignal<ApplicationOutput[]>([])
  createEffect(() => {
    const stream = logStream()
    if (!stream) {
      setStreamedLog([])
      return
    }

    const iterate = async () => {
      try {
        for await (const log of stream) {
          setStreamedLog((prev) => prev.concat(log))
        }
      } catch (err) {
        // ignore abort error
        const isAbortErr = err instanceof ConnectError && err.code === Code.Canceled
        if (!isAbortErr) {
          console.trace(err)
          return
        }
      }
    }

    void iterate()
  })
  onCleanup(() => {
    logStreamAbort.abort()
  })

  const streamedLogOldest = createMemo(() => {
    const logs = streamedLog()
    if (logs.length === 0) return
    return logs.reduce((acc, log) => minTimestamp(acc, log.time), Timestamp.now())
  })
  createEffect(() => {
    if (!streamedLogOldest()) return
    if (lessTimestamp(streamedLogOldest(), loadedUntil())) {
      setLoadedUntil(streamedLogOldest())
    }
  })

  let logRef: Ref<HTMLDivElement>
  createEffect(() => {
    streamedLog()
    const ref = logRef as HTMLDivElement
    if (!ref) return
    if (atBottom()) {
      ref.scrollTop = ref.scrollHeight
    }
  })

  const [atBottom, setAtBottom] = createSignal(true)
  const onScroll = (e: { target: Element }) => setAtBottom(isScrolledToBottom(e.target))

  return (
    <LogContainer ref={logRef} overflowX="scroll" onScroll={onScroll}>
      {/* cannot distinguish zero log and loading (but should be enough for most use-cases) */}
      <Show when={streamedLog().length > 0}>
        <LoadMoreContainer>
          Loaded until {loadedUntil().toDate().toLocaleString()}
          <Show when={!loadDisabled()} fallback={<span>(reached load limit)</span>}>
            <Show when={!loading()} fallback={<span>Loading...</span>}>
              <LoadMoreButton onClick={load}>Load more</LoadMoreButton>
            </Show>
          </Show>
        </LoadMoreContainer>
      </Show>
      <For each={olderLogs()}>{(log) => <code innerHTML={formatLogLine(log, props.showTimestamp)} />}</For>
      <For each={streamedLog()}>{(log) => <code innerHTML={formatLogLine(log, props.showTimestamp)} />}</For>
    </LogContainer>
  )
}

const formatLogLine = (log: ApplicationOutput, withTimestamp: boolean): string => {
  return (withTimestamp ? `${log.time.toDate().toLocaleString()} ` : '') + toWithAnsi(log.log)
}
