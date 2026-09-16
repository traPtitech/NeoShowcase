import { Title } from '@solidjs/meta'
import { type RouteSectionProps, useMatch, useNavigate, useParams } from '@solidjs/router'
import { Show, Suspense, startTransition } from 'solid-js'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import { WithNav } from '/@/components/layouts/WithNav'
import { RepositoryNav } from '/@/components/templates/repo/RepositoryNav'
import { TabRound } from '/@/components/UI/TabRound'
import { useRepositoryData } from '/@/routes'

export default (props: RouteSectionProps) => {
  const { repo, refetchRepo } = useRepositoryData()

  // Route params, not the resource: these memos are owned by the component, outside the boundary it
  // renders, so reading a failed resource here would escape to the boundary above instead of this page's.
  const params = useParams()
  const matchIndexPage = useMatch(() => `/repos/${params.id}/`)
  const matchSettingsPage = useMatch(() => `/repos/${params.id}/settings/*`)

  const navigator = useNavigate()
  const navigate = (path: string) => startTransition(() => navigator(path))

  return (
    <PageBoundary onRetry={refetchRepo}>
      <Show when={repo()}>
        <WithNav.Container>
          <Title>{`${repo()!.name} - Repository - NeoShowcase`}</Title>
          <WithNav.Navs>
            <RepositoryNav repository={repo()!} />
            <WithNav.Tabs>
              <TabRound
                onClick={() => navigate(`/repos/${repo()!.id}`)}
                state={matchIndexPage() ? 'active' : 'default'}
              >
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
            <Suspense>{props.children}</Suspense>
          </WithNav.Body>
        </WithNav.Container>
      </Show>
    </PageBoundary>
  )
}
