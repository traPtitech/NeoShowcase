import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { AiOutlineBranches } from 'solid-icons/ai'
import { type Component, For, Show } from 'solid-js'
import { type Application, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import Code from '/@/components/UI/Code'
import { ToolTip } from '/@/components/UI/ToolTip'
import { type CommitsMap, systemInfo } from '/@/libs/api'
import { ApplicationState, deploymentState } from '/@/libs/application'
import { diffHuman, shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { List } from '../List'

const DataRow = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '4px',
    alignItems: 'center',
  },
})

const DataRows = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    alignContent: 'left',
  },
})

const greyText = style({
  color: colorVars.semantic.text.grey,
  ...textVars.caption.regular,
})

const AppInfoLists: Component<{
  app: Application
  commits?: CommitsMap
  refreshCommit: () => void
  disableRefreshCommit: boolean
  hasPermission: boolean
}> = (props) => {
  const sshAccessCommand = () => `ssh -p ${systemInfo()?.ssh?.port} ${props.app.id}@${systemInfo()?.ssh?.host}`

  const commit = () => props.commits?.[props.app.commit]
  const commitDisplay = () => {
    const c = commit()
    if (!c || !c.commitDate) {
      return (
        <DataRow>
          <AiOutlineBranches />
          <div>{`${props.app.refName} (${shortSha(props.app.commit)})`}</div>
        </DataRow>
      )
    }

    const { diff, localeString } = diffHuman(c.commitDate.toDate())
    return (
      <DataRows>
        <DataRow>
          <AiOutlineBranches />
          <div>{props.app.refName}</div>
        </DataRow>
        <For each={c.message.split('\n')}>{(line) => <div>{line}</div>}</For>
        <div class={greyText}>
          {c.authorName}
          <span>, </span>
          <ToolTip props={{ content: localeString }}>
            <span>{diff}</span>
          </ToolTip>
          <span>, </span>
          {shortSha(c.hash)}
        </div>
      </DataRows>
    )
  }

  return (
    <>
      <List.Container>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>Branch (Commit)</List.RowTitle>
            <List.RowData>{commitDisplay()}</List.RowData>
          </List.RowContent>
          <Show when={props.hasPermission}>
            <Button
              variants="ghost"
              size="medium"
              onClick={props.refreshCommit}
              disabled={props.disableRefreshCommit}
              tooltip={{
                props: {
                  content: (
                    <>
                      <div>リポジトリの最新コミットを取得</div>
                      <div>Pro Tip: Webhookを設定することで</div>
                      <div>同期して更新されます</div>
                    </>
                  ),
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
