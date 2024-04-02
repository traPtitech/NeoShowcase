import {
  Navigate,
  Route,
  type RouteLoadFunc,
  type RouteSectionProps,
  Router,
  createAsync,
  useParams,
} from '@solidjs/router'
import { type Component, createMemo, lazy } from 'solid-js'
import ErrorView from './components/layouts/ErrorView'
import {
  getApplication,
  getBuild,
  getRepository,
  getRepositoryApps,
  hasApplicationPermission,
  hasRepositoryPermission,
  revalidateApplication,
  revalidateBuild,
  revalidateRepository,
} from './libs/api'

const loadApplicationData: RouteLoadFunc = ({ params }) => {
  getApplication(params.id).then((app) => {
    void getRepository(app.repositoryId)
  })
}

export const useApplicationData = () => {
  const params = useParams()
  const app = createAsync(() => getApplication(params.id))
  const repo = createAsync(async () => app() && (await getRepository(app()?.repositoryId)))
  const refetchApp = () => revalidateApplication(params.id)
  const hasPermission = createMemo(() => hasApplicationPermission(app))
  return {
    app,
    repo,
    refetchApp,
    hasPermission,
  }
}

const loadRepositoryData: RouteLoadFunc = ({ params }) => {
  void getRepository(params.id)
  void getRepositoryApps(params.id)
}

export const useRepositoryData = () => {
  const params = useParams()
  const repo = createAsync(() => getRepository(params.id))
  const apps = createAsync(() => getRepositoryApps(params.id))
  const refetchRepo = () => revalidateRepository(params.id)
  const hasPermission = createMemo(() => hasRepositoryPermission(repo))
  return {
    repo,
    apps,
    refetchRepo,
    hasPermission,
  }
}

const loadBuildData: RouteLoadFunc = ({ params }) => {
  void getApplication(params.id)
  void getBuild(params.buildID)
}

export const useBuildData = () => {
  const params = useParams()
  const app = createAsync(() => getApplication(params.id))
  const build = createAsync(() => getBuild(params.buildID))
  const refetchApp = () => revalidateApplication(params.id)
  const refetchBuild = () => revalidateBuild(params.buildID)
  const hasPermission = () => hasApplicationPermission(app)
  return {
    app,
    build,
    refetchApp,
    refetchBuild,
    hasPermission,
  }
}

declare module '@solidjs/router' {
  type RouteProps<S extends string> = {
    // Invalid component type? workaround
    component?: Component<RouteSectionProps>
  }
}

export const Routes: Component<{ root: Component<RouteSectionProps> }> = (props) => (
  <Router root={props.root}>
    <Route path="/" component={() => <Navigate href="/apps" />} />
    <Route path="/apps" component={lazy(() => import('/@/pages/apps'))} />
    <Route path="/apps/:id" component={lazy(() => import('/@/pages/apps/[id]'))} load={loadApplicationData}>
      <Route path="/" component={lazy(() => import('/@/pages/apps/[id]/index'))} />
      <Route path="/builds" component={lazy(() => import('/@/pages/apps/[id]/builds'))} />
      <Route
        path="/builds/:buildID"
        component={lazy(() => import('/@/pages/apps/[id]/builds/[id]'))}
        load={loadBuildData}
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
    <Route path="/repos/:id" component={lazy(() => import('/@/pages/repos/[id]'))} load={loadRepositoryData}>
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
