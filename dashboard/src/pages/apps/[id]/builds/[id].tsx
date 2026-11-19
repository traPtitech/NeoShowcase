import { Title } from '@solidjs/meta'
import { useParams } from '@solidjs/router'
import { createMemo, For, Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { SectionBoundary } from '/@/components/layouts/SectionBoundary'
import { ArtifactRow } from '/@/components/templates/build/ArtifactRow'
import { BuildLog } from '/@/components/templates/build/BuildLog'
import BuildStatusTable from '/@/components/templates/build/BuildStatusTable'
import { RuntimeImageRow } from '/@/components/templates/build/RuntimeImageRow'
import { List } from '/@/components/templates/List'
import { useApplicationData } from '/@/features/application/data'
import { getBuild, getRepositoryCommits, revalidateBuild } from '/@/libs/api'

export default () => {
  const params = useParams<{ buildID: string }>()
  const { app, repo, hasPermission, refetch: refetchApp } = useApplicationData()

  const build = createMemo(() => getBuild(params.buildID))
  const commits = createMemo(() => getRepositoryCommits([build().commit]))

  const refetch = () => {
    revalidateBuild(params.buildID)
    refetchApp()
  }

  return (
    <MainViewContainer>
      <Title>{`${app().name} - Build - NeoShowcase`}</Title>
      <div class="flex w-full flex-col gap-8">
        <DataTable.Container>
          <DataTable.Title>Build Status</DataTable.Title>
          <BuildStatusTable
            app={app()}
            repo={repo()}
            build={build()}
            commit={commits()[build().commit]}
            refetch={refetch}
            hasPermission={hasPermission()}
          />
        </DataTable.Container>
        <Show when={build().artifacts.length > 0}>
          <DataTable.Container>
            <DataTable.Title>Artifacts</DataTable.Title>
            <List.Container>
              <For each={build().artifacts}>{(artifact) => <ArtifactRow artifact={artifact} />}</For>
            </List.Container>
          </DataTable.Container>
        </Show>
        <Show when={build().runtimeImage}>
          {(image) => (
            <DataTable.Container>
              <DataTable.Title>Artifact</DataTable.Title>
              <List.Container>
                <RuntimeImageRow image={image()} />
              </List.Container>
            </DataTable.Container>
          )}
        </Show>
        <Show when={hasPermission()}>
          <SectionBoundary title="Build Log">
            <div class="border-ui-border w-full rounded-lg border px-5 py-4">
              <BuildLog buildID={build().id} finished={build().finishedAt?.valid ?? false} refetch={refetch} />
            </div>
          </SectionBoundary>
        </Show>
      </div>
    </MainViewContainer>
  )
}
