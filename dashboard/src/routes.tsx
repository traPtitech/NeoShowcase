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
])
