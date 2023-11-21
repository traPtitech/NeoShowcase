import { Application, Build, BuildStatus, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { client, handleAPIError } from '/@/libs/api'
import { buildStatusStr, buildTypeStr } from '/@/libs/application'
import { shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import toast from 'solid-toast'
import { BuildStatusIcon } from '../UI/BuildStatusIcon'
import { Button } from '../UI/Button'
import JumpButton from '../UI/JumpButton'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { List } from './List'

const PlaceHolder = styled('div', {
  base: {
    width: '100%',
    height: '400px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
    alignItems: 'center',
    justifyContent: 'center',

    color: colorVars.semantic.text.black,
    ...textVars.h4.medium,
  },
})
const BuildStatusRow = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
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
  refetchApp: () => void
  repo: Repository
  refreshRepo: () => void
  disableRefresh: () => boolean
  latestBuild?: Build
  refetchLatestBuild: () => void
}> = (props) => {
  const startApp = async () => {
    try {
      await client.startApplication({ id: props.app.id })
      await props.refetchApp()
      toast.success('アプリケーションを起動しました')
    } catch (e) {
      handleAPIError(e, 'アプリケーションの再起動に失敗しました')
    }
  }
  const rebuild = async () => {
    try {
      await client.retryCommitBuild({
        applicationId: props.app.id,
        commit: props.app.commit,
      })
      await props.refetchLatestBuild()
    } catch (e) {
      handleAPIError(e, '再ビルドに失敗しました')
    }
  }
  const cancelBuild = async () => {
    try {
      await client.cancelBuild({
        buildId: props.latestBuild?.id,
      })
      await props.refetchLatestBuild()
    } catch (e) {
      handleAPIError(e, 'ビルドのキャンセルに失敗しました')
    }
  }

  return (
    <Show
      when={props.latestBuild}
      fallback={
        <List.Container>
          <PlaceHolder>
            <MaterialSymbols displaySize={80}>deployed_code</MaterialSymbols>
            No Builds
            <Button
              variants="primary"
              size="medium"
              onClick={startApp}
              disabled={props.disableRefresh()}
              leftIcon={<MaterialSymbols>add</MaterialSymbols>}
            >
              Build and Start App
            </Button>
          </PlaceHolder>
        </List.Container>
      }
    >
      {(nonNullLatestBuild) => (
        <List.Container>
          <BuildStatusRow>
            <BuildStatusLabel>
              <BuildStatusIcon state={nonNullLatestBuild().status} size={24} />
              {buildStatusStr[nonNullLatestBuild().status]}
            </BuildStatusLabel>
            <Show when={!nonNullLatestBuild().retriable}>
              <Button
                variants="borderError"
                size="small"
                onClick={rebuild}
                disabled={props.disableRefresh()}
                tooltip={{
                  props: {
                    content: '同じコミットで再ビルド',
                  },
                }}
              >
                Rebuild
              </Button>
            </Show>
            <Show when={nonNullLatestBuild().status === BuildStatus.BUILDING}>
              <Button variants="borderError" size="small" onClick={cancelBuild} disabled={props.disableRefresh()}>
                Cancel Build
              </Button>
            </Show>
          </BuildStatusRow>
          <List.Row>
            <List.RowContent>
              <List.RowTitle>Latest Build ID</List.RowTitle>
              <List.RowData>{nonNullLatestBuild().id}</List.RowData>
            </List.RowContent>
            <JumpButton href={`/apps/${props.app.id}/builds/${nonNullLatestBuild().id}`} />
          </List.Row>
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
            <Button
              variants="ghost"
              size="medium"
              onClick={props.refreshRepo}
              disabled={props.disableRefresh()}
              tooltip={{
                props: {
                  content: 'リポジトリの最新コミットを取得',
                },
              }}
            >
              Refresh Commit
            </Button>
          </List.Row>
          <List.Row>
            <List.RowContent>
              <List.RowTitle>Build Type</List.RowTitle>
              <List.RowData>{buildTypeStr[props.app.config.buildConfig.case]}</List.RowData>
            </List.RowContent>
            <JumpButton href={`/apps/${props.app.id}/settings/build`} />
          </List.Row>
        </List.Container>
      )}
    </Show>
  )
}

export default BuildStatusTable
