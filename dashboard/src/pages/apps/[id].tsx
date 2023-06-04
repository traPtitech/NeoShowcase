import { A, useNavigate, useParams } from '@solidjs/router'
import { createEffect, createMemo, createResource, createSignal, For, onCleanup, Ref, Show } from 'solid-js'
import { client, handleAPIError } from '/@/libs/api'
import { Header } from '/@/components/Header'
import { applicationState, buildTypeStr, getWebsiteURL, providerToIcon } from '/@/libs/application'
import { StatusIcon } from '/@/components/StatusIcon'
import { titleCase } from '/@/libs/casing'
import {
  Application_ContainerState,
  ApplicationConfig,
  ApplicationOutput,
  DeployType,
  RuntimeConfig,
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
  CardsContainer,
  CardTitle,
} from '/@/components/Card'
import { LogContainer } from '/@/components/Log'
import { sleep } from '/@/libs/sleep'
import { Timestamp } from '@bufbuild/protobuf'
import { Code, ConnectError } from '@bufbuild/connect'
import useModal from '/@/libs/useModal'
import { ModalButtonsContainer, ModalContainer, ModalText } from '/@/components/Modal'
import toast from 'solid-toast'

interface RuntimeConfigInfoProps {
  config: RuntimeConfig
}

const RuntimeConfigInfo = (props: RuntimeConfigInfoProps) => {
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

interface ApplicationConfigInfoProps {
  config: ApplicationConfig
}

const ApplicationConfigInfo = (props: ApplicationConfigInfoProps) => {
  const c = props.config.buildConfig
  switch (c.case) {
    case 'runtimeBuildpack':
      return (
        <>
          <RuntimeConfigInfo config={c.value.runtimeConfig} />
          <CardItem>
            <CardItemTitle>Context</CardItemTitle>
            <CardItemContent>{c.value.context}</CardItemContent>
          </CardItem>
        </>
      )
    case 'runtimeCmd':
      return (
        <>
          <RuntimeConfigInfo config={c.value.runtimeConfig} />
          <CardItem>
            <CardItemTitle>Base Image</CardItemTitle>
            <CardItemContent>{c.value.baseImage || 'Scratch'}</CardItemContent>
          </CardItem>
          {c.value.buildCmd && (
            <CardItem>
              <CardItemTitle>Build Command{c.value.buildCmdShell && ' (Shell)'}</CardItemTitle>
              <CardItemContent>{c.value.buildCmd}</CardItemContent>
            </CardItem>
          )}
        </>
      )
    case 'runtimeDockerfile':
      return (
        <>
          <RuntimeConfigInfo config={c.value.runtimeConfig} />
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
    case 'staticCmd':
      return (
        <>
          <CardItem>
            <CardItemTitle>Base Image</CardItemTitle>
            <CardItemContent>{c.value.baseImage || 'Scratch'}</CardItemContent>
          </CardItem>
          {c.value.buildCmd && (
            <CardItem>
              <CardItemTitle>Build Command{c.value.buildCmdShell && ' (Shell)'}</CardItemTitle>
              <CardItemContent>{c.value.buildCmd}</CardItemContent>
            </CardItem>
          )}
          <CardItem>
            <CardItemTitle>Artifact Path</CardItemTitle>
            <CardItemContent>{c.value.artifactPath}</CardItemContent>
          </CardItem>
        </>
      )
    case 'staticDockerfile':
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
          <CardItem>
            <CardItemTitle>Artifact Path</CardItemTitle>
            <CardItemContent>{c.value.artifactPath}</CardItemContent>
          </CardItem>
        </>
      )
  }
}

export default () => {
  const navigate = useNavigate()
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

  const now = new Date()
  const [log] = createResource(
    () => app()?.deployType === DeployType.RUNTIME && app()?.id,
    (id) => client.getOutput({ applicationId: id, before: Timestamp.fromDate(now) }),
  )
  const reversedLog = createMemo(() => log() && ([...log().outputs].reverse() as ApplicationOutput[]))
  const logStreamAbort = new AbortController()
  const [logStream] = createResource(
    () => app()?.deployType === DeployType.RUNTIME && app()?.id,
    (id) =>
      client.getOutputStream({ applicationId: id, after: Timestamp.fromDate(now) }, { signal: logStreamAbort.signal }),
  )
  const [streamedLog, setStreamedLog] = createSignal<string[]>([])
  createEffect(() => {
    const stream = logStream()
    if (!stream) {
      setStreamedLog([])
      return
    }

    const iterate = async () => {
      try {
        for await (const log of stream) {
          setStreamedLog((prev) => prev.concat(log.log))
        }
      } catch (err) {
        // ignore abort error
        const isAbortErr = err instanceof ConnectError && err.code === Code.Canceled
        if (!isAbortErr) {
          console.trace(err)
          return
        }
      }
      await sleep(1000)
      await refetchApp()
    }
    void iterate()
  })
  onCleanup(() => {
    logStreamAbort.abort()
  })

  let logRef: Ref<HTMLDivElement>
  createEffect(() => {
    if (!log() || !streamedLog()) return
    const ref = logRef as HTMLDivElement
    if (!ref) return
    setTimeout(() => {
      ref.scrollTop = ref.scrollHeight
    })
  })

  return (
    <Container>
      <Header />
      <Show when={loaded()}>
        <AppNav appID={app().id} appName={app().name} repoName={repo().name} />
        <CardsContainer>
          <Card>
            <CardTitle>Actions</CardTitle>
            <Button color='black1' size='large' onclick={refreshRepo} disabled={disableRefresh()}>
              Refresh Commit
            </Button>
            <Button
              onclick={openDeleteAppModal}
              color='black1'
              size='large'
              disabled={app().running}
              title={
                app().running ? 'アプリケーションを削除します' : 'アプリケーションが起動しているため削除できません'
              }
            >
              Delete Application
            </Button>
            <DeleteAppModal>
              <ModalContainer>
                <ModalText>本当に削除しますか?</ModalText>
                <ModalButtonsContainer>
                  <Button onclick={closeDeleteAppModal} color='black1' size='large'>
                    キャンセル
                  </Button>
                  <Button onclick={deleteApp} color='black1' size='large'>
                    削除
                  </Button>
                </ModalButtonsContainer>
              </ModalContainer>
            </DeleteAppModal>
            <Show when={!app().running}>
              <Button color='black1' size='large' onclick={startApp}>
                Start App
              </Button>
            </Show>
            <Show when={app().running}>
              <Button color='black1' size='large' onclick={startApp}>
                Restart App
              </Button>
            </Show>
            <Show when={app().running}>
              <Button color='black1' size='large' onclick={stopApp}>
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
                    <For each={app().websites}>
                      {(website) => (
                        <URLText href={getWebsiteURL(website)} target='_blank' rel='noreferrer'>
                          {getWebsiteURL(website)}
                        </URLText>
                      )}
                    </For>
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
                    <CenterInline>{providerToIcon('GitHub', 20)}</CenterInline>
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
                  <Show
                    when={app().currentCommit !== app().wantCommit}
                    fallback={<div>{shortSha(app().currentCommit)}</div>}
                  >
                    <div>
                      {shortSha(app().currentCommit)} → {shortSha(app().wantCommit)} (Deploying)
                    </div>
                  </Show>
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
          <Show when={log()}>
            <Card>
              <CardTitle>Container Log</CardTitle>
              <LogContainer ref={logRef} overflowX='scroll'>
                <For each={reversedLog()}>{(l) => <span>{l.log}</span>}</For>
                <For each={streamedLog()}>{(line) => <span>{line}</span>}</For>
              </LogContainer>
            </Card>
          </Show>
        </CardsContainer>
      </Show>
    </Container>
  )
}
