import { Title } from '@solidjs/meta'
import { A } from '@solidjs/router'
import { createVirtualizer } from '@tanstack/solid-virtual'
import Fuse from 'fuse.js'
import { type Component, For, Suspense, createMemo, createResource, createSignal, useTransition } from 'solid-js'
import { GetApplicationsRequest_Scope, GetRepositoriesRequest_Scope } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '/@/components/styled-components'
import type { SelectOption } from '/@/components/templates/Select'
import { client, getRepositoryCommits, user } from '/@/libs/api'
import { ApplicationState, type RepositoryOrigin, useApplicationsFilter, RepoWithApp } from '/@/libs/application'
import { createSessionSignal } from '/@/libs/localStore'
import { Button } from '../components/UI/Button'
import { TabRound } from '../components/UI/TabRound'
import { TextField } from '../components/UI/TextField'
import SuspenseContainer from '../components/layouts/SuspenseContainer'
import { WithNav } from '../components/layouts/WithNav'
import { AppsNav } from '../components/templates/AppsNav'
import { List, RepositoryList } from '../components/templates/List'
import AppsFilter from '../components/templates/app/AppsFilter'

const FilterContainer = styled('div', 'z-1 flex w-full items-center justify-center pt-10 pb-8')

export const sortItems: { [k in 'desc' | 'asc']: SelectOption<k> } = {
  desc: { value: 'desc', label: 'Newest' },
  asc: { value: 'asc', label: 'Oldest' },
}

const scopeItems = (admin: boolean | undefined) => {
  const items: SelectOption<GetRepositoriesRequest_Scope>[] = [
    { value: GetRepositoriesRequest_Scope.MINE, label: 'My Apps' },
    { value: GetRepositoriesRequest_Scope.PUBLIC, label: 'All Apps' },
  ]
  if (admin) {
    items.push({
      value: GetRepositoriesRequest_Scope.ALL,
      label: 'All Apps (admin)',
    })
  }
  return items
}

export const allStatuses: SelectOption<ApplicationState>[] = [
  { label: 'Idle', value: ApplicationState.Idle },
  { label: 'Deploying', value: ApplicationState.Deploying },
  { label: 'Running', value: ApplicationState.Running },
  { label: 'Sleeping', value: ApplicationState.Sleeping },
  { label: 'Serving', value: ApplicationState.Serving },
  { label: 'Error', value: ApplicationState.Error },
]
export const allOrigins: SelectOption<RepositoryOrigin>[] = [
  { label: 'GitHub', value: 'GitHub' },
  { label: 'Gitea', value: 'Gitea' },
  { label: 'Others', value: 'Others' },
]

const AppsList: Component<{
  repoWithApps: RepoWithApp[]
  query: string
  parentRef: HTMLDivElement
}> = (props) => {
  const hashes = () => props.repoWithApps.flatMap((r) => r.apps.map((a) => a.commit))
  const [commits] = createResource(
    () => hashes(),
    (hashes) => getRepositoryCommits(hashes),
  )

  const fuse = createMemo(() => {
    return new Fuse(props.repoWithApps, {
      keys: ['repo.name', 'apps.name'],
    })
  })
  const filteredRepos = createMemo(() => {
    if (props.query === '') return props.repoWithApps
    return fuse()
      .search(props.query)
      .map((r) => r.item)
  })

  const virtualizer = createMemo(() =>
    createVirtualizer({
      count: filteredRepos().length,
      getScrollElement: () => props.parentRef,
      estimateSize: (i) => 76 + 16 + filteredRepos()[i].apps.length * 80,
      // scrollParentRef内に高さ120pxのFilterContainerが存在するため、この分を設定
      scrollMargin: 120,
      paddingEnd: 72,
      gap: 16,
    }),
  )

  const items = () => virtualizer().getVirtualItems()

  return (
    <div
      class="relative w-full"
      style={{
        height: `${virtualizer().getTotalSize()}px`,
        // scrollParentRef内に高さ120pxのFilterContainerが存在するため、この分を減算
        transform: `translateY(${-virtualizer().options.scrollMargin}px)`,
      }}
    >
      <For
        each={items() ?? []}
        fallback={
          <List.Container>
            <List.PlaceHolder>
              <div class="i-material-symbols:search shrink-0 text-20/20" />
              No Apps Found
            </List.PlaceHolder>
          </List.Container>
        }
      >
        {(vRow) => (
          <div
            data-index={vRow.index}
            ref={(el) => queueMicrotask(() => virtualizer().measureElement(el))}
            class="absolute top-0 left-0 w-full"
            style={{
              height: `${items()[vRow.index]}px`,
              transform: `translateY(${vRow.start}px)`,
            }}
          >
            <RepositoryList
              repository={filteredRepos()[vRow.index].repo}
              apps={filteredRepos()[vRow.index].apps}
              commits={commits()}
            />
          </div>
        )}
      </For>
    </div>
  )
}

export default () => {
  const [scope, _setScope] = createSessionSignal('apps-scope', GetRepositoriesRequest_Scope.MINE)
  const [isPending, start] = useTransition()

  const setScope = (scope: GetRepositoriesRequest_Scope) => {
    start(() => {
      _setScope(scope)
    })
  }

  const [statuses, setStatuses] = createSessionSignal(
    'apps-statuses-v1',
    allStatuses.map((s) => s.value),
  )
  const [origin, setOrigin] = createSessionSignal<RepositoryOrigin[]>('apps-repository-origin', [
    'GitHub',
    'Gitea',
    'Others',
  ])
  const [query, setQuery] = createSessionSignal('apps-query', '')
  const [sort, setSort] = createSessionSignal<keyof typeof sortItems>('apps-sort', sortItems.desc.value)
  const [includeNoApp, setIncludeNoApp] = createSessionSignal('apps-include-no-app', false)

  const [scrollParentRef, setScrollParentRef] = createSignal<HTMLDivElement>()

  const appScope = () => {
    const mine = scope() === GetRepositoriesRequest_Scope.MINE
    return mine ? GetApplicationsRequest_Scope.MINE : GetApplicationsRequest_Scope.ALL
  }
  const [repos] = createResource(
    () => scope(),
    (scope) => client.getRepositories({ scope }),
  )
  const [apps] = createResource(
    () => appScope(),
    (scope) => client.getApplications({ scope }),
  )
  const repoWithApps = () =>
    useApplicationsFilter(
      repos()?.repositories ?? [],
      apps()?.applications ?? [],
      statuses(),
      origin(),
      includeNoApp(),
      sort(),
    )

  return (
    <div class="h-full overflow-y-auto" ref={setScrollParentRef}>
      <WithNav.Container>
        <Title>Apps - NeoShowcase</Title>
        <WithNav.Navs>
          <AppsNav />
          <WithNav.Tabs>
            <For each={scopeItems(user()?.admin)}>
              {(s) => (
                <TabRound state={s.value === scope() ? 'active' : 'default'} onClick={() => setScope(s.value)}>
                  <div class="i-material-symbols:deployed-code-outline shrink-0 text-2xl/6" />
                  {s.label}
                </TabRound>
              )}
            </For>
            <A href="/apps/new" class="ml-auto">
              <Button
                variants="primary"
                size="medium"
                leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
              >
                Add New App
              </Button>
            </A>
          </WithNav.Tabs>
        </WithNav.Navs>
        <WithNav.Body>
          <div class="relative flex h-full w-full flex-col bg-ui-background px-[max(calc(50%-500px),32px)] max-md:px-4">
            <FilterContainer>
              <TextField
                placeholder="Search"
                value={query()}
                onInput={(e) => setQuery(e.currentTarget.value)}
                leftIcon={<div class="i-material-symbols:search shrink-0 text-2xl/6" />}
                rightIcon={
                  <AppsFilter
                    statuses={statuses()}
                    setStatues={setStatuses}
                    origin={origin()}
                    setOrigin={setOrigin}
                    sort={sort()}
                    setSort={setSort}
                    includeNoApp={includeNoApp()}
                    setIncludeNoApp={setIncludeNoApp}
                  />
                }
              />
            </FilterContainer>
            <Suspense
              fallback={
                <div class="flex w-full flex-col gap-4">
                  <RepositoryList apps={[undefined]} />
                  <RepositoryList apps={[undefined]} />
                  <RepositoryList apps={[undefined]} />
                  <RepositoryList apps={[undefined]} />
                </div>
              }
            >
              <SuspenseContainer isPending={isPending()}>
                <AppsList repoWithApps={repoWithApps()} query={query()} parentRef={scrollParentRef()!} />
              </SuspenseContainer>
            </Suspense>
          </div>
        </WithNav.Body>
      </WithNav.Container>
    </div>
  )
}
