import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TabRound } from '/@/components/UI/TabRound'
import { WithNav } from '/@/components/layouts/WithNav'
import { RepositoryNav } from '/@/components/templates/RepositoryNav'
import { useRepositoryData } from '/@/routes'

export default () => {
  const navigate = useNavigate()
  const { repo } = useRepositoryData()

  const matchIndexPage = useMatch(() => `/repos/${repo()?.id}/`)
  const matchSettingsPage = useMatch(() => `/repos/${repo()?.id}/settings/*`)

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
            <Outlet />
          </WithNav.Body>
        </WithNav.Container>
      )}
    </Show>
  )
}
