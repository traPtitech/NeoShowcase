import { createMemo, createResource } from 'solid-js'
import { client } from '/@/libs/api'
import { useParams } from '@solidjs/router'
import { Container } from '/@/libs/layout'
import { Header } from '/@/components/Header'
import { Show } from 'solid-js'
import { AppNav } from '/@/components/AppNav'
import { BuildList } from '/@/components/BuildList'

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
  const [builds] = createResource(
    () => params.id,
    (id) => client.getBuilds({ id }),
  )
  const loaded = () => !!(app() && repo() && builds())

  const sortedBuilds = createMemo(
    () =>
      builds() &&
      [...builds().builds].sort((b1, b2) => {
        return b2.queuedAt.toDate().getTime() - b1.queuedAt.toDate().getTime()
      }),
  )

  return (
    <Container>
      <Header />
      <Show when={loaded()}>
        <AppNav repo={repo()} app={app()} />
        <BuildList builds={sortedBuilds()} showAppID={false} />
      </Show>
    </Container>
  )
}
