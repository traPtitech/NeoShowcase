import { Title } from '@solidjs/meta'
import { type RouteSectionProps, useMatch, useNavigate, useParams } from '@solidjs/router'
import { Show, startTransition } from 'solid-js'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import { WithNav } from '/@/components/layouts/WithNav'
import { AppNav } from '/@/components/templates/app/AppNav'
import { TabRound } from '/@/components/UI/TabRound'
import { useApplicationData } from '/@/routes'

export default (props: RouteSectionProps) => {
  const { app, repo } = useApplicationData()
  const loaded = () => !!(app() && repo())

  // Route params, not app(): see the note in pages/repos/[id].tsx.
  const params = useParams()
  const matchIndexPage = useMatch(() => `/apps/${params.id}/`)
  const matchBuildsPage = useMatch(() => `/apps/${params.id}/builds/*`)
  const matchSettingsPage = useMatch(() => `/apps/${params.id}/settings/*`)

  const navigator = useNavigate()
  const navigate = (path: string) => startTransition(() => navigator(path))

  return (
    <PageBoundary>
      <WithNav.Container>
        <Show when={loaded()}>
          <Title>{`${app()?.name} - Application - NeoShowcase`}</Title>
          <WithNav.Navs>
            <AppNav app={app()!} repository={repo()!} />
            <WithNav.Tabs>
              <TabRound onClick={() => navigate(`/apps/${app()?.id}`)} state={matchIndexPage() ? 'active' : 'default'}>
                <div class="i-material-symbols:insert-chart-outline shrink-0 text-2xl/6" />
                Info
              </TabRound>
              <TabRound
                onClick={() => navigate(`/apps/${app()?.id}/builds`)}
                state={matchBuildsPage() ? 'active' : 'default'}
              >
                <div class="i-material-symbols:history shrink-0 text-2xl/6" />
                Build History
              </TabRound>
              <TabRound
                onClick={() => navigate(`/apps/${app()?.id}/settings`)}
                state={matchSettingsPage() ? 'active' : 'default'}
              >
                <div class="i-material-symbols:settings-outline shrink-0 text-2xl/6" />
                Settings
              </TabRound>
            </WithNav.Tabs>
          </WithNav.Navs>
        </Show>
        <WithNav.Body>{props.children}</WithNav.Body>
      </WithNav.Container>
    </PageBoundary>
  )
}
