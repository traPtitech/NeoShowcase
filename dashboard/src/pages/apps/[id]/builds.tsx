import { styled } from '@macaron-css/solid'
import { createMemo, createResource } from 'solid-js'
import { Show } from 'solid-js'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { BuildList, List } from '/@/components/templates/List'
import { client } from '/@/libs/api'
import { useApplicationData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'

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

export default () => {
  const { app } = useApplicationData()
  const [builds] = createResource(
    () => app()?.id,
    (id) => client.getBuilds({ id }),
  )
  const loaded = () => !!(app() && builds())

  const sortedBuilds = createMemo(
    () =>
      builds() &&
      [...builds().builds].sort((b1, b2) => {
        return b2.queuedAt.toDate().getTime() - b1.queuedAt.toDate().getTime()
      }),
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
              <PlaceHolder>
                <MaterialSymbols displaySize={80}>deployed_code</MaterialSymbols>
                No Builds
              </PlaceHolder>
            </List.Container>
          </Show>
        </DataTable.Container>
      </Show>
    </MainViewContainer>
  )
}
