import { Component, For, createSignal } from 'solid-js'
import toast from 'solid-toast'
import { Application, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { client, handleAPIError } from '/@/libs/api'
import { ApplicationState, deploymentState } from '/@/libs/application'
import { List } from '../List'
import { BuildRow } from '../build/BuildRow'

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

  // 最新5件のビルド
  const latestBuilds = () => props.sortedBuilds.slice(0, 4)

  return (
    <List.Container>
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
        {(build) => {
          const deployState = deploymentState(props.app)
          const isCurrentBuild = build.id === props.app.currentBuild
          const isDeploying = isCurrentBuild && deployState === ApplicationState.Deploying
          const isDeployed =
            isCurrentBuild && (deployState === ApplicationState.Running || deployState === ApplicationState.Static)

          return <BuildRow build={build} isDeployed={isDeployed ? 'deployed' : isDeploying ? 'deploying' : undefined} />
        }}
      </For>
    </List.Container>
  )
}

export default AppLatestBuilds
