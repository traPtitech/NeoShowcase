import { timestampDate } from '@bufbuild/protobuf/wkt'
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
      <div class="flex w-full items-center gap-1">
        <AiOutlineBranches size={20} />
        <div class="flex w-full items-center gap-1">
          {`${props.app.refName} → `}
          <Show when={isErrorCommit} fallback={shortSha(props.app.commit)}>
            <AppStatusIcon state={ApplicationState.Error} size={20} />
            Error
          </Show>
        </div>
      </div>
    )
    if (!c || !c.commitDate) {
      return base
    }

    const diff = diffHuman(timestampDate(c.commitDate))
    const localeString = timestampDate(c.commitDate).toLocaleString()
    return (
      <div class="flex w-full flex-col gap-1">
        {base}
        <div>
          <For each={c.message.split('\n')}>{(line) => <div class="truncate">{line}</div>}</For>
          <div class="caption-regular text-text-grey">
            {c.authorName}
            <span>, </span>
            <ToolTip props={{ content: localeString }}>
              <span>{diff()}</span>
            </ToolTip>
            <span>, </span>
            {shortSha(c.hash)}
          </div>
        </div>
      </div>
    )
  }

  return (
    <>
      <List.Container>
        <List.Row>
          <List.RowContent class="min-w-0 shrink-0 grow-1 basis-0">{commitDisplay()}</List.RowContent>
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
