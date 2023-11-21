import { Application, Build, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import AppDeployInfoTable from '/@/components/templates/AppDeployInfoTable'
import AppInfoTable from '/@/components/templates/AppInfoTable'
import { AppMetrics } from '/@/components/templates/AppMetrics'
import BuildStatusTable from '/@/components/templates/BuildStatusTable'
import { ContainerLog } from '/@/components/templates/ContainerLog'
import { availableMetrics, client, systemInfo } from '/@/libs/api'
import { useApplicationData } from '/@/routes'
import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, For, Show, createResource, createSignal, onCleanup } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    overflowY: 'auto',
    scrollbarGutter: 'stable',
  },
})
const MainViewContainer = styled('div', {
  base: {
    width: '100%',
    padding: '40px 32px 72px',

    '@media': {
      'screen and (max-width: 768px)': {
        padding: '40px 16px 72px',
      },
    },
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
              variants="text"
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

const getLatestBuild = (appId: Application['id']): Promise<Build | undefined> =>
  client
    .getBuilds({ id: appId })
    .then((res) => res.builds.sort((b1, b2) => b2.queuedAt.toDate().getTime() - b1.queuedAt.toDate().getTime())[0])

export default () => {
  const { app, refetchApp, repo, hasPermission } = useApplicationData()

  const [latestBuild, { refetch: refetchLatestBuild }] = createResource(
    () => app().id,
    (appId) => getLatestBuild(appId),
  )

  const loaded = () => !!(systemInfo() && app() && repo())

  const [disableRefresh, setDisableRefresh] = createSignal(false)
  const refreshRepo = async () => {
    setDisableRefresh(true)
    setTimeout(() => setDisableRefresh(false), 3000)
    await client.refreshRepository({ repositoryId: repo().id })
    await refetchApp()
  }

  const refetchAppTimer = setInterval(refetchApp, 10000)
  onCleanup(() => clearInterval(refetchAppTimer))

  const refetchLatestBuildTimer = setInterval(refetchLatestBuild, 10000)
  onCleanup(() => clearInterval(refetchLatestBuildTimer))

  return (
    <Container>
      <Show when={loaded()}>
        <MainViewContainer gray>
          <MainView>
            <DataTable.Container>
              <DataTable.Title>Deployment</DataTable.Title>
              <AppDeployInfoTable
                app={app()}
                refetchApp={refetchApp}
                repo={repo()}
                refreshRepo={refreshRepo}
                disableRefresh={disableRefresh}
                latestBuildId={latestBuild()?.id}
                hasPermission={hasPermission()}
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
                latestBuild={latestBuild()}
                refetchLatestBuild={refetchLatestBuild}
                hasPermission={hasPermission()}
              />
            </DataTable.Container>
            <DataTable.Container>
              <DataTable.Title>Information</DataTable.Title>
              <AppInfoTable app={app()} />
            </DataTable.Container>
            <Show when={app().deployType === DeployType.RUNTIME && hasPermission()}>
              <DataTable.Container>
                <DataTable.Title>Usage</DataTable.Title>
                <Metrics app={app()} />
              </DataTable.Container>
            </Show>
            <Show when={app().deployType === DeployType.RUNTIME && hasPermission()}>
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
