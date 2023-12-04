import { Title } from '@solidjs/meta'
import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { ErrorBoundary, Show, Suspense, startTransition } from 'solid-js'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TabRound } from '/@/components/UI/TabRound'
import ErrorView from '/@/components/layouts/ErrorView'
import { WithNav } from '/@/components/layouts/WithNav'
import { AppNav } from '/@/components/templates/app/AppNav'
import { useApplicationData } from '/@/routes'

export default () => {
  const { app, repo } = useApplicationData()
  const loaded = () => !!(app() && repo())

  const matchIndexPage = useMatch(() => `/apps/${app()?.id}/`)
  const matchBuildsPage = useMatch(() => `/apps/${app()?.id}/builds/*`)
  const matchSettingsPage = useMatch(() => `/apps/${app()?.id}/settings/*`)

  const navigator = useNavigate()
  const navigate = (path: string) => startTransition(() => navigator(path))

  return (
    <WithNav.Container>
      <Show when={loaded()}>
        <Title>{`${app()?.name} - Application - NeoShowcase`}</Title>
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
        <ErrorBoundary fallback={(props) => <ErrorView {...props} />}>
          <Suspense>
            <Outlet />
          </Suspense>
        </ErrorBoundary>
      </WithNav.Body>
    </WithNav.Container>
  )
}
