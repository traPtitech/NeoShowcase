import { Timestamp } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, Show } from 'solid-js'
import toast from 'solid-toast'
import { Application, Build, BuildStatus, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import JumpButton from '/@/components/UI/JumpButton'
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

const BuildStatusTable: Component<{
  app: Application
  repo: Repository
  refreshRepo?: () => void
  disableRefresh?: () => boolean
  build: Build
  refetchBuild: () => void
  hasPermission: boolean
}> = (props) => {
  const rebuild = async () => {
    try {
      await client.retryCommitBuild({
        applicationId: props.app.id,
        commit: props.build.commit,
      })
      await props.refetchBuild()
      toast.success('再ビルドを開始しました')
    } catch (e) {
      handleAPIError(e, '再ビルドに失敗しました')
    }
  }
  const cancelBuild = async () => {
    try {
      await client.cancelBuild({
        buildId: props.build?.id,
      })
      await props.refetchBuild()
      toast.success('ビルドをキャンセルしました')
    } catch (e) {
      handleAPIError(e, 'ビルドのキャンセルに失敗しました')
    }
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
            disabled={props.disableRefresh?.()}
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
          <Button variants="borderError" size="small" onClick={cancelBuild} disabled={props.disableRefresh?.()}>
            Cancel Build
          </Button>
        </Show>
      </BuildStatusRow>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Repository</List.RowTitle>
          <List.RowData>{props.repo.name}</List.RowData>
        </List.RowContent>
        <JumpButton href={`/repos/${props.repo.id}`} />
      </List.Row>
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Source Commit</List.RowTitle>
          <List.RowData>{shortSha(props.build.commit)}</List.RowData>
        </List.RowContent>
        <Show when={props.refreshRepo && props.hasPermission}>
          <Button
            variants="ghost"
            size="medium"
            onClick={props.refreshRepo}
            disabled={props.disableRefresh?.()}
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
                  props.build.finishedAt?.timestamp?.toDate().getTime() -
                    props.build.startedAt?.timestamp?.toDate().getTime(),
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
