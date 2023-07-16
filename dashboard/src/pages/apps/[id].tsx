import { A, useNavigate, useParams } from '@solidjs/router'
import { Component, createResource, createSignal, For, onCleanup, Show } from 'solid-js'
import { availableMetrics, client, handleAPIError, systemInfo } from '/@/libs/api'
import { Header } from '/@/components/Header'
import {
  applicationState,
  buildTypeStr,
  getWebsiteURL,
  providerToIcon,
  repositoryURLToProvider,
} from '/@/libs/application'
import { StatusIcon } from '/@/components/StatusIcon'
import { titleCase } from '/@/libs/casing'
import {
  Application_ContainerState,
  ApplicationConfig,
  DeployType,
  RuntimeConfig,
  StaticConfig,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DiffHuman, shortSha } from '/@/libs/format'
import { CenterInline, Container } from '/@/libs/layout'
import { URLText } from '/@/components/URLText'
import { Button } from '/@/components/Button'
import { AppNav } from '/@/components/AppNav'
import {
  Card,
  CardItem,
  CardItemContent,
  CardItems,
  CardItemTitle,
  CardRowsContainer,
  CardsRow,
  CardTitle,
} from '/@/components/Card'
import useModal from '/@/libs/useModal'
import { ModalButtonsContainer, ModalContainer, ModalText } from '/@/components/Modal'
import toast from 'solid-toast'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { unreachable } from '/@/libs/unreachable'
import { ContainerLog } from '/@/components/ContainerLog'
import { Checkbox } from '/@/components/Checkbox'
import { AppMetrics } from '/@/components/AppMetrics'

const RuntimeConfigInfo: Component<{ config: RuntimeConfig }> = (props) => {
  return (
    <>
      <CardItem>
        <CardItemTitle>Use MariaDB</CardItemTitle>
        <CardItemContent>{`${props.config.useMariadb}`}</CardItemContent>
      </CardItem>
      <CardItem>
        <CardItemTitle>Use MongoDB</CardItemTitle>
        <CardItemContent>{`${props.config.useMongodb}`}</CardItemContent>
      </CardItem>
      <Show when={props.config.entrypoint !== ''}>
        <CardItem>
          <CardItemTitle>Entrypoint</CardItemTitle>
          <CardItemContent>{props.config.entrypoint}</CardItemContent>
        </CardItem>
      </Show>
      <Show when={props.config.command !== ''}>
        <CardItem>
          <CardItemTitle>Command</CardItemTitle>
          <CardItemContent>{props.config.command}</CardItemContent>
        </CardItem>
      </Show>
    </>
  )
}

const StaticConfigInfo: Component<{ config: StaticConfig }> = (props) => {
  return (
    <>
      <CardItem>
        <CardItemTitle>Artifact Path</CardItemTitle>
        <CardItemContent>{props.config.artifactPath}</CardItemContent>
      </CardItem>
      <CardItem>
        <CardItemTitle>Single Page Application</CardItemTitle>
        <CardItemContent>{props.config.spa}</CardItemContent>
      </CardItem>
    </>
  )
}

const ApplicationConfigInfo: Component<{ config: ApplicationConfig }> = (props) => {
  const c = props.config.buildConfig
  switch (c.case) {
    case 'runtimeBuildpack':
      return (
        <>
          <CardItem>
            <CardItemTitle>Context</CardItemTitle>
            <CardItemContent>{c.value.context}</CardItemContent>
          </CardItem>
          <RuntimeConfigInfo config={c.value.runtimeConfig} />
        </>
      )
    case 'runtimeCmd':
      return (
        <>
          <CardItem>
            <CardItemTitle>Base Image</CardItemTitle>
            <CardItemContent>{c.value.baseImage || 'Scratch'}</CardItemContent>
          </CardItem>
          {c.value.buildCmd && (
            <>
              <CardItem>
                <CardItemTitle>Build Command</CardItemTitle>
              </CardItem>
              <CardItem>
                <pre>{c.value.buildCmd}</pre>
              </CardItem>
            </>
          )}
          <RuntimeConfigInfo config={c.value.runtimeConfig} />
        </>
      )
    case 'runtimeDockerfile':
      return (
        <>
          <CardItem>
            <CardItemTitle>Dockerfile</CardItemTitle>
            <CardItemContent>{c.value.dockerfileName}</CardItemContent>
          </CardItem>
          <CardItem>
            <CardItemTitle>Context</CardItemTitle>
            <CardItemContent>{c.value.context}</CardItemContent>
          </CardItem>
          <RuntimeConfigInfo config={c.value.runtimeConfig} />
        </>
      )
    case 'staticBuildpack':
      return (
        <>
          <StaticConfigInfo config={c.value.staticConfig} />
          <CardItem>
            <CardItemTitle>Context</CardItemTitle>
            <CardItemContent>{c.value.context}</CardItemContent>
          </CardItem>
        </>
      )
    case 'staticCmd':
      return (
        <>
          <StaticConfigInfo config={c.value.staticConfig} />
          <CardItem>
            <CardItemTitle>Base Image</CardItemTitle>
            <CardItemContent>{c.value.baseImage || 'Scratch'}</CardItemContent>
          </CardItem>
          {c.value.buildCmd && (
            <CardItem>
              <CardItemTitle>Build Command</CardItemTitle>
              <CardItemContent>{c.value.buildCmd}</CardItemContent>
            </CardItem>
          )}
        </>
      )
    case 'staticDockerfile':
      return (
        <>
          <StaticConfigInfo config={c.value.staticConfig} />
          <CardItem>
            <CardItemTitle>Dockerfile</CardItemTitle>
            <CardItemContent>{c.value.dockerfileName}</CardItemContent>
          </CardItem>
          <CardItem>
            <CardItemTitle>Context</CardItemTitle>
            <CardItemContent>{c.value.context}</CardItemContent>
          </CardItem>
        </>
      )
  }
  return unreachable(c)
}

const URLsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
})

const SSHCode = styled('code', {
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
  const navigate = useNavigate()
  const params = useParams()
  const id = params.id
  const [app, { refetch: refetchApp }] = createResource(
    () => id,
    (id) => client.getApplication({ id }),
  )
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const loaded = () => !!(systemInfo() && app() && repo())

  const refetchTimer = setInterval(refetchApp, 10000)
  onCleanup(() => clearInterval(refetchTimer))

  const [disableRefresh, setDisableRefresh] = createSignal(false)
  const refreshRepo = async () => {
    setDisableRefresh(true)
    setTimeout(() => setDisableRefresh(false), 3000)
    await client.refreshRepository({ repositoryId: repo().id })
    await refetchApp()
  }
  const startApp = async () => {
    await client.startApplication({ id: app().id })
    await refetchApp()
  }
  const stopApp = async () => {
    await client.stopApplication({ id: app().id })
    await refetchApp()
  }
  const deleteApp = async () => {
    try {
      await client.deleteApplication({ id: app().id })
    } catch (e) {
      handleAPIError(e, 'アプリケーションの削除に失敗しました')
      return
    }
    toast.success('アプリケーションを削除しました')
    navigate('/apps')
  }
  const { Modal: DeleteAppModal, open: openDeleteAppModal, close: closeDeleteAppModal } = useModal()

  const [showLogTimestamp, setShowLogTimestamp] = createSignal(true)

  return (
    <Container>
      <Header />
      <Show when={loaded()}>
        <AppNav repo={repo()} app={app()} />
        <CardRowsContainer>
          <CardsRow>
            <Card>
              <CardTitle>Actions</CardTitle>
              <Button color='black1' size='large' width='full' onclick={refreshRepo} disabled={disableRefresh()}>
                Refresh Commit
              </Button>
              <Button
                onclick={openDeleteAppModal}
                color='black1'
                size='large'
                width='full'
                disabled={app().running}
                tooltip={
                  app().running ? 'アプリケーションが起動しているため削除できません' : 'アプリケーションを削除します'
                }
              >
                Delete Application
              </Button>
              <DeleteAppModal>
                <ModalContainer>
                  <ModalText>本当に削除しますか?</ModalText>
                  <ModalButtonsContainer>
                    <Button onclick={closeDeleteAppModal} color='black1' size='large' width='full'>
                      キャンセル
                    </Button>
                    <Button onclick={deleteApp} color='black1' size='large' width='full'>
                      削除
                    </Button>
                  </ModalButtonsContainer>
                </ModalContainer>
              </DeleteAppModal>
              <Show when={!app().running}>
                <Button color='black1' size='large' width='full' onclick={startApp}>
                  Start App
                </Button>
              </Show>
              <Show when={app().running}>
                <Button color='black1' size='large' width='full' onclick={startApp}>
                  Restart App
                </Button>
              </Show>
              <Show when={app().running}>
                <Button color='black1' size='large' width='full' onclick={stopApp}>
                  Stop App
                </Button>
              </Show>
            </Card>
            <Card>
              <CardTitle>Overall</CardTitle>
              <CardItems>
                <CardItem>
                  <CardItemTitle>状態</CardItemTitle>
                  <CardItemContent>
                    <StatusIcon state={applicationState(app())} size={24} />
                    {applicationState(app())}
                  </CardItemContent>
                </CardItem>
                <Show when={app().deployType === DeployType.RUNTIME}>
                  <CardItem>
                    <CardItemTitle>コンテナの状態</CardItemTitle>
                    <CardItemContent>{app() && titleCase(Application_ContainerState[app().container])}</CardItemContent>
                  </CardItem>
                </Show>
                <CardItem>
                  <CardItemTitle>起動時刻</CardItemTitle>
                  <CardItemContent>
                    <Show when={app().running} fallback={'-'}>
                      <DiffHuman target={app().updatedAt.toDate()} />
                    </Show>
                  </CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>作成日</CardItemTitle>
                  <CardItemContent>
                    <DiffHuman target={app().createdAt.toDate()} />
                  </CardItemContent>
                </CardItem>
                <Show when={app().websites.length > 0}>
                  <CardItem>
                    <CardItemTitle>URLs</CardItemTitle>
                  </CardItem>
                  <CardItem>
                    <CardItemTitle>
                      <URLsContainer>
                        <For each={app().websites}>
                          {(website) => (
                            <URLText href={getWebsiteURL(website)} target='_blank' rel='noreferrer'>
                              {getWebsiteURL(website)}
                            </URLText>
                          )}
                        </For>
                      </URLsContainer>
                    </CardItemTitle>
                  </CardItem>
                </Show>
              </CardItems>
            </Card>
            <Card>
              <CardTitle>Info</CardTitle>
              <CardItems>
                <CardItem>
                  <CardItemTitle>ID</CardItemTitle>
                  <CardItemContent>{app().id}</CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>Name</CardItemTitle>
                  <CardItemContent>{app().name}</CardItemContent>
                </CardItem>
                <A href={`/repos/${repo().id}`}>
                  <CardItem>
                    <CardItemTitle>Repository</CardItemTitle>
                    <CardItemContent>
                      <CenterInline>{providerToIcon(repositoryURLToProvider(repo().url), 20)}</CenterInline>
                      {repo().name}
                    </CardItemContent>
                  </CardItem>
                </A>
                <CardItem>
                  <CardItemTitle>Git ref (short)</CardItemTitle>
                  <CardItemContent>{app().refName}</CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>Deploy type</CardItemTitle>
                  <CardItemContent>{titleCase(DeployType[app().deployType])}</CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>Commit</CardItemTitle>
                  <CardItemContent>
                    <div>{shortSha(app().commit)}</div>
                  </CardItemContent>
                </CardItem>
                <CardItem>
                  <CardItemTitle>Deployed build</CardItemTitle>
                  <CardItemContent>
                    <div>{shortSha(app().currentBuild)}</div>
                  </CardItemContent>
                </CardItem>
              </CardItems>
            </Card>
            <Card>
              <CardTitle>Config</CardTitle>
              <CardItems>
                <CardItem>
                  <CardItemTitle>Build Type</CardItemTitle>
                  <CardItemContent>{buildTypeStr[app().config.buildConfig.case]}</CardItemContent>
                </CardItem>
                <ApplicationConfigInfo config={app().config} />
              </CardItems>
            </Card>
          </CardsRow>
          <Show when={app().deployType === DeployType.RUNTIME && availableMetrics()}>
            <CardsRow>
              <For each={availableMetrics().metricsNames || []}>
                {(name) => (
                  <Card>
                    <CardTitle>{name}</CardTitle>
                    <CardItems>
                      <AppMetrics appID={id} metricsName={name} />
                    </CardItems>
                  </Card>
                )}
              </For>
            </CardsRow>
          </Show>
          <Show when={app().deployType === DeployType.RUNTIME}>
            <CardsRow>
              <Card>
                <CardTitle>SSH Access</CardTitle>
                <CardItems>
                  <Show
                    when={app().running}
                    fallback={<CardItem>アプリケーションが起動している間のみSSHでアクセス可能です</CardItem>}
                  >
                    <CardItem>
                      <SSHCode>{`ssh -p ${systemInfo().ssh.port} ${app().id}@${systemInfo().ssh.host}`}</SSHCode>
                    </CardItem>
                  </Show>
                </CardItems>
              </Card>
            </CardsRow>
          </Show>
          <Show when={app().deployType === DeployType.RUNTIME}>
            <CardsRow>
              <Card>
                <CardTitle>Container Log</CardTitle>
                <Checkbox selected={showLogTimestamp()} setSelected={setShowLogTimestamp}>
                  Show Timestamps
                </Checkbox>
                <ContainerLog appID={app().id} showTimestamp={showLogTimestamp()} />
              </Card>
            </CardsRow>
          </Show>
        </CardRowsContainer>
      </Show>
    </Container>
  )
}
