import { Title } from '@solidjs/meta'
import { createMemo, For, isPending, Loading, Show } from 'solid-js'
import { GetApplicationsRequest_Scope, GetRepositoriesRequest_Scope } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { MainViewContainer } from '/@/components/layouts/MainView'
import PendingContainer from '/@/components/layouts/PendingContainer'
import { WithNav } from '/@/components/layouts/WithNav'
import { List } from '/@/components/templates/List'
import { Nav } from '/@/components/templates/Nav'
import ReposFilter from '/@/components/templates/repo/ReposFilter'
import { RepositoryRow } from '/@/components/templates/repo/RepositoryRow'
import type { SelectOption } from '/@/components/templates/Select'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import { TabRound } from '/@/components/UI/TabRound'
import { TextField } from '/@/components/UI/TextField'
import { client } from '/@/libs/api'
import { allOrigins, type RepositoryOrigin, repositoryURLToOrigin } from '/@/libs/application'
import { me } from '/@/libs/session'
import { createSessionSignal } from '/@/libs/storage'

const scopeItems = (admin: boolean) => {
  const items: SelectOption<GetRepositoriesRequest_Scope>[] = [
    { value: GetRepositoriesRequest_Scope.MINE, label: 'My Repositories' },
    { value: GetRepositoriesRequest_Scope.PUBLIC, label: 'All Repositories' },
  ]
  if (admin) {
    items.push({ value: GetRepositoriesRequest_Scope.ALL, label: 'All Repositories (admin)' })
  }
  return items
}

export default () => {
  const [scope, setScope] = createSessionSignal('repos-scope', GetRepositoriesRequest_Scope.MINE)
  const [origin, setOrigin] = createSessionSignal<RepositoryOrigin[]>(
    'repos-origin',
    allOrigins.map((o) => o.value),
  )
  const [query, setQuery] = createSessionSignal('repos-query', '')

  const repos = createMemo(() =>
    client
      .getRepositories({ scope: scope() })
      .then((res) => res.repositories.toSorted((a, b) => a.name.localeCompare(b.name))),
  )
  const appCounts = createMemo(() =>
    client.getApplications({ scope: GetApplicationsRequest_Scope.ALL }).then((res) => {
      const counts = new Map<string, number>()
      for (const app of res.applications) counts.set(app.repositoryId, (counts.get(app.repositoryId) ?? 0) + 1)
      return counts
    }),
  )
  const filteredRepos = createMemo(() => {
    const q = query().toLowerCase()
    return repos().filter(
      (repo) => origin().includes(repositoryURLToOrigin(repo.url)) && repo.name.toLowerCase().includes(q),
    )
  })

  return (
    <WithNav.Container>
      <Title>Repositories - NeoShowcase</Title>
      <WithNav.Navs>
        <Nav title="Repositories" />
        <WithNav.Tabs>
          <For each={scopeItems(me().admin)}>
            {(s) => (
              <TabRound state={s.value === scope() ? 'active' : 'default'} onClick={() => setScope(s.value)}>
                <div class="i-mdi:source-repository shrink-0 text-2xl/6" />
                {s.label}
              </TabRound>
            )}
          </For>
          <a href="/repos/new" class="ml-auto">
            <Button
              variants="primary"
              size="medium"
              leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
            >
              Register Repository
            </Button>
          </a>
        </WithNav.Tabs>
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey">
          <div class="flex w-full flex-col gap-8">
            <TextField
              placeholder="Search"
              value={query()}
              onInput={(e) => setQuery(e.currentTarget.value)}
              leftIcon={<div class="i-material-symbols:search shrink-0 text-2xl/6" />}
              rightIcon={<ReposFilter origin={origin()} setOrigin={setOrigin} />}
            />
            <Loading fallback={<Skeleton height={200} />}>
              <PendingContainer isPending={isPending(() => repos())}>
                <List.Container>
                  <Show
                    when={filteredRepos().length > 0}
                    fallback={
                      <List.PlaceHolder>
                        <div class="i-mdi:source-repository text-20/20 shrink-0" />
                        No Repositories
                      </List.PlaceHolder>
                    }
                  >
                    <For each={filteredRepos()}>
                      {(repo) => <RepositoryRow repository={repo} appCount={appCounts().get(repo.id) ?? 0} />}
                    </For>
                  </Show>
                </List.Container>
              </PendingContainer>
            </Loading>
          </div>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
