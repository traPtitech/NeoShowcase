import { describe, it } from 'vitest'
import { idle, scenario, waitForText } from '/@/testing/harness'

/** Renders the real <App/> against the in-memory backend for every route and checks for errors. */

const routes: { url: string; texts: string[] }[] = [
  { url: '/apps', texts: ['runtime-app', 'static-app', 'broken-app'] },
  { url: '/apps/new', texts: ['Repository'] },
  { url: '/apps/new?repositoryID=repo-github', texts: ['github-repo'] },
  { url: '/apps/app-runtime', texts: ['runtime-app', 'feat: add feature A'] },
  { url: '/apps/app-runtime/builds', texts: ['Build History', 'feat: add feature A'] },
  { url: '/apps/app-runtime/builds/build-ok', texts: ['feat: add feature A'] },
  { url: '/apps/app-static/builds/build-running', texts: ['streamed build line 1'] },
  { url: '/apps/app-runtime/settings', texts: ['runtime-app'] },
  { url: '/apps/app-runtime/settings/build', texts: ['runtime-app'] },
  { url: '/apps/app-runtime/settings/urls', texts: ['.trap.show', '8080'] },
  { url: '/apps/app-runtime/settings/portForwarding', texts: ['5000'] },
  { url: '/apps/app-runtime/settings/envVars', texts: ['FOO', 'bar'] },
  { url: '/apps/app-runtime/settings/owners', texts: ['admin-user'] },
  { url: '/apps/app-static', texts: ['static-app'] },
  { url: '/apps/app-static/settings/build', texts: ['static-app'] },
  { url: '/apps/app-error', texts: ['broken-app'] },
  { url: '/repos/repo-github', texts: ['github-repo', 'runtime-app'] },
  { url: '/repos/repo-gitea', texts: ['gitea-repo', 'broken-app'] },
  { url: '/repos/repo-github/settings', texts: ['github-repo'] },
  { url: '/repos/repo-github/settings/authorization', texts: ['github-repo'] },
  { url: '/repos/repo-github/settings/owners', texts: ['admin-user'] },
  { url: '/repos/new', texts: ['Register Repository'] },
  { url: '/groups', texts: ['web-service'] },
  { url: '/groups/new', texts: ['Create Group'] },
  { url: '/groups/group-web', texts: ['web-service', 'static-app'] },
  { url: '/groups/group-web/settings', texts: ['web-service', 'Delete Group'] },
  { url: '/groups/group-web/settings/owners', texts: ['admin-user'] },
  { url: '/apps/app-static/settings/group', texts: ['web-service'] },
  { url: '/apps/app-static/settings/owners', texts: ['web-service', 'admin-user'] },
  { url: '/builds', texts: ['feat: add feature A'] },
  { url: '/settings', texts: ['laptop'] },
  { url: '/does/not/exist', texts: ['Not Found'] },
]

describe('route smoke tests', () => {
  it.each(routes)('renders $url', async ({ url, texts }) => {
    await scenario(
      url,
      url,
      async (ctx) => {
        await waitForText(ctx, ...texts)
        await idle(ctx)
      },
      { expectErrorView: url === '/does/not/exist' },
    )
  })
})
