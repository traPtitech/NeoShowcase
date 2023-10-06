import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TabRound } from '/@/components/UI/TabRound'
import { AppNav } from '/@/components/templates/AppNav'
import { Header } from '/@/components/templates/Header'
import { useApplicationData } from '/@/routes'
import { colorVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
})
const NavTabContainer = styled('div', {
  base: {
    width: '100%',
    padding: '0 32px 16px',
    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
  },
})
const NavTabs = styled('div', {
  base: {
    width: '100%',
    maxWidth: '1000px',
    margin: '0 auto',
    display: 'flex',
    gap: '8px',
  },
})

export default () => {
  const navigate = useNavigate()
  const { app, repo } = useApplicationData()
  const loaded = () => !!(app() && repo())

  const matchIndexPage = useMatch(() => `/apps/${app()?.id}/`)
  const matchBuildsPage = useMatch(() => `/apps/${app()?.id}/builds/*`)
  const matchSettingsPage = useMatch(() => `/apps/${app()?.id}/settings/*`)

  return (
    <Container>
      <Header />
      <Show when={loaded()}>
        <AppNav app={app()} repository={repo()} />
        <NavTabContainer>
          <NavTabs>
            <TabRound onClick={() => navigate(`/apps/${app()?.id}`)} state={matchIndexPage() ? 'active' : 'default'}>
              <MaterialSymbols>insert_chart</MaterialSymbols>
              App
            </TabRound>
            <TabRound
              onClick={() => navigate(`/apps/${app()?.id}/builds`)}
              state={matchBuildsPage() ? 'active' : 'default'}
            >
              <MaterialSymbols>deployed_code</MaterialSymbols>
              Builds
            </TabRound>
            <TabRound
              onClick={() => navigate(`/apps/${app()?.id}/settings`)}
              state={matchSettingsPage() ? 'active' : 'default'}
            >
              <MaterialSymbols>settings</MaterialSymbols>
              Settings
            </TabRound>
          </NavTabs>
        </NavTabContainer>
        <Outlet />
      </Show>
    </Container>
  )
}
