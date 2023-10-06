import { Navigate, RouteDataFunc, useRouteData, useRoutes } from '@solidjs/router'
import { Resource, createResource, lazy } from 'solid-js'
import { Application, GetApplicationsRequest_Scope, Repository } from './api/neoshowcase/protobuf/gateway_pb'
import { client } from './libs/api'

const RepositoryData: RouteDataFunc<
  unknown,
  {
    repo: Resource<Repository>
    refetchRepo: () => void
    apps: Resource<Application[]>
  }
> = ({ params }) => {
  const [repo, { refetch: refetchRepo }] = createResource(
    () => params.id,
    (id) => client.getRepository({ repositoryId: id }),
  )
  const [apps] = createResource(repo, async (repo) => {
    const allAppsRes = await client.getApplications({
      scope: GetApplicationsRequest_Scope.ALL,
    })
    return allAppsRes.applications.filter((app) => app.repositoryId === repo.id)
  })
  return {
    repo,
    refetchRepo,
    apps,
  }
}
export const useRepositoryData = () => useRouteData<ReturnType<typeof RepositoryData>>()

const ApplicationData: RouteDataFunc<
  unknown,
  {
    app: Resource<Application>
    refetchApp: () => void
    repo: Resource<Repository>
  }
> = ({ params }) => {
  const [app, { refetch: refetchApp }] = createResource(
    () => params.id,
    (id) => client.getApplication({ id }),
  )
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )
  return {
    app,
    refetchApp,
    repo,
  }
}
export const useApplicationData = () => useRouteData<ReturnType<typeof ApplicationData>>()

export default useRoutes([
  {
    path: '/',
    component: () => <Navigate href="/apps" />,
  },
  {
    path: '/apps',
    component: lazy(() => import('/@/pages/apps')),
  },
  {
    path: '/apps/:id',
    component: lazy(() => import('/@/pages/apps/[id]')),
    data: ApplicationData,
    children: [
      {
        path: '/',
        component: lazy(() => import('/@/pages/apps/[id]/index')),
      },
      {
        path: '/builds',
        component: lazy(() => import('/@/pages/apps/[id]/builds')),
      },
      {
        path: '/builds/:buildID',
        component: lazy(() => import('/@/pages/apps/[id]/builds/[id]')),
      },
      {
        path: '/settings',
        component: lazy(() => import('/@/pages/apps/[id]/settings')),
        children: [
          {
            path: '/',
            component: lazy(() => import('/@/pages/apps/[id]/settings/general')),
          },
          {
            path: '/build',
            component: lazy(() => import('/@/pages/apps/[id]/settings/build')),
          },
          {
            path: '/domains',
            component: lazy(() => import('/@/pages/apps/[id]/settings/domains')),
          },
          {
            path: '/portForwarding',
            component: lazy(() => import('/@/pages/apps/[id]/settings/portForwarding')),
          },
          {
            path: '/envVars',
            component: lazy(() => import('/@/pages/apps/[id]/settings/envVars')),
          },
          {
            path: '/owner',
            component: lazy(() => import('/@/pages/apps/[id]/settings/owner')),
          },
        ],
      },
    ],
  },
  {
    path: '/apps/new',
    component: lazy(() => import('/@/pages/apps/new')),
  },
  {
    path: '/builds',
    component: lazy(() => import('/@/pages/builds')),
  },
  {
    path: '/repos/:id',
    component: lazy(() => import('/@/pages/repos/[id]')),
    data: RepositoryData,
    children: [
      {
        path: '/',
        component: lazy(() => import('/@/pages/repos/[id]/index')),
      },
      {
        path: '/settings',
        component: lazy(() => import('/@/pages/repos/[id]/settings')),
        children: [
          {
            path: '/',
            component: lazy(() => import('/@/pages/repos/[id]/settings/general')),
          },
          {
            path: '/authorization',
            component: lazy(() => import('/@/pages/repos/[id]/settings/authorization')),
          },
          {
            path: '/owner',
            component: lazy(() => import('/@/pages/repos/[id]/settings/owner')),
          },
        ],
      },
    ],
  },
  {
    path: '/repos/new',
    component: lazy(() => import('/@/pages/repos/new')),
  },
  {
    path: '/settings',
    component: lazy(() => import('/@/pages/settings')),
  },
])
