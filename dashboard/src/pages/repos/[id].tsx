import { createResource, For, JSX, Show } from 'solid-js'
import toast from 'solid-toast'
import { useNavigate, useParams } from '@solidjs/router'
import { GetApplicationsRequest_Scope } from '/@/api/neoshowcase/protobuf/gateway_pb'
import AppRow from '/@/components/AppRow'
import { Button } from '/@/components/Button'
import { Card, CardItem, CardItemContent, CardItems, CardItemTitle, CardsRow, CardTitle } from '/@/components/Card'
import { Header } from '/@/components/Header'
import { URLText } from '/@/components/URLText'
import { client, handleAPIError } from '/@/libs/api'
import { Container } from '/@/libs/layout'
import useModal from '/@/libs/useModal'
import { ModalButtonsContainer, ModalContainer, ModalText } from '/@/components/Modal'
import RepositoryNav from '/@/components/RepositoryNav'

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
      const allAppsRes = await client.getApplications({ scope: GetApplicationsRequest_Scope.ALL })
      return repo() ? allAppsRes.applications.filter((app) => app.repositoryId === repo()?.id) : []
    },
  )

  const { Modal: DeleteRepoModal, open: openDeleteRepoModal, close: closeDeleteRepoModal } = useModal()

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
      handleAPIError(e, 'リポジトリの削除に失敗しました')
    }
  }

  const handleCreateApplication = async () => {
    navigate(`/apps/new?repositoryID=${repo()?.id}`)
  }

  return (
    <Container>
      <Header />
      <Show when={repo()}>
        <RepositoryNav repository={repo()} />
        <CardsRow>
          <Card>
            <CardTitle>Actions</CardTitle>
            <CardItems>
              <Button
                onclick={handleCreateApplication}
                color='black1'
                size='large'
                width='full'
                tooltip='このリポジトリからアプリケーションを作成します'
              >
                Create New Application
              </Button>
              <Button
                onclick={openDeleteRepoModal}
                color='black1'
                size='large'
                width='full'
                disabled={!canDeleteRepository()}
                tooltip={
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
                  <Button onclick={closeDeleteRepoModal} color='black1' size='large' width='full'>
                    キャンセル
                  </Button>
                  <Button onclick={handleDeleteRepository} color='black1' size='large' width='full'>
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
          <Show when={apps()?.length > 0}>
            <Card
              style={{
                width: '100%',
              }}
            >
              <CardTitle>Apps</CardTitle>
              <div>
                <For each={apps()}>{(app) => <AppRow app={app} />}</For>
              </div>
            </Card>
          </Show>
        </CardsRow>
      </Show>
    </Container>
  )
}
