import { styled } from '@macaron-css/solid'
import { type Component, For, Show, createSignal, onCleanup, useTransition } from 'solid-js'
import toast from 'solid-toast'
import { type Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import AppBranchResolution from '/@/components/templates/app/AppBranchResolution'
import AppDeployInfo from '/@/components/templates/app/AppDeployInfo'
import AppInfoLists from '/@/components/templates/app/AppInfoLists'
import AppLatestBuilds from '/@/components/templates/app/AppLatestBuilds'
import { AppMetrics } from '/@/components/templates/app/AppMetrics'
import { ContainerLog } from '/@/components/templates/app/ContainerLog'
import { availableMetrics, client, handleAPIError } from '/@/libs/api'
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
      <ContainerLog appID={props.app.id} />
    </LogContainer>
  )
}

export default () => {
  const { app, builds, commits, refetch, repo, hasPermission } = useApplicationData()

  const sortedBuilds = () =>
    builds()?.sort((b1, b2) => {
      return (b2.queuedAt?.toDate().getTime() ?? 0) - (b1.queuedAt?.toDate().getTime() ?? 0)
    })
  const deployedBuild = () => sortedBuilds()?.find((b) => b.id === app()?.currentBuild)
  const latestBuild = () => sortedBuilds()?.[0]

  const loaded = () => !!(app() && repo())

  const startApp = async () => {
    const wasRunning = app()?.running
    try {
      await client.startApplication({ id: app()?.id })
      await refetch()
      // 非同期でビルドが開始されるので1秒程度待ってから再度リロード
      setTimeout(refetch, 1000)
      toast.success(`アプリケーションを${wasRunning ? '再' : ''}起動しました`)
    } catch (e) {
      handleAPIError(e, `アプリケーションの${wasRunning ? '再' : ''}起動に失敗しました`)
    }
  }

  const [disableRefreshCommit, setDisableRefreshCommit] = createSignal(false)
  const refreshCommit = async () => {
    setDisableRefreshCommit(true)
    await client.refreshRepository({ repositoryId: repo()?.id })
    setTimeout(() => {
      // バックエンド側で非同期で取得されるので1秒程度待ってからリロード
      setDisableRefreshCommit(false)
      void refetch()
    }, 1000)
  }

  const refetchTimer = setInterval(refetch, 10000)
  onCleanup(() => clearInterval(refetchTimer))

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
                  refetch={refetch}
                  repo={repo()!}
                  startApp={startApp}
                  deployedBuild={deployedBuild()}
                  latestBuildId={latestBuild()?.id}
                  hasPermission={hasPermission()}
                />
              </DataTable.Container>
            </MainView>
          </MainViewContainer>
          <MainViewContainer>
            <MainView>
              <DataTable.Container>
                <DataTable.Title>Branch Resolution</DataTable.Title>
                <AppBranchResolution
                  app={app()!}
                  commits={commits()}
                  refreshCommit={refreshCommit}
                  disableRefreshCommit={disableRefreshCommit()}
                  hasPermission={hasPermission()}
                />
              </DataTable.Container>
              <DataTable.Container>
                <Show when={builds()}>
                  <DataTable.Title>Latest Builds</DataTable.Title>
                  <AppLatestBuilds
                    app={app()!}
                    refetch={refetch}
                    repo={repo()!}
                    startApp={startApp}
                    hasPermission={hasPermission()}
                    sortedBuilds={sortedBuilds()!}
                  />
                </Show>
              </DataTable.Container>
              <DataTable.Container>
                <DataTable.Title>Information</DataTable.Title>
                <AppInfoLists
                  app={app()!}
                  commits={commits()}
                  refreshCommit={refreshCommit}
                  disableRefreshCommit={disableRefreshCommit()}
                  hasPermission={hasPermission()}
                />
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
