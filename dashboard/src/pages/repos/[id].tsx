import ChartIcon from '/@/assets/icons/24/insert_chart.svg'
import SettingsIcon from '/@/assets/icons/24/settings.svg'
import { TabRound } from '/@/components/UI/TabRound'
import { Header } from '/@/components/templates/Header'
import { RepositoryNav } from '/@/components/templates/RepositoryNav'
import { useRepositoryData } from '/@/routes'
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
  const { repo } = useRepositoryData()

  const matchIndexPage = useMatch(() => `/repos/${repo()?.id}/`)
  const matchSettingsPage = useMatch(() => `/repos/${repo()?.id}/settings/*`)

  return (
    <Container>
      <Header />
      <Show when={repo()}>
        <RepositoryNav repository={repo()} />
        <NavTabContainer>
          <NavTabs>
            <TabRound onClick={() => navigate(`/repos/${repo()?.id}`)} state={matchIndexPage() ? 'active' : 'default'}>
              <ChartIcon />
              Project
            </TabRound>
            <TabRound
              onClick={() => navigate(`/repos/${repo()?.id}/settings`)}
              state={matchSettingsPage() ? 'active' : 'default'}
            >
              <SettingsIcon />
              Settings
            </TabRound>
          </NavTabs>
        </NavTabContainer>
        <Outlet />
      </Show>
    </Container>
  )
}
