import type { Timestamp } from '@bufbuild/protobuf'
import { useNavigate } from '@solidjs/router'
import { type Component, For, Show } from 'solid-js'
import toast from 'solid-toast'
import {
  type Application,
  type Build,
  BuildStatus,
  type Repository,
  type SimpleCommit,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { ToolTip } from '/@/components/UI/ToolTip'
import { client, handleAPIError } from '/@/libs/api'
import { buildStatusStr } from '/@/libs/application'
import { diffHuman, durationHuman, shortSha } from '/@/libs/format'
import { List } from '../List'
import { BuildStatusIcon } from './BuildStatusIcon'

const BuildStatusTable: Component<{
  app: Application
  repo: Repository
  build: Build
  commit?: SimpleCommit
  refetch: () => Promise<void>
  hasPermission: boolean
}> = (props) => {
  const navigate = useNavigate()

  const rebuild = async () => {
    try {
      await client.retryCommitBuild({
        applicationId: props.app.id,
        commit: props.build.commit,
      })
      toast.success('再ビルドを開始しました')
      // 非同期でビルドが開始されるので1秒程度待ってから遷移
      setTimeout(() => {
        void props.refetch()
        navigate(`/apps/${props.app.id}`)
      }, 1000)
    } catch (e) {
      handleAPIError(e, '再ビルドに失敗しました')
    }
  }
  const cancelBuild = async () => {
    try {
      await client.cancelBuild({
        buildId: props.build?.id,
      })
      await props.refetch()
      toast.success('ビルドをキャンセルしました')
    } catch (e) {
      handleAPIError(e, 'ビルドのキャンセルに失敗しました')
    }
  }

  const commitDisplay = () => {
    const c = props.commit
    if (!c || !c.commitDate) {
      return shortSha(props.build.commit)
    }

    const diff = diffHuman(c.commitDate.toDate())
    const localeString = c.commitDate.toDate().toLocaleString()
    return (
      <div class="flex flex-col">
        <For each={c.message.split('\n')}>{(line) => <div>{line}</div>}</For>
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
    )
  }

  return (
    <List.Container>
      <div class="flex w-full items-center gap-2 bg-ui-secondary px-5 py-4">
        <div class="text flex w-full items-center gap-1 text-black text-medium">
          <BuildStatusIcon state={props.build.status} size={24} />
          {buildStatusStr[props.build.status]}
        </div>
        <Show when={!props.build.retriable && props.hasPermission}>
          <Button
            variants="borderError"
            size="small"
            onClick={rebuild}
            tooltip={{
              props: {
                content: '同じコミットで再ビルド',
              },
            }}
          >
            Rebuild
          </Button>
        </Show>
        <Show when={props.build.status === BuildStatus.BUILDING && props.hasPermission}>
          <Button variants="borderError" size="small" onClick={cancelBuild}>
            Cancel Build
          </Button>
        </Show>
      </div>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Source Commit</List.RowTitle>
          <List.RowData>{commitDisplay()}</List.RowData>
        </List.RowContent>
      </List.Row>
      <List.Columns>
        <Show when={props.build.queuedAt}>
          {(nonNullQueuedAt) => {
            const diff = diffHuman(nonNullQueuedAt().toDate())
            const localeString = nonNullQueuedAt().toDate().toLocaleString()
            return (
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>キュー登録時刻</List.RowTitle>
                  <ToolTip props={{ content: localeString }}>
                    <List.RowData>{diff()}</List.RowData>
                  </ToolTip>
                </List.RowContent>
              </List.Row>
            )
          }}
        </Show>
        <Show when={props.build.startedAt?.valid && props.build.startedAt} fallback={'-'}>
          {(nonNullStartedAt) => {
            const ts = (nonNullStartedAt().timestamp as Timestamp).toDate()
            const diff = diffHuman(ts)
            const localeString = ts.toLocaleString()
            return (
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>ビルド開始時刻</List.RowTitle>
                  <ToolTip props={{ content: localeString }}>
                    <List.RowData>{diff()}</List.RowData>
                  </ToolTip>
                </List.RowContent>
              </List.Row>
            )
          }}
        </Show>
      </List.Columns>
      <List.Columns>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>ビルド終了時刻</List.RowTitle>
            <Show when={props.build.finishedAt?.valid && props.build.finishedAt} fallback={'-'}>
              {(nonNullFinishedAt) => {
                const ts = (nonNullFinishedAt().timestamp as Timestamp).toDate()
                const diff = diffHuman(ts)
                const localeString = ts.toLocaleString()
                return (
                  <ToolTip props={{ content: localeString }}>
                    <List.RowData>{diff()}</List.RowData>
                  </ToolTip>
                )
              }}
            </Show>
          </List.RowContent>
        </List.Row>
        <List.Row>
          <List.RowContent>
            <List.RowTitle>ビルド時間</List.RowTitle>
            <Show when={props.build.finishedAt?.valid && props.build.startedAt?.valid} fallback={'-'}>
              <List.RowData>
                {durationHuman(
                  props.build.finishedAt!.timestamp!.toDate().getTime() -
                    props.build.startedAt!.timestamp!.toDate().getTime(),
                )}
              </List.RowData>
            </Show>
          </List.RowContent>
        </List.Row>
      </List.Columns>
    </List.Container>
  )
}

export default BuildStatusTable
