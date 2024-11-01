import { Title } from '@solidjs/meta'
import { type RouteSectionProps, useMatch, useNavigate } from '@solidjs/router'
import { ErrorBoundary, Show, startTransition } from 'solid-js'
import { TabRound } from '/@/components/UI/TabRound'
import ErrorView from '/@/components/layouts/ErrorView'
import { WithNav } from '/@/components/layouts/WithNav'
import { AppNav } from '/@/components/templates/app/AppNav'
import { useApplicationData } from '/@/routes'

export default (props: RouteSectionProps) => {
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
              <div class="i-material-symbols:insert-chart-outline text-2xl/6" />
              Info
            </TabRound>
            <TabRound
              onClick={() => navigate(`/apps/${app()?.id}/builds`)}
              state={matchBuildsPage() ? 'active' : 'default'}
            >
              <div class="i-material-symbols:history text-2xl/6" />
              Build History
            </TabRound>
            <TabRound
              onClick={() => navigate(`/apps/${app()?.id}/settings`)}
              state={matchSettingsPage() ? 'active' : 'default'}
            >
              <div class="i-material-symbols:settings-outline text-2xl/6" />
              Settings
            </TabRound>
          </WithNav.Tabs>
        </WithNav.Navs>
      </Show>
      <WithNav.Body>
        <ErrorBoundary fallback={(props) => <ErrorView {...props} />}>{props.children}</ErrorBoundary>
      </WithNav.Body>
    </WithNav.Container>
  )
}
