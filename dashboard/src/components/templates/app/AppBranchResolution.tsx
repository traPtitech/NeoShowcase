import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { AiOutlineBranches } from 'solid-icons/ai'
import { type Component, For, Show } from 'solid-js'
import type { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import JumpButton from '/@/components/UI/JumpButton'
import { ToolTip } from '/@/components/UI/ToolTip'
import { List } from '/@/components/templates/List'
import { AppStatusIcon } from '/@/components/templates/app/AppStatusIcon'
import type { CommitsMap } from '/@/libs/api'
import { ApplicationState, errorCommit } from '/@/libs/application'
import { diffHuman, shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'

const DataRows = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignContent: 'left',
    gap: '4px',
  },
})

const DataRow = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    gap: '4px',
    alignItems: 'center',
  },
})

const greyText = style({
  color: colorVars.semantic.text.grey,
  ...textVars.caption.regular,
})

const noOverflow = style({
  overflowX: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const shrinkFirst = style({
  flexShrink: 0,
  flexGrow: 1,
  flexBasis: 0,
  minWidth: 0,
})

const AppBranchResolution: Component<{
  app: Application
  commits?: CommitsMap
  refreshCommit: () => void
  disableRefreshCommit: boolean
  hasPermission: boolean
}> = (props) => {
  const commit = () => props.commits?.[props.app.commit]
  const commitDisplay = () => {
    const c = commit()
    const isErrorCommit = props.app.commit === errorCommit
    const base = (
      <DataRow>
        <AiOutlineBranches size={20} />
        <DataRow>
          {`${props.app.refName} → `}
          <Show when={isErrorCommit} fallback={shortSha(props.app.commit)}>
            <AppStatusIcon state={ApplicationState.Error} size={20} />
            Error
          </Show>
        </DataRow>
      </DataRow>
    )
    if (!c || !c.commitDate) {
      return base
    }

    const { diff, localeString } = diffHuman(c.commitDate.toDate())
    return (
      <DataRows>
        {base}
        <div>
          <For each={c.message.split('\n')}>{(line) => <div class={noOverflow}>{line}</div>}</For>
          <div class={greyText}>
            {c.authorName}
            <span>, </span>
            <ToolTip props={{ content: localeString }}>
              <span>{diff}</span>
            </ToolTip>
            <span>, </span>
            {shortSha(c.hash)}
          </div>
        </div>
      </DataRows>
    )
  }

  return (
    <>
      <List.Container>
        <List.Row>
          <List.RowContent class={shrinkFirst}>{commitDisplay()}</List.RowContent>
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
          <JumpButton href={`/apps/${props.app.id}/settings`} tooltip="設定を変更" />
        </List.Row>
      </List.Container>
    </>
  )
}

export default AppBranchResolution
