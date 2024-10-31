import { Title } from '@solidjs/meta'
import { type RouteSectionProps, useMatch, useNavigate } from '@solidjs/router'
import { ErrorBoundary, Show, Suspense, startTransition } from 'solid-js'
import { TabRound } from '/@/components/UI/TabRound'
import ErrorView from '/@/components/layouts/ErrorView'
import { WithNav } from '/@/components/layouts/WithNav'
import { RepositoryNav } from '/@/components/templates/repo/RepositoryNav'
import { useRepositoryData } from '/@/routes'

export default (props: RouteSectionProps) => {
  const { repo } = useRepositoryData()

  const matchIndexPage = useMatch(() => `/repos/${repo()?.id}/`)
  const matchSettingsPage = useMatch(() => `/repos/${repo()?.id}/settings/*`)

  const navigator = useNavigate()
  const navigate = (path: string) => startTransition(() => navigator(path))

  return (
    <Show when={repo()}>
      <WithNav.Container>
        <Title>{`${repo()!.name} - Repository - NeoShowcase`}</Title>
        <WithNav.Navs>
          <RepositoryNav repository={repo()!} />
          <WithNav.Tabs>
            <TabRound onClick={() => navigate(`/repos/${repo()!.id}`)} state={matchIndexPage() ? 'active' : 'default'}>
              <span class="i-material-symbols:insert-chart-outline text-2xl/6" />
              Info
            </TabRound>
            <TabRound
              onClick={() => navigate(`/repos/${repo()!.id}/settings`)}
              state={matchSettingsPage() ? 'active' : 'default'}
            >
              <span class="i-material-symbols:settings-outline text-2xl/6" />
              Settings
            </TabRound>
          </WithNav.Tabs>
        </WithNav.Navs>
        <WithNav.Body>
          <ErrorBoundary fallback={(props) => <ErrorView {...props} />}>
            <Suspense>{props.children}</Suspense>
          </ErrorBoundary>
        </WithNav.Body>
      </WithNav.Container>
    </Show>
  )
}
