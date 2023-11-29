import { Component, Show } from 'solid-js'
import { Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'
import { diffHuman } from '/@/libs/format'
import Code from '../../UI/Code'
import { ToolTip } from '../../UI/ToolTip'
import { List } from '../List'

const AppInfoLists: Component<{ app: Application }> = (props) => {
  const sshAccessCommand = () => `ssh -p ${systemInfo().ssh.port} ${props.app.id}@${systemInfo().ssh.host}`

  return (
    <>
      <List.Container>
        <List.Columns>
          <Show when={props.app.updatedAt}>
            {(nonNullUpdatedAt) => {
              const { diff, localeString } = diffHuman(nonNullUpdatedAt().toDate())

              return (
                <List.Row>
                  <List.RowContent>
                    <List.RowTitle>起動時刻</List.RowTitle>
                    <ToolTip props={{ content: localeString }}>
                      <List.RowData>{diff}</List.RowData>
                    </ToolTip>
                  </List.RowContent>
                </List.Row>
              )
            }}
          </Show>
          <Show when={props.app.createdAt}>
            {(nonNullCreatedAt) => {
              const { diff, localeString } = diffHuman(nonNullCreatedAt().toDate())
              return (
                <List.Row>
                  <List.RowContent>
                    <List.RowTitle>作成日</List.RowTitle>
                    <ToolTip props={{ content: localeString }}>
                      <List.RowData>{diff}</List.RowData>
                    </ToolTip>
                  </List.RowContent>
                </List.Row>
              )
            }}
          </Show>
        </List.Columns>
        <Show when={props.app.deployType === DeployType.RUNTIME}>
          <List.Row>
            <List.RowContent>
              <List.RowTitle>SSH Access</List.RowTitle>
              <Show
                when={props.app.running}
                fallback={<List.RowData>アプリケーションが起動している間のみSSHでアクセス可能です</List.RowData>}
              >
                <List.RowData>
                  <Code value={sshAccessCommand()} copyable />
                </List.RowData>
              </Show>
            </List.RowContent>
          </List.Row>
        </Show>
      </List.Container>
    </>
  )
}
export default AppInfoLists
