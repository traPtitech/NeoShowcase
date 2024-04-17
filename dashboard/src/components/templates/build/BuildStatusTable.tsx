import type { Timestamp } from '@bufbuild/protobuf'
import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
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
import { colorVars, textVars } from '/@/theme'
import { List } from '../List'
import { BuildStatusIcon } from './BuildStatusIcon'

const BuildStatusRow = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

    background: colorVars.semantic.ui.secondary,
  },
})
const BuildStatusLabel = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '4px',

    color: colorVars.semantic.text.black,
    ...textVars.text.medium,
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

    const { diff, localeString } = diffHuman(c.commitDate.toDate())
    return (
      <DataRows>
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
    <List.Container>
      <BuildStatusRow>
        <BuildStatusLabel>
          <BuildStatusIcon state={props.build.status} size={24} />
          {buildStatusStr[props.build.status]}
        </BuildStatusLabel>
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
      </BuildStatusRow>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Source Commit</List.RowTitle>
          <List.RowData>{commitDisplay()}</List.RowData>
        </List.RowContent>
      </List.Row>
      <List.Columns>
        <Show when={props.build.queuedAt}>
          {(nonNullQueuedAt) => {
            const { diff, localeString } = diffHuman(nonNullQueuedAt().toDate())
            return (
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>キュー登録時刻</List.RowTitle>
                  <ToolTip props={{ content: localeString }}>
                    <List.RowData>{diff}</List.RowData>
                  </ToolTip>
                </List.RowContent>
              </List.Row>
            )
          }}
        </Show>
        <Show when={props.build.startedAt?.valid && props.build.startedAt} fallback={'-'}>
          {(nonNullStartedAt) => {
            const { diff, localeString } = diffHuman((nonNullStartedAt().timestamp as Timestamp).toDate())
            return (
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>ビルド開始時刻</List.RowTitle>
                  <ToolTip props={{ content: localeString }}>
                    <List.RowData>{diff}</List.RowData>
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
                const { diff, localeString } = diffHuman((nonNullFinishedAt().timestamp as Timestamp).toDate())
                return (
                  <ToolTip props={{ content: localeString }}>
                    <List.RowData>{diff}</List.RowData>
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
