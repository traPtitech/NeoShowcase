import { createResource, JSX, onCleanup, Show } from 'solid-js'
import toast from 'solid-toast'
import { ConnectError } from '@bufbuild/connect'
import { styled } from '@macaron-css/solid'
import { useNavigate, useParams } from '@solidjs/router'
import { Button } from '/@/components/Button'
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
import useModal from '/@/libs/useModal'
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

const ModalContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
  },
})
const ModalButtonsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '16px',
    justifyContent: 'center',
  },
})
const ModalText = styled('div', {
  base: {
    fontSize: '16px',
    fontWeight: 'bold',
    color: vars.text.black1,
    textAlign: 'center',
  },
})

export default () => {
  const navigate = useNavigate()
  const params = useParams()
  const [repo] = createResource(
    () => params.id,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const [apps] = createResource(
    () => params.id,
    async () => {
      const allAppsRes = await client.getApplications({})
      return repo() ? allAppsRes.applications.filter((app) => app.repositoryId === repo()?.id) : []
    },
  )

  const { Modal: DeleteRepoModal, open: openDeleteRepoModal, close: closeDeleteRepoModal } = useModal()

  onCleanup(() => {
    closeDeleteRepoModal()
  })

  // リポジトリに紐づくアプリケーションが存在するかどうか
  const canDeleteRepository = (): boolean => apps()?.length === 0

  // リポジトリの削除処理
  const handleDeleteRepository: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async () => {
    try {
      await client.deleteRepository({ repositoryId: repo()?.id })
      toast.success('リポジトリを削除しました')
      closeDeleteRepoModal()
      // アプリ一覧ページに遷移
      navigate('/apps')
    } catch (e) {
      console.error(e)
      // gRPCエラー
      if (e instanceof ConnectError) {
        toast.error('リポジトリの削除に失敗しました\n' + e.message)
      }
    }
  }

  return (
    <Container>
      <Header />
      <Show when={repo()}>
        <RepoTitleContainer>
          <CenterInline>{providerToIcon(repositoryURLToProvider(repo().url), 36)}</CenterInline>
          {repo()?.name}
        </RepoTitleContainer>
        <CardsContainer>
          <Card>
            <CardTitle>Actions</CardTitle>
            <CardItems>
              <Button
                onclick={() => {
                  openDeleteRepoModal()
                }}
                color='black1'
                size='large'
                disabled={!canDeleteRepository()}
                title={
                  canDeleteRepository()
                    ? 'リポジトリを削除します'
                    : 'リポジトリに紐づくアプリケーションが存在するため削除できません'
                }
              >
                Delete Repository
              </Button>
            </CardItems>
            <DeleteRepoModal>
              <ModalContainer>
                <ModalText>本当に削除しますか?</ModalText>
                <ModalButtonsContainer>
                  <Button
                    onclick={() => {
                      closeDeleteRepoModal()
                    }}
                    color='black1'
                    size='large'
                  >
                    キャンセル
                  </Button>
                  <Button onclick={handleDeleteRepository} color='black1' size='large'>
                    削除
                  </Button>
                </ModalButtonsContainer>
              </ModalContainer>
            </DeleteRepoModal>
          </Card>
          <Card>
            <CardTitle>Info</CardTitle>
            <CardItems>
              <CardItem>
                <CardItemTitle>ID</CardItemTitle>
                <CardItemContent>{repo()?.id}</CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>Name</CardItemTitle>
                <CardItemContent>{repo()?.name}</CardItemContent>
              </CardItem>
              <CardItem>
                <CardItemTitle>URL</CardItemTitle>
                <CardItemContent>
                  <URLText href={repo()?.url} target='_blank' rel='noreferrer'>
                    {repo()?.url}
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
