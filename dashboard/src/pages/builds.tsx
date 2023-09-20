import { BuildList } from '/@/components/BuildList'
import { Header } from '/@/components/Header'
import { client } from '/@/libs/api'
import { Container, PageTitle } from '/@/libs/layout'
import { styled } from '@macaron-css/solid'
import { createResource } from 'solid-js'
import { Show } from 'solid-js'

const PageContainer = styled('div', {
  base: {
    paddingTop: '24px',
  },
})

export default () => {
  const [builds] = createResource(() => client.getAllBuilds({ page: 0, limit: 100 }))

  return (
    <Container>
      <Header />
      <PageTitle>Build Queue</PageTitle>
      <Show when={builds()}>
        <PageContainer>
          <BuildList builds={builds().builds} showAppID={false} />
        </PageContainer>
      </Show>
    </Container>
  )
}
