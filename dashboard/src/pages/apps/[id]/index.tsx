import { styled } from '@macaron-css/solid'
import { Component, For, Show, createResource, createSignal, onCleanup, useTransition } from 'solid-js'

import { Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import AppDeployInfo from '/@/components/templates/app/AppDeployInfo'
import AppInfoLists from '/@/components/templates/app/AppInfoLists'
import AppLatestBuilds from '/@/components/templates/app/AppLatestBuilds'
import { AppMetrics } from '/@/components/templates/app/AppMetrics'
import { ContainerLog } from '/@/components/templates/app/ContainerLog'
import { availableMetrics, client } from '/@/libs/api'
import { useApplicationData } from '/@/routes'
import { colorVars, media } from '/@/theme'

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
      [media.mobile]: {
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
  const metricsNames = () => availableMetrics()?.metricsNames ?? []
  const [currentView, setCurrentView] = createSignal(metricsNames()[0])

  return (
    <MetricsContainer>
      <MetricsTypeButtons>
        <For each={metricsNames()}>
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
        <For each={metricsNames()}>
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
  const { app, refetchApp, repo, hasPermission } = useApplicationData()

  const [builds, { refetch: refetchBuilds }] = createResource(
    () => app()?.id,
    (id) => client.getBuilds({ id }),
  )
  const sortedBuilds = () =>
    builds()?.builds.sort((b1, b2) => {
      return (b2.queuedAt?.toDate().getTime() ?? 0) - (b1.queuedAt?.toDate().getTime() ?? 0)
    })
  const latestBuild = () => sortedBuilds()?.[0]

  const loaded = () => !!(app() && repo())

  const [disableRefresh, setDisableRefresh] = createSignal(false)
  const refreshRepo = async () => {
    setDisableRefresh(true)
    setTimeout(() => setDisableRefresh(false), 3000)
    await client.refreshRepository({ repositoryId: repo()?.id })
    await refetchApp()
  }

  const refetchAppTimer = setInterval(refetchApp, 10000)
  onCleanup(() => clearInterval(refetchAppTimer))

  const refetchBuildsTimer = setInterval(refetchBuilds, 10000)
  onCleanup(() => clearInterval(refetchBuildsTimer))

  const [isPending] = useTransition()

  return (
    <SuspenseContainer isPending={isPending()}>
      <Container>
        <Show when={loaded()}>
          <MainViewContainer gray>
            <MainView>
              <DataTable.Container>
                <DataTable.Title>Deployment</DataTable.Title>
                <AppDeployInfo
                  app={app()!}
                  refetchApp={refetchApp}
                  repo={repo()!}
                  refreshRepo={refreshRepo}
                  disableRefresh={disableRefresh}
                  isLatestBuild={latestBuild()?.id === app()?.currentBuild}
                  hasPermission={hasPermission()}
                />
              </DataTable.Container>
            </MainView>
          </MainViewContainer>
          <MainViewContainer>
            <MainView>
              <DataTable.Container>
                <Show when={builds()}>
                  <DataTable.Title>Latest Builds</DataTable.Title>
                  <AppLatestBuilds
                    app={app()!}
                    refetchApp={refetchApp}
                    repo={repo()!}
                    hasPermission={hasPermission()}
                    sortedBuilds={sortedBuilds()!}
                    refetchBuilds={refetchBuilds}
                  />
                </Show>
              </DataTable.Container>
              <DataTable.Container>
                <DataTable.Title>Information</DataTable.Title>
                <AppInfoLists app={app()!} />
              </DataTable.Container>
              <Show when={app()?.deployType === DeployType.RUNTIME && hasPermission()}>
                <DataTable.Container>
                  <DataTable.Title>Usage</DataTable.Title>
                  <Metrics app={app()!} />
                </DataTable.Container>
              </Show>
              <Show when={app()?.deployType === DeployType.RUNTIME && hasPermission()}>
                <DataTable.Container>
                  <DataTable.Title>Container Log</DataTable.Title>
                  <Logs app={app()!} />
                </DataTable.Container>
              </Show>
            </MainView>
          </MainViewContainer>
        </Show>
      </Container>
    </SuspenseContainer>
  )
}
