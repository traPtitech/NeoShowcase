import { createSignal, For, onSettled, Show } from 'solid-js'
import { type Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { SectionBoundary } from '/@/components/layouts/SectionBoundary'
import { styled } from '/@/components/styled-components'
import AppBranchResolution from '/@/components/templates/app/AppBranchResolution'
import AppDeployInfo from '/@/components/templates/app/AppDeployInfo'
import AppLatestBuilds from '/@/components/templates/app/AppLatestBuilds'
import { AppMetrics } from '/@/components/templates/app/AppMetrics'
import { ContainerLog } from '/@/components/templates/app/ContainerLog'
import { Button } from '/@/components/UI/Button'
import { useApplicationData } from '/@/features/application/data'
import { client, handleAPIError } from '/@/libs/api'
import { availableMetrics } from '/@/libs/session'
import { toast } from '/@/libs/toast'

const MainViewContainer = styled('div', 'w-full bg-ui-primary px-8 pt-10 pb-18 max-md:px-4')

const MainView = styled('div', 'mx-auto flex w-full max-w-250 flex-col gap-8')

const Metrics = (props: { app: Application }) => {
  const [selected, setSelected] = createSignal<string>()
  const currentView = () => selected() ?? availableMetrics()[0]

  return (
    <div class="border-ui-border flex w-full flex-col gap-4 rounded-lg border px-5 py-4">
      <div class="flex">
        <For each={availableMetrics()}>
          {(metrics) => (
            <Button
              variants="text"
              size="medium"
              onClick={() => setSelected(metrics)}
              active={currentView() === metrics}
            >
              {metrics}
            </Button>
          )}
        </For>
      </div>
      <div class="aspect-ratio-[960/464] bg-ui-secondary h-auto w-full rounded-lg">
        <Show when={currentView()} keyed>
          {(metrics) => <AppMetrics appID={props.app.id} metricsName={metrics} />}
        </Show>
      </div>
    </div>
  )
}

export default () => {
  const { app, builds, commits, refetch, repo, hasPermission } = useApplicationData()

  const deployedBuild = () => builds().find((b) => b.id === app().currentBuild)

  const startApp = async () => {
    const wasRunning = app().running
    try {
      await client.startApplication({ id: app().id })
      refetch()
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
    await client.refreshRepository({ repositoryId: repo().id })
    // バックエンド側で非同期で取得されるので1秒程度待ってからリロード
    setTimeout(() => {
      setDisableRefreshCommit(false)
      refetch()
    }, 1000)
  }

  onSettled(() => {
    const timer = setInterval(refetch, 10000)
    return () => clearInterval(timer)
  })

  return (
    <div class="h-full w-full overflow-y-auto">
      <MainViewContainer class="bg-ui-background">
        <MainView>
          <SectionBoundary title="Deployment">
            <AppDeployInfo
              app={app()}
              refetch={refetch}
              startApp={startApp}
              deployedBuild={deployedBuild()}
              latestBuildId={builds()[0]?.id}
              hasPermission={hasPermission()}
            />
          </SectionBoundary>
        </MainView>
      </MainViewContainer>
      <MainViewContainer>
        <MainView>
          <SectionBoundary title="Branch Resolution">
            <AppBranchResolution
              app={app()}
              commits={commits()}
              refreshCommit={refreshCommit}
              disableRefreshCommit={disableRefreshCommit()}
              hasPermission={hasPermission()}
            />
          </SectionBoundary>
          <SectionBoundary title="Latest Builds">
            <AppLatestBuilds app={app()} startApp={startApp} hasPermission={hasPermission()} sortedBuilds={builds()} />
          </SectionBoundary>
          <Show when={app().deployType === DeployType.RUNTIME && hasPermission()}>
            <SectionBoundary title="Usage">
              <Metrics app={app()} />
            </SectionBoundary>
            <SectionBoundary title="Container Log">
              <div class="border-ui-border w-full rounded-lg border px-5 py-4">
                <ContainerLog appID={app().id} />
              </div>
            </SectionBoundary>
          </Show>
        </MainView>
      </MainViewContainer>
    </div>
  )
}
