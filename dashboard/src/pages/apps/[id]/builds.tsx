import { createMemo, onSettled, Show } from 'solid-js'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { BuildList, List } from '/@/components/templates/List'
import { useApplicationData } from '/@/features/application/data'

export default () => {
  const { app, builds, commits, refetch } = useApplicationData()

  onSettled(() => {
    const timer = setInterval(refetch, 10000)
    return () => clearInterval(timer)
  })

  const rows = createMemo(() => builds().map((build) => ({ build })))

  return (
    <MainViewContainer>
      <DataTable.Container>
        <DataTable.Title>Builds</DataTable.Title>
        <Show
          when={rows().length > 0}
          fallback={
            <List.Container>
              <List.PlaceHolder>
                <div class="i-material-symbols:deployed-code-outline text-20/20 shrink-0" />
                No Builds
              </List.PlaceHolder>
            </List.Container>
          }
        >
          <BuildList builds={rows()} commits={commits()} currentBuild={app().currentBuild} />
        </Show>
      </DataTable.Container>
    </MainViewContainer>
  )
}
