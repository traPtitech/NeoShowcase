import AppsPlaceholder from '/@/assets/icons/apps_placeholder.svg'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { URLText } from '/@/components/UI/URLText'
import { AppsList, ListContainer } from '/@/components/templates/List'
import { useRepositoryData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    padding: '40px 32px',
    overflowY: 'auto',
  },
})
const MainView = styled('div', {
  base: {
    width: '100%',
    maxWidth: '1000px',
    margin: '0 auto',
    display: 'flex',
    flexDirection: 'column',
    gap: '32px',
  },
})
const DataContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
  },
})
const DataTitle = styled('h2', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    color: colorVars.semantic.text.black,
    ...textVars.h2.medium,
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
const TableRow = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',

    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
    selectors: {
      '&:last-child': {
        borderBottom: 'none',
      },
    },
  },
})
const TableTitle = styled('h3', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.text.medium,
  },
})
const TableData = styled('div', {
  base: {
    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
  },
})

export default () => {
  const { repo, apps } = useRepositoryData()
  const loaded = () => !!(repo() && apps())
  const navigator = useNavigate()
  const showPlaceHolder = () => apps()?.length === 0

  return (
    <Container>
      <MainView>
        <Show when={loaded()}>
          <DataContainer>
            <DataTitle>
              Apps
              <Show when={!showPlaceHolder()}>
                <Button
                  color="primary"
                  size="medium"
                  leftIcon={<MaterialSymbols>add</MaterialSymbols>}
                  onClick={() => {
                    navigator(`/apps/new?repositoryID=${repo()?.id}`)
                  }}
                >
                  Add New App
                </Button>
              </Show>
            </DataTitle>
            <Show when={showPlaceHolder()} fallback={<AppsList apps={apps()} />}>
              <ListContainer>
                <PlaceHolder>
                  <AppsPlaceholder />
                  No Apps
                  <Button
                    color="primary"
                    size="medium"
                    leftIcon={<MaterialSymbols>add</MaterialSymbols>}
                    onClick={() => {
                      navigator(`/apps/new?repositoryID=${repo()?.id}`)
                    }}
                  >
                    Add New App
                  </Button>
                </PlaceHolder>
              </ListContainer>
            </Show>
          </DataContainer>
          <DataContainer>
            <DataTitle>Information</DataTitle>
            <ListContainer>
              <TableRow>
                <TableTitle>ID</TableTitle>
                <TableData>{repo()?.id}</TableData>
              </TableRow>
              <TableRow>
                <TableTitle>Name</TableTitle>
                <TableData>{repo()?.name}</TableData>
              </TableRow>
              <TableRow>
                <TableTitle>URL</TableTitle>
                <TableData>
                  <URLText text={repo()?.url} href={repo()?.htmlUrl} />
                </TableData>
              </TableRow>
            </ListContainer>
          </DataContainer>
        </Show>
      </MainView>
    </Container>
  )
}
