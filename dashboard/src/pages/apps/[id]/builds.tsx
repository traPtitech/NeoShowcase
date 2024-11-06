import { createMemo, onCleanup, useTransition } from 'solid-js'
import { Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import SuspenseContainer from '/@/components/layouts/SuspenseContainer'
import { BuildList, List } from '/@/components/templates/List'
import { useApplicationData } from '/@/routes'

export default () => {
  const { app, builds, commits, refetch } = useApplicationData()
  const loaded = () => !!(app() && builds())

  const refetchTimer = setInterval(refetch, 10000)
  onCleanup(() => clearInterval(refetchTimer))

  const sortedBuilds = createMemo(
    () =>
      builds()
        ?.sort((b1, b2) => {
          return (b2.queuedAt?.toDate().getTime() ?? 0) - (b1.queuedAt?.toDate().getTime() ?? 0)
        })
        ?.map((b) => ({ build: b })) ?? [],
  )
  const showPlaceHolder = () => builds()?.length === 0

  const [isPending] = useTransition()

  return (
    <SuspenseContainer isPending={isPending()}>
      <MainViewContainer>
        <Show when={loaded()}>
          <DataTable.Container>
            <DataTable.Title>Builds</DataTable.Title>
            <Show
              when={showPlaceHolder()}
              fallback={<BuildList builds={sortedBuilds()} commits={commits()} currentBuild={app()?.currentBuild} />}
            >
              <List.Container>
                <List.PlaceHolder>
                  <div class="i-material-symbols:deployed-code-outline shrink-0 text-20/20" />
                  No Builds
                </List.PlaceHolder>
              </List.Container>
            </Show>
          </DataTable.Container>
        </Show>
      </MainViewContainer>
    </SuspenseContainer>
  )
}
