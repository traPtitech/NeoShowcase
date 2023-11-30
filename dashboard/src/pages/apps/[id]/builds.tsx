import { createMemo, createResource } from 'solid-js'
import { Show } from 'solid-js'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { BuildList, List } from '/@/components/templates/List'
import { client } from '/@/libs/api'
import { useApplicationData } from '/@/routes'

export default () => {
  const { app } = useApplicationData()
  const [builds] = createResource(
    () => app()?.id,
    (id) => client.getBuilds({ id }),
  )
  const loaded = () => !!(app.state === 'ready' && builds.state === 'ready')

  const sortedBuilds = createMemo(() =>
    builds.state === 'ready'
      ? [...builds().builds].sort((b1, b2) => {
          return (b2.queuedAt?.toDate().getTime() ?? 0) - (b1.queuedAt?.toDate().getTime() ?? 0)
        })
      : [],
  )
  const showPlaceHolder = () => builds()?.builds.length === 0

  return (
    <MainViewContainer>
      <Show when={loaded()}>
        <DataTable.Container>
          <DataTable.Title>Builds</DataTable.Title>
          <Show
            when={showPlaceHolder()}
            fallback={<BuildList builds={sortedBuilds()} showAppID={false} deployedBuild={app()?.currentBuild} />}
          >
            <List.Container>
              <List.PlaceHolder>
                <MaterialSymbols displaySize={80}>deployed_code</MaterialSymbols>
                No Builds
              </List.PlaceHolder>
            </List.Container>
          </Show>
        </DataTable.Container>
      </Show>
    </MainViewContainer>
  )
}
