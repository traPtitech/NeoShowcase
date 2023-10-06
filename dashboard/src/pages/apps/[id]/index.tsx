import {
  Application,
  ApplicationConfig,
  BuildStatus,
  DeployType,
  Repository,
  RuntimeConfig,
  StaticConfig,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { StatusIcon } from '/@/components/StatusIcon'
import { BuildStatusIcon } from '/@/components/UI/BuildStatusIcon'
import { Button } from '/@/components/UI/Button'
import { URLText } from '/@/components/UI/URLText'
import { DataTable } from '/@/components/layouts/DataTable'
import { AppMetrics } from '/@/components/templates/AppMetrics'
import { ContainerLog } from '/@/components/templates/ContainerLog'
import { List } from '/@/components/templates/List'
import { availableMetrics, client, handleAPIError, systemInfo } from '/@/libs/api'
import { applicationState, buildStatusStr, buildTypeStr, getWebsiteURL } from '/@/libs/application'
import { titleCase } from '/@/libs/casing'
import { diffHuman, shortSha } from '/@/libs/format'
import { useApplicationData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, For, Match, Show, Switch, createSignal, onCleanup } from 'solid-js'
import { tippy as tippyDir } from 'solid-tippy'
import toast from 'solid-toast'

// https://github.com/solidjs/solid/discussions/845
const tippy = tippyDir

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    overflowY: 'auto',
  },
})
const MainViewContainer = styled('div', {
  base: {
    width: '100%',
    padding: '40px 32px 72px 32px',
  },
  variants: {
    gray: {
      true: {
        background: colorVars.semantic.ui.background,
      },
      false: {
        background: colorVars.semantic.ui.primary,
      },
    },
  },
})
const MainView = styled('div', {
  base: {
    width: '100%',
    maxWidth: '1000px',
    margin: '0 auto',
    display: 'flex',
    flexDirection: 'column',
    gap: '32px',
  },
})

const DeploymentContainer = styled('div', {
  base: {
    width: '100%',
    display: 'grid',
    gridTemplateColumns: '32% 1fr 1fr',
    gridTemplateRows: 'repeat(4, auto)',
    gap: '1px',

    background: colorVars.semantic.ui.border,
    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
    overflow: 'hidden',
  },
})
const AppStateContainer = styled('div', {
  base: {
    gridArea: '1 / 1 / 5 / 2',
    width: '100%',
    display: 'grid',
    gridTemplateRows: '1fr 2fr 1fr',
    justifyItems: 'center',

    cursor: 'pointer',
    background: colorVars.semantic.ui.primary,
    color: colorVars.semantic.text.black,
    ...textVars.h3.medium,
  },
})
const AppState = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    gap: '8px',
  },
})
const ActionButtons = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
  },
})
const DeployInfoContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',

    background: colorVars.semantic.ui.primary,
  },
  variants: {
    long: {
      true: {
        gridColumn: '2 / 4',
      },
    },
  },
})
const UrlCount = styled('div', {
  base: {
    height: '20px',
    padding: '0 8px',
    borderRadius: '9999px',

    background: colorVars.primitive.blackAlpha[200],
    color: colorVars.semantic.text.black,
    ...textVars.caption.regular,
  },
})
const DeploymentInfo: Component<{
  app: Application
  refetchApp: () => void
  repo: Repository
  refreshRepo: () => void
  disableRefresh: () => boolean
}> = (props) => {
  const [showActions, setShowActions] = createSignal(false)

  const restartApp = async () => {
    try {
      await client.startApplication({ id: props.app.id })
      await props.refetchApp()
      toast.success('アプリケーションを再起動しました')
    } catch (e) {
      handleAPIError(e, 'アプリケーションの再起動に失敗しました')
    }
  }
  const stopApp = async () => {
    try {
      await client.stopApplication({ id: props.app.id })
      await props.refetchApp()
      toast.success('アプリケーションを停止しました')
    } catch (e) {
      handleAPIError(e, 'アプリケーションの停止に失敗しました')
    }
  }

  return (
    <DeploymentContainer>
      <AppStateContainer onMouseEnter={() => setShowActions(true)} onMouseLeave={() => setShowActions(false)}>
        <div />
        <AppState>
          <StatusIcon state={applicationState(props.app)} size={80} />
          {applicationState(props.app)}
        </AppState>
        <Show when={showActions()}>
          <ActionButtons>
            <Button color="borderError" size="small" onClick={restartApp} disabled={props.disableRefresh()}>
              Restart App
            </Button>
            <Button color="borderError" size="small" onClick={stopApp} disabled={props.disableRefresh()}>
              Stop App
            </Button>
          </ActionButtons>
        </Show>
      </AppStateContainer>
      <DeployInfoContainer>
        <List.RowContent>
          <List.RowTitle>Repository</List.RowTitle>
          <List.RowData>
            <A href={`/repos/${props.repo.id}`}>{props.repo.name}</A>
          </List.RowData>
        </List.RowContent>
      </DeployInfoContainer>
      <DeployInfoContainer>
        <List.RowContent>
          <List.RowTitle>Source Branch</List.RowTitle>
          <List.RowData>{props.app.refName}</List.RowData>
        </List.RowContent>
      </DeployInfoContainer>
      <DeployInfoContainer>
        <List.RowContent>
          <List.RowTitle>Deployed Build</List.RowTitle>
          <List.RowData>{shortSha(props.app.currentBuild)}</List.RowData>
        </List.RowContent>
      </DeployInfoContainer>
      <DeployInfoContainer>
        <List.RowContent>
          <List.RowTitle>Deploy Type</List.RowTitle>
          <List.RowData>{titleCase(DeployType[props.app.deployType])}</List.RowData>
        </List.RowContent>
      </DeployInfoContainer>
      <DeployInfoContainer long>
        <List.RowContent>
          <List.RowTitle>Source Commit</List.RowTitle>
          <List.RowData>{shortSha(props.app.commit)}</List.RowData>
        </List.RowContent>
        <Button color="ghost" size="medium" onClick={props.refreshRepo} disabled={props.disableRefresh()}>
          Refresh Commit
        </Button>
      </DeployInfoContainer>
      <DeployInfoContainer long>
        <List.RowContent>
          <List.RowTitle>
            Domains
            <UrlCount>{props.app.websites.length}</UrlCount>
          </List.RowTitle>
          <For each={props.app.websites.map(getWebsiteURL)}>
            {(url) => (
              <List.RowData>
                <URLText text={url} href={url} />
              </List.RowData>
            )}
          </For>
        </List.RowContent>
      </DeployInfoContainer>
    </DeploymentContainer>
  )
}

const BuildStatusRow = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
    background: colorVars.semantic.ui.secondary,
  },
})
const BuildStatusLabel = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '4px',

    color: colorVars.semantic.text.black,
    ...textVars.text.medium,
  },
})

const BuildStatusTable: Component<{
  app: Application
  refetchApp: () => void
  repo: Repository
  refreshRepo: () => void
  disableRefresh: () => boolean
}> = (props) => {
  const canRebuild = () =>
    props.app.latestBuildStatus === BuildStatus.SUCCEEDED ||
    props.app.latestBuildStatus === BuildStatus.FAILED ||
    props.app.latestBuildStatus === BuildStatus.CANCELLED

  const rebuild = async () => {
    try {
      await client.retryCommitBuild({
        applicationId: props.app.id,
        commit: props.app.commit,
      })
      await props.refetchApp()
    } catch (e) {
      handleAPIError(e, '再ビルドに失敗しました')
    }
  }
  // const cancelBuild = async () => {
  //   try {
  //     await client.cancelBuild({
  //       buildId: props.app.currentBuild,
  //     });
  //     await props.refetchApp();
  //   } catch (e) {
  //     handleAPIError(e, "再ビルドに失敗しました");
  //   }
  // };

  return (
    <List.Container>
      <BuildStatusRow>
        <BuildStatusLabel>
          <BuildStatusIcon state={props.app.latestBuildStatus} size={24} />
          {buildStatusStr[props.app.latestBuildStatus]}
        </BuildStatusLabel>
        <Show when={canRebuild()}>
          <Button color="borderError" size="small" onClick={rebuild} disabled={props.disableRefresh()}>
            Rebuild
          </Button>
        </Show>
        {/* <Show when={props.app.latestBuildStatus === BuildStatus.BUILDING}>
          <Button
            color="borderError"
            size="small"
            onClick={cancelBuild}
            disabled={disableRefresh()}
          >
            Cancel Build
          </Button>
        </Show> */}
      </BuildStatusRow>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Repository</List.RowTitle>
          <List.RowData>
            <A href={`/repos/${props.repo.id}`}>{props.repo.name}</A>
          </List.RowData>
        </List.RowContent>
      </List.Row>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Source Branch (Commit)</List.RowTitle>
          <List.RowData>
            {props.app.refName} ({shortSha(props.app.commit)})
          </List.RowData>
        </List.RowContent>
        <Button color="ghost" size="medium" onClick={props.refreshRepo} disabled={props.disableRefresh()}>
          Refresh Commit
        </Button>
      </List.Row>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Build Type</List.RowTitle>
          <List.RowData>{buildTypeStr[props.app.config.buildConfig.case]}</List.RowData>
        </List.RowContent>
      </List.Row>
    </List.Container>
  )
}

const RuntimeConfigInfo: Component<{ config: RuntimeConfig }> = (props) => {
  return (
    <>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Use MariaDB</List.RowTitle>
          <List.RowData>{`${props.config.useMariadb}`}</List.RowData>
        </List.RowContent>
      </List.Row>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Use MongoDB</List.RowTitle>
          <List.RowData>{`${props.config.useMongodb}`}</List.RowData>
        </List.RowContent>
      </List.Row>
      <Show when={props.config.entrypoint !== ''}>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>Entrypoint</List.RowTitle>
            <List.RowData code>{props.config.entrypoint}</List.RowData>
          </List.RowContent>
        </List.Row>
      </Show>
      <Show when={props.config.command !== ''}>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>Command</List.RowTitle>
            <List.RowData code>{props.config.command}</List.RowData>
          </List.RowContent>
        </List.Row>
      </Show>
    </>
  )
}
const StaticConfigInfo: Component<{ config: StaticConfig }> = (props) => {
  return (
    <>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Artifact Path</List.RowTitle>
          <List.RowData>{props.config.artifactPath}</List.RowData>
        </List.RowContent>
      </List.Row>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Single Page Application</List.RowTitle>
          <List.RowData>{`${props.config.spa}`}</List.RowData>
        </List.RowContent>
      </List.Row>
    </>
  )
}
const ApplicationConfigInfo: Component<{ config: ApplicationConfig }> = (props) => {
  const c = props.config.buildConfig
  return (
    <Switch>
      <Match when={c.case === 'runtimeBuildpack' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Context</List.RowTitle>
                <List.RowData>{c().value.context}</List.RowData>
              </List.RowContent>
            </List.Row>
            <RuntimeConfigInfo config={c().value.runtimeConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'runtimeCmd' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Base Image</List.RowTitle>
                <List.RowData>{c().value.baseImage}</List.RowData>
              </List.RowContent>
            </List.Row>
            <Show when={c().value.buildCmd !== ''}>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>Build Command</List.RowTitle>
                  <List.RowData code>{c().value.buildCmd}</List.RowData>
                </List.RowContent>
              </List.Row>
            </Show>
            <RuntimeConfigInfo config={c().value.runtimeConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'runtimeDockerfile' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Dockerfile</List.RowTitle>
                <List.RowData>{c().value.dockerfileName}</List.RowData>
              </List.RowContent>
            </List.Row>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Context</List.RowTitle>
                <List.RowData>{c().value.context}</List.RowData>
              </List.RowContent>
            </List.Row>
            <RuntimeConfigInfo config={c().value.runtimeConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'staticBuildpack' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Context</List.RowTitle>
                <List.RowData>{c().value.context}</List.RowData>
              </List.RowContent>
            </List.Row>
            <StaticConfigInfo config={c().value.staticConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'staticCmd' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Base Image</List.RowTitle>
                <List.RowData>{c().value.baseImage}</List.RowData>
              </List.RowContent>
            </List.Row>
            <Show when={c().value.buildCmd !== ''}>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>Build Command</List.RowTitle>
                  <List.RowData code>{c().value.buildCmd}</List.RowData>
                </List.RowContent>
              </List.Row>
            </Show>
            <StaticConfigInfo config={c().value.staticConfig} />
          </>
        )}
      </Match>
      <Match when={c.case === 'staticDockerfile' && c}>
        {(c) => (
          <>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Dockerfile</List.RowTitle>
                <List.RowData>{c().value.dockerfileName}</List.RowData>
              </List.RowContent>
            </List.Row>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>Context</List.RowTitle>
                <List.RowData>{c().value.context}</List.RowData>
              </List.RowContent>
            </List.Row>
            <StaticConfigInfo config={c().value.staticConfig} />
          </>
        )}
      </Match>
    </Switch>
  )
}

const Information: Component<{ app: Application }> = (props) => {
  const [showDetails, setShowDetails] = createSignal(false)
  const sshAccessCommand = () => `ssh -p ${systemInfo().ssh.port} ${props.app.id}@${systemInfo().ssh.host}`

  return (
    <>
      <List.Container>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>ID</List.RowTitle>
            <List.RowData>{props.app.id}</List.RowData>
          </List.RowContent>
        </List.Row>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>Name</List.RowTitle>
            <List.RowData>{props.app.name}</List.RowData>
          </List.RowContent>
        </List.Row>
        <Show when={props.app.updatedAt}>
          {(nonNullUpdatedAt) => {
            const { diff, localeString } = diffHuman(nonNullUpdatedAt().toDate())

            return (
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>起動時刻</List.RowTitle>
                  <List.RowData>
                    <span
                      use:tippy={{
                        props: { content: localeString, maxWidth: 1000 },
                        hidden: true,
                      }}
                    >
                      {diff}
                    </span>
                  </List.RowData>
                </List.RowContent>
              </List.Row>
            )
          }}
        </Show>
        <Show when={props.app.createdAt}>
          {(nonNullCreatedAt) => {
            const { diff, localeString } = diffHuman(nonNullCreatedAt().toDate())
            return (
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>作成日</List.RowTitle>
                  <List.RowData>
                    <span
                      use:tippy={{
                        props: { content: localeString, maxWidth: 1000 },
                        hidden: true,
                      }}
                    >
                      {diff}
                    </span>
                  </List.RowData>
                </List.RowContent>
              </List.Row>
            )
          }}
        </Show>
      </List.Container>
      <Show when={!showDetails()}>
        <Button
          color="ghost"
          size="small"
          onClick={() => setShowDetails(true)}
          style={{
            margin: '0 auto',
          }}
        >
          Show More
        </Button>
      </Show>
      <Show when={showDetails()}>
        <List.Container>
          <List.Row>
            <List.RowContent>
              <List.RowTitle>Build Type</List.RowTitle>
              <List.RowData>{buildTypeStr[props.app.config.buildConfig.case]}</List.RowData>
            </List.RowContent>
          </List.Row>
          <ApplicationConfigInfo config={props.app.config} />
        </List.Container>
        <Show when={props.app.deployType === DeployType.RUNTIME}>
          <List.Container>
            <List.Row>
              <List.RowContent>
                <List.RowTitle>SSH Access</List.RowTitle>
                <Show
                  when={props.app.running}
                  fallback={<List.RowData>アプリケーションが起動している間のみSSHでアクセス可能です</List.RowData>}
                >
                  <List.RowData code>{sshAccessCommand()}</List.RowData>
                </Show>
              </List.RowContent>
            </List.Row>
          </List.Container>
        </Show>
      </Show>
    </>
  )
}

const MetricsContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',

    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
  },
})
const MetricsTypeButtons = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
  },
})
const ChartContainer = styled('div', {
  base: {
    width: '100%',
    borderRadius: '8px',
    height: 'auto',
    aspectRatio: '960 / 464',

    background: colorVars.semantic.ui.secondary,
  },
})

const Metrics: Component<{ app: Application }> = (props) => {
  const { metricsNames } = availableMetrics()
  const [currentView, setCurrentView] = createSignal(metricsNames[0])

  return (
    <MetricsContainer>
      <MetricsTypeButtons>
        <For each={metricsNames}>
          {(metrics) => (
            <Button
              color="text"
              size="medium"
              onClick={() => setCurrentView(metrics)}
              active={currentView() === metrics}
            >
              {metrics}
            </Button>
          )}
        </For>
      </MetricsTypeButtons>
      <ChartContainer>
        <For each={metricsNames}>
          {(metrics) => (
            <Show when={currentView() === metrics}>
              <AppMetrics appID={props.app.id} metricsName={metrics} />
            </Show>
          )}
        </For>
      </ChartContainer>
    </MetricsContainer>
  )
}

const LogContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',

    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
  },
})

const Logs: Component<{ app: Application }> = (props) => {
  return (
    <LogContainer>
      <ContainerLog appID={props.app.id} showTimestamp={true} />
    </LogContainer>
  )
}

export default () => {
  const { app, refetchApp, repo } = useApplicationData()
  const loaded = () => !!(systemInfo() && app() && repo())

  const [disableRefresh, setDisableRefresh] = createSignal(false)
  const refreshRepo = async () => {
    setDisableRefresh(true)
    setTimeout(() => setDisableRefresh(false), 3000)
    await client.refreshRepository({ repositoryId: repo().id })
    await refetchApp()
  }

  const refetchTimer = setInterval(refetchApp, 10000)
  onCleanup(() => clearInterval(refetchTimer))

  return (
    <Container>
      <Show when={loaded()}>
        <MainViewContainer gray>
          <MainView>
            <DataTable.Container>
              <DataTable.Title>Deployment</DataTable.Title>
              <DeploymentInfo
                app={app()}
                refetchApp={refetchApp}
                repo={repo()}
                refreshRepo={refreshRepo}
                disableRefresh={disableRefresh}
              />
            </DataTable.Container>
          </MainView>
        </MainViewContainer>
        <MainViewContainer>
          <MainView>
            <DataTable.Container>
              <DataTable.Title>Build Status</DataTable.Title>
              <BuildStatusTable
                app={app()}
                refetchApp={refetchApp}
                repo={repo()}
                refreshRepo={refreshRepo}
                disableRefresh={disableRefresh}
              />
            </DataTable.Container>
            <DataTable.Container>
              <DataTable.Title>Information</DataTable.Title>
              <Information app={app()} />
            </DataTable.Container>
            <Show when={app().deployType === DeployType.RUNTIME}>
              <DataTable.Container>
                <DataTable.Title>Usage</DataTable.Title>
                <Metrics app={app()} />
              </DataTable.Container>
            </Show>
            <Show when={app().deployType === DeployType.RUNTIME}>
              <DataTable.Container>
                <DataTable.Title>Container Log</DataTable.Title>
                <Logs app={app()} />
              </DataTable.Container>
            </Show>
          </MainView>
        </MainViewContainer>
      </Show>
    </Container>
  )
}
