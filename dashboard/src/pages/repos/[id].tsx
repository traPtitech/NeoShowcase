import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TabRound } from '/@/components/UI/TabRound'
import { WithNav } from '/@/components/layouts/WithNav'
import { RepositoryNav } from '/@/components/templates/RepositoryNav'
import { useRepositoryData } from '/@/routes'
import { Outlet, useMatch, useNavigate } from '@solidjs/router'
import { Show } from 'solid-js'

export default () => {
  const navigate = useNavigate()
  const { repo } = useRepositoryData()

  const matchIndexPage = useMatch(() => `/repos/${repo()?.id}/`)
  const matchSettingsPage = useMatch(() => `/repos/${repo()?.id}/settings/*`)

  return (
    <Show when={repo()}>
      <WithNav.Container>
        <WithNav.Navs>
          <RepositoryNav repository={repo()} />
          <WithNav.Tabs>
            <TabRound onClick={() => navigate(`/repos/${repo()?.id}`)} state={matchIndexPage() ? 'active' : 'default'}>
              <MaterialSymbols>insert_chart</MaterialSymbols>
              Project
            </TabRound>
            <TabRound
              onClick={() => navigate(`/repos/${repo()?.id}/settings`)}
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
