import { styled } from '@macaron-css/solid'
import { Title } from '@solidjs/meta'
import { A } from '@solidjs/router'
import { createVirtualizer } from '@tanstack/solid-virtual'
import Fuse from 'fuse.js'
import { Component, For, Suspense, createMemo, createResource, createSignal, useTransition } from 'solid-js'
import {
  Application,
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  Repository,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { SelectOption } from '/@/components/templates/Select'
import { client, user } from '/@/libs/api'
import { ApplicationState, Provider, applicationState, repositoryURLToProvider } from '/@/libs/application'
import { createLocalSignal } from '/@/libs/localStore'
import { Button } from '../components/UI/Button'
import { MaterialSymbols } from '../components/UI/MaterialSymbols'
import { TabRound } from '../components/UI/TabRound'
import { TextField } from '../components/UI/TextField'
import SuspenseContainer from '../components/layouts/SuspenseContainer'
import { WithNav } from '../components/layouts/WithNav'
import { AppsNav } from '../components/templates/AppsNav'
import { List, RepositoryList } from '../components/templates/List'
import AppsFilter from '../components/templates/app/AppsFilter'
import { colorVars, media } from '../theme'

const MainView = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    height: '100%',
    overflowY: 'auto',
    padding: '0 max(calc(50% - 500px), 32px)',
    background: colorVars.semantic.ui.background,

    '@media': {
      [media.mobile]: {
        padding: '0 16px',
      },
    },
  },
})
const FilterContainer = styled('div', {
  base: {
    position: 'sticky',
    width: '100%',
    top: '0',
    left: '0',
    padding: '40px 0 32px',
    zIndex: 1,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: `linear-gradient(0deg, rgba(255,255,255,0), ${colorVars.semantic.ui.background} 20px)`,
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

export const allStatuses: SelectOption<ApplicationState>[] = [
  { label: 'Idle', value: ApplicationState.Idle },
  { label: 'Deploying', value: ApplicationState.Deploying },
  { label: 'Running', value: ApplicationState.Running },
  { label: 'Serving', value: ApplicationState.Serving },
  { label: 'Error', value: ApplicationState.Error },
]
export const allProviders: SelectOption<Provider>[] = [
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
  includeNoApp: boolean
  parentRef: HTMLDivElement
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
    const res = filteredReposByProvider().reduce<RepoWithApp[]>((acc, repo) => {
      if (!props.includeNoApp && !appsMap[repo.id]) return acc
      acc.push({ repo, apps: appsMap[repo.id] || [] })
      return acc
    }, [])
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

  const virtualizer = createMemo(() =>
    createVirtualizer({
      count: filteredRepos().length,
      getScrollElement: () => props.parentRef,
      estimateSize: (i) => 76 + 16 + filteredRepos()[i].apps.length * 80,
      paddingEnd: 72,
    }),
  )

  const items = () => virtualizer().getVirtualItems()

  return (
    <div
      style={{
        width: '100%',
        height: `${virtualizer().getTotalSize()}px`,
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
        <For
          each={items() ?? []}
          fallback={
            <List.Container>
              <List.PlaceHolder>
                <MaterialSymbols displaySize={80}>search</MaterialSymbols>
                No Apps Found
              </List.PlaceHolder>
            </List.Container>
          }
        >
          {(vRow) => (
            <div ref={virtualizer().measureElement}>
              <div style={{ 'padding-bottom': '16px' }}>
                <RepositoryList repository={filteredRepos()[vRow.index].repo} apps={filteredRepos()[vRow.index].apps} />
              </div>
            </div>
          )}
        </For>
      </div>
    </div>
  )
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
    'apps-statuses-v1',
    allStatuses.map((s) => s.value),
  )
  const [provider, setProvider] = createLocalSignal<Provider[]>('apps-provider', ['GitHub', 'GitLab', 'Gitea'])
  const [query, setQuery] = createLocalSignal('apps-query', '')
  const [sort, setSort] = createLocalSignal<keyof typeof sortItems>('apps-sort', sortItems.desc.value)
  const [includeNoApp, setIncludeNoApp] = createLocalSignal('apps-include-no-app', false)

  const [scrollParentRef, setScrollParentRef] = createSignal<HTMLDivElement>()

  return (
    <WithNav.Container>
      <Title>Apps - NeoShowcase</Title>
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
          <A href="/apps/new" style={{ 'margin-left': 'auto' }}>
            <Button variants="primary" size="medium" leftIcon={<MaterialSymbols>add</MaterialSymbols>}>
              Add New App
            </Button>
          </A>
        </WithNav.Tabs>
      </WithNav.Navs>
      <WithNav.Body>
        <MainView ref={setScrollParentRef}>
          <FilterContainer>
            <TextField
              placeholder="Search"
              value={query()}
              onInput={(e) => setQuery(e.currentTarget.value)}
              leftIcon={<MaterialSymbols>search</MaterialSymbols>}
              rightIcon={
                <AppsFilter
                  statuses={statuses()}
                  setStatues={setStatuses}
                  provider={provider()}
                  setProvider={setProvider}
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
              <Repositories>
                <RepositoryList apps={[undefined]} />
                <RepositoryList apps={[undefined]} />
                <RepositoryList apps={[undefined]} />
                <RepositoryList apps={[undefined]} />
              </Repositories>
            }
          >
            <SuspenseContainer isPending={isPending()}>
              <AppsList
                scope={scope()}
                statuses={statuses()}
                provider={provider()}
                query={query()}
                sort={sort()}
                includeNoApp={includeNoApp()}
                parentRef={scrollParentRef()!}
              />
            </SuspenseContainer>
          </Suspense>
        </MainView>
      </WithNav.Body>
    </WithNav.Container>
  )
}
