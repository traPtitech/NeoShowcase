import { useParams } from '@solidjs/router'
import { Component, createEffect, createMemo, createResource, createSignal, For, JSX, onMount, Show } from 'solid-js'
import { client, handleAPIError } from '/@/libs/api'
import { Container } from '/@/libs/layout'
import { Header } from '/@/components/Header'
import { AppNav } from '/@/components/AppNav'
import { Button } from '/@/components/Button'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { createStore } from 'solid-js/store'
import {
  ApplicationEnvVar,
  CreateWebsiteRequest,
  DeployType,
  PortPublication,
  RuntimeConfig,
  User,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { BuildConfig, BuildConfigMethod, BuildConfigs } from '/@/components/BuildConfigs'
import toast from 'solid-toast'
import { WebsiteSettings } from '/@/components/WebsiteSettings'
import { InputBar, InputLabel } from '/@/components/Input'
import { FormTextBig } from '/@/components/AppsNew'
import { PortPublicationSettings } from '/@/components/PortPublications'
import { userFromId, users } from '/@/libs/useAllUsers'
import { UserSearch } from '/@/components/UserSearch'
import useModal from '/@/libs/useModal'
import { style } from '@macaron-css/core'
import { ModalButtonsContainer, ModalContainer, ModalText } from '/@/components/Modal'
import { PlainMessage } from '@bufbuild/protobuf'
import { getPlainMessage } from '../../../libs/getPlainMessage'

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
const EnvVarContainerClass = style({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr 1fr',
  gap: '16px',
})
const EnvVarsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
})
const EnvVarKeyCode = styled('code', {
  base: {
    padding: '8px 12px',
    borderRadius: '4px',
    fontSize: '14px',
    marginLeft: '4px',

    width: '100%',
    display: 'flex',
    alignItems: 'center',
  },
})
const EnvVarButtonContainer = styled('div', {
  base: {
    display: 'flex',
    justifyContent: 'space-between',
    gap: '16px',
  },
})

export default () => {
  const params = useParams()
  const [app, { refetch: refetchApp }] = createResource(
    () => params.id,
    (id) => client.getApplication({ id }),
  )
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const loaded = () => !!(app() && repo())

  const GeneralConfigsContainer: Component = () => {
    // 現在の設定で初期化
    const [generalConfig, setGeneralConfig] = createStore({
      name: app().name,
      refName: app().refName,
    })
    let formContainer: HTMLFormElement

    const updateGeneralSettings: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
      // prevent default form submit (reload page)
      e.preventDefault()

      // validate form
      if (!formContainer.reportValidity()) {
        return
      }

      try {
        await client.updateApplication({
          id: app().id,
          name: generalConfig.name,
          refName: generalConfig.refName,
        })
        toast.success('アプリ設定を更新しました')
        refetchApp()
      } catch (e) {
        handleAPIError(e, 'アプリ設定の更新に失敗しました')
      }
    }

    return (
      <form ref={formContainer}>
        <SettingFieldSet>
          <FormTextBig id='general-settings'>General settings</FormTextBig>
          <div>
            <InputLabel>Application Name</InputLabel>
            <InputBar
              placeholder='my-app'
              value={generalConfig.name}
              onChange={(e) => setGeneralConfig('name', e.currentTarget.value)}
              required
            />
          </div>
          <div>
            <InputLabel>Branch Name</InputLabel>
            <InputBar
              placeholder='main'
              value={generalConfig.refName}
              onChange={(e) => setGeneralConfig('refName', e.currentTarget.value)}
              required
            />
          </div>
          <Button color='black1' size='large' width='auto' onclick={updateGeneralSettings} type='submit'>
            Save
          </Button>
        </SettingFieldSet>
      </form>
    )
  }

  const BuildConfigsContainer: Component = () => {
    const [runtimeConfig, setRuntimeConfig] = createStore<PlainMessage<RuntimeConfig>>({
      command: '',
      entrypoint: '',
      useMariadb: false,
      useMongodb: false,
    })
    const [buildConfigMethod, setBuildConfigMethod] = createSignal<BuildConfigMethod>(app().config.buildConfig.case)
    const [buildConfig, setBuildConfig] = createStore<BuildConfig>({
      runtimeBuildpack: {
        case: 'runtimeBuildpack',
        value: {
          context: '',
          runtimeConfig: runtimeConfig,
        },
      },
      runtimeCmd: {
        case: 'runtimeCmd',
        value: {
          baseImage: '',
          buildCmd: '',
          buildCmdShell: false,
          runtimeConfig: runtimeConfig,
        },
      },
      runtimeDockerfile: {
        case: 'runtimeDockerfile',
        value: {
          context: '',
          dockerfileName: '',
          runtimeConfig: runtimeConfig,
        },
      },
      staticCmd: {
        case: 'staticCmd',
        value: {
          artifactPath: '',
          baseImage: '',
          buildCmd: '',
          buildCmdShell: false,
        },
      },
      staticDockerfile: {
        case: 'staticDockerfile',
        value: {
          artifactPath: '',
          context: '',
          dockerfileName: '',
        },
      },
    })

    // 現在のビルド設定を反映
    onMount(() => {
      const conf = getPlainMessage(app().config)
      setBuildConfigMethod(conf.buildConfig.case)
      setBuildConfig({
        [conf.buildConfig.case]: {
          case: conf.buildConfig.case,
          value: conf.buildConfig.value,
        },
      })
      switch (conf.buildConfig.case) {
        case 'runtimeBuildpack':
        case 'runtimeCmd':
        case 'runtimeDockerfile':
          setRuntimeConfig(conf.buildConfig.value.runtimeConfig)
      }
    })

    const updateBuildSettings: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
      // prevent default form submit (reload page)
      e.preventDefault()

      try {
        await client.updateApplication({
          id: app().id,
          config: {
            buildConfig: buildConfig[buildConfigMethod()],
          },
        })
        toast.success('ビルド設定を更新しました')
        refetchApp()
      } catch (e) {
        handleAPIError(e, 'ビルド設定の更新に失敗しました')
      }
    }

    return (
      <SettingFieldSet>
        <FormTextBig id='build-settings'>Build Settings</FormTextBig>
        <BuildConfigs
          setBuildConfig={setBuildConfig}
          buildConfig={buildConfig}
          runtimeConfig={runtimeConfig}
          setRuntimeConfig={setRuntimeConfig}
          buildConfigMethod={buildConfigMethod()}
          setBuildConfigMethod={setBuildConfigMethod}
        />
        <Button color='black1' size='large' width='auto' onclick={updateBuildSettings} type='submit'>
          Save
        </Button>
      </SettingFieldSet>
    )
  }

  const WebsitesConfigContainer: Component = () => {
    const [websites, setWebsites] = createStore<PlainMessage<CreateWebsiteRequest>[]>([])
    // 現在のウェブサイト設定を反映 (`onMount`ではrefetch時に反映されないので`createEffect`を使用)
    createEffect(() => {
      // classのまま渡すとリアクティビティが失われるのでオブジェクトに変換して渡す
      setWebsites(app().websites.map(getPlainMessage))
    })

    const updateWebsites = async () => {
      console.log(`websites: ${websites.length}`)
      console.log(websites)

      try {
        await client.updateApplication({
          id: app().id,
          websites: { websites },
        })
        toast.success('ウェブサイト設定を更新しました')
        refetchApp()
      } catch (e) {
        handleAPIError(e, 'ウェブサイト設定の更新に失敗しました')
      }
    }

    return (
      <SettingFieldSet>
        <FormTextBig id='website-settings'>Website Settings</FormTextBig>
        <WebsiteSettings
          runtime={app().deployType === DeployType.RUNTIME}
          websiteConfigs={websites}
          setWebsiteConfigs={setWebsites}
        />
        <Button color='black1' size='large' width='auto' onclick={updateWebsites} type='submit'>
          Save
        </Button>
      </SettingFieldSet>
    )
  }

  const PortPublicationConfigContainer: Component = () => {
    const [portPublications, setPortPublications] = createStore<PlainMessage<PortPublication>[]>([])
    // 現在のポート設定を反映
    onMount(() => {
      // classのまま渡すとリアクティビティが失われるのでオブジェクトに変換して渡す
      setPortPublications(app().portPublications.map(getPlainMessage))
    })

    const updatePortPublications = async () => {
      try {
        await client.updateApplication({
          id: app().id,
          portPublications: { portPublications },
        })
        toast.success('ポート設定を更新しました')
        refetchApp()
      } catch (e) {
        handleAPIError(e, 'ポート設定の更新に失敗しました')
      }
    }

    return (
      <SettingFieldSet>
        <FormTextBig id='port-settings'>Port Publication Settings</FormTextBig>
        <PortPublicationSettings portPublications={portPublications} setPortPublications={setPortPublications} />
        <Button color='black1' size='large' width='auto' onclick={updatePortPublications} type='submit'>
          Save
        </Button>
      </SettingFieldSet>
    )
  }

  const OwnerConfigContainer: Component = () => {
    const { Modal, open } = useModal()

    const nonOwnerUsers = createMemo(() => {
      return users()?.filter((user) => !app().ownerIds.includes(user.id)) ?? []
    })

    const handleAddOwner = async (user: User) => {
      try {
        await client.updateApplication({
          id: app().id,
          ownerIds: { ownerIds: app().ownerIds.concat(user.id) },
        })
        toast.success('アプリオーナーを追加しました')
        refetchApp()
      } catch (e) {
        handleAPIError(e, 'アプリオーナーの追加に失敗しました')
      }
    }
    const handleDeleteOwner = async (owner: User) => {
      try {
        await client.updateApplication({
          id: app().id,
          ownerIds: { ownerIds: app().ownerIds.filter((id) => id !== owner.id) },
        })
        toast.success('アプリのオーナーを削除しました')
        refetchApp()
      } catch (e) {
        handleAPIError(e, 'アプリのオーナーの削除に失敗しました')
      }
    }

    return (
      <>
        <SettingFieldSet>
          <FormTextBig id='owner-settings'>Owner Settings</FormTextBig>
          <Button color='black1' size='large' width='auto' onclick={open}>
            アプリオーナーを追加する
          </Button>
          <UserSearch users={app().ownerIds.map((userId) => userFromId(userId))}>
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

  const EnvVarConfigContainer: Component = () => {
    const [envVars, { refetch: refetchEnvVar }] = createResource(
      () => app().id,
      (id) => client.getEnvVars({ id }),
    )

    const EditEnvVarContainer: Component<{
      envVar: ApplicationEnvVar
    }> = (props) => {
      const [isEditing, setIsEditing] = createSignal(false)
      let formRef: HTMLFormElement
      let valueInputRef: HTMLInputElement
      const { Modal: DeleteEnvVarModal, open: openDeleteEnvVarModal, close: closeDeleteEnvVarModal } = useModal()

      const handleUpdateEnvVar: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
        // prevent default form submit (reload page)
        e.preventDefault()

        // validate form
        if (!formRef.reportValidity()) {
          return
        }

        try {
          await client.setEnvVar({
            applicationId: app().id,
            key: props.envVar.key,
            value: valueInputRef.value,
          })
          toast.success('環境変数を更新しました')
          refetchEnvVar()
          setIsEditing(false)
        } catch (e) {
          handleAPIError(e, '環境変数の更新に失敗しました')
        }
      }

      const handleDeleteEnvVar: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
        // prevent default form submit (reload page)
        e.preventDefault()

        try {
          await client.deleteEnvVar({
            applicationId: app().id,
            key: props.envVar.key,
          })
          toast.success('環境変数を削除しました')
          refetchEnvVar()
          setIsEditing(false)
        } catch (e) {
          handleAPIError(e, '環境変数の削除に失敗しました')
        }
      }

      return (
        <form class={EnvVarContainerClass} ref={formRef}>
          <EnvVarKeyCode>{props.envVar.key}</EnvVarKeyCode>
          <InputBar
            type='text'
            disabled={!isEditing()}
            required
            placeholder='VALUE'
            ref={valueInputRef}
            value={props.envVar.value}
          />
          <EnvVarButtonContainer>
            <Show
              when={!isEditing()}
              fallback={
                <Button color='black1' size='large' width='full' type='submit' onclick={handleUpdateEnvVar}>
                  Save
                </Button>
              }
            >
              <Button
                color='black1'
                size='large'
                width='full'
                type='button'
                onclick={() => setIsEditing(true)}
                disabled={props.envVar.system}
                title={props.envVar.system ? 'システム環境変数は編集できません' : undefined}
              >
                Edit
              </Button>
              <Button
                color='black1'
                size='large'
                width='full'
                type='button'
                onclick={openDeleteEnvVarModal}
                disabled={props.envVar.system}
                title={props.envVar.system ? 'システム環境変数は削除できません' : undefined}
              >
                Delete
              </Button>
              <DeleteEnvVarModal>
                <ModalContainer>
                  <ModalText>本当に削除しますか?</ModalText>
                  <ModalButtonsContainer>
                    <Button onclick={closeDeleteEnvVarModal} color='black1' size='large' width='full'>
                      キャンセル
                    </Button>
                    <Button onclick={handleDeleteEnvVar} color='black1' size='large' width='full'>
                      削除
                    </Button>
                  </ModalButtonsContainer>
                </ModalContainer>
              </DeleteEnvVarModal>
            </Show>
          </EnvVarButtonContainer>
        </form>
      )
    }

    const AddEnvVarContainer: Component = () => {
      let formRef: HTMLFormElement
      let keyInputRef: HTMLInputElement
      let valueInputRef: HTMLInputElement

      const handleAddEnvVar: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
        // prevent default form submit (reload page)
        e.preventDefault()

        // validate form
        if (!formRef.reportValidity()) {
          return
        }

        try {
          await client.setEnvVar({
            applicationId: app().id,
            key: keyInputRef.value,
            value: valueInputRef.value,
          })
          toast.success('環境変数を追加しました')
          keyInputRef.value = ''
          valueInputRef.value = ''
          refetchEnvVar()
        } catch (e) {
          handleAPIError(e, '環境変数の追加に失敗しました')
        }
      }

      return (
        <form class={EnvVarContainerClass} ref={formRef}>
          <InputBar type='text' required placeholder='KEY' ref={keyInputRef} />
          <InputBar type='text' required placeholder='VALUE' ref={valueInputRef} />
          <Button color='black1' size='large' width='full' type='submit' onclick={handleAddEnvVar}>
            Add
          </Button>
        </form>
      )
    }

    return (
      <SettingFieldSet>
        <FormTextBig id='env-var-settings'>Environment Variable Settings</FormTextBig>
        <div>
          <EnvVarsContainer>
            <div class={EnvVarContainerClass}>
              <div>key</div>
              <div>value</div>
            </div>
            <For each={envVars()?.variables}>
              {(envVar) => {
                return <EditEnvVarContainer envVar={envVar} />
              }}
            </For>
            <AddEnvVarContainer />
          </EnvVarsContainer>
        </div>
      </SettingFieldSet>
    )
  }

  return (
    <Container>
      <Header />
      <Show when={loaded()}>
        <AppNav repoName={repo().name} appName={app().name} appID={app().id} />
        <ContentContainer>
          <div>
            <SidebarContainer>
              <SidebarOptions>
                <SidebarNavAnchor href='#general-settings'>General</SidebarNavAnchor>
                <SidebarNavAnchor href='#build-settings'>Build</SidebarNavAnchor>
                <SidebarNavAnchor href='#website-settings'>Website</SidebarNavAnchor>
                <SidebarNavAnchor href='#port-settings'>Port Publication</SidebarNavAnchor>
                <SidebarNavAnchor href='#owner-settings'>Owner</SidebarNavAnchor>
                <SidebarNavAnchor href='#env-var-settings'>Environment Variable</SidebarNavAnchor>
              </SidebarOptions>
            </SidebarContainer>
          </div>
          <ConfigsContainer>
            <GeneralConfigsContainer />
            <BuildConfigsContainer />
            <WebsitesConfigContainer />
            <PortPublicationConfigContainer />
            <OwnerConfigContainer />
            <EnvVarConfigContainer />
          </ConfigsContainer>
        </ContentContainer>
      </Show>
    </Container>
  )
}
