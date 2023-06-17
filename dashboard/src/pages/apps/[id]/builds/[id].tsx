import { useNavigate, useParams } from '@solidjs/router'
import { createEffect, createResource, createSignal, For, onCleanup, Ref } from 'solid-js'
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
  CardRowsContainer,
  CardsRow,
  CardTitle,
} from '/@/components/Card'
import { Button } from '/@/components/Button'
import { DiffHuman, durationHuman, shortSha } from '/@/libs/format'
import { BuildStatusIcon } from '/@/components/BuildStatusIcon'
import { buildStatusStr } from '/@/libs/application'
import { concatBuffers, toUTF8WithAnsi } from '/@/libs/buffers'
import { sleep } from '/@/libs/sleep'
import { LogContainer } from '/@/components/Log'
import { Code, ConnectError } from '@bufbuild/connect'
import { Build_BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '@macaron-css/solid'
import { ArtifactRow } from '/@/components/ArtifactRow'

const ArtifactsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
  },
})

export default () => {
  const navigate = useNavigate()
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
  const logStreamAbort = new AbortController()
  const [buildLogStream] = createResource(
    () => !buildFinished() && build()?.id,
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
      await refetchBuild() // refetch build on stream end
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
    setTimeout(() => {
      ref.scrollTop = ref.scrollHeight
    })
  })

  const retryBuild = async () => {
    await client.retryCommitBuild({ applicationId: params.id, commit: build().commit })
    navigate(`/apps/${app().id}/builds`)
  }

  const cancelBuild = async () => {
    await client.cancelBuild({ buildId: build().id })
    await refetchBuild()
  }

  return (
    <Container>
      <Header />
      <Show when={loaded()}>
        <AppNav appID={app().id} appName={app().name} repoName={repo().name} />
        <CardRowsContainer>
          <CardsRow>
            <Card>
              <CardTitle>Actions</CardTitle>
              <Button
                color='black1'
                size='large'
                width='full'
                onclick={retryBuild}
                disabled={build().retriable}
                title={build().retriable ? '既に再ビルドが行われています' : '同じコミットで再ビルドします'}
              >
                Retry build
              </Button>
              <Show when={build().status === Build_BuildStatus.BUILDING}>
                <Button color='black1' size='large' width='full' onclick={cancelBuild}>
                  Cancel build
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
                  <CardItemTitle>Queued at</CardItemTitle>
                  <CardItemContent>
                    <DiffHuman target={build().queuedAt.toDate()} />
                  </CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>Started at</CardItemTitle>
                  <CardItemContent>
                    <Show when={build().startedAt.valid} fallback={'-'}>
                      <DiffHuman target={build().startedAt.timestamp.toDate()} />
                    </Show>
                  </CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>Finished at</CardItemTitle>
                  <CardItemContent>
                    <Show when={build().finishedAt.valid} fallback={'-'}>
                      <DiffHuman target={build().finishedAt.timestamp.toDate()} />
                    </Show>
                  </CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>Duration</CardItemTitle>
                  <CardItemContent>
                    <Show when={build().startedAt.valid && build().finishedAt.valid} fallback={'-'}>
                      {durationHuman(
                        build().finishedAt.timestamp.toDate().getTime() -
                          build().startedAt.timestamp.toDate().getTime(),
                      )}
                    </Show>
                  </CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>Retried</CardItemTitle>
                  <CardItemContent>{build().retriable ? 'Yes' : 'No'}</CardItemContent>
                </CardItem>
              </CardItems>
            </Card>
          </CardsRow>
          <Show when={build().artifacts.length > 0}>
            <CardsRow>
              <Card>
                <CardTitle>Artifacts</CardTitle>
                <ArtifactsContainer>
                  <For each={build().artifacts || []}>{(artifact) => <ArtifactRow artifact={artifact} />}</For>
                </ArtifactsContainer>
              </Card>
            </CardsRow>
          </Show>
          <CardsRow>
            <Card>
              <CardTitle>Build Log</CardTitle>
              <Show when={buildLog()}>
                <LogContainer innerHTML={toUTF8WithAnsi(buildLog().log)} ref={logRef} />
              </Show>
              <Show when={!buildLog() && buildLogStream()}>
                <LogContainer innerHTML={toUTF8WithAnsi(streamedLog())} ref={streamLogRef} />
              </Show>
            </Card>
          </CardsRow>
        </CardRowsContainer>
      </Show>
    </Container>
  )
}
