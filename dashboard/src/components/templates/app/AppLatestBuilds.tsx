import { styled } from '@macaron-css/solid'
import { Component, For, createSignal } from 'solid-js'
import { Show } from 'solid-js'
import toast from 'solid-toast'
import { Application, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { client, handleAPIError } from '/@/libs/api'
import { ApplicationState, applicationState } from '/@/libs/application'
import { colorVars } from '/@/theme'
import { List } from '../List'
import { BuildRow } from '../build/BuildRow'

const Bordered = styled('div', {
  base: {
    borderBottom: `2px solid ${colorVars.semantic.ui.border}`,
  },
})

const AppLatestBuilds: Component<{
  app: Application
  refetchApp: () => void
  repo: Repository
  sortedBuilds: Build[]
  refetchBuilds: () => void
  hasPermission: boolean
}> = (props) => {
  const [disabled, setDisabled] = createSignal(false)

  const StartApp = async () => {
    try {
      setDisabled(true)
      await client.startApplication({ id: props.app.id })
      await Promise.all([props.refetchApp(), props.refetchBuilds()])
      toast.success('アプリケーションを再起動しました')
    } catch (e) {
      handleAPIError(e, 'アプリケーションの再起動に失敗しました')
      setDisabled(false)
    }
  }

  // 最新4件のビルド
  const latestBuilds = () => props.sortedBuilds.slice(0, 4)

  // 最新4件のビルドにデプロイ済みビルドが含まれるか
  const hasDeployedBuild = () => latestBuilds()?.some((build) => build.id === props.app.currentBuild)

  return (
    <List.Container>
      <Show when={!hasDeployedBuild() && props.sortedBuilds.find((b) => b.id === props.app.currentBuild)}>
        {(deployedBuild) => (
          <Bordered>
            <BuildRow build={deployedBuild()} isDeployed="deployed" />
          </Bordered>
        )}
      </Show>
      <For
        each={latestBuilds()}
        fallback={
          <List.PlaceHolder>
            <MaterialSymbols displaySize={80}>deployed_code</MaterialSymbols>
            No Builds
            <Button
              variants="primary"
              size="medium"
              onClick={StartApp}
              leftIcon={<MaterialSymbols>add</MaterialSymbols>}
              loading={disabled()}
            >
              Build and Start App
            </Button>
          </List.PlaceHolder>
        }
      >
        {(build, i) => {
          const isCurrentBuild = build.id === props.app.currentBuild
          const isDeploying = isCurrentBuild && applicationState(props.app) === ApplicationState.Deploying
          const isDeployed =
            isCurrentBuild &&
            (applicationState(props.app) === ApplicationState.Running ||
              applicationState(props.app) === ApplicationState.Static)

          return (
            <BuildRow
              build={build}
              isDeployed={isDeployed ? 'deployed' : isDeploying ? 'deploying' : undefined}
              isLatest={i() === 0}
            />
          )
        }}
      </For>
    </List.Container>
  )
}

export default AppLatestBuilds
