import { type Component, For, createSignal } from 'solid-js'

import type { Application, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'

import { List } from '../List'
import { BuildRow } from '../build/BuildRow'

const AppLatestBuilds: Component<{
  app: Application
  refetch: () => Promise<void>
  repo: Repository
  startApp: () => Promise<void>
  sortedBuilds: Build[]
  hasPermission: boolean
}> = (props) => {
  const [disabled, setDisabled] = createSignal(false)

  const startApp = async () => {
    setDisabled(true)
    await props.startApp()
    setDisabled(false)
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
              onClick={startApp}
              leftIcon={<MaterialSymbols>add</MaterialSymbols>}
              loading={disabled()}
            >
              Start App to Trigger Builds
            </Button>
          </List.PlaceHolder>
        }
      >
        {(build) => <BuildRow build={build} isCurrent={build.id === props.app.currentBuild} />}
      </For>
    </List.Container>
  )
}

export default AppLatestBuilds
