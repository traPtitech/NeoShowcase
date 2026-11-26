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
/** How often a build whose live log is lost is re-checked for completion. */
const LOST_LOG_POLL_MS = 10000
/**
 * Sent (as FailedPrecondition) for a build marked as building that no instance streams: the controller holding
 * its in-memory log restarted, so retrying the stream is pointless (pkg/usecase/apiserver/app_build_service.go).
 */
const LOG_STREAM_LOST = 'build log stream no longer available'

/** The log of one finished build: `buildID` is fixed for the lifetime of the component. */
const FinishedBuildLog: Component<{ buildID: string }> = (props) => {
  const buildLog = createMemo(() => client.getBuildLog({ buildId: props.buildID }))

  let logRef!: HTMLElement
  // the memo has no reactive input, so this runs once: when the log is first displayed
  createEffect(buildLog, () => {
    logRef.scrollTop = logRef.scrollHeight
  })

  return <LogContainer innerHTML={toUTF8WithAnsi(buildLog().log)} ref={logRef} />
}

/** The live log of one unfinished build: `buildID` is fixed for the lifetime of the component. */
const StreamingBuildLog: Component<{ buildID: string; refetch: () => void }> = (props) => {
  const [log, setLog] = createSignal<Uint8Array>(new Uint8Array())
  const [lost, setLost] = createSignal(false)

  onSettled(() => {
    const buildID = props.buildID
    const abort = new AbortController()
    const refetch = () => {
      if (!abort.signal.aborted) props.refetch()
    }
    void (async () => {
      let started = false
      while (!abort.signal.aborted) {
        try {
          for await (const chunk of client.getBuildLogStream({ buildId: buildID }, { signal: abort.signal })) {
            if (!started) {
              started = true
              refetch() // the build left the queue: show its new status
            }
            setLog((prev) => concatBuffers(prev, chunk.log))
          }
        } catch (err) {
          if (err instanceof ConnectError && err.code === Code.Canceled) return
          if (err instanceof ConnectError && err.code === Code.FailedPrecondition) {
            if (err.rawMessage === LOG_STREAM_LOST) {
              // The live log is gone for good: only wait for the build to finish (the parent then shows its archived log).
              setLost(true)
              while (!abort.signal.aborted) {
                await sleep(LOST_LOG_POLL_MS)
                refetch()
              }
              return
            }
            // The build is still queued ("build not started yet"): keep watching it until it starts or is finished.
            await sleep(QUEUED_RETRY_MS)
            refetch()
            continue
          }
          console.trace(err)
          return
        }
        await sleep(1000)
        refetch() // stream ended: build is finished
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
    <Show
      when={!lost()}
      fallback={
        <div class="caption-regular text-text-grey">
          ビルドログを表示できません。ビルドを実行していたコントローラーが再起動したため、実行中のログは失われました。
        </div>
      }
    >
      <LogContainer
        innerHTML={toUTF8WithAnsi(log())}
        ref={logRef}
        onScroll={(e) => {
          atBottom = isScrolledToBottom(e.currentTarget)
        }}
      />
    </Show>
  )
}

export const BuildLog: Component<BuildLogProps> = (props) => {
  // Primitive memos: the props come from the build query, and reading it directly would re-run on every revalidation.
  const buildID = createMemo(() => props.buildID)
  const finished = createMemo(() => props.finished)
  // Keyed by build id: the build page is reused when moving to another build, and each log belongs to one build.
  return (
    <Show when={buildID()} keyed>
      {(id) => (
        <Show when={finished()} fallback={<StreamingBuildLog buildID={id} refetch={props.refetch} />}>
          <FinishedBuildLog buildID={id} />
        </Show>
      )}
    </Show>
  )
}
