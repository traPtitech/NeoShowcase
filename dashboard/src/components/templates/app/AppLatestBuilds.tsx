import { type Component, createSignal, For } from 'solid-js'

import type { Application, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'

import { useApplicationData } from '/@/routes'
import { BuildRow } from '../build/BuildRow'
import { List } from '../List'

const AppLatestBuilds: Component<{
  app: Application
  refetch: () => Promise<void>
  repo: Repository
  startApp: () => Promise<void>
  sortedBuilds: Build[]
  hasPermission: boolean
}> = (props) => {
  const { commits } = useApplicationData()

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
            <div class="i-material-symbols:deployed-code-outline shrink-0 text-20/20" />
            No Builds
            <Button
              variants="primary"
              size="medium"
              onClick={startApp}
              leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
              loading={disabled()}
            >
              Start App to Trigger Builds
            </Button>
          </List.PlaceHolder>
        }
      >
        {(build) => <BuildRow build={build} commits={commits()} isCurrent={build.id === props.app.currentBuild} />}
      </For>
    </List.Container>
  )
}

export default AppLatestBuilds
