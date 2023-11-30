import { styled } from '@macaron-css/solid'
import { createVirtualizer } from '@tanstack/solid-virtual'
import Fuse from 'fuse.js'
import { Component, For, Show, Suspense, createMemo, createResource, useTransition } from 'solid-js'
import {
  Application,
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  Repository,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { MultiSelect, SelectOption, SingleSelect } from '/@/components/templates/Select'
import { client, systemInfo, user } from '/@/libs/api'
import { ApplicationState, Provider, applicationState, repositoryURLToProvider } from '/@/libs/application'
import { createLocalSignal } from '/@/libs/localStore'
import { MaterialSymbols } from '../components/UI/MaterialSymbols'
import { TabRound } from '../components/UI/TabRound'
import { TextField } from '../components/UI/TextField'
import { MainViewContainer } from '../components/layouts/MainView'
import SuspenseContainer from '../components/layouts/SuspenseContainer'
import { WithNav } from '../components/layouts/WithNav'
import { AppsNav } from '../components/templates/AppsNav'
import { RepositoryList } from '../components/templates/List'
import { media } from '../theme'

const MainView = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'grid',
    gridTemplateRows: 'auto 1fr',
    flexDirection: 'column',
    gap: '32px',
  },
})
const SortContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '16px',

    '@media': {
      [media.mobile]: {
        flexDirection: 'column',
        gap: '8px',
      },
    },
  },
})
const SortSelects = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
})
const Repositories = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
  },
})

const sortItems: { [k in 'desc' | 'asc']: SelectOption<k> } = {
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
interface RepoWithApp {
  repo: Repository
  apps: Application[]
}

const newestAppDate = (apps: Application[]): number =>
  Math.max(0, ...apps.map((a) => a.updatedAt?.toDate().getTime() ?? 0))
const compareRepoWithApp =
  (sort: 'asc' | 'desc') =>
  (a: RepoWithApp, b: RepoWithApp): number => {
    // Sort by apps updated at
    if (a.apps.length > 0 && b.apps.length > 0) {
      if (sort === 'asc') {
        return newestAppDate(a.apps) - newestAppDate(b.apps)
      } else {
        return newestAppDate(b.apps) - newestAppDate(a.apps)
      }
    }
    // Bring up repositories with 1 or more apps at top
    if ((a.apps.length > 0 && b.apps.length === 0) || (a.apps.length === 0 && b.apps.length > 0)) {
      return b.apps.length - a.apps.length
    }
    // Fallback to sort by repository id
    return a.repo.id.localeCompare(b.repo.id)
  }

const allStatuses: SelectOption<ApplicationState>[] = [
  { label: 'Idle', value: ApplicationState.Idle },
  { label: 'Deploying', value: ApplicationState.Deploying },
  { label: 'Running', value: ApplicationState.Running },
  { label: 'Static', value: ApplicationState.Static },
  { label: 'Error', value: ApplicationState.Error },
]
const allProviders: SelectOption<Provider>[] = [
  { label: 'GitHub', value: 'GitHub' },
  { label: 'GitLab', value: 'GitLab' },
  { label: 'Gitea', value: 'Gitea' },
]

const AppsList: Component<{
  scope: GetRepositoriesRequest_Scope
  statuses: ApplicationState[]
  provider: Provider[]
  query: string
  sort: keyof typeof sortItems
}> = (props) => {
  const appScope = () => {
    const mine = props.scope === GetRepositoriesRequest_Scope.MINE
    return mine ? GetApplicationsRequest_Scope.MINE : GetApplicationsRequest_Scope.ALL
  }
  const [repos] = createResource(
    () => props.scope,
    (scope) => client.getRepositories({ scope }),
  )
  const [apps] = createResource(
    () => appScope(),
    (scope) => client.getApplications({ scope }),
  )

  const filteredReposByProvider = createMemo(() => {
    const p = props.provider
    return repos()?.repositories.filter((r) => p.includes(repositoryURLToProvider(r.url))) ?? []
  })
  const filteredApps = createMemo(() => {
    const s = props.statuses
    return apps()?.applications.filter((a) => s.includes(applicationState(a))) ?? []
  })
  const repoWithApps = createMemo(() => {
    const appsMap = {} as Record<string, Application[]>
    for (const app of filteredApps()) {
      if (!appsMap[app.repositoryId]) appsMap[app.repositoryId] = []
      appsMap[app.repositoryId].push(app)
    }
    const res = filteredReposByProvider().map((repo): RepoWithApp => ({ repo, apps: appsMap[repo.id] || [] }))
    res.sort(compareRepoWithApp(props.sort))
    return res
  })

  const fuse = createMemo(() => {
    return new Fuse(repoWithApps(), {
      keys: ['repo.name', 'apps.name'],
    })
  })
  const filteredRepos = createMemo(() => {
    if (props.query === '') return repoWithApps()
    return fuse()
      .search(props.query)
      .map((r) => r.item)
  })

  let scrollParentRef: HTMLDivElement | undefined
  const virtualizer = createMemo(() =>
    createVirtualizer({
      count: filteredRepos().length,
      getScrollElement: () => scrollParentRef,
      estimateSize: (i) => 76 + 16 + filteredRepos()[i].apps.length * 80,
    }),
  )

  const items = () => virtualizer().getVirtualItems()

  return (
    <div
      ref={scrollParentRef}
      class="List"
      style={{
        width: `100%`,
        height: `100%`,
        'padding-bottom': '16px',
        overflow: 'auto',
      }}
    >
      <div
        style={{
          height: `${virtualizer().getTotalSize()}px`,
          width: '100%',
          position: 'relative',
        }}
      >
        <div
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            transform: `translateY(${items()?.[0]?.start ?? 0}px)`,
          }}
        >
          <For each={items() ?? []}>
            {(vRow) => (
              <div ref={vRow?.measureElement}>
                <div style={{ 'padding-bottom': '16px' }}>
                  <RepositoryList
                    repository={filteredRepos()[vRow.index].repo}
                    apps={filteredRepos()[vRow.index].apps}
                  />
                </div>
              </div>
            )}
          </For>
        </div>
      </div>
    </div>
  )
  // return (
  //   <Repositories>
  //     <For each={filteredRepos()}>{(r) => <RepositoryList repository={r.repo} apps={r.apps} />}</For>
  //   </Repositories>
  // )
}

export default () => {
  const [scope, _setScope] = createLocalSignal('apps-scope', GetRepositoriesRequest_Scope.MINE)
  const [isPending, start] = useTransition()

  const setScope = (scope: GetRepositoriesRequest_Scope) => {
    start(() => {
      _setScope(scope)
    })
  }

  const [statuses, setStatuses] = createLocalSignal(
    'apps-statuses',
    allStatuses.map((s) => s.value),
  )
  const [provider, setProvider] = createLocalSignal<Provider[]>('apps-provider', ['GitHub', 'GitLab', 'Gitea'])
  const [query, setQuery] = createLocalSignal('apps-query', '')
  const [sort, setSort] = createLocalSignal<keyof typeof sortItems>('apps-sort', sortItems.asc.value)

  return (
    <WithNav.Container>
      <WithNav.Navs>
        <AppsNav />
        <WithNav.Tabs>
          <For each={scopeItems(user()?.admin)}>
            {(s) => (
              <TabRound state={s.value === scope() ? 'active' : 'default'} onClick={() => setScope(s.value)}>
                <MaterialSymbols>deployed_code</MaterialSymbols>
                {s.label}
              </TabRound>
            )}
          </For>
          <Show when={systemInfo()?.adminerUrl}>
            <a
              href={systemInfo()?.adminerUrl}
              target="_blank"
              rel="noopener noreferrer"
              style={{
                'margin-left': 'auto',
              }}
            >
              <TabRound variant="ghost">
                Adminer
                <MaterialSymbols>open_in_new</MaterialSymbols>
              </TabRound>
            </a>
          </Show>
        </WithNav.Tabs>
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer background="grey" scrollable={false}>
          <MainView>
            <SortContainer>
              <TextField
                placeholder="Search"
                value={query()}
                onInput={(e) => setQuery(e.currentTarget.value)}
                leftIcon={<MaterialSymbols>search</MaterialSymbols>}
              />
              <SortSelects>
                <MultiSelect options={allStatuses} placeholder="Status" value={statuses()} setValue={setStatuses} />
                <MultiSelect options={allProviders} placeholder="Provider" value={provider()} setValue={setProvider} />
                <SingleSelect options={Object.values(sortItems)} placeholder="Sort" value={sort()} setValue={setSort} />
              </SortSelects>
            </SortContainer>
            <div style={{ height: '100%', 'overflow-y': 'hidden' }}>
              <Suspense
                fallback={
                  <Repositories>
                    <RepositoryList apps={[undefined]} />
                    <RepositoryList apps={[undefined]} />
                    <RepositoryList apps={[undefined]} />
                    <RepositoryList apps={[undefined]} />
                  </Repositories>
                }
              >
                <SuspenseContainer isPending={isPending()}>
                  <AppsList scope={scope()} statuses={statuses()} provider={provider()} query={query()} sort={sort()} />
                </SuspenseContainer>
              </Suspense>
            </div>
          </MainView>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
