import { type Component, Show } from 'solid-js'
import { type Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import Code from '/@/components/UI/Code'
import { type CommitsMap, systemInfo } from '/@/libs/api'
import { ApplicationState, deploymentState } from '/@/libs/application'
import { List } from '../List'

const AppInfoLists: Component<{
  app: Application
  commits?: CommitsMap
  refreshCommit: () => void
  disableRefreshCommit: boolean
  hasPermission: boolean
}> = (props) => {
  const sshAccessCommand = () => `ssh -p ${systemInfo()?.ssh?.port} ${props.app.id}@${systemInfo()?.ssh?.host}`

  return (
    <>
      <List.Container>
        <Show when={props.app.deployType === DeployType.RUNTIME}>
          <List.Row>
            <List.RowContent>
              <List.RowTitle>SSH Access</List.RowTitle>
              <Code value={sshAccessCommand()} copyable />
              <Show when={deploymentState(props.app) !== ApplicationState.Running}>
                <List.RowData>現在アプリが起動していないためSSHアクセスはできません</List.RowData>
              </Show>
            </List.RowContent>
          </List.Row>
        </Show>
      </List.Container>
    </>
  )
}
export default AppInfoLists
