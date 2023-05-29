import { createResource, Show } from 'solid-js'
import { styled } from '@macaron-css/solid'
import { useParams } from '@solidjs/router'
import {
  Card,
  CardItem,
  CardItemContent,
  CardItems,
  CardItemTitle,
  CardsContainer,
  CardTitle,
} from '/@/components/Card'
import { Header } from '/@/components/Header'
import { URLText } from '/@/components/URLText'
import { client } from '/@/libs/api'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { CenterInline, Container } from '/@/libs/layout'
import { vars } from '/@/theme'

// copy from AppTitleContainer in AppNav.tsx
const RepoTitleContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '14px',
    alignContent: 'center',

    marginTop: '48px',
    fontSize: '32px',
    fontWeight: 'bold',
    color: vars.text.black1,
  },
})

export default () => {
  const params = useParams()
  const [repo] = createResource(
    () => params.id,
    (id) => client.getRepository({ repositoryId: id }),
  )

  return (
    <Container>
      <Header />
      <Show when={repo()}>
        <RepoTitleContainer>
          <CenterInline>{providerToIcon(repositoryURLToProvider(repo().url), 36)}</CenterInline>
          {repo().name}
        </RepoTitleContainer>
        <CardsContainer>
          <Card>
            <CardTitle>Info</CardTitle>
            <CardItems>
              <CardItem>
                <CardItemTitle>ID</CardItemTitle>
                <CardItemContent>{repo().id}</CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>Name</CardItemTitle>
                <CardItemContent>{repo().name}</CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>URL</CardItemTitle>
                <CardItemContent>
                  <URLText href={repo().url} target='_blank' rel='noreferrer'>
                    {repo().url}
                  </URLText>
                </CardItemContent>
              </CardItem>
            </CardItems>
          </Card>
        </CardsContainer>
      </Show>
    </Container>
  )
}
