import Fuse from 'fuse.js'
import { Component, createEffect, createResource, createSignal, For, JSX, onCleanup, Show } from 'solid-js'
import toast from 'solid-toast'
import { ConnectError } from '@bufbuild/connect'
import { styled } from '@macaron-css/solid'
import { useNavigate, useParams } from '@solidjs/router'
import { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import AppRow from '/@/components/AppRow'
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
import useAllUsers from '/@/libs/useAllUsers'
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

const UserContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
})
const UserRowLeft = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
    alignItems: 'center',
  },
})
const UserAvatar = styled('img', {
  base: {
    width: '32px',
    height: '32px',
    borderRadius: '50%',
  },
})
const UserName = styled('div', {
  base: {
    fontSize: '16px',
    color: vars.text.black1,
  },
})

const OwnerEditorContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
    height: '480px',
  },
})

const OwnerEditorUserList = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
    overflowY: 'auto',
  },
})

export default () => {
  const navigate = useNavigate()
  const params = useParams()
  const [repo, { refetch: refetchRepository }] = createResource(
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
  const [users] = useAllUsers()

  const userFromId = (userId: string): User | undefined => {
    return users()?.find((user) => user.id === userId)
  }

  const { Modal: DeleteRepoModal, open: openDeleteRepoModal, close: closeDeleteRepoModal } = useModal()

  onCleanup(closeDeleteRepoModal)

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

  const handleCreateApplication: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async () => {
    navigate(`/apps/new?repositoryID=${repo()?.id}`)
  }

  // リポジトリのオーナー編集処理
  const { Modal: EditOwnerModal, open: openEditOwnerModal } = useModal()

  const [userSearchQuery, setUserSearchQuery] = createSignal('')
  const [userSearchResults, setUserSearchResults] = createSignal<User[]>([])

  // ユーザー検索
  // - users()の更新時にFuseインスタンスを再生成する
  // - userSearchQuery()の更新時に検索を実行する
  // ため二重にcreateEffectを使用している
  createEffect(() => {
    const fuse = new Fuse(users() ?? [], {
      keys: ['name'],
    })

    createEffect(() => {
      // 検索クエリが空の場合は全ユーザーを表示する
      if (userSearchQuery() === '') {
        setUserSearchResults(users() ?? [])
      } else {
        setUserSearchResults(fuse.search(userSearchQuery()).map((result) => result.item))
      }
    })
  })

  const handleAddOwner = async (user: User): Promise<void> => {
    try {
      await client.updateRepository({
        id: repo()?.id,
        ownerIds: repo()?.ownerIds.concat(user.id),
      })
      refetchRepository()
      toast.success('リポジトリのオーナーを追加しました')
    } catch (e) {
      // gRPCエラー
      if (e instanceof ConnectError) {
        toast.error('リポジトリのオーナーの追加に失敗しました\n' + e.message)
      }
    }
  }

  const handleDeleteOwner = async (user: User): Promise<void> => {
    try {
      await client.updateRepository({
        id: repo()?.id,
        ownerIds: repo()?.ownerIds.filter((id) => id !== user.id),
      })
      refetchRepository()
      toast.success('リポジトリのオーナーを削除しました')
    } catch (e) {
      console.error(e)
      // gRPCエラー
      if (e instanceof ConnectError) {
        toast.error('リポジトリのオーナーの削除に失敗しました\n' + e.message)
      }
    }
  }

  const OwnerSuggestions: Component<{
    users: User[]
  }> = (props) => {
    return (
      <For each={props.users}>
        {(user) => {
          return (
            <UserContainer>
              <UserRowLeft>
                <UserAvatar src={user.avatarUrl} />
                <UserName>{user.name}</UserName>
              </UserRowLeft>
              <Button
                color='black1'
                size='large'
                onclick={() => {
                  handleAddOwner(user)
                }}
              >
                追加
              </Button>
            </UserContainer>
          )
        }}
      </For>
    )
  }

  const OwnerRow: Component<{
    user: User
  }> = (props) => {
    return (
      <UserContainer>
        <UserRowLeft>
          <UserAvatar src={props.user.avatarUrl} />
          <UserName>{props.user.name}</UserName>
        </UserRowLeft>
        <Button
          color='black1'
          size='large'
          onclick={() => {
            handleDeleteOwner(props.user)
          }}
        >
          削除
        </Button>
      </UserContainer>
    )
  }

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
            <CardTitle>Actions</CardTitle>
            <CardItems>
              <Button
                onclick={handleCreateApplication}
                color='black1'
                size='large'
                title='このリポジトリからアプリケーションを作成します'
              >
                Create New Application
              </Button>
              <Button
                onclick={openDeleteRepoModal}
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
                  <Button onclick={closeDeleteRepoModal} color='black1' size='large'>
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
          <Show when={users()}>
            <Card>
              <CardTitle>Owners</CardTitle>
              <Button onclick={openEditOwnerModal} color='black1' size='large'>
                リポジトリ所有者を追加する
              </Button>
              <EditOwnerModal>
                <OwnerEditorContainer>
                  ユーザーを検索して追加
                  <input
                    type="text"
                    value={userSearchQuery()}
                    placeholder='ユーザー名'
                    oninput={(e) => {
                      setUserSearchQuery(e.currentTarget.value)
                    }}
                  />
                  <OwnerEditorUserList>
                    <OwnerSuggestions users={userSearchResults()} />
                  </OwnerEditorUserList>
                </OwnerEditorContainer>
              </EditOwnerModal>
              <div>
                <For each={repo().ownerIds}>
                  {(ownerId) => {
                    const user = userFromId(ownerId)
                    return <OwnerRow user={user} />
                  }}
                </For>
              </div>
            </Card>
          </Show>
          <Show when={apps()}>
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
        </CardsContainer>
      </Show>
    </Container>
  )
}
