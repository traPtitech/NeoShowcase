import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TabRound } from '/@/components/UI/TabRound'
import { WithNav } from '/@/components/layouts/WithNav'
import { useApplicationData } from '/@/routes'
import { AppNav } from '../../components/templates/app/AppNav'

export default () => {
  const navigate = useNavigate()
  const { app, repo } = useApplicationData()
  const loaded = () => app.state === 'ready' && repo.state === 'ready'

  const matchIndexPage = useMatch(() => `/apps/${app()?.id}/`)
  const matchBuildsPage = useMatch(() => `/apps/${app()?.id}/builds/*`)
  const matchSettingsPage = useMatch(() => `/apps/${app()?.id}/settings/*`)

  return (
    <WithNav.Container>
      <Show when={loaded()}>
        <WithNav.Navs>
          <AppNav app={app()!} repository={repo()!} />
          <WithNav.Tabs>
            <TabRound onClick={() => navigate(`/apps/${app()?.id}`)} state={matchIndexPage() ? 'active' : 'default'}>
              <MaterialSymbols>insert_chart</MaterialSymbols>
              Info
            </TabRound>
            <TabRound
              onClick={() => navigate(`/apps/${app()?.id}/builds`)}
              state={matchBuildsPage() ? 'active' : 'default'}
            >
              <MaterialSymbols>deployed_code</MaterialSymbols>
              Build History
            </TabRound>
            <TabRound
              onClick={() => navigate(`/apps/${app()?.id}/settings`)}
              state={matchSettingsPage() ? 'active' : 'default'}
            >
              <MaterialSymbols>settings</MaterialSymbols>
              Settings
            </TabRound>
          </WithNav.Tabs>
        </WithNav.Navs>
      </Show>
      <WithNav.Body>
        <Outlet />
      </WithNav.Body>
    </WithNav.Container>
  )
}
