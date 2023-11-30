import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { ErrorBoundary, Show, Suspense, startTransition } from 'solid-js'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TabRound } from '/@/components/UI/TabRound'
import ErrorView from '/@/components/layouts/ErrorView'
import { WithNav } from '/@/components/layouts/WithNav'
import { useRepositoryData } from '/@/routes'
import { RepositoryNav } from '../../components/templates/repo/RepositoryNav'

export default () => {
  const { repo } = useRepositoryData()

  const matchIndexPage = useMatch(() => `/repos/${repo()?.id}/`)
  const matchSettingsPage = useMatch(() => `/repos/${repo()?.id}/settings/*`)

  const navigator = useNavigate()
  const navigate = (path: string) => startTransition(() => navigator(path))

  return (
    <Show when={repo()}>
      {(nonNullRepo) => (
        <WithNav.Container>
          <WithNav.Navs>
            <RepositoryNav repository={nonNullRepo()} />
            <WithNav.Tabs>
              <TabRound
                onClick={() => navigate(`/repos/${nonNullRepo().id}`)}
                state={matchIndexPage() ? 'active' : 'default'}
              >
                <MaterialSymbols>insert_chart</MaterialSymbols>
                Info
              </TabRound>
              <TabRound
                onClick={() => navigate(`/repos/${nonNullRepo().id}/settings`)}
                state={matchSettingsPage() ? 'active' : 'default'}
              >
                <MaterialSymbols>settings</MaterialSymbols>
                Settings
              </TabRound>
            </WithNav.Tabs>
          </WithNav.Navs>
          <WithNav.Body>
            <ErrorBoundary fallback={(props) => <ErrorView {...props} />}>
              <Suspense>
                <Outlet />
              </Suspense>
            </ErrorBoundary>
          </WithNav.Body>
        </WithNav.Container>
      )}
    </Show>
  )
}
