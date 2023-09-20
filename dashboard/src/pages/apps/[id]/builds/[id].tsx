import { Build_BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { AppNav } from '/@/components/AppNav'
import { ArtifactRow } from '/@/components/ArtifactRow'
import { BuildLog } from '/@/components/BuildLog'
import { BuildStatusIcon } from '/@/components/BuildStatusIcon'
import { Button } from '/@/components/Button'
import {
  Card,
  CardItem,
  CardItemContent,
  CardItemTitle,
  CardItems,
  CardRowsContainer,
  CardTitle,
  CardsRow,
} from '/@/components/Card'
import { Header } from '/@/components/Header'
import { client } from '/@/libs/api'
import { buildStatusStr } from '/@/libs/application'
import { DiffHuman, durationHuman, shortSha } from '/@/libs/format'
import { Container } from '/@/libs/layout'
import { styled } from '@macaron-css/solid'
import { useNavigate, useParams } from '@solidjs/router'
import { For, Ref, createEffect, createResource, createSignal, onCleanup } from 'solid-js'
import { Show } from 'solid-js'

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
        <AppNav repo={repo()} app={app()} />
        <CardRowsContainer>
          <CardsRow>
            <Card>
              <CardTitle>Actions</CardTitle>
              <Button
                color="black1"
                size="large"
                width="full"
                onclick={retryBuild}
                disabled={build().retriable}
                tooltip={build().retriable ? '既に再ビルドが行われています' : '同じコミットで再ビルドします'}
              >
                Retry build
              </Button>
              <Show when={build().status === Build_BuildStatus.BUILDING}>
                <Button color="black1" size="large" width="full" onclick={cancelBuild}>
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
              <BuildLog buildID={build().id} finished={buildFinished()} refetchBuild={refetchBuild} />
            </Card>
          </CardsRow>
        </CardRowsContainer>
      </Show>
    </Container>
  )
}
