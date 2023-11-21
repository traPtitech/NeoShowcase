import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { URLText } from '/@/components/UI/URLText'
import { DataTable } from '/@/components/layouts/DataTable'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { AppsList, List } from '/@/components/templates/List'
import { useRepositoryData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import { Show, createMemo } from 'solid-js'

const MainView = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '32px',
  },
})
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
  const { repo, apps, hasPermission } = useRepositoryData()
  const loaded = () => !!(repo() && apps())
  const navigator = useNavigate()
  const showPlaceHolder = createMemo(() => apps()?.length === 0)

  const AddNewAppButton = () => (
    <Button
      variants="primary"
      size="medium"
      leftIcon={<MaterialSymbols>add</MaterialSymbols>}
      onClick={() => {
        navigator(`/apps/new?repositoryID=${repo()?.id}`)
      }}
      tooltip={{
        props: {
          content: hasPermission() ? (
            'このリポジトリからアプリケーションを作成します'
          ) : (
            <>
              <div>アプリケーションを作成するには</div>
              <div>リポジトリのオーナーになる必要があります</div>
            </>
          ),
        },
        style: 'center',
      }}
      disabled={!hasPermission()}
    >
      Add New App
    </Button>
  )

  return (
    <MainViewContainer>
      <MainView>
        <Show when={loaded()}>
          <DataTable.Container>
            <DataTable.Title>
              Apps
              <Show when={!showPlaceHolder()}>
                <AddNewAppButton />
              </Show>
            </DataTable.Title>
            <Show when={showPlaceHolder()} fallback={<AppsList apps={apps()} />}>
              <List.Container>
                <PlaceHolder>
                  <MaterialSymbols displaySize={80}>deployed_code</MaterialSymbols>
                  No Apps
                  <AddNewAppButton />
                </PlaceHolder>
              </List.Container>
            </Show>
          </DataTable.Container>
          <DataTable.Container>
            <DataTable.Title>Information</DataTable.Title>
            <List.Container>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>ID</List.RowTitle>
                  <List.RowData>{repo()?.id}</List.RowData>
                </List.RowContent>
              </List.Row>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>Name</List.RowTitle>
                  <List.RowData>{repo()?.name}</List.RowData>
                </List.RowContent>
              </List.Row>
              <List.Row>
                <List.RowContent>
                  <List.RowTitle>URL</List.RowTitle>
                  <List.RowData>
                    <URLText text={repo()?.url} href={repo()?.htmlUrl} />
                  </List.RowData>
                </List.RowContent>
              </List.Row>
            </List.Container>
          </DataTable.Container>
        </Show>
      </MainView>
    </MainViewContainer>
  )
}
