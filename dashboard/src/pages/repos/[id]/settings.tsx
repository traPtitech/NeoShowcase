import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { useRepositoryData } from '/@/routes'
import { styled } from '@macaron-css/solid'
import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    padding: '40px 32px 72px 32px',
    overflowY: 'auto',
  },
})
const MainView = styled('div', {
  base: {
    width: '100%',
    maxWidth: '1000px',
    margin: '0 auto',
    display: 'grid',
    gridTemplateColumns: '235px 1fr',
    gap: '48px',
  },
})
const SideMenuContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
  },
})
const SideMenu = styled('div', {
  base: {
    position: 'sticky',
    width: '100%',
    top: '0',
    display: 'flex',
    flexDirection: 'column',
  },
})

export default () => {
  const { repo } = useRepositoryData()
  const loaded = () => !!repo()
  const navigator = useNavigate()
  const matchGeneralPage = useMatch(() => `/repos/${repo()?.id}/settings/`)
  const matchAuthPage = useMatch(() => `/repos/${repo()?.id}/settings/authorization`)
  const matchOwnerPage = useMatch(() => `/repos/${repo()?.id}/settings/owner`)

  return (
    <Container>
      <Show when={loaded()}>
        <MainView>
          <SideMenuContainer>
            <SideMenu>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchGeneralPage()}
                onclick={() => {
                  navigator(`/repos/${repo()?.id}/settings/`)
                }}
                leftIcon={<MaterialSymbols>browse_activity</MaterialSymbols>}
              >
                General
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchAuthPage()}
                onclick={() => {
                  navigator(`/repos/${repo()?.id}/settings/authorization`)
                }}
                leftIcon={<MaterialSymbols>conversion_path</MaterialSymbols>}
              >
                Authorization
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchOwnerPage()}
                onclick={() => {
                  navigator(`/repos/${repo()?.id}/settings/owner`)
                }}
                leftIcon={<MaterialSymbols>person</MaterialSymbols>}
              >
                Owner
              </Button>
            </SideMenu>
          </SideMenuContainer>
          <Outlet />
        </MainView>
      </Show>
    </Container>
  )
}
