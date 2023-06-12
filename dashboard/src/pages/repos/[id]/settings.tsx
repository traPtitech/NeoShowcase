import { useParams } from '@solidjs/router'
import { Component, JSX, Show, Switch, createEffect, createMemo, createResource, createSignal } from 'solid-js'
import { client, handleAPIError } from '/@/libs/api'
import { Container } from '/@/libs/layout'
import { Header } from '/@/components/Header'
import { Button } from '/@/components/Button'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { createStore } from 'solid-js/store'
import {
  CreateRepositoryAuth,
  CreateRepositoryAuthBasic,
  CreateRepositoryAuthSSH,
  Repository_AuthMethod,
  UpdateRepositoryRequest,
  User,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import toast from 'solid-toast'
import { InputLabel } from '/@/components/Input'
import { InputBar } from '/@/components/Input'
import { FormTextBig } from '/@/components/AppsNew'
import { userFromId, users } from '/@/libs/useAllUsers'
import { UserSearch } from '/@/components/UserSearch'
import useModal from '/@/libs/useModal'
import RepositoryNav from '/@/components/RepositoryNav'
import { Empty } from '@bufbuild/protobuf'
import { Match } from 'solid-js'
import { Radio } from '/@/components/Radio'
import { extractRepositoryNameFromURL } from '/@/libs/application'

const ContentContainer = styled('div', {
  base: {
    marginTop: '24px',
    display: 'grid',
    gridTemplateColumns: '380px 1fr',
    gap: '40px',
    position: 'relative',
  },
})
const SidebarContainer = styled('div', {
  base: {
    position: 'sticky',
    top: '64px',
    padding: '24px 40px',
    backgroundColor: vars.bg.white1,
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
  },
})
const SidebarOptions = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '12px',

    fontSize: '20px',
    color: vars.text.black1,
  },
})
const SidebarNavAnchor = styled('a', {
  base: {
    color: vars.text.black2,
    textDecoration: 'none',
    selectors: {
      '&:hover': {
        color: vars.text.black1,
      },
    },
  },
})
const ConfigsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
  },
})
const SettingFieldSet = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
    padding: '24px',
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    background: vars.bg.white1,
  },
})
const SshDetails = styled('div', {
  base: {
    color: vars.text.black2,
    marginBottom: '4px',
  },
})
const PublicKeyCode = styled('code', {
  base: {
    display: 'block',
    padding: '8px 12px',
    fontFamily: 'monospace',
    fontSize: '14px',
    background: vars.bg.white2,
    color: vars.text.black1,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
  },
})

export default () => {
  const params = useParams()
  const [repo, { refetch: refetchRepo }] = createResource(
    () => params.id,
    (repositoryId) => client.getRepository({ repositoryId }),
  )

  const GeneralConfigsContainer: Component = () => {
    let formContainer: HTMLFormElement
    // 現在の設定で初期化
    const [generalConfig, setGeneralConfig] = createStore({
      name: repo().name,
      url: repo().url,
    })

    const [updateAuthConfig, setUpdateAuthConfig] = createSignal(false)

    // 認証方法 ("none" | "ssh" | "basic")
    type AuthMethod = CreateRepositoryAuth['auth']['case']
    const [authMethod, setAuthMethod] = createSignal<AuthMethod>('none')

    const [systemPublicKey] = createResource(() => client.getSystemPublicKey({}))
    const [useTmpKey, setUseTmpKey] = createSignal(false)
    const [tmpKey] = createResource(
      () => (useTmpKey() ? true : undefined),
      () => client.generateKeyPair({}),
    )
    createEffect(() => {
      if (!tmpKey()) return
      setAuthConfig('ssh', 'value', 'keyId', tmpKey().keyId)
    })
    const publicKey = () => (useTmpKey() ? tmpKey()?.publicKey : systemPublicKey()?.publicKey)

    // 認証情報
    // 認証方法の切り替え時に情報を保持するために、storeを使用して3種類の認証情報を保持する
    const [authConfig, setAuthConfig] = createStore<{
      [K in AuthMethod]: Extract<CreateRepositoryAuth['auth'], { case: K }>
    }>({
      none: {
        case: 'none',
        value: new Empty(),
      },
      basic: {
        case: 'basic',
        value: new CreateRepositoryAuthBasic(),
      },
      ssh: {
        case: 'ssh',
        value: new CreateRepositoryAuthSSH(),
      },
    })

    // URLからリポジトリ名を自動入力
    createEffect(() => {
      const repositoryName = extractRepositoryNameFromURL(generalConfig.url)
      setGeneralConfig('name', repositoryName)
    })

    const updateGeneralSettings: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
      // prevent default form submit (reload page)
      e.preventDefault()

      // validate form
      if (!formContainer.reportValidity()) {
        return
      }

      const updateRepositoryRequest = new UpdateRepositoryRequest({
        id: repo().id,
        name: generalConfig.name,
        url: generalConfig.url,
        auth: updateAuthConfig()
          ? {
              auth: authConfig[authMethod()],
            }
          : undefined,
      })

      try {
        await client.updateRepository(updateRepositoryRequest)
        toast.success('リポジトリ設定を更新しました')
        refetchRepo()
      } catch (e) {
        handleAPIError(e, 'リポジトリ設定の更新に失敗しました')
      }
    }

    return (
      <form ref={formContainer}>
        <SettingFieldSet>
          <FormTextBig id='general-settings'>General settings</FormTextBig>
          <div>
            <InputLabel>Repository Name</InputLabel>
            <InputBar
              placeholder='my-app'
              value={generalConfig.name}
              onChange={(e) => setGeneralConfig('name', e.currentTarget.value)}
              required
            />
          </div>
          <div>
            <InputLabel>Repository URL</InputLabel>
            <InputBar
              // SSH URLはURLとしては不正なのでtypeを変更
              type={repo().authMethod === Repository_AuthMethod.SSH ? 'text' : 'url'}
              placeholder='https://example.com/my-app.git'
              value={generalConfig.url}
              onChange={(e) => setGeneralConfig('url', e.currentTarget.value)}
              required
            />
          </div>
          <Show
            when={!updateAuthConfig()}
            fallback={
              <Button color='black1' size='large' width='auto' onclick={() => setUpdateAuthConfig(false)}>
                認証方法の更新をキャンセルする
              </Button>
            }
          >
            <Button color='black1' size='large' width='auto' onclick={() => setUpdateAuthConfig(true)}>
              認証方法を更新する
            </Button>
          </Show>
          <Show when={updateAuthConfig()}>
            <div>
              <InputLabel>認証方法</InputLabel>
              <Radio
                items={[
                  { title: '認証を使用しない', value: 'none' },
                  { title: 'Basic認証を使用', value: 'basic' },
                  { title: 'SSH認証を使用', value: 'ssh' },
                ]}
                selected={authMethod()}
                setSelected={setAuthMethod}
              />
              <Switch>
                <Match when={authMethod() === 'basic'}>
                  <InputLabel>ユーザー名</InputLabel>
                  <InputBar
                    // SSH URLはURLとしては不正なのでtypeを変更
                    value={authConfig.basic.value.username}
                    onInput={(e) => setAuthConfig('basic', 'value', 'username', e.currentTarget.value)}
                  />
                  <InputLabel>パスワード</InputLabel>
                  <InputBar
                    // SSH URLはURLとしては不正なのでtypeを変更
                    type='password'
                    value={authConfig.basic.value.password}
                    onInput={(e) => setAuthConfig('basic', 'value', 'password', e.currentTarget.value)}
                  />
                </Match>
                <Match when={authMethod() === 'ssh'}>
                  <SshDetails>
                    以下のSSH公開鍵{!useTmpKey() && ' (システムデフォルト) '}をリポジトリに登録してください。
                  </SshDetails>
                  <PublicKeyCode>{publicKey()}</PublicKeyCode>
                  <Show when={!useTmpKey()}>
                    <Button color='black1' size='large' width='auto' onclick={() => setUseTmpKey(true)} type='submit'>
                      新たなSSH鍵を生成する (for github.com)
                    </Button>
                  </Show>
                </Match>
              </Switch>
            </div>
          </Show>
          <Button color='black1' size='large' width='auto' onclick={updateGeneralSettings} type='submit'>
            Save
          </Button>
        </SettingFieldSet>
      </form>
    )
  }

  const OwnerConfigContainer: Component = () => {
    const { Modal, open } = useModal()

    const nonOwnerUsers = createMemo(() => {
      return users()?.filter((user) => !repo().ownerIds.includes(user.id)) ?? []
    })

    const handleAddOwner = async (user: User) => {
      const updateApplicationRequest = new UpdateRepositoryRequest({
        id: repo().id,
        ownerIds: repo().ownerIds.concat(user.id),
      })

      try {
        await client.updateRepository(updateApplicationRequest)
        toast.success('リポジトリオーナーを追加しました')
        refetchRepo()
      } catch (e) {
        handleAPIError(e, 'リポジトリオーナーの追加に失敗しました')
      }
    }
    const handleDeleteOwner = async (owner: User) => {
      const updateApplicationRequest = new UpdateRepositoryRequest({
        id: repo().id,
        ownerIds: repo().ownerIds.filter((id) => id !== owner.id),
      })

      try {
        await client.updateRepository(updateApplicationRequest)
        toast.success('リポジトリのオーナーを削除しました')
        refetchRepo()
      } catch (e) {
        handleAPIError(e, 'リポジトリのオーナーの削除に失敗しました')
      }
    }

    return (
      <>
        <SettingFieldSet>
          <FormTextBig id='owner-settings'>Owner Settings</FormTextBig>
          <Button color='black1' size='large' width='auto' onclick={open}>
            リポジトリオーナーを追加する
          </Button>
          <UserSearch users={repo().ownerIds.map((userId) => userFromId(userId))}>
            {(user) => (
              <Button
                color='black1'
                size='large'
                width='auto'
                onclick={() => {
                  handleDeleteOwner(user)
                }}
              >
                削除
              </Button>
            )}
          </UserSearch>
        </SettingFieldSet>
        <Modal>
          <UserSearch users={nonOwnerUsers()}>
            {(user) => (
              <Button
                color='black1'
                size='large'
                width='auto'
                onclick={() => {
                  handleAddOwner(user)
                }}
              >
                追加
              </Button>
            )}
          </UserSearch>
        </Modal>
      </>
    )
  }

  return (
    <Container>
      <Header />
      <Show when={repo()}>
        <RepositoryNav repository={repo()} />
        <ContentContainer>
          <div>
            <SidebarContainer>
              <SidebarOptions>
                <SidebarNavAnchor href='#general-settings'>General</SidebarNavAnchor>
                <SidebarNavAnchor href='#owner-settings'>Owner</SidebarNavAnchor>
              </SidebarOptions>
            </SidebarContainer>
          </div>
          <ConfigsContainer>
            <GeneralConfigsContainer />
            <OwnerConfigContainer />
          </ConfigsContainer>
        </ContentContainer>
      </Show>
    </Container>
  )
}
