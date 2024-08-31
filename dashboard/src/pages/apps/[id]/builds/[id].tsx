import { styled } from '@macaron-css/solid'
import { Title } from '@solidjs/meta'
import { For, Show, createResource } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { List } from '/@/components/templates/List'
import { ArtifactRow } from '/@/components/templates/build/ArtifactRow'
import { BuildLog } from '/@/components/templates/build/BuildLog'
import BuildStatusTable from '/@/components/templates/build/BuildStatusTable'
import { RuntimeImageRow } from '/@/components/templates/build/RuntimeImageRow'
import { client } from '/@/libs/api'
import { useBuildData } from '/@/routes'
import { colorVars } from '/@/theme'

const MainView = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '32px',
  },
})
const LogContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',

    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
  },
})

export default () => {
  const { app, build, commit, refetch, hasPermission } = useBuildData()
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const loaded = () => !!(app() && repo() && build())

  const buildFinished = () => build()?.finishedAt?.valid ?? false

  return (
    <MainViewContainer>
      <Show when={loaded()}>
        <Title>{`${app()!.name} - Build - NeoShowcase`}</Title>
        <MainView>
          <DataTable.Container>
            <DataTable.Title>Build Status</DataTable.Title>
            <BuildStatusTable
              app={app()!}
              repo={repo()!}
              build={build()!}
              commit={commit()}
              refetch={refetch}
              hasPermission={hasPermission()}
            />
          </DataTable.Container>
          <Show when={build()!.artifacts.length > 0}>
            <DataTable.Container>
              <DataTable.Title>Artifacts</DataTable.Title>
              <List.Container>
                <For each={build()!.artifacts}>{(artifact) => <ArtifactRow artifact={artifact} />}</For>
              </List.Container>
            </DataTable.Container>
          </Show>
          <Show when={build()!.runtimeImage != null}>
            <DataTable.Container>
              <DataTable.Title>Runtime Image</DataTable.Title>
              <List.Container>
                <RuntimeImageRow image={build()!.runtimeImage!} />
              </List.Container>
            </DataTable.Container>
          </Show>
          <Show when={hasPermission()}>
            <DataTable.Container>
              <DataTable.Title>Build Log</DataTable.Title>
              <LogContainer>
                <BuildLog buildID={build()!.id} finished={buildFinished()} refetch={refetch} />
              </LogContainer>
            </DataTable.Container>
          </Show>
        </MainView>
      </Show>
    </MainViewContainer>
  )
}
