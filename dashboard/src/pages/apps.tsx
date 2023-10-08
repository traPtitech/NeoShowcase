import {
  Application,
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  Repository,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { MultiSelect, SelectItem, SingleSelect } from '/@/components/templates/Select'
import { client, user } from '/@/libs/api'
import { ApplicationState, applicationState, repositoryURLToProvider } from '/@/libs/application'
import { createLocalSignal } from '/@/libs/localStore'
import { styled } from '@macaron-css/solid'
import Fuse from 'fuse.js'
import { For, Show, createMemo, createResource } from 'solid-js'
import { Provider } from '../components/RepositoryRow'
import { MaterialSymbols } from '../components/UI/MaterialSymbols'
import { TabRound } from '../components/UI/TabRound'
import { TextInput } from '../components/UI/TextInput'
import { WithNav } from '../components/layouts/WithNav'
import { AppsNav } from '../components/templates/AppsNav'
import { RepositoryList } from '../components/templates/List'
import { colorVars } from '../theme'

const MainViewContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    padding: '40px 32px 72px 32px',
    overflowY: 'auto',
    background: colorVars.semantic.ui.background,
  },
})
const MainView = styled('div', {
  base: {
    width: '100%',
    maxWidth: '1000px',
    margin: '0 auto',
    display: 'flex',
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
      'screen and (max-width: 768px)': {
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

const sortItems: { [k in 'desc' | 'asc']: SelectItem<k> } = {
  desc: { value: 'desc', title: 'Newest' },
  asc: { value: 'asc', title: 'Oldest' },
}

const scopeItems = (admin: boolean) => {
  const items: SelectItem<GetRepositoriesRequest_Scope>[] = [
    { value: GetRepositoriesRequest_Scope.MINE, title: 'My Apps' },
    { value: GetRepositoriesRequest_Scope.PUBLIC, title: 'All Apps' },
  ]
  if (admin) {
    items.push({
      value: GetRepositoriesRequest_Scope.ALL,
      title: 'All Apps (admin)',
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
const compareRepoWithApp = (sort: 'asc' | 'desc') => (a: RepoWithApp, b: RepoWithApp): number => {
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

const allStatuses: SelectItem<ApplicationState>[] = [
  { title: 'Idle', value: ApplicationState.Idle },
  { title: 'Deploying', value: ApplicationState.Deploying },
  { title: 'Running', value: ApplicationState.Running },
  { title: 'Static', value: ApplicationState.Static },
  { title: 'Error', value: ApplicationState.Error },
]
const allProviders: SelectItem<Provider>[] = [
  { title: 'GitHub', value: 'GitHub' },
  { title: 'GitLab', value: 'GitLab' },
  { title: 'Gitea', value: 'Gitea' },
]

export default () => {
  const [statuses, setStatuses] = createLocalSignal(
    'apps-statuses',
    allStatuses.map((s) => s.value),
  )
  const [scope, setScope] = createLocalSignal('apps-scope', GetRepositoriesRequest_Scope.MINE)
  const [provider, setProvider] = createLocalSignal<Provider[]>('apps-provider', ['GitHub', 'GitLab', 'Gitea'])
  const appScope = () => {
    const mine = scope() === GetRepositoriesRequest_Scope.MINE
    return mine ? GetApplicationsRequest_Scope.MINE : GetApplicationsRequest_Scope.ALL
  }
  const [query, setQuery] = createLocalSignal('apps-query', '')
  const [sort, setSort] = createLocalSignal<keyof typeof sortItems>('apps-sort', sortItems.asc.value)

  const [repos] = createResource(
    () => scope(),
    (scope) => client.getRepositories({ scope }),
  )
  const [apps] = createResource(
    () => appScope(),
    (scope) => client.getApplications({ scope }),
  )
  const loaded = () => !!(user() && repos() && apps())

  const filteredReposByProvider = createMemo(() => {
    if (!repos()) return
    const p = provider()
    return repos()?.repositories.filter((r) => p.includes(repositoryURLToProvider(r.url)))
  })
  const filteredApps = createMemo(() => {
    if (!apps()) return
    const s = statuses()
    return apps()?.applications.filter((a) => s.includes(applicationState(a)))
  })
  const repoWithApps = createMemo(() => {
    if (!filteredReposByProvider() || !filteredApps()) return
    const appsMap = {} as Record<string, Application[]>
    for (const app of filteredApps()) {
      if (!appsMap[app.repositoryId]) appsMap[app.repositoryId] = []
      appsMap[app.repositoryId].push(app)
    }
    const res = filteredReposByProvider().map((repo): RepoWithApp => ({ repo, apps: appsMap[repo.id] || [] }))
    res.sort(compareRepoWithApp(sort()))
    return res
  })

  const fuse = createMemo(() => {
    if (!repoWithApps()) return
    return new Fuse(repoWithApps(), {
      keys: ['repo.name', 'apps.name'],
    })
  })
  const filteredRepos = createMemo(() => {
    if (!repoWithApps()) return
    if (query() === '') return repoWithApps()
    return fuse()
      .search(query())
      .map((r) => r.item)
  })

  return (
    <Show when={loaded()}>
      <WithNav.Container>
        <WithNav.Navs>
          <AppsNav />
          <WithNav.Tabs>
            <For each={scopeItems(user()?.admin)}>
              {(s) => (
                <TabRound state={s.value === scope() ? 'active' : 'default'} onClick={() => setScope(s.value)}>
                  <MaterialSymbols>deployed_code</MaterialSymbols>
                  {s.title}
                </TabRound>
              )}
            </For>
          </WithNav.Tabs>
        </WithNav.Navs>
        <WithNav.Body>
          <MainViewContainer>
            <MainView>
              <SortContainer>
                <TextInput
                  placeholder="Search"
                  value={query()}
                  onInput={(e) => setQuery(e.target.value)}
                  leftIcon={<MaterialSymbols>search</MaterialSymbols>}
                />
                <SortSelects>
                  <MultiSelect
                    placeHolder="Status"
                    items={allStatuses}
                    selected={statuses()}
                    setSelected={setStatuses}
                  />
                  <MultiSelect
                    placeHolder="Provider"
                    items={allProviders}
                    selected={provider()}
                    setSelected={setProvider}
                  />
                  <SingleSelect
                    placeHolder="Sort"
                    items={Object.values(sortItems)}
                    selected={sort()}
                    setSelected={setSort}
                  />
                </SortSelects>
              </SortContainer>
              <Repositories>
                <For each={filteredRepos()}>{(r) => <RepositoryList repository={r.repo} apps={r.apps} />}</For>
              </Repositories>
            </MainView>
          </MainViewContainer>
        </WithNav.Body>
      </WithNav.Container>
    </Show>
  )
}
