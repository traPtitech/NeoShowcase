import { Navigate, RouteDataFunc, useRouteData, useRoutes } from '@solidjs/router'
import { Resource, createResource, lazy } from 'solid-js'
import { Application, GetApplicationsRequest_Scope, Repository } from './api/neoshowcase/protobuf/gateway_pb'
import { client } from './libs/api'

const RepositoryData: RouteDataFunc<
  unknown,
  {
    repo: Resource<Repository>
    apps: Resource<Application[]>
  }
> = ({ params }) => {
  const [repo] = createResource(
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
    apps,
  }
}
export const useRepositoryData = () => useRouteData<ReturnType<typeof RepositoryData>>()

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
  },
  {
    path: '/apps/:id/builds',
    component: lazy(() => import('/@/pages/apps/[id]/builds')),
  },
  {
    path: '/apps/:id/builds/:buildID',
    component: lazy(() => import('/@/pages/apps/[id]/builds/[id]')),
  },
  {
    path: '/apps/:id/settings',
    component: lazy(() => import('/@/pages/apps/[id]/settings')),
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
