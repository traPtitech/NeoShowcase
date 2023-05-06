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
])
