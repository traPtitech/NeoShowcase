import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { MainViewContainer } from '/@/components/layouts/MainView'
import { useApplicationData } from '/@/routes'
import { styled } from '@macaron-css/solid'
import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'

const MainView = styled('div', {
  base: {
    width: '100%',
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
  const { app } = useApplicationData()
  const loaded = () => !!app()
  const navigator = useNavigate()
  const matchGeneralPage = useMatch(() => `/apps/${app()?.id}/settings/`)
  const matchBuildPage = useMatch(() => `/apps/${app()?.id}/settings/build`)
  const matchDomainsPage = useMatch(() => `/apps/${app()?.id}/settings/domains`)
  const matchPortPage = useMatch(() => `/apps/${app()?.id}/settings/portForwarding`)
  const matchEnvVarsPage = useMatch(() => `/apps/${app()?.id}/settings/envVars`)
  const matchOwnerPage = useMatch(() => `/apps/${app()?.id}/settings/owner`)

  return (
    <MainViewContainer>
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
                  navigator(`/apps/${app()?.id}/settings/`)
                }}
                leftIcon={<MaterialSymbols>browse_activity</MaterialSymbols>}
              >
                General
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchBuildPage()}
                onclick={() => {
                  navigator(`/apps/${app()?.id}/settings/build`)
                }}
                leftIcon={<MaterialSymbols>deployed_code</MaterialSymbols>}
              >
                Build
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchDomainsPage()}
                onclick={() => {
                  navigator(`/apps/${app()?.id}/settings/domains`)
                }}
                leftIcon={<MaterialSymbols>language</MaterialSymbols>}
              >
                Domain
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchPortPage()}
                onclick={() => {
                  navigator(`/apps/${app()?.id}/settings/portForwarding`)
                }}
                leftIcon={<MaterialSymbols>lan</MaterialSymbols>}
              >
                Port Forwarding
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchEnvVarsPage()}
                onclick={() => {
                  navigator(`/apps/${app()?.id}/settings/envVars`)
                }}
                leftIcon={<MaterialSymbols>password</MaterialSymbols>}
              >
                Environment Variables
              </Button>
              <Button
                color="text"
                size="medium"
                full
                active={!!matchOwnerPage()}
                onclick={() => {
                  navigator(`/apps/${app()?.id}/settings/owner`)
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
    </MainViewContainer>
  )
}
