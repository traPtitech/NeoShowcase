import { useParams } from '@solidjs/router'
import { createEffect, createResource, createSignal, Ref } from 'solid-js'
import { client } from '/@/libs/api'
import { Container } from '/@/libs/layout'
import { Header } from '/@/components/Header'
import { Show } from 'solid-js'
import { AppNav } from '/@/components/AppNav'
import {
  Card,
  CardItem,
  CardItemContent,
  CardItems,
  CardItemTitle,
  CardsContainer,
  CardTitle,
} from '/@/components/Card'
import { Button } from '/@/components/Button'
import { DiffHuman, durationHuman, shortSha } from '/@/libs/format'
import { BuildStatusIcon } from '/@/components/BuildStatusIcon'
import { buildStatusStr } from '/@/libs/application'
import { concatBuffers, toUTF8WithAnsi } from '/@/libs/buffers'
import { sleep } from '/@/libs/sleep'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'

const LogContainer = styled('div', {
  base: {
    backgroundColor: vars.bg.black1,
    padding: '10px',
    color: vars.text.white1,
    borderRadius: '4px',

    whiteSpace: 'pre-wrap',
    overflowWrap: 'anywhere',
    maxHeight: '500px',
    overflowY: 'scroll',
  }
})

export default () => {
  const params = useParams()
  const [app] = createResource(
    () => params.id,
    (id) => client.getApplication({ id }),
  )
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const [build, { refetch: refetchBuild }] = createResource(
    () => params.buildID,
    (id) => client.getBuild({ buildId: id }),
  )
  const loaded = () => !!(app() && repo() && build())

  const buildFinished = () => build()?.finishedAt.valid
  const [buildLog] = createResource(
    () => buildFinished() && build()?.id,
    (id) => client.getBuildLog({ buildId: id }),
  )
  const [buildLogStream] = createResource(
    () => !buildFinished() && build()?.id,
    (id) => client.getBuildLogStream({ buildId: id }),
  )
  const [streamedLog, setStreamedLog] = createSignal(new Uint8Array())
  createEffect(() => {
    const stream = buildLogStream()
    if (!stream) return

    const iterate = async () => {
      for await (const log of stream) {
        setStreamedLog(prev => concatBuffers(prev, log.log))
      }
      await sleep(1000)
      await refetchBuild() // refetch build on stream end
    }
    void iterate()
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
    setTimeout(() => {
      ref.scrollTop = ref.scrollHeight
    })
  })

  const retryBuild = async () => {
    await client.retryCommitBuild({ applicationId: params.id, commit: build().commit })
    await refetchBuild()
  }

  return (
    <Container>
      <Header />
      <Show when={loaded()}>
        <AppNav appID={app().id} appName={app().name} repoName={repo().name} />
        <CardsContainer>
          <Card>
            <CardTitle>Actions</CardTitle>
            <Show when={!build().retriable}>
              <Button color='black1' size='large' onclick={retryBuild}>
                Retry build
              </Button>
            </Show>
          </Card>
          <Card>
            <CardTitle>Info</CardTitle>
            <CardItems>
              <CardItem>
                <CardItemTitle>ID</CardItemTitle>
                <CardItemContent>{build().id}</CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>Commit</CardItemTitle>
                <CardItemContent>{shortSha(build().commit)}</CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>Status</CardItemTitle>
                <CardItemContent>
                  <BuildStatusIcon state={build().status} size={24} />
                  {buildStatusStr[build().status]}
                </CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>Started at</CardItemTitle>
                <CardItemContent>
                  <Show when={build().startedAt.valid} fallback={'-'}>
                    <DiffHuman target={build().startedAt.timestamp.toDate()}></DiffHuman>
                  </Show>
                </CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>Finished at</CardItemTitle>
                <CardItemContent>
                  <Show when={build().finishedAt.valid} fallback={'-'}>
                    <DiffHuman target={build().finishedAt.timestamp.toDate()}></DiffHuman>
                  </Show>
                </CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>Duration</CardItemTitle>
                <CardItemContent>
                  <Show when={build().startedAt.valid && build().finishedAt.valid} fallback={'-'}>
                    {durationHuman(build().finishedAt.timestamp.toDate().getTime() - build().startedAt.timestamp.toDate().getTime())}
                  </Show>
                </CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>Retried</CardItemTitle>
                <CardItemContent>
                  {build().retriable ? 'Yes' : 'No'}
                </CardItemContent>
              </CardItem>
            </CardItems>
          </Card>
          <Card>
            <CardTitle>Build Log</CardTitle>
            <Show when={buildLog()}>
              <LogContainer innerHTML={toUTF8WithAnsi(buildLog().log)} ref={logRef} />
            </Show>
            <Show when={!buildLog() && buildLogStream()}>
              <LogContainer innerHTML={toUTF8WithAnsi(streamedLog())} ref={streamLogRef} />
            </Show>
          </Card>
        </CardsContainer>
      </Show>
    </Container>
  )
}
