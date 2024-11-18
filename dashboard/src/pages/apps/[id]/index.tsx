import { type Component, For, Show, createSignal, onCleanup, useTransition } from 'solid-js'
import toast from 'solid-toast'
import { type Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import { styled } from '/@/components/styled-components'
import AppBranchResolution from '/@/components/templates/app/AppBranchResolution'
import AppDeployInfo from '/@/components/templates/app/AppDeployInfo'
import AppLatestBuilds from '/@/components/templates/app/AppLatestBuilds'
import { AppMetrics } from '/@/components/templates/app/AppMetrics'
import { ContainerLog } from '/@/components/templates/app/ContainerLog'
import { availableMetrics, client, handleAPIError } from '/@/libs/api'
import { useApplicationData } from '/@/routes'

const MainViewContainer = styled('div', 'w-full bg-ui-primary px-8 pt-10 pb-18 max-md:px-4')

const MainView = styled('div', 'mx-auto flex w-full max-w-250 flex-col gap-8')

const Metrics: Component<{ app: Application }> = (props) => {
  const metricsNames = () => availableMetrics()?.metricsNames ?? []
  const [currentView, setCurrentView] = createSignal(metricsNames()[0])

  return (
    <div class="flex w-full flex-col gap-4 rounded-lg border border-ui-border px-5 py-4">
      <div class="flex">
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
      </div>
      <div class="aspect-ratio-[960/464] h-auto w-full rounded-lg bg-ui-secondary">
        <For each={metricsNames()}>
          {(metrics) => (
            <Show when={currentView() === metrics}>
              <AppMetrics appID={props.app.id} metricsName={metrics} />
            </Show>
          )}
        </For>
      </div>
    </div>
  )
}

const Logs: Component<{ app: Application }> = (props) => {
  return (
    <div class="w-full rounded-lg border border-ui-border px-5 py-4">
      <ContainerLog appID={props.app.id} />
    </div>
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
      <div class="h-full w-full overflow-y-auto">
        <Show when={loaded()}>
          <MainViewContainer class="bg-ui-background">
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
      </div>
    </SuspenseContainer>
  )
}
