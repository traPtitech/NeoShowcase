import { Code, ConnectError } from '@connectrpc/connect'
import { type Component, createEffect, createMemo, createSignal, onSettled, Show } from 'solid-js'
import { LogContainer } from '/@/components/UI/LogContainer'
import { client } from '/@/libs/api'
import { concatBuffers, toUTF8WithAnsi } from '/@/libs/buffers'
import { isScrolledToBottom } from '/@/libs/scroll'
import { sleep } from '/@/libs/sleep'

export interface BuildLogProps {
  buildID: string
  finished: boolean
  /** Called when the log stream ends, so that the parent can reload the build. */
  refetch: () => void
}

/** How often a still-queued build is re-checked (its log stream cannot be opened before it starts). */
const QUEUED_RETRY_MS = 3000

const FinishedBuildLog: Component<{ buildID: string }> = (props) => {
  const buildLog = createMemo(() => client.getBuildLog({ buildId: props.buildID }))

  let logRef!: HTMLElement
  createEffect(buildLog, () => {
    logRef.scrollTop = logRef.scrollHeight
  })

  return <LogContainer innerHTML={toUTF8WithAnsi(buildLog().log)} ref={logRef} />
}

const StreamingBuildLog: Component<{ buildID: string; refetch: () => void }> = (props) => {
  const [log, setLog] = createSignal<Uint8Array>(new Uint8Array())

  onSettled(() => {
    const abort = new AbortController()
    void (async () => {
      let started = false
      while (!abort.signal.aborted) {
        try {
          for await (const chunk of client.getBuildLogStream({ buildId: props.buildID }, { signal: abort.signal })) {
            if (!started) {
              started = true
              props.refetch() // the build left the queue: show its new status
            }
            setLog((prev) => concatBuffers(prev, chunk.log))
          }
        } catch (err) {
          if (err instanceof ConnectError && err.code === Code.Canceled) return
          if (err instanceof ConnectError && err.code === Code.FailedPrecondition) {
            // The build is still queued ("build not started yet"): keep watching it until it starts or is finished.
            await sleep(QUEUED_RETRY_MS)
            if (!abort.signal.aborted) props.refetch()
            continue
          }
          console.trace(err)
          return
        }
        await sleep(1000)
        if (!abort.signal.aborted) props.refetch() // stream ended: build is finished
        return
      }
    })()
    return () => abort.abort()
  })

  let logRef!: HTMLElement
  let atBottom = true
  createEffect(log, () => {
    if (atBottom) logRef.scrollTop = logRef.scrollHeight
  })

  return (
    <LogContainer
      innerHTML={toUTF8WithAnsi(log())}
      ref={logRef}
      onScroll={(e) => {
        atBottom = isScrolledToBottom(e.currentTarget)
      }}
    />
  )
}

export const BuildLog: Component<BuildLogProps> = (props) => (
  <Show when={props.finished} fallback={<StreamingBuildLog buildID={props.buildID} refetch={props.refetch} />}>
    <FinishedBuildLog buildID={props.buildID} />
  </Show>
)
