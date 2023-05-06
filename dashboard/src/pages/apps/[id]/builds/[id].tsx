import { useParams } from '@solidjs/router'
import { createResource } from 'solid-js'
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
        </CardsContainer>
      </Show>
    </Container>
  )
}
