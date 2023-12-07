import { Component, Show } from 'solid-js'
import { Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import Code from '/@/components/UI/Code'
import { ToolTip } from '/@/components/UI/ToolTip'
import { systemInfo } from '/@/libs/api'
import { diffHuman, shortSha } from '/@/libs/format'
import { List } from '../List'

import { Button } from '/@/components/UI/Button'

const AppInfoLists: Component<{
  app: Application
  refreshCommit: () => void
  disableRefreshCommit: boolean
  hasPermission: boolean
}> = (props) => {
  const sshAccessCommand = () => `ssh -p ${systemInfo()?.ssh?.port} ${props.app.id}@${systemInfo()?.ssh?.host}`

  return (
    <>
      <List.Container>
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
        <List.Row>
          <List.RowContent>
            <List.RowTitle>Branch (Commit)</List.RowTitle>
            <List.RowData>{`${props.app.refName} (${shortSha(props.app.commit)})`}</List.RowData>
          </List.RowContent>
          <Show when={props.hasPermission}>
            <Button
              variants="ghost"
              size="medium"
              onClick={props.refreshCommit}
              disabled={props.disableRefreshCommit}
              tooltip={{
                props: {
                  content: 'リポジトリの最新コミットを取得',
                },
              }}
            >
              Refresh Commit
            </Button>
          </Show>
        </List.Row>
        <Show when={props.app.deployType === DeployType.RUNTIME}>
          <List.Row>
            <List.RowContent>
              <List.RowTitle>SSH Access</List.RowTitle>
              <Code value={sshAccessCommand()} copyable />
              <List.RowData>アプリケーションが起動している間のみアクセス可能です</List.RowData>
            </List.RowContent>
          </List.Row>
        </Show>
      </List.Container>
    </>
  )
}
export default AppInfoLists
