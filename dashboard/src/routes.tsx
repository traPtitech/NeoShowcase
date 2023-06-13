import { Navigate, useRoutes } from '@solidjs/router'
import { lazy } from 'solid-js'

export default useRoutes([
  {
    path: '/',
    component: () => <Navigate href='/apps' />,
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
    path: '/repos/:id',
    component: lazy(() => import('/@/pages/repos/[id]')),
  },
  {
    path: '/repos/:id/settings',
    component: lazy(() => import('/@/pages/repos/[id]/settings')),
  },
  {
    path: '/repos/new',
    component: lazy(() => import('/@/pages/repos/new')),
  },
])
