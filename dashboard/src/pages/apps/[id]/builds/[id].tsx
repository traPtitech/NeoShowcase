import { Timestamp } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { For, Show, VoidComponent, createResource } from 'solid-js'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { ToolTip } from '/@/components/UI/ToolTip'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { ArtifactRow } from '/@/components/templates/ArtifactRow'
import { BuildLog } from '/@/components/templates/BuildLog'
import BuildStatusTable from '/@/components/templates/BuildStatusTable'
import { Bordered, List } from '/@/components/templates/List'
import { client } from '/@/libs/api'
import { diffHuman, durationHuman } from '/@/libs/format'
import { useBuildData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'

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
  const { app, build, refetchBuild, hasPermission } = useBuildData()
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const loaded = () => !!(app() && repo() && build())

  const buildFinished = () => build()?.finishedAt?.valid

  return (
    <MainViewContainer>
      <Show when={loaded()}>
        <MainView>
          <DataTable.Container>
            <DataTable.Title>Build Status</DataTable.Title>
            <BuildStatusTable
              app={app()}
              repo={repo()}
              build={build()}
              refetchBuild={refetchBuild}
              hasPermission={hasPermission()}
            />
          </DataTable.Container>
          <DataTable.Container>
            <DataTable.Title>Information</DataTable.Title>
            <List.Container>
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
          <Show when={hasPermission()}>
            <DataTable.Container>
              <DataTable.Title>Build Log</DataTable.Title>
              <LogContainer>
                <BuildLog buildID={build()?.id} finished={buildFinished()} refetchBuild={refetchBuild} />
              </LogContainer>
            </DataTable.Container>
          </Show>
        </MainView>
      </Show>
    </MainViewContainer>
  )
}
