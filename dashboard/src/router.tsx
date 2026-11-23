import { createRouter, intentPreload, type RoutePreloadFunc } from '@solidjs/router'
import { lazy } from 'solid-js'
import {
  getApplication,
  getApplicationGroup,
  getApplicationGroupApps,
  getBuild,
  getBuilds,
  getRepository,
  getRepositoryApps,
  getRepositoryCommits,
} from '/@/libs/api'

// Preloads only warm the query cache; the route components read the same queries.
const preloadApplication: RoutePreloadFunc = ({ params }) => {
  if (!params.id) return
  void getBuilds(params.id)
  void getApplication(params.id).then((app) => {
    void getRepository(app.repositoryId)
    void getRepositoryCommits([app.commit])
  })
}

const preloadRepository: RoutePreloadFunc = ({ params }) => {
  if (!params.id) return
  void getRepository(params.id)
  void getRepositoryApps(params.id).then((apps) => getRepositoryCommits(apps.map((app) => app.commit)))
}

const preloadApplicationGroup: RoutePreloadFunc = ({ params }) => {
  if (!params.id) return
  void getApplicationGroup(params.id)
  void getApplicationGroupApps(params.id).then((apps) => getRepositoryCommits(apps.map((app) => app.commit)))
}

const preloadBuild: RoutePreloadFunc = ({ params }) => {
  if (!params.buildID) return
  void getBuild(params.buildID).then((build) => getRepositoryCommits([build.commit]))
}

export const Router = createRouter({
  preloadLinks: intentPreload(),
  routes: [
    { path: '/', component: lazy(() => import('/@/pages/index')) },
    { path: '/apps', component: lazy(() => import('/@/pages/apps')) },
    { path: '/apps/new', component: lazy(() => import('/@/pages/apps/new')) },
    {
      path: '/apps/:id',
      component: lazy(() => import('/@/pages/apps/[id]')),
      preload: preloadApplication,
      children: [
        { path: '/', component: lazy(() => import('/@/pages/apps/[id]/index')) },
        { path: '/builds', component: lazy(() => import('/@/pages/apps/[id]/builds')) },
        {
          path: '/builds/:buildID',
          component: lazy(() => import('/@/pages/apps/[id]/builds/[id]')),
          preload: preloadBuild,
        },
        {
          path: '/settings',
          component: lazy(() => import('/@/pages/apps/[id]/settings')),
          children: [
            { path: '/', component: lazy(() => import('/@/pages/apps/[id]/settings/general')) },
            { path: '/build', component: lazy(() => import('/@/pages/apps/[id]/settings/build')) },
            { path: '/urls', component: lazy(() => import('/@/pages/apps/[id]/settings/urls')) },
            { path: '/portForwarding', component: lazy(() => import('/@/pages/apps/[id]/settings/portForwarding')) },
            { path: '/envVars', component: lazy(() => import('/@/pages/apps/[id]/settings/envVars')) },
            { path: '/group', component: lazy(() => import('/@/pages/apps/[id]/settings/group')) },
            { path: '/owners', component: lazy(() => import('/@/pages/apps/[id]/settings/owners')) },
          ],
        },
      ],
    },
    { path: '/repos/new', component: lazy(() => import('/@/pages/repos/new')) },
    {
      path: '/repos/:id',
      component: lazy(() => import('/@/pages/repos/[id]')),
      preload: preloadRepository,
      children: [
        { path: '/', component: lazy(() => import('/@/pages/repos/[id]/index')) },
        {
          path: '/settings',
          component: lazy(() => import('/@/pages/repos/[id]/settings')),
          children: [
            { path: '/', component: lazy(() => import('/@/pages/repos/[id]/settings/general')) },
            { path: '/authorization', component: lazy(() => import('/@/pages/repos/[id]/settings/authorization')) },
            { path: '/owners', component: lazy(() => import('/@/pages/repos/[id]/settings/owners')) },
          ],
        },
      ],
    },
    { path: '/groups', component: lazy(() => import('/@/pages/groups')) },
    { path: '/groups/new', component: lazy(() => import('/@/pages/groups/new')) },
    {
      path: '/groups/:id',
      component: lazy(() => import('/@/pages/groups/[id]')),
      preload: preloadApplicationGroup,
      children: [
        { path: '/', component: lazy(() => import('/@/pages/groups/[id]/index')) },
        {
          path: '/settings',
          component: lazy(() => import('/@/pages/groups/[id]/settings')),
          children: [
            { path: '/', component: lazy(() => import('/@/pages/groups/[id]/settings/general')) },
            { path: '/owners', component: lazy(() => import('/@/pages/groups/[id]/settings/owners')) },
          ],
        },
      ],
    },
    { path: '/builds', component: lazy(() => import('/@/pages/builds')) },
    { path: '/settings', component: lazy(() => import('/@/pages/settings')) },
    { path: '*404', component: lazy(() => import('/@/pages/[...404]')) },
  ],
})

export const { paths } = Router
