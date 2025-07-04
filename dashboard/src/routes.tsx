import {
  createAsync,
  Navigate,
  Route,
  type RouteLoadFunc,
  Router,
  type RouteSectionProps,
  useParams,
} from '@solidjs/router'
import { type Component, createMemo, lazy } from 'solid-js'
import ErrorView from './components/layouts/ErrorView'
import {
  getApplication,
  getBuild,
  getBuilds,
  getRepository,
  getRepositoryApps,
  getRepositoryCommits,
  hasApplicationPermission,
  hasRepositoryPermission,
  revalidateApplication,
  revalidateBuild,
  revalidateBuilds,
  revalidateRepository,
} from './libs/api'

const loadApplicationData: RouteLoadFunc = ({ params }) => {
  getApplication(params.id).then((app) => {
    void getRepository(app.repositoryId)
    void getRepositoryCommits([app.commit])
  })
}

export const useApplicationData = () => {
  const params = useParams()
  const app = createAsync(() => getApplication(params.id))
  const repo = createAsync(async () => app() && (await getRepository(app()?.repositoryId)))
  const builds = createAsync(() => getBuilds(params.id))
  const hashes = () => {
    const a = app()
    const b = builds()
    if (a && b) return [a.commit, ...b.map((b) => b.commit)]
    return undefined
  }
  const commits = createAsync(async () => {
    const h = hashes()
    if (!h) return undefined
    return getRepositoryCommits(h)
  })
  const refetch = async () => {
    await Promise.all([revalidateApplication(params.id), revalidateBuilds(params.id)])
  }
  const hasPermission = createMemo(() => hasApplicationPermission(app))
  return {
    app,
    repo,
    builds,
    commits,
    refetch,
    hasPermission,
  }
}

const loadRepositoryData: RouteLoadFunc = ({ params }) => {
  void getRepository(params.id)
  getRepositoryApps(params.id).then((apps) => {
    const hashes = apps.map((app) => app.commit)
    void getRepositoryCommits(hashes)
  })
}

export const useRepositoryData = () => {
  const params = useParams()
  const repo = createAsync(() => getRepository(params.id))
  const apps = createAsync(() => getRepositoryApps(params.id))
  const commits = createAsync(async () => {
    const a = apps()
    if (!a) return undefined
    return getRepositoryCommits(a.map((a) => a.commit))
  })
  const refetchRepo = () => revalidateRepository(params.id)
  const hasPermission = createMemo(() => hasRepositoryPermission(repo))
  return {
    repo,
    apps,
    commits,
    refetchRepo,
    hasPermission,
  }
}

const loadBuildData: RouteLoadFunc = ({ params }) => {
  void getApplication(params.id)
  getBuild(params.buildID).then((build) => {
    void getRepositoryCommits([build.commit])
  })
}

export const useBuildData = () => {
  const params = useParams()
  const app = createAsync(() => getApplication(params.id))
  const build = createAsync(() => getBuild(params.buildID))
  const commit = createAsync(async () => {
    const hash = build()?.commit
    if (!hash) return undefined
    return getRepositoryCommits([hash]).then((c) => c[hash])
  })
  const refetch = async () => {
    await Promise.all([revalidateApplication(params.id), revalidateBuild(params.buildID), revalidateBuilds(params.id)])
  }
  const hasPermission = () => hasApplicationPermission(app)
  return {
    app,
    build,
    commit,
    refetch,
    hasPermission,
  }
}

export const Routes: Component<{ root: Component<RouteSectionProps> }> = (props) => (
  <Router root={props.root}>
    <Route path="/" component={() => <Navigate href="/apps" />} />
    <Route path="/apps" component={lazy(() => import('/@/pages/apps'))} />
    <Route path="/apps/:id" component={lazy(() => import('/@/pages/apps/[id]'))} preload={loadApplicationData}>
      <Route path="/" component={lazy(() => import('/@/pages/apps/[id]/index'))} />
      <Route path="/builds" component={lazy(() => import('/@/pages/apps/[id]/builds'))} />
      <Route
        path="/builds/:buildID"
        component={lazy(() => import('/@/pages/apps/[id]/builds/[id]'))}
        preload={loadBuildData}
      />
      <Route path="/settings" component={lazy(() => import('/@/pages/apps/[id]/settings'))}>
        <Route path="/" component={lazy(() => import('/@/pages/apps/[id]/settings/general'))} />
        <Route path="/build" component={lazy(() => import('/@/pages/apps/[id]/settings/build'))} />
        <Route path="/urls" component={lazy(() => import('/@/pages/apps/[id]/settings/urls'))} />
        <Route path="/portForwarding" component={lazy(() => import('/@/pages/apps/[id]/settings/portForwarding'))} />
        <Route path="/envVars" component={lazy(() => import('/@/pages/apps/[id]/settings/envVars'))} />
        <Route path="/owners" component={lazy(() => import('/@/pages/apps/[id]/settings/owners'))} />
      </Route>
      <Route path="/" component={lazy(() => import('/@/pages/apps/[id]/index'))} />
      <Route path="/" component={lazy(() => import('/@/pages/apps/[id]/index'))} />
    </Route>
    <Route path="/apps/new" component={lazy(() => import('/@/pages/apps/new'))} />
    <Route path="/repos/:id" component={lazy(() => import('/@/pages/repos/[id]'))} preload={loadRepositoryData}>
      <Route path="/" component={lazy(() => import('/@/pages/repos/[id]/index'))} />
      <Route path="/settings" component={lazy(() => import('/@/pages/repos/[id]/settings'))}>
        <Route path="/" component={lazy(() => import('/@/pages/repos/[id]/settings/general'))} />
        <Route path="/authorization" component={lazy(() => import('/@/pages/repos/[id]/settings/authorization'))} />
        <Route path="/owners" component={lazy(() => import('/@/pages/repos/[id]/settings/owners'))} />
      </Route>
    </Route>
    <Route path="/repos/new" component={lazy(() => import('/@/pages/repos/new'))} />
    <Route path="/settings" component={lazy(() => import('/@/pages/settings'))} />
    <Route path="/builds" component={lazy(() => import('/@/pages/builds'))} />
    <Route path="/*" component={() => <ErrorView error={new Error('Not Found')} />} />
  </Router>
)
