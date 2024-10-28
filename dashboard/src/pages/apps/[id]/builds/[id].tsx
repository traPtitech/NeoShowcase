import { Title } from '@solidjs/meta'
import { For, Show, createResource } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { List } from '/@/components/templates/List'
import { ArtifactRow } from '/@/components/templates/build/ArtifactRow'
import { BuildLog } from '/@/components/templates/build/BuildLog'
import BuildStatusTable from '/@/components/templates/build/BuildStatusTable'
import { client } from '/@/libs/api'
import { useBuildData } from '/@/routes'

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
        <div class="flex w-full flex-col gap-8">
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
          <Show when={hasPermission()}>
            <DataTable.Container>
              <DataTable.Title>Build Log</DataTable.Title>
              <div class="w-full rounded-lg border border-ui-border px-5 py-4">
                <BuildLog buildID={build()!.id} finished={buildFinished()} refetch={refetch} />
              </div>
            </DataTable.Container>
          </Show>
        </div>
      </Show>
    </MainViewContainer>
  )
}
