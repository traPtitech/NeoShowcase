import { type Timestamp, timestampDate, timestampNow } from '@bufbuild/protobuf/wkt'
import { Code, ConnectError } from '@connectrpc/connect'
import { type Component, For, Show, createEffect, createSignal, onCleanup } from 'solid-js'
import type { ApplicationOutput } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { LogContainer } from '/@/components/UI/LogContainer'
import { ContainerLogExport } from '/@/components/templates/app/ContainerLogExport'
import { client, handleAPIError } from '/@/libs/api'
import { toWithAnsi } from '/@/libs/buffers'
import { isScrolledToBottom } from '/@/libs/scroll'
import { addTimestamp, lessTimestamp, minTimestamp } from '/@/libs/timestamp'

const logLinesLimit = 1000
const loadLimitSeconds = 7 * 86400
const loadDuration = 86400n

const stripLogLines = (logs: ApplicationOutput[]): ApplicationOutput[] => {
  if (logs.length >= logLinesLimit) {
    return logs.slice(logs.length - logLinesLimit, logs.length)
  }
  return logs
}

const loadLogChunk = async (appID: string, before: Timestamp, limit: number): Promise<ApplicationOutput[]> => {
  const res = await client.getOutput({
    applicationId: appID,
    before: before,
    limit: limit,
  })
  return res.outputs
}

const oldestTimestamp = (ts: ApplicationOutput[]): Timestamp =>
  ts.reduce((acc, t) => (t.time ? minTimestamp(acc, t.time) : acc), timestampNow())
const sortByTimestamp = (ts: ApplicationOutput[]) =>
  ts.sort((a, b) => (a.time && b.time ? (lessTimestamp(a.time, b.time) ? -1 : 1) : 0))

export interface ContainerLogProps {
  appID: string
}

export const ContainerLog: Component<ContainerLogProps> = (props) => {
  // TODO: show timestamps toggle button?
  const [showTimestamps] = createSignal(true)

  const componentLoadTime = timestampNow()
  const [loadedUntil, setLoadedUntil] = createSignal(componentLoadTime)
  const [logs, setLogs] = createSignal<ApplicationOutput[]>([])

  const loadDisabled = () =>
    timestampNow().seconds - loadedUntil().seconds >= loadLimitSeconds || logs().length >= logLinesLimit
  const [loading, setLoading] = createSignal(false)
  const load = async () => {
    setLoading(true)
    try {
      const loadedOlderLogs = await loadLogChunk(props.appID, loadedUntil(), 100)
      if (loadedOlderLogs.length === 0) {
        setLoadedUntil(addTimestamp(loadedUntil(), -loadDuration))
      } else {
        setLoadedUntil(oldestTimestamp(loadedOlderLogs))
      }
      sortByTimestamp(loadedOlderLogs)
      setLogs((prev) => stripLogLines(loadedOlderLogs.concat(prev)))
    } catch (e) {
      handleAPIError(e, 'ログの読み込み中にエラーが発生しました')
    }
    setLoading(false)
  }
  // Load logs before component load time
  void load()

  // Stream logs beginning from component load time
  const logStreamAbort = new AbortController()
  const logStream = client.getOutputStream(
    {
      applicationId: props.appID,
      begin: componentLoadTime,
    },
    { signal: logStreamAbort.signal },
  )
  const iterate = async () => {
    try {
      for await (const log of logStream) {
        setLogs((prev) => stripLogLines(prev.concat(log)))
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

  onCleanup(() => {
    logStreamAbort.abort()
  })

  let logRef!: HTMLDivElement
  createEffect(() => {
    logs() // on change to (streamed) logs
    const ref = logRef
    if (!ref) return
    if (atBottom()) {
      ref.scrollTop = ref.scrollHeight
    }
  })

  const [atBottom, setAtBottom] = createSignal(true)
  const onScroll = (e: { target: Element }) => setAtBottom(isScrolledToBottom(e.target))

  return (
    <div class="flex flex-col gap-2.5">
      <ContainerLogExport currentLogs={logs()} />
      <LogContainer ref={logRef!} overflowX="scroll" onScroll={onScroll}>
        <div class="mb-1.5 flex items-center gap-2">
          Loaded until {timestampDate(loadedUntil()).toLocaleString()}
          <Show when={!loadDisabled()} fallback={<span>(reached load limit)</span>}>
            <Button variants="ghost" size="small" onClick={load} disabled={loading()}>
              {loading() ? 'Loading...' : 'Load more'}
            </Button>
          </Show>
        </div>
        <For each={logs()}>{(log) => <code innerHTML={formatLogLine(log, showTimestamps())} />}</For>
      </LogContainer>
    </div>
  )
}

const formatLogLine = (log: ApplicationOutput, withTimestamp: boolean): string => {
  return (
    (withTimestamp ? `${log.time ? timestampDate(log.time).toLocaleString() : undefined} ` : '') + toWithAnsi(log.log)
  )
}
