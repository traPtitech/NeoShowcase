import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TabRound } from '/@/components/UI/TabRound'
import { WithNav } from '/@/components/layouts/WithNav'
import { AppNav } from '/@/components/templates/AppNav'
import { useApplicationData } from '/@/routes'
import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'

export default () => {
  const navigate = useNavigate()
  const { app, repo } = useApplicationData()
  const loaded = () => !!(app() && repo())

  const matchIndexPage = useMatch(() => `/apps/${app()?.id}/`)
  const matchBuildsPage = useMatch(() => `/apps/${app()?.id}/builds/*`)
  const matchSettingsPage = useMatch(() => `/apps/${app()?.id}/settings/*`)

  return (
    <Show when={loaded()}>
      <WithNav.Container>
        <WithNav.Navs>
          <AppNav app={app()} repository={repo()} />
          <WithNav.Tabs>
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
          </WithNav.Tabs>
        </WithNav.Navs>
        <WithNav.Body>
          <Outlet />
        </WithNav.Body>
      </WithNav.Container>
    </Show>
  )
}
