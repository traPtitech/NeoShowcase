import { LogContainer } from '/@/components/Log'
import { Component, createEffect, createResource, createSignal, For, onCleanup, Ref } from 'solid-js'
import { toWithAnsi } from '/@/libs/buffers'
import { client } from '/@/libs/api'
import { Code, ConnectError } from '@bufbuild/connect'
import { sleep } from '/@/libs/sleep'
import { isScrolledToBottom } from '/@/libs/scroll'

export interface ContainerLogProps {
  appID: string
}

export const ContainerLog: Component<ContainerLogProps> = (props) => {
  const logStreamAbort = new AbortController()
  const [logStream] = createResource(
    () => props.appID,
    (id) => client.getOutputStream({ id }, { signal: logStreamAbort.signal }),
  )
  const [streamedLog, setStreamedLog] = createSignal<string[]>([])
  createEffect(() => {
    const stream = logStream()
    if (!stream) {
      setStreamedLog([])
      return
    }

    const iterate = async () => {
      try {
        for await (const log of stream) {
          setStreamedLog((prev) => prev.concat(log.log))
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
    }
    void iterate()
  })
  onCleanup(() => {
    logStreamAbort.abort()
  })

  let logRef: Ref<HTMLDivElement>
  createEffect(() => {
    if (!streamedLog()) return
    const ref = logRef as HTMLDivElement
    if (!ref) return
    if (atBottom()) {
      ref.scrollTop = ref.scrollHeight
    }
  })

  const [atBottom, setAtBottom] = createSignal(true)
  const onScroll = (e: { target: Element }) => setAtBottom(isScrolledToBottom(e.target))

  return (
    <LogContainer ref={logRef} overflowX='scroll' onScroll={onScroll}>
      <For each={streamedLog()}>{(line) => <code innerHTML={toWithAnsi(line)} />}</For>
    </LogContainer>
  )
}
