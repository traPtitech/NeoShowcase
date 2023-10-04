import {
  CreateRepositoryAuth,
  Repository_AuthMethod,
  UpdateRepositoryRequest,
  User,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import GeneralIcon from '/@/assets/icons/24/browse_activity.svg'
import AuthIcon from '/@/assets/icons/24/conversion_path.svg'
import OwnerIcon from '/@/assets/icons/24/person.svg'
import { FormTextBig } from '/@/components/AppsNew'
import { InfoTooltip } from '/@/components/InfoTooltip'
import { InputBar, InputLabel } from '/@/components/Input'
import { RepositoryAuthSettings } from '/@/components/RepositoryAuthSettings'
import { Button } from '/@/components/UI/Button'
import { UserSearch } from '/@/components/UserSearch'
import { client, handleAPIError } from '/@/libs/api'
import { extractRepositoryNameFromURL } from '/@/libs/application'
import { userFromId, users } from '/@/libs/useAllUsers'
import useModal from '/@/libs/useModal'
import { vars } from '/@/theme'
import { PartialMessage, PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Outlet, useMatch, useNavigate, useParams } from '@solidjs/router'
import { Component, JSX, Show, createEffect, createMemo, createResource, createSignal } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    padding: '40px 32px',
    overflowY: 'auto',
  },
})
const MainView = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    maxWidth: '1000px',
    margin: '0 auto',
    display: 'grid',
    gridTemplateColumns: '235px 1fr',
    gap: '48px',
  },
})
const SideMenu = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
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

export default () => {
  const params = useParams()
  const [repo, { refetch: refetchRepo }] = createResource(
    () => params.id,
    (repositoryId) => client.getRepository({ repositoryId }),
  )
  const loaded = () => !!(users() && repo())
  const navigator = useNavigate()
  const matchGeneralPage = useMatch(() => `/repos/${repo()?.id}/settings/`)
  const matchAuthPage = useMatch(() => `/repos/${repo()?.id}/settings/authorization`)
  const matchOwnerPage = useMatch(() => `/repos/${repo()?.id}/settings/owner`)

  const update = async (req: PlainMessage<UpdateRepositoryRequest>) => {
    try {
      await client.updateRepository(req)
      toast.success('リポジトリ設定を更新しました')
      refetchRepo()
    } catch (e) {
      handleAPIError(e, 'リポジトリ設定の更新に失敗しました')
    }
  }

  const GeneralConfigsContainer: Component = () => {
    let formContainer: HTMLFormElement

    // 現在の設定で初期化
    const [generalConfig, setGeneralConfig] = createStore<PlainMessage<UpdateRepositoryRequest>>({
      id: repo().id,
      name: repo().name,
      url: repo().url,
    })
    const [updateAuthConfig, setUpdateAuthConfig] = createSignal(false)
    const mapAuthMethod = (authMethod: Repository_AuthMethod): PlainMessage<CreateRepositoryAuth>['auth'] => {
      switch (authMethod) {
        case Repository_AuthMethod.NONE:
          return { case: 'none', value: {} }
        case Repository_AuthMethod.BASIC:
          return { case: 'basic', value: { username: '', password: '' } }
        case Repository_AuthMethod.SSH:
          return { case: 'ssh', value: { keyId: '' } }
      }
    }
    const [authConfig, setAuthConfig] = createStore<PlainMessage<CreateRepositoryAuth>>({
      auth: mapAuthMethod(repo().authMethod),
    })

    // URLからリポジトリ名を自動入力
    createEffect(() => {
      const repositoryName = extractRepositoryNameFromURL(generalConfig.url)
      setGeneralConfig('name', repositoryName)
    })

    const onClickSave: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
      e.preventDefault()
      if (!formContainer.reportValidity()) {
        return
      }
      return update({ ...generalConfig, auth: authConfig })
    }

    return (
      <form ref={formContainer}>
        <SettingFieldSet>
          <FormTextBig id="general-settings">General Settings</FormTextBig>
          <div>
            <InputLabel>Repository Name</InputLabel>
            <InputBar
              placeholder="my-app"
              value={generalConfig.name}
              onChange={(e) => setGeneralConfig('name', e.currentTarget.value)}
              required
            />
          </div>
          <div>
            <InputLabel>Repository URL</InputLabel>
            <InputBar
              // SSH URLはURLとしては不正なのでtypeを変更
              type={authConfig.auth.case === 'ssh' ? 'text' : 'url'}
              placeholder="https://example.com/my-app.git"
              value={generalConfig.url}
              onChange={(e) => setGeneralConfig('url', e.currentTarget.value)}
              required
            />
          </div>
          <Show
            when={!updateAuthConfig()}
            fallback={
              <Button color="black1" size="large" width="auto" onclick={() => setUpdateAuthConfig(false)}>
                認証方法の更新をキャンセルする
              </Button>
            }
          >
            <Button color="black1" size="large" width="auto" onclick={() => setUpdateAuthConfig(true)}>
              認証方法を更新する
            </Button>
          </Show>
          <Show when={updateAuthConfig()}>
            <RepositoryAuthSettings authConfig={authConfig} setAuthConfig={setAuthConfig} />
          </Show>
          <Button color="black1" size="large" width="auto" onclick={onClickSave} type="submit">
            Save
          </Button>
        </SettingFieldSet>
      </form>
    )
  }

  const OwnerConfigContainer: Component = () => {
    const { Modal, open } = useModal()

    const nonOwnerUsers = createMemo(() => {
      return users().filter((user) => !repo().ownerIds.includes(user.id)) ?? []
    })

    const handleAddOwner = async (user: User) => {
      const updateApplicationRequest: PartialMessage<UpdateRepositoryRequest> = {
        id: repo().id,
        ownerIds: {
          ownerIds: repo().ownerIds.concat(user.id),
        },
      }

      try {
        await client.updateRepository(updateApplicationRequest)
        toast.success('リポジトリオーナーを追加しました')
        refetchRepo()
      } catch (e) {
        handleAPIError(e, 'リポジトリオーナーの追加に失敗しました')
      }
    }
    const handleDeleteOwner = async (owner: User) => {
      const updateApplicationRequest: PartialMessage<UpdateRepositoryRequest> = {
        id: repo().id,
        ownerIds: {
          ownerIds: repo().ownerIds.filter((id) => id !== owner.id),
        },
      }

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
          <FormTextBig id="owner-settings">
            Owners
            <InfoTooltip
              tooltip={['オーナーは以下が可能になります', 'リポジトリの設定を変更']}
              style="bullets-with-title"
            />
          </FormTextBig>
          <Button color="black1" size="large" width="auto" onclick={open}>
            リポジトリオーナーを追加する
          </Button>
          <UserSearch users={repo().ownerIds.map((userId) => userFromId(userId))}>
            {(user) => (
              <Button
                color="black1"
                size="large"
                width="auto"
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
                color="black1"
                size="large"
                width="auto"
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
      <Show when={loaded()}>
        <MainView>
          <SideMenu>
            <Button
              color="text"
              size="medium"
              full
              active={!!matchGeneralPage()}
              onclick={() => {
                navigator(`/repos/${repo()?.id}/settings/`)
              }}
              leftIcon={<GeneralIcon />}
            >
              General
            </Button>
            <Button
              color="text"
              size="medium"
              full
              active={!!matchAuthPage()}
              onclick={() => {
                navigator(`/repos/${repo()?.id}/settings/authorization`)
              }}
              leftIcon={<AuthIcon />}
            >
              Authorization
            </Button>
            <Button
              color="text"
              size="medium"
              full
              active={!!matchOwnerPage()}
              onclick={() => {
                navigator(`/repos/${repo()?.id}/settings/owner`)
              }}
              leftIcon={<OwnerIcon />}
            >
              Owner
            </Button>
          </SideMenu>
          <Outlet />
        </MainView>
      </Show>
    </Container>
  )
}
