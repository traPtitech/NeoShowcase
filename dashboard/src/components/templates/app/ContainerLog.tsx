import { type Timestamp, timestampDate, timestampNow } from '@bufbuild/protobuf/wkt'
import { Code, ConnectError } from '@connectrpc/connect'
import { type Component, createEffect, createMemo, createSignal, For, Show } from 'solid-js'
import type { ApplicationOutput } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { ContainerLogExport } from '/@/components/templates/app/ContainerLogExport'
import { Button } from '/@/components/UI/Button'
import { LogContainer } from '/@/components/UI/LogContainer'
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

const loadLogChunk = async (
  appID: string,
  before: Timestamp,
  limit: number,
  signal: AbortSignal,
): Promise<ApplicationOutput[]> => {
  const res = await client.getOutput({ applicationId: appID, before, limit }, { signal })
  return res.outputs
}

const oldestTimestamp = (ts: ApplicationOutput[]): Timestamp =>
  ts.reduce((acc, t) => (t.time ? minTimestamp(acc, t.time) : acc), timestampNow())
const sortByTimestamp = (ts: ApplicationOutput[]) =>
  ts.sort((a, b) => (a.time && b.time ? (lessTimestamp(a.time, b.time) ? -1 : 1) : 0))

const streamRetryMs = 3000

const sleep = (ms: number, signal: AbortSignal) =>
  new Promise<void>((resolve) => {
    const timer = setTimeout(resolve, ms)
    signal.addEventListener(
      'abort',
      () => {
        clearTimeout(timer)
        resolve()
      },
      { once: true },
    )
  })

/**
 * Streams the logs of `appID` from `begin` into `onLog` until `signal` is aborted, reconnecting after the stream
 * ends or fails. A reconnect resumes from the last line received; the server's `begin` is inclusive and only
 * second-precise, so lines it sends again are skipped. Lines are compared only with received ones (server time),
 * never with `begin` (client time), so clock skew cannot drop lines.
 */
const streamLogs = async (
  appID: string,
  begin: Timestamp,
  signal: AbortSignal,
  onLog: (log: ApplicationOutput) => void,
) => {
  const key = (l: ApplicationOutput) => `${l.time?.seconds}.${l.time?.nanos} ${l.log}`
  let last: Timestamp | undefined
  let seenInLatestSecond = new Set<string>()
  while (!signal.aborted) {
    try {
      for await (const log of client.getOutputStream({ applicationId: appID, begin: last ?? begin }, { signal })) {
        if (log.time) {
          if (last && (lessTimestamp(log.time, last) || seenInLatestSecond.has(key(log)))) continue
          if (!last || log.time.seconds > last.seconds) seenInLatestSecond = new Set()
          last = log.time
          seenInLatestSecond.add(key(log))
        }
        onLog(log)
      }
    } catch (err) {
      const isAbortErr = err instanceof ConnectError && err.code === Code.Canceled
      if (signal.aborted || isAbortErr) return
      console.trace(err)
    }
    await sleep(streamRetryMs, signal)
  }
}

export interface ContainerLogProps {
  appID: string
}

export const ContainerLog: Component<ContainerLogProps> = (props) => {
  // TODO: show timestamps toggle button?
  const [showTimestamps] = createSignal(true)

  const [loadedUntil, setLoadedUntil] = createSignal(timestampNow())
  const [logs, setLogs] = createSignal<ApplicationOutput[]>([])
  const [loading, setLoading] = createSignal(false)

  const loadDisabled = () =>
    timestampNow().seconds - loadedUntil().seconds >= loadLimitSeconds || logs().length >= logLinesLimit

  // Aborted when the app changes or the component goes away. It is per visit, not per app id: after A → B → A,
  // a response requested during the first visit of A must not land in the second one.
  let visit = new AbortController()

  const load = async (appID: string, before: Timestamp, signal: AbortSignal) => {
    setLoading(true)
    try {
      const loadedOlderLogs = await loadLogChunk(appID, before, 100, signal)
      if (signal.aborted) return
      if (loadedOlderLogs.length === 0) {
        setLoadedUntil(addTimestamp(before, -loadDuration))
      } else {
        setLoadedUntil(oldestTimestamp(loadedOlderLogs))
      }
      sortByTimestamp(loadedOlderLogs)
      setLogs((prev) => stripLogLines(loadedOlderLogs.concat(prev)))
    } catch (e) {
      if (signal.aborted) return
      handleAPIError(e, 'ログの読み込み中にエラーが発生しました')
    }
    setLoading(false)
  }

  // (Re)start for each app: load logs before now, and stream logs from now on.
  // `props.appID` comes from the app query, which goes pending on every revalidation (the page polls it). An effect
  // whose compute reads that directly re-runs when it settles even with the same id, which would wipe the loaded logs
  // and restart the stream; the memo only propagates an actual change of id.
  const appID = createMemo(() => props.appID)
  createEffect(appID, (appID) => {
    const abort = new AbortController()
    visit = abort
    const begin = timestampNow()
    setLogs([])
    setLoadedUntil(begin)
    setLoading(false)
    void load(appID, begin, abort.signal)

    void streamLogs(appID, begin, abort.signal, (log) => setLogs((prev) => stripLogLines(prev.concat(log))))
    return () => abort.abort()
  })

  let logRef!: HTMLElement
  let atBottom = true
  createEffect(logs, () => {
    if (atBottom) logRef.scrollTop = logRef.scrollHeight
  })

  return (
    <div class="flex flex-col gap-2.5">
      <ContainerLogExport currentLogs={logs()} />
      <LogContainer
        ref={logRef}
        overflowX="scroll"
        onScroll={(e) => {
          atBottom = isScrolledToBottom(e.currentTarget)
        }}
      >
        <div class="mb-1.5 flex items-center gap-2">
          Loaded until {timestampDate(loadedUntil()).toLocaleString()}
          <Show when={!loadDisabled()} fallback={<span>(reached load limit)</span>}>
            <Button
              variants="ghost"
              size="small"
              onClick={() => load(appID(), loadedUntil(), visit.signal)}
              disabled={loading()}
            >
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
