import { Title } from '@solidjs/meta'
import { createMemo, For, isPending, Loading, Show } from 'solid-js'
import {
  type Application,
  GetApplicationGroupsRequest_Scope,
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import PendingContainer from '/@/components/layouts/PendingContainer'
import { WithNav } from '/@/components/layouts/WithNav'
import { styled } from '/@/components/styled-components'
import { AppsNav } from '/@/components/templates/AppsNav'
import AppsFilter from '/@/components/templates/app/AppsFilter'
import { AppsList } from '/@/components/templates/List'
import type { SelectOption } from '/@/components/templates/Select'
import { Button } from '/@/components/UI/Button'
import { TabRound } from '/@/components/UI/TabRound'
import { TextField } from '/@/components/UI/TextField'
import { client, getRepositoryCommits } from '/@/libs/api'
import {
  allOrigins,
  allStatuses,
  applicationState,
  groupApplications,
  type RepositoryOrigin,
  repositoryURLToOrigin,
  sortItems,
} from '/@/libs/application'
import { me } from '/@/libs/session'
import { createSessionSignal } from '/@/libs/storage'

const FilterContainer = styled('div', 'z-1 flex w-full items-center justify-center pt-10 pb-8')

const scopeItems = (admin: boolean) => {
  const items: SelectOption<GetRepositoriesRequest_Scope>[] = [
    { value: GetRepositoriesRequest_Scope.MINE, label: 'My Apps' },
    { value: GetRepositoriesRequest_Scope.PUBLIC, label: 'All Apps' },
  ]
  if (admin) {
    items.push({ value: GetRepositoriesRequest_Scope.ALL, label: 'All Apps (admin)' })
  }
  return items
}

const Skeletons = () => <AppsList apps={[undefined, undefined, undefined, undefined]} />

export default () => {
  const [scope, setScope] = createSessionSignal('apps-scope', GetRepositoriesRequest_Scope.MINE)
  const [statuses, setStatuses] = createSessionSignal(
    'apps-statuses-v1',
    allStatuses.map((s) => s.value),
  )
  const [origin, setOrigin] = createSessionSignal<RepositoryOrigin[]>(
    'apps-repository-origin',
    allOrigins.map((o) => o.value),
  )
  const [query, setQuery] = createSessionSignal('apps-query', '')
  const [sort, setSort] = createSessionSignal<keyof typeof sortItems>('apps-sort', sortItems.desc.value)

  // Repositories the user can see: "All Apps" lists only their apps, and the apps show their names.
  const repos = createMemo(() =>
    client
      .getRepositories({
        scope:
          scope() === GetRepositoriesRequest_Scope.ALL
            ? GetRepositoriesRequest_Scope.ALL
            : GetRepositoriesRequest_Scope.PUBLIC,
      })
      .then((res) => Object.fromEntries(res.repositories.map((r) => [r.id, r]))),
  )
  const apps = createMemo(() =>
    client
      .getApplications({
        scope:
          scope() === GetRepositoriesRequest_Scope.MINE
            ? GetApplicationsRequest_Scope.MINE
            : GetApplicationsRequest_Scope.ALL,
      })
      .then((res) => res.applications),
  )
  const groups = createMemo(() =>
    client
      .getApplicationGroups({ scope: GetApplicationGroupsRequest_Scope.ALL })
      .then((res) => Object.fromEntries(res.groups.map((g) => [g.id, g]))),
  )
  // Commits of all apps (not only the filtered ones) so that changing filters does not refetch.
  const commits = createMemo(() => getRepositoryCommits(apps().map((a) => a.commit)))

  const scopedApps = createMemo(() =>
    scope() === GetRepositoriesRequest_Scope.PUBLIC ? apps().filter((a) => repos()[a.repositoryId]) : apps(),
  )
  const originOf = (app: Application): RepositoryOrigin => {
    const repo = repos()[app.repositoryId]
    return repo ? repositoryURLToOrigin(repo.url) : 'Others'
  }
  const matchesQuery = (app: Application, q: string) =>
    [app.name, repos()[app.repositoryId]?.name, app.groupId && groups()[app.groupId]?.name].some((text) =>
      text?.toLowerCase().includes(q),
    )
  const entries = createMemo(() => {
    const q = query().toLowerCase()
    const filtered = scopedApps().filter(
      (app) => statuses().includes(applicationState(app)) && origin().includes(originOf(app)) && matchesQuery(app, q),
    )
    return groupApplications(filtered, sort())
  })
  const scopeChanging = () => isPending(() => apps())

  return (
    <div class="h-full overflow-y-auto">
      <WithNav.Container>
        <Title>Apps - NeoShowcase</Title>
        <WithNav.Navs>
          <AppsNav />
          <WithNav.Tabs>
            <For each={scopeItems(me().admin)}>
              {(s) => (
                <TabRound state={s.value === scope() ? 'active' : 'default'} onClick={() => setScope(s.value)}>
                  <div class="i-material-symbols:deployed-code-outline shrink-0 text-2xl/6" />
                  {s.label}
                </TabRound>
              )}
            </For>
            <a href="/apps/new" class="ml-auto">
              <Button
                variants="primary"
                size="medium"
                leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
              >
                Add New App
              </Button>
            </a>
          </WithNav.Tabs>
        </WithNav.Navs>
        <WithNav.Body>
          <div class="bg-ui-background flex h-full w-full flex-col px-[max(calc(50%-500px),32px)] max-md:px-4">
            <FilterContainer>
              <TextField
                placeholder="Search"
                value={query()}
                onInput={(e) => setQuery(e.currentTarget.value)}
                leftIcon={<div class="i-material-symbols:search shrink-0 text-2xl/6" />}
                rightIcon={
                  <Loading>
                    <AppsFilter
                      apps={scopedApps()}
                      originOf={originOf}
                      statuses={statuses()}
                      setStatues={setStatuses}
                      origin={origin()}
                      setOrigin={setOrigin}
                      sort={sort()}
                      setSort={setSort}
                    />
                  </Loading>
                }
              />
            </FilterContainer>
            <Loading fallback={<Skeletons />}>
              <PendingContainer isPending={scopeChanging()}>
                <Show
                  when={entries().length > 0}
                  fallback={
                    <div class="text-text-black mt-32 flex flex-col items-center justify-center gap-4">
                      <div class="i-material-symbols:search text-20/20 shrink-0" />
                      <span class="font-bold">No Apps Found</span>
                    </div>
                  }
                >
                  {/* `content-visibility: auto` lets the browser skip layout/paint of off-screen entries,
                      which is what the former JS virtualizer was for, without breaking find-in-page or scroll anchoring. */}
                  <For each={entries()}>
                    {(entry) => (
                      <div class="mb-10 [contain-intrinsic-size:auto_300px] [content-visibility:auto]">
                        <AppsList
                          group={entry.groupId === undefined ? undefined : groups()[entry.groupId]}
                          apps={entry.apps}
                          repositories={repos()}
                          commits={commits()}
                        />
                      </div>
                    )}
                  </For>
                </Show>
              </PendingContainer>
            </Loading>
          </div>
        </WithNav.Body>
      </WithNav.Container>
    </div>
  )
}
