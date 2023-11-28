import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import toast from 'solid-toast'
import { Application, Build, BuildStatus, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client, handleAPIError } from '/@/libs/api'
import { buildStatusStr, buildTypeStr } from '/@/libs/application'
import { shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { Button } from '../../UI/Button'
import JumpButton from '../../UI/JumpButton'
import { MaterialSymbols } from '../../UI/MaterialSymbols'
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
  showJumpButton?: boolean
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
        <Show when={props.showJumpButton}>
          <A href={`/apps/${props.app.id}/builds/${props.build.id}`}>
            <Button
              size="small"
              variants="border"
              rightIcon={<MaterialSymbols opticalSize={20}>arrow_outward</MaterialSymbols>}
            >
              View Details
            </Button>
          </A>
        </Show>
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
          <List.RowTitle>Source Branch (Commit)</List.RowTitle>
          <List.RowData>
            {props.app.refName} ({shortSha(props.app.commit)})
          </List.RowData>
        </List.RowContent>
        <JumpButton href={`/apps/${props.app.id}/settings`} />
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
      <List.Row>
        <List.RowContent>
          <List.RowTitle>Build Type</List.RowTitle>
          <List.RowData>{buildTypeStr[props.app.config.buildConfig.case]}</List.RowData>
        </List.RowContent>
        <JumpButton href={`/apps/${props.app.id}/settings/build`} />
      </List.Row>
    </List.Container>
  )
}

export default BuildStatusTable
