import { BuildStatus } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { BuildStatusIcon } from '/@/components/UI/BuildStatusIcon'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { ToolTip } from '/@/components/UI/ToolTip'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { ArtifactRow } from '/@/components/templates/ArtifactRow'
import { BuildLog } from '/@/components/templates/BuildLog'
import { Bordered, List } from '/@/components/templates/List'
import { client, handleAPIError } from '/@/libs/api'
import { buildStatusStr, buildTypeStr } from '/@/libs/application'
import { diffHuman, durationHuman, shortSha } from '/@/libs/format'
import { useBuildData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { Timestamp } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { A, useNavigate } from '@solidjs/router'
import { For, Show, VoidComponent, createResource, createSignal } from 'solid-js'

const MainView = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '32px',
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
const JumpButtonContainer = styled('div', {
  base: {
    width: '32px',
    height: '32px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',

    flexShrink: 0,
    background: 'none',
    border: 'none',
    borderRadius: '6px',
    cursor: 'pointer',
    color: colorVars.semantic.text.black,
    selectors: {
      '&:hover': {
        background: colorVars.semantic.transparent.primaryHover,
      },
      '&:active, &[data-active="true"]': {
        color: colorVars.semantic.primary.main,
        background: colorVars.semantic.transparent.primarySelected,
      },
      '&:disabled': {
        cursor: 'not-allowed',
        border: 'none !important',
        color: `${colorVars.semantic.text.black} !important`,
        background: `${colorVars.semantic.text.disabled} !important`,
      },
    },
  },
})
const JumpButton: VoidComponent<{ href: string }> = (props) => (
  <A href={props.href}>
    <JumpButtonContainer>
      <MaterialSymbols opticalSize={20}>arrow_outward</MaterialSymbols>
    </JumpButtonContainer>
  </A>
)
const LogContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',

    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
  },
})

export default () => {
  const navigate = useNavigate()
  const { app, refetchApp, build, refetchBuild } = useBuildData()
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const [disableRefresh, setDisableRefresh] = createSignal(false)
  const refreshRepo = async () => {
    setDisableRefresh(true)
    setTimeout(() => setDisableRefresh(false), 3000)
    await client.refreshRepository({ repositoryId: repo()?.id })
    await refetchApp()
  }

  const loaded = () => !!(app() && build())

  const buildFinished = () => build()?.finishedAt?.valid

  const rebuild = async () => {
    try {
      await client.retryCommitBuild({
        applicationId: app()?.id,
        commit: build()?.commit,
      })
      navigate(`/apps/${app()?.id}/builds`)
    } catch (e) {
      handleAPIError(e, '再ビルドに失敗しました')
    }
  }

  const cancelBuild = async () => {
    try {
      await client.cancelBuild({
        buildId: build()?.id,
      })
      await refetchBuild()
    } catch (e) {
      handleAPIError(e, 'ビルドのキャンセルに失敗しました')
    }
  }

  return (
    <MainViewContainer>
      <Show when={loaded()}>
        <MainView>
          <DataTable.Container>
            <DataTable.Title>Build Status</DataTable.Title>
            <List.Container>
              <BuildStatusRow>
                <BuildStatusLabel>
                  <BuildStatusIcon state={build()?.status} size={24} />
                  {buildStatusStr[build()?.status]}
                </BuildStatusLabel>
                <Show when={!build()?.retriable}>
                  <Button
                    color="borderError"
                    size="small"
                    onClick={rebuild}
                    disabled={disableRefresh()}
                    tooltip={{
                      props: {
                        content: '同じコミットで再ビルド',
                      },
                    }}
                  >
                    Rebuild
                  </Button>
                </Show>
                <Show when={build()?.status === BuildStatus.BUILDING}>
                  <Button color="borderError" size="small" onClick={cancelBuild} disabled={disableRefresh()}>
                    Cancel Build
                  </Button>
                </Show>
              </BuildStatusRow>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>Source Branch (Commit)</List.RowTitle>
                  <List.RowData>
                    {app()?.refName} ({shortSha(app()?.commit)})
                  </List.RowData>
                </List.RowContent>
                <JumpButton href={`/apps/${app()?.id}/settings`} />
                <Button
                  color="ghost"
                  size="medium"
                  onClick={refreshRepo}
                  disabled={disableRefresh()}
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
                  <List.RowData>{buildTypeStr[app()?.config?.buildConfig?.case]}</List.RowData>
                </List.RowContent>
                <JumpButton href={`/apps/${app()?.id}/settings/build`} />
              </List.Row>
            </List.Container>
          </DataTable.Container>
          <DataTable.Container>
            <DataTable.Title>Information</DataTable.Title>
            <List.Container>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>ID</List.RowTitle>
                  <List.RowData>{build()?.id}</List.RowData>
                </List.RowContent>
              </List.Row>
              <Show when={build()?.queuedAt}>
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
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>ビルド開始時刻</List.RowTitle>
                  <Show when={build()?.startedAt?.valid && build()?.startedAt} fallback={'-'}>
                    {(nonNullStartedAt) => {
                      const { diff, localeString } = diffHuman((nonNullStartedAt().timestamp as Timestamp).toDate())
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
                  <List.RowTitle>ビルド終了時刻</List.RowTitle>
                  <Show when={build()?.finishedAt?.valid && build()?.finishedAt} fallback={'-'}>
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
                  <Show when={build()?.finishedAt?.valid && build()?.startedAt?.valid} fallback={'-'}>
                    <List.RowData>
                      {durationHuman(
                        build()?.finishedAt?.timestamp?.toDate().getTime() -
                          build()?.startedAt?.timestamp?.toDate().getTime(),
                      )}
                    </List.RowData>
                  </Show>
                </List.RowContent>
              </List.Row>
            </List.Container>
          </DataTable.Container>
          <Show when={build()?.artifacts?.length > 0}>
            <DataTable.Container>
              <DataTable.Title>Artifacts</DataTable.Title>
              <List.Container>
                <For each={build()?.artifacts}>
                  {(artifact) => (
                    <Bordered>
                      <ArtifactRow artifact={artifact} />
                    </Bordered>
                  )}
                </For>
              </List.Container>
            </DataTable.Container>
          </Show>
          <DataTable.Container>
            <DataTable.Title>Build Log</DataTable.Title>
            <LogContainer>
              <BuildLog buildID={build()?.id} finished={buildFinished()} refetchBuild={refetchBuild} />
            </LogContainer>
          </DataTable.Container>
        </MainView>
      </Show>
    </MainViewContainer>
  )
}
