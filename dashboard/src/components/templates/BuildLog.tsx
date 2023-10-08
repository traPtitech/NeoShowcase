import { client } from '/@/libs/api'
import { concatBuffers, toUTF8WithAnsi } from '/@/libs/buffers'
import { isScrolledToBottom } from '/@/libs/scroll'
import { sleep } from '/@/libs/sleep'
import { Code, ConnectError } from '@bufbuild/connect'
import { Component, Ref, Show, createEffect, createResource, createSignal, onCleanup } from 'solid-js'
import { LogContainer } from '../UI/LogContainer'

export interface BuildLogProps {
  buildID: string
  finished: boolean
  refetchBuild: () => void
}

export const BuildLog: Component<BuildLogProps> = (props) => {
  const [buildLog] = createResource(
    () => props.finished && props.buildID,
    (id) => client.getBuildLog({ buildId: id }),
  )
  const logStreamAbort = new AbortController()
  const [buildLogStream] = createResource(
    () => !props.finished && props.buildID,
    (id) => client.getBuildLogStream({ buildId: id }, { signal: logStreamAbort.signal }),
  )
  const [streamedLog, setStreamedLog] = createSignal(new Uint8Array())
  createEffect(() => {
    const stream = buildLogStream()
    if (!stream) return

    const iterate = async () => {
      try {
        for await (const log of stream) {
          setStreamedLog((prev) => concatBuffers(prev, log.log))
        }
      } catch (err) {
        // ignore abort error
        const isAbortErr = err instanceof ConnectError && err.code === Code.Canceled
        if (!isAbortErr) {
          console.trace(err)
          return
        }
      }
      await sleep(1000)
      props.refetchBuild() // refetch build on stream end
    }
    void iterate()
  })
  onCleanup(() => {
    logStreamAbort.abort()
  })

  let logRef: Ref<HTMLDivElement>
  let streamLogRef: Ref<HTMLDivElement>
  createEffect(() => {
    if (!buildLog()) return
    const ref = logRef as HTMLDivElement
    if (!ref) return
    setTimeout(() => {
      ref.scrollTop = ref.scrollHeight
    })
  })
  createEffect(() => {
    if (!streamedLog()) return
    const ref = streamLogRef as HTMLDivElement
    if (!ref) return
    if (atBottom()) {
      ref.scrollTop = ref.scrollHeight
    }
  })

  const [atBottom, setAtBottom] = createSignal(true)
  const onScroll = (e: { target: Element }) => setAtBottom(isScrolledToBottom(e.target))

  return (
    <>
      <Show when={buildLog()}>
        <LogContainer innerHTML={toUTF8WithAnsi(buildLog().log)} ref={logRef} />
      </Show>
      <Show when={!buildLog() && buildLogStream()}>
        <LogContainer innerHTML={toUTF8WithAnsi(streamedLog())} ref={streamLogRef} onScroll={onScroll} />
      </Show>
    </>
  )
}
