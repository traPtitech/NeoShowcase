import { Title } from '@solidjs/meta'
import { type RouteSectionProps, useMatch, useNavigate } from '@solidjs/router'
import { Show, Suspense, startTransition } from 'solid-js'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import { WithNav } from '/@/components/layouts/WithNav'
import { RepositoryNav } from '/@/components/templates/repo/RepositoryNav'
import { TabRound } from '/@/components/UI/TabRound'
import { useRepositoryData } from '/@/routes'

export default (props: RouteSectionProps) => {
  const { repo, refetchRepo } = useRepositoryData()

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
              <div class="i-material-symbols:insert-chart-outline shrink-0 text-2xl/6" />
              Info
            </TabRound>
            <TabRound
              onClick={() => navigate(`/repos/${repo()!.id}/settings`)}
              state={matchSettingsPage() ? 'active' : 'default'}
            >
              <div class="i-material-symbols:settings-outline shrink-0 text-2xl/6" />
              Settings
            </TabRound>
          </WithNav.Tabs>
        </WithNav.Navs>
        <WithNav.Body>
          <PageBoundary onRetry={refetchRepo}>
            <Suspense>{props.children}</Suspense>
          </PageBoundary>
        </WithNav.Body>
      </WithNav.Container>
    </Show>
  )
}
