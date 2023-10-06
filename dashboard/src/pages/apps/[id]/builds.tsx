import { BuildList } from '/@/components/BuildList'
import { client } from '/@/libs/api'
import { styled } from '@macaron-css/solid'
import { useParams } from '@solidjs/router'
import { createMemo, createResource } from 'solid-js'
import { Show } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    padding: '40px 32px 72px 32px',
    overflowY: 'auto',
  },
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
      <Show when={loaded()}>
        <BuildList builds={sortedBuilds()} showAppID={false} />
      </Show>
    </Container>
  )
}
