import { Header } from '/@/components/Header'
import { Checkbox } from '/@/components/Checkbox'
import { createMemo, createResource, For, Show } from 'solid-js'
import { Radio, RadioItem } from '/@/components/Radio'
import { client, user } from '/@/libs/api'
import {
  Application,
  GetApplicationsRequest_Scope,
  GetRepositoriesRequest_Scope,
  Repository,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { RepositoryRow } from '/@/components/RepositoryRow'
import { applicationState, ApplicationState } from '/@/libs/application'
import { AppStatus } from '/@/components/AppStatus'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { Container, PageTitle } from '/@/libs/layout'
import { Button } from '/@/components/Button'
import { useNavigate } from '@solidjs/router'
import Fuse from 'fuse.js'
import { unique } from '/@/libs/unique'
import { createLocalSignal } from '/@/libs/localStore'

const sortItems: RadioItem<'asc' | 'desc'>[] = [
  { value: 'desc', title: '最新順' },
  { value: 'asc', title: '古い順' },
]

const scopeItems = (admin: boolean) => {
  const items: RadioItem<GetRepositoriesRequest_Scope>[] = [
    { value: GetRepositoriesRequest_Scope.MINE, title: '自分のアプリ' },
    { value: GetRepositoriesRequest_Scope.PUBLIC, title: 'すべてのアプリ' },
  ]
  if (admin) {
    items.push({ value: GetRepositoriesRequest_Scope.ALL, title: 'すべてのアプリ (admin)' })
  }
  return items
}

const ContentContainer = styled('div', {
  base: {
    marginTop: '24px',
    display: 'grid',
    gridTemplateColumns: '380px 1fr',
    gap: '40px',
  },
})

const SidebarContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '22px',

    padding: '24px 40px',
    backgroundColor: vars.bg.white1,
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
  },
})

const SidebarSection = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
  },
})

const SidebarTitle = styled('div', {
  base: {
    fontSize: '24px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const SidebarOptions = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '12px',

    fontSize: '20px',
    color: vars.text.black1,
  },
})

const MainContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
})

const SearchBarContainer = styled('div', {
  base: {
    display: 'grid',
    gridTemplateColumns: '1fr 180px',
    gap: '20px',
    height: '44px',
  },
})

const SearchBar = styled('input', {
  base: {
    padding: '12px 20px',
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
    fontSize: '14px',

    '::placeholder': {
      color: vars.text.black3,
    },
  },
})

const RepositoriesContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
})

interface RepoWithApp {
  repo: Repository
  apps: Application[]
}

const newestAppDate = (apps: Application[]): number => Math.max(0, ...apps.map((a) => a.updatedAt.toDate().getTime()))
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

const allStatuses = [
  ApplicationState.Idle,
  ApplicationState.Deploying,
  ApplicationState.Running,
  ApplicationState.Static,
]

export default () => {
  const navigate = useNavigate()

  const [statuses, setStatuses] = createLocalSignal('apps-statuses', [...allStatuses])
  const checkStatus = (status: ApplicationState, checked: boolean) => {
    if (checked) {
      setStatuses((statuses) => unique([status, ...statuses]))
    } else {
      setStatuses((statuses) => statuses.filter((s) => s !== status))
    }
  }

  const [scope, setScope] = createLocalSignal('apps-scope', GetRepositoriesRequest_Scope.MINE)
  const appScope = () => {
    const mine = scope() === GetRepositoriesRequest_Scope.MINE
    return mine ? GetApplicationsRequest_Scope.MINE : GetApplicationsRequest_Scope.ALL
  }
  const [query, setQuery] = createLocalSignal('apps-query', '')
  const [sort, setSort] = createLocalSignal('apps-sort', sortItems[0].value)

  const [repos] = createResource(
    () => scope(),
    (scope) => client.getRepositories({ scope }),
  )
  const [apps] = createResource(
    () => appScope(),
    (scope) => client.getApplications({ scope }),
  )
  const loaded = () => !!(user() && repos() && apps())

  const filteredApps = createMemo(() => {
    if (!apps()) return
    const s = statuses()
    return apps().applications.filter((a) => s.includes(applicationState(a)))
  })
  const repoWithApps = createMemo(() => {
    if (!repos() || !filteredApps()) return
    const appsMap = {} as Record<string, Application[]>
    for (const app of filteredApps()) {
      if (!appsMap[app.repositoryId]) appsMap[app.repositoryId] = []
      appsMap[app.repositoryId].push(app)
    }
    const res = repos().repositories.map((repo): RepoWithApp => ({ repo, apps: appsMap[repo.id] || [] }))
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
    <Container>
      <Header />
      <PageTitle>Apps</PageTitle>
      <Show when={loaded()}>
        <ContentContainer>
          <SidebarContainer>
            <SidebarSection>
              <SidebarTitle>Status</SidebarTitle>
              <SidebarOptions>
                <For each={allStatuses}>
                  {(status) => (
                    <Checkbox selected={statuses().includes(status)} setSelected={(s) => checkStatus(status, s)}>
                      <AppStatus apps={apps().applications} state={status} />
                    </Checkbox>
                  )}
                </For>
              </SidebarOptions>
            </SidebarSection>
            <SidebarSection>
              <SidebarTitle>Scope</SidebarTitle>
              <SidebarOptions>
                <Radio items={scopeItems(user().admin)} selected={scope()} setSelected={setScope} />
              </SidebarOptions>
            </SidebarSection>
            <SidebarOptions>
              <SidebarTitle>Sort</SidebarTitle>
              <Radio items={sortItems} selected={sort()} setSelected={setSort} />
            </SidebarOptions>
          </SidebarContainer>
          <MainContainer>
            <SearchBarContainer>
              <SearchBar value={query()} onInput={(e) => setQuery(e.target.value)} placeholder='Search...' />
              <Button color='black1' size='large' width='full' onclick={() => navigate('/repos/new')}>
                + New Repository
              </Button>
            </SearchBarContainer>
            <RepositoriesContainer>
              <For each={filteredRepos()}>{(r) => <RepositoryRow repo={r.repo} apps={r.apps} />}</For>
            </RepositoriesContainer>
          </MainContainer>
        </ContentContainer>
      </Show>
    </Container>
  )
}
