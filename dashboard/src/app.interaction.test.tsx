import { Code } from '@connectrpc/connect'
import { describe, expect, it, vi } from 'vitest'
import { Repository_AuthMethod } from '/@/api/neoshowcase/protobuf/gateway_pb'
import {
  button,
  buttonContaining,
  click,
  dialog,
  field,
  go,
  idle,
  link,
  scenario,
  settle,
  type,
  waitForNoText,
  waitForText,
} from '/@/testing/harness'

/** User flows against the in-memory backend, where Solid 2 / Kobalte 2 / Formisch integration is most fragile. */

describe('interactions', () => {
  it('create app: walks the three steps and calls createApplication', async () => {
    await scenario('create app', '/apps/new', async (ctx) => {
      await waitForText(ctx, 'github-repo')
      click(ctx, buttonContaining(ctx, 'github-repo'))
      await settle(ctx, () => expect(window.location.search).toBe('?repositoryID=repo-github'))
      await waitForText(ctx, 'Create Application from')

      type(ctx, field(ctx, 'Application Name'), 'new-app')
      // branch: type to open the combobox, then pick the option
      const branch = await settle(ctx, () => expect(ctx.screen.querySelector('[role=combobox]')).not.toBeNull()).then(
        () => ctx.screen.querySelector<HTMLInputElement>('[role=combobox]') as HTMLInputElement,
      )
      type(ctx, branch, 'mai')
      await settle(ctx, () => expect(ctx.screen.querySelector('[role=option]')?.textContent).toBe('main'))
      const option = ctx.screen.querySelector<HTMLElement>('[role=option]') as HTMLElement
      for (const e of ['pointerdown', 'pointerup', 'click']) {
        option.dispatchEvent(new MouseEvent(e, { bubbles: true, button: 0 }))
      }
      await settle(ctx, () => expect(branch.value).toBe('main'))

      const radio = (name: string, value: string) =>
        ctx.screen.querySelector<HTMLInputElement>(`input[type=radio][name='${name}'][value='${value}']`) as HTMLElement
      click(ctx, radio('["form","config","deployConfig","type"]', 'runtime'))
      await settle(ctx, () => expect(ctx.screen.textContent).toContain('Build Type'))
      click(ctx, radio('["form","config","buildConfig","type"]', 'buildpack'))

      // step 2 -> 3
      click(ctx, button(ctx, 'Next'))
      await settle(ctx, () => expect(button(ctx, 'Create Application')).toBeTruthy())
      // the form was already submitted by "Next", yet a freshly added URL row must not start out with errors
      click(ctx, button(ctx, 'Add URL'))
      const subdomain = await settle(ctx, () =>
        expect(ctx.screen.querySelector(`input[name$='"subdomain"]']`)).not.toBeNull(),
      ).then(() => ctx.screen.querySelector<HTMLInputElement>(`input[name$='"subdomain"]']`) as HTMLInputElement)
      await idle(ctx, 200)
      expect(ctx.screen.textContent).not.toContain('Please Enter Subdomain Name')
      type(ctx, subdomain, 'new-app')
      click(ctx, button(ctx, 'Create Application'))
      await settle(ctx, () =>
        expect(ctx.backend.requestsTo('createApplication')).toEqual([
          expect.objectContaining({ name: 'new-app', refName: 'main', repositoryId: 'repo-github' }),
        ]),
      )
      // navigates to the created app
      await settle(ctx, () => expect(window.location.pathname).toBe('/apps/app-new'))
      await waitForText(ctx, 'new-app')
    })
  })

  it('settings: add an SSH key through the modal', async () => {
    await scenario('add ssh key', '/settings', async (ctx) => {
      await waitForText(ctx, 'laptop')
      click(ctx, button(ctx, 'Add New SSH Key'))
      await settle(ctx, () => expect(document.querySelector('[role=dialog]')).not.toBeNull())

      // validation: submitting empty fields shows errors and does not call the API
      click(ctx, button(ctx, 'Add'))
      await waitForText(ctx, 'Enter Name', 'Enter SSH Public Key')
      expect(ctx.backend.requestsTo('createUserKey')).toHaveLength(0)

      type(ctx, field(ctx, 'Name'), 'desktop')
      type(ctx, field(ctx, 'Key'), 'ssh-ed25519 AAAAdesktop')
      // the errors are cleared by revalidation, which re-enables the (invalid-disabled) submit button
      await settle(ctx, () => expect(button(ctx, 'Add').disabled).toBe(false))
      click(ctx, button(ctx, 'Add'))

      await settle(ctx, () =>
        expect(ctx.backend.requestsTo('createUserKey')).toEqual([
          expect.objectContaining({ name: 'desktop', publicKey: 'ssh-ed25519 AAAAdesktop' }),
        ]),
      )
      // dialog closes and the list is revalidated
      await settle(ctx, () => expect(document.querySelector('[role=dialog]')).toBeNull())
      await waitForText(ctx, 'desktop')
    })
  })

  it('settings: cancelling the modal closes it', async () => {
    await scenario('cancel ssh key', '/settings', async (ctx) => {
      await waitForText(ctx, 'laptop')
      click(ctx, button(ctx, 'Add New SSH Key'))
      await settle(ctx, () => expect(document.querySelector('[role=dialog]')).not.toBeNull())
      click(ctx, button(ctx, 'Cancel'))
      await settle(ctx, () => expect(document.querySelector('[role=dialog]')).toBeNull())
      expect(ctx.backend.requestsTo('createUserKey')).toHaveLength(0)
    })
  })

  it('app general settings: rename enables Save and calls updateApplication', async () => {
    await scenario('rename app', '/apps/app-runtime/settings', async (ctx) => {
      await waitForText(ctx, 'Branch')
      await idle(ctx)
      expect(button(ctx, 'Save').disabled).toBe(true)

      await settle(ctx, () => field(ctx, 'Application Name'))
      type(ctx, field(ctx, 'Application Name'), 'renamed-app')
      await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
      click(ctx, button(ctx, 'Save'))

      await settle(ctx, () =>
        expect(ctx.backend.requestsTo('updateApplication')).toEqual([
          expect.objectContaining({ id: 'app-runtime', name: 'renamed-app' }),
        ]),
      )
      // refetch brings the new name into the header and resets the form (Save disabled again)
      await waitForText(ctx, 'renamed-app')
      await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(true))
    })
  })

  it('app env vars: typing into the trailing empty row appends another empty row', async () => {
    await scenario('env vars rows', '/apps/app-runtime/settings/envVars', async (ctx) => {
      await waitForText(ctx, 'FOO', 'bar')
      const keyInputs = () => [...ctx.screen.querySelectorAll<HTMLInputElement>(`input[name$='"key"]']`)]
      await idle(ctx)
      const before = keyInputs().length // FOO, NS_APP_ID (system) and the trailing empty row

      type(ctx, keyInputs().at(-1) as HTMLInputElement, 'NEW_KEY')
      await settle(ctx, () => expect(keyInputs()).toHaveLength(before + 1))
      expect(keyInputs().at(-1)?.value).toBe('')

      await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
      click(ctx, button(ctx, 'Save'))
      await settle(ctx, () =>
        expect(ctx.backend.requestsTo('setEnvVar')).toEqual([
          expect.objectContaining({ applicationId: 'app-runtime', key: 'NEW_KEY' }),
        ]),
      )
    })
  })

  it('apps list: the search box narrows the list', async () => {
    await scenario('apps filter', '/apps', async (ctx) => {
      await waitForText(ctx, 'runtime-app', 'static-app')
      const search = ctx.screen.querySelector<HTMLInputElement>('input[placeholder="Search"]')
      if (!search) throw new Error('no search box')
      type(ctx, search, 'static')
      await waitForNoText(ctx, 'broken-app')
      await waitForText(ctx, 'static-app')
      expect(search.value).toBe('static')

      type(ctx, search, 'zzz')
      await waitForText(ctx, 'No Apps Found')

      type(ctx, search, '')
      await waitForText(ctx, 'runtime-app', 'static-app')
    })
  })

  it('tooltip: shows its content on focus of a disabled button', async () => {
    // "Save" is disabled for a non-owner; admin-user owns runtime-app, but not broken-app (admins still pass isOwner),
    // so use the env var system row, whose inputs carry a tooltip.
    await scenario('tooltip', '/apps/app-runtime/settings/envVars', async (ctx) => {
      await waitForText(ctx, 'NS_APP_ID')
      const wrapper = ctx.screen.querySelector<HTMLElement>('input[disabled]')?.closest('span')
      if (!wrapper) throw new Error('no tooltip trigger around the disabled system variable input')
      wrapper.focus()
      wrapper.dispatchEvent(new FocusEvent('focus'))
      wrapper.dispatchEvent(new PointerEvent('pointermove', { bubbles: true, pointerType: 'mouse' }))
      await settle(ctx, () =>
        expect(document.querySelector('[role=tooltip]')?.textContent).toContain('システム環境変数'),
      )
    })
  })

  it('navigation: clicking a link changes the route client-side', async () => {
    await scenario('navigate', '/apps', async (ctx) => {
      await waitForText(ctx, 'runtime-app')
      click(ctx, link(ctx, '/apps/app-runtime'))
      await settle(ctx, () => expect(window.location.pathname).toBe('/apps/app-runtime'))
      await waitForText(ctx, 'Branch Resolution')
    })
  })

  it('app page: container logs are loaded and streamed', async () => {
    await scenario('container log', '/apps/app-runtime', async (ctx) => {
      await waitForText(ctx, 'previous log line', 'streamed line 1', 'streamed line 2')
      // leaving the page aborts the stream without errors (checked by scenario)
      click(ctx, link(ctx, '/apps'))
      await settle(ctx, () => expect(window.location.pathname).toBe('/apps'))
      await waitForText(ctx, 'static-app')
    })
  })

  it('app page: "Load more" prepends older container logs and keeps the loaded ones', async () => {
    await scenario('container log load more', '/apps/app-runtime', async (ctx) => {
      await waitForText(ctx, 'previous log line', 'streamed line 1', 'streamed line 2')
      click(ctx, button(ctx, 'Load more'))
      await waitForText(ctx, 'older log line 1', 'previous log line', 'streamed line 1', 'streamed line 2')
      click(ctx, button(ctx, 'Load more'))
      await waitForText(ctx, 'older log line 2', 'older log line 1', 'previous log line', 'streamed line 2')
      // the page revalidates the app every 10s; the loaded logs must survive it
      // (imported after mount: the harness resets modules per scenario, so a top-level import targets a stale cache;
      // the short timeout makes sure this revalidation, not the page's own 10s poll, is what we observe)
      const { revalidateApplication } = await import('/@/libs/api')
      const fetched = ctx.backend.requestsTo('getApplication').length
      revalidateApplication('app-runtime')
      await settle(ctx, () => expect(ctx.backend.requestsTo('getApplication').length).toBeGreaterThan(fetched), 2_000)
      await idle(ctx, 200)
      for (const t of ['older log line 2', 'older log line 1', 'previous log line', 'streamed line 1']) {
        expect(ctx.screen.textContent).toContain(t)
      }
      // each request continues from the oldest line loaded so far
      const befores = ctx.backend.requestsTo('getOutput').map((r) => (r.before as { seconds: bigint }).seconds)
      expect(befores).toHaveLength(3)
      expect(befores[1]).toBe(befores[0] - 60n)
      expect(befores[2]).toBe(befores[1] - 60n)
    })
  })

  it('app page: a "Load more" response of an earlier visit does not leak into a later visit of the same app', async () => {
    await scenario('container log stale load more', '/apps/app-runtime', async (ctx) => {
      // client-side navigation that reuses the app page (as Back/Forward between apps does)
      const go = (path: string) => {
        window.history.pushState(null, '', path)
        window.dispatchEvent(new PopStateEvent('popstate'))
      }
      await waitForText(ctx, 'previous log line', 'streamed line 1')
      const release = ctx.backend.holdNext('getOutput')
      click(ctx, button(ctx, 'Load more'))
      await settle(ctx, () => expect(ctx.backend.requestsTo('getOutput')).toHaveLength(2))

      go('/apps/app-error')
      await waitForText(ctx, 'broken-app', 'older log line 2')
      go('/apps/app-runtime')
      await waitForText(ctx, 'runtime-app', 'older log line 3')

      release() // the held response ('older log line 1', for the first visit) arrives now
      await idle(ctx, 200)
      expect(ctx.screen.textContent).not.toContain('older log line 1')
      expect(button(ctx, 'Load more').disabled).toBe(false)
    })
  })

  it('app page: the container log stream reconnects after it breaks, without repeating lines', async () => {
    await scenario(
      'container log reconnect',
      '/apps/app-runtime',
      async (ctx) => {
        vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] })
        try {
          ctx.backend.dropOutputStream(1)
          await waitForText(ctx, 'previous log line', 'streamed line 2')
          await settle(ctx, () => expect(ctx.backend.requestsTo('getOutputStream')).toHaveLength(1))
          await vi.advanceTimersByTimeAsync(3000)
          await settle(ctx, () => expect(ctx.backend.requestsTo('getOutputStream')).toHaveLength(2))
          await vi.advanceTimersByTimeAsync(0)
          ctx.flush()
        } finally {
          vi.useRealTimers()
        }
        await idle(ctx, 100)
        // the second connection resends both lines (the mock ignores `begin`): each is shown once
        const count = (text: string) => (ctx.screen.textContent ?? '').split(text).length - 1
        expect(count('streamed line 1')).toBe(1)
        expect(count('streamed line 2')).toBe(1)
        // it resumes from the last line received
        const second = ctx.backend.requestsTo('getOutputStream')[1] as { begin: { seconds: bigint } }
        expect(second.begin.seconds).toBe(BigInt(Date.parse('2024-01-02T00:00:02Z') / 1000))
      },
      { expectedErrors: ['stream broke'] },
    )
  })

  it('build page: a running build streams its log', async () => {
    await scenario('build log stream', '/apps/app-static/builds/build-running', async (ctx) => {
      await waitForText(ctx, 'streamed build line 1')
    })
  })

  it('build page: a still-queued build keeps retrying its log stream until it starts', async () => {
    await scenario(
      'queued build log',
      '/apps/app-static/builds/build-running',
      async (ctx) => {
        // "build not started yet" is expected for a queued build and must not surface as an error
        await waitForText(ctx, 'streamed build line 1')
        expect(ctx.backend.requestsTo('getBuildLogStream').length).toBeGreaterThanOrEqual(2)
      },
      { setup: (b) => b.failMethod('getBuildLogStream', 'build not started yet', 1, Code.FailedPrecondition) },
    )
  })

  it('repo settings: rename calls updateRepository', async () => {
    await scenario('rename repo', '/repos/repo-github/settings', async (ctx) => {
      await waitForText(ctx, 'github-repo')
      await idle(ctx)
      await settle(ctx, () => field(ctx, 'Repository Name'))
      const name = field(ctx, 'Repository Name')
      type(ctx, name, 'renamed-repo')
      await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
      click(ctx, button(ctx, 'Save'))
      await settle(ctx, () =>
        expect(ctx.backend.requestsTo('updateRepository')).toEqual([
          expect.objectContaining({ id: 'repo-github', name: 'renamed-repo' }),
        ]),
      )
      await waitForText(ctx, 'renamed-repo')
    })
  })

  it('app urls: "Add More" appends a website row and Save sends all websites', async () => {
    await scenario('add website', '/apps/app-runtime/settings/urls', async (ctx) => {
      await waitForText(ctx, '8080')
      await idle(ctx)
      const rows = () => ctx.screen.querySelectorAll(`input[name$='"subdomain"]']`).length
      const before = rows()
      click(ctx, button(ctx, 'Add More'))
      await settle(ctx, () => expect(rows()).toBe(before + 1))
      const inputs = [...ctx.screen.querySelectorAll<HTMLInputElement>(`input[name$='"subdomain"]']`)]
      type(ctx, inputs.at(-1) as HTMLInputElement, 'second')
      await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
      click(ctx, button(ctx, 'Save'))
      await settle(ctx, () => expect(ctx.backend.requestsTo('updateApplication')).toHaveLength(1))
      const [req] = ctx.backend.requestsTo('updateApplication') as { websites?: { websites: { fqdn: string }[] } }[]
      expect(req.websites?.websites.map((w) => w.fqdn)).toEqual(['runtime-app.trap.show', 'second.trap.show'])
    })
  })

  describe('state bound to the previous resource after the route param changes', () => {
    // Route components are reused across param changes, so anything captured for one resource must not
    // carry over to the next.
    it('app general settings: an open delete confirmation does not delete another app', async () => {
      await scenario('delete app retarget', '/apps/app-runtime/settings', async (ctx) => {
        await settle(ctx, () => button(ctx, 'Delete Application'))
        await idle(ctx)
        click(ctx, button(ctx, 'Delete Application'))
        await settle(ctx, () => expect(dialog(ctx)).not.toBeNull())

        go(ctx, '/apps/app-static/settings')
        await settle(ctx, () => expect(field(ctx, 'Application Name').value).toBe('static-app'))
        await idle(ctx)
        // the confirmation was for runtime-app: it is gone, and nothing deletes static-app in its place
        const confirm = [...ctx.screen.querySelectorAll<HTMLElement>('button')].find(
          (b) => b.textContent?.trim() === 'Yes, Delete',
        )
        if (confirm) click(ctx, confirm)
        await idle(ctx)
        expect(ctx.backend.requestsTo('deleteApplication')).toEqual([])
        expect(dialog(ctx)).toBeNull()
      })
    })

    it('repo general settings: an open delete confirmation does not delete another repo', async () => {
      await scenario(
        'delete repo retarget',
        '/repos/repo-github/settings',
        async (ctx) => {
          await settle(ctx, () => expect(button(ctx, 'Delete Repository').disabled).toBe(false))
          click(ctx, button(ctx, 'Delete Repository'))
          await settle(ctx, () => expect(dialog(ctx)).not.toBeNull())

          go(ctx, '/repos/repo-gitea/settings')
          await settle(ctx, () => expect(field(ctx, 'Repository Name').value).toBe('gitea-repo'))
          await idle(ctx)
          const confirm = [...ctx.screen.querySelectorAll<HTMLElement>('button')].find(
            (b) => b.textContent?.trim() === 'Yes, Delete',
          )
          if (confirm) click(ctx, confirm)
          await idle(ctx)
          expect(ctx.backend.requestsTo('deleteRepository')).toEqual([])
          expect(dialog(ctx)).toBeNull()
        },
        // a repository can only be deleted once it has no applications
        { setup: (b) => b.fixtures.applications.splice(0) },
      )
    })

    it('repo authorization: a key generated for one repo does not land in another repo', async () => {
      await scenario(
        'generated key retarget',
        '/repos/repo-gitea/settings/authorization',
        async (ctx) => {
          await settle(ctx, () => buttonContaining(ctx, '専用公開鍵を生成する'))
          await idle(ctx)
          const release = ctx.backend.holdMethod('generateKeyPair')
          click(ctx, buttonContaining(ctx, '専用公開鍵を生成する'))
          await settle(ctx, () => expect(ctx.backend.requestsTo('generateKeyPair')).toHaveLength(1))

          // the response for gitea-repo arrives after moving to github-repo
          go(ctx, '/repos/repo-github/settings/authorization')
          await settle(ctx, () => expect(field(ctx, 'Repository URL').value).toContain('NeoShowcase'))
          release()
          await idle(ctx, 100)
          expect(ctx.screen.textContent).not.toContain('(このリポジトリ専用)')
          expect(field(ctx, 'Repository URL').closest('form')?.textContent).toContain('NeoShowcase全体共通')
          expect(button(ctx, 'Save').disabled).toBe(true)

          // a key generated for github-repo stays with github-repo
          click(ctx, buttonContaining(ctx, '専用公開鍵を生成する'))
          await waitForText(ctx, 'ssh-ed25519 AAAAgenerated')
          go(ctx, '/repos/repo-gitea/settings/authorization')
          await settle(ctx, () => expect(field(ctx, 'Repository URL').value).toContain('gitea-repo'))
          await idle(ctx)
          expect(ctx.screen.textContent).not.toContain('(このリポジトリ専用)')
          await waitForNoText(ctx, 'ssh-ed25519 AAAAgenerated')
        },
        {
          setup: (b) => {
            b.fixtures.repositories[0].authMethod = b.fixtures.repositories[1].authMethod
          },
        },
      )
    })

    it('build page: switching between two running builds shows the log of the selected one', async () => {
      await scenario('switch running builds', '/apps/app-static/builds/build-running', async (ctx) => {
        await waitForText(ctx, 'streamed build line 1 of build-running')
        go(ctx, '/apps/app-static/builds/build-running-2')
        await waitForText(ctx, 'streamed build line 1 of build-running-2')
        expect(ctx.screen.textContent).not.toMatch(/of build-running(?!-2)/)
        expect(ctx.backend.requestsTo('getBuildLogStream').map((r) => r.buildId)).toEqual([
          'build-running',
          'build-running-2',
        ])
      })
    })

    it('build page: Rebuild does not navigate back once the user has left', async () => {
      await scenario('rebuild then leave', '/apps/app-static/builds/build-running', async (ctx) => {
        await waitForText(ctx, 'streamed build line 1')
        vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] })
        try {
          click(ctx, button(ctx, 'Rebuild'))
          await settle(ctx, () => expect(ctx.backend.requestsTo('retryCommitBuild')).toHaveLength(1))
          await vi.advanceTimersByTimeAsync(0)
          click(ctx, link(ctx, '/apps'))
          await settle(ctx, () => expect(window.location.pathname).toBe('/apps'))
          await vi.advanceTimersByTimeAsync(2000)
          ctx.flush()
        } finally {
          vi.useRealTimers()
        }
        await idle(ctx)
        expect(window.location.pathname).toBe('/apps')
      })
    })
  })

  describe('list rows survive a refetch', () => {
    // Each refetch returns new objects; rows keyed by object identity would be recreated and lose their state.
    it('build history: refreshing the list keeps the focused build link', async () => {
      await scenario('builds refetch focus', '/apps/app-runtime/builds', async (ctx) => {
        await settle(ctx, () => link(ctx, '/apps/app-runtime/builds/build-ok'))
        await idle(ctx)
        const focused = link(ctx, '/apps/app-runtime/builds/build-ok')
        focused.focus()
        const before = ctx.backend.requestsTo('getBuilds').length
        const { revalidateBuilds } = await import('/@/libs/api')
        revalidateBuilds('app-runtime')
        await settle(ctx, () => expect(ctx.backend.requestsTo('getBuilds').length).toBeGreaterThan(before))
        await idle(ctx, 100)
        expect(focused.isConnected).toBe(true)
        expect(document.activeElement).toBe(focused)
      })
    })

    it('settings: refreshing the SSH keys keeps an open delete confirmation', async () => {
      await scenario('ssh keys refetch dialog', '/settings', async (ctx) => {
        await waitForText(ctx, 'laptop')
        click(ctx, button(ctx, 'Delete'))
        await waitForText(ctx, 'Delete SSH Key')
        // the kept row still shows the refetched data
        ctx.backend.fixtures.userKeys[0].name = 'renamed-laptop'
        const { revalidateUserKeys } = await import('/@/libs/api')
        revalidateUserKeys()
        await waitForText(ctx, 'renamed-laptop')
        expect(ctx.screen.textContent).toContain('Delete SSH Key')
      })
    })

    it('owners: retrying the session data keeps an open owner delete confirmation', async () => {
      await scenario('owners refetch dialog', '/apps/app-runtime/settings/owners', async (ctx) => {
        await waitForText(ctx, 'Add Owners')
        click(ctx, await settle(ctx, () => button(ctx, 'Delete')).then(() => button(ctx, 'Delete')))
        await waitForText(ctx, 'Delete Owner')
        const before = ctx.backend.requestsTo('getUsers').length
        const { retryAll } = await import('/@/libs/session')
        await retryAll()
        await settle(ctx, () => expect(ctx.backend.requestsTo('getUsers').length).toBeGreaterThan(before))
        await idle(ctx, 100)
        expect(ctx.screen.textContent).toContain('Delete Owner')
      })
    })
  })

  it('build page: a build whose live log is lost shows that, without retrying the stream', async () => {
    await scenario(
      'lost build log',
      '/apps/app-static/builds/build-running',
      async (ctx) => {
        await waitForText(ctx, 'ビルドログを表示できません')
        vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout', 'setInterval', 'clearInterval'] })
        try {
          await vi.advanceTimersByTimeAsync(10_000)
          ctx.flush()
        } finally {
          vi.useRealTimers()
        }
        expect(ctx.backend.requestsTo('getBuildLogStream')).toHaveLength(1)
      },
      {
        setup: (b) =>
          b.failMethod('getBuildLogStream', 'build log stream no longer available', 1, Code.FailedPrecondition),
      },
    )
  })

  it('build page: revalidating a finished build neither re-fetches its log nor scrolls it', async () => {
    await scenario('finished build revalidation', '/apps/app-runtime/builds/build-ok', async (ctx) => {
      await waitForText(ctx, 'step 2: done')
      await idle(ctx)
      const log = [...ctx.screen.querySelectorAll<HTMLElement>('code')].find((el) =>
        el.textContent?.includes('step 2: done'),
      ) as HTMLElement
      const scrolls = vi.fn()
      Object.defineProperty(log, 'scrollTop', { configurable: true, get: () => 0, set: scrolls })

      const { revalidateBuild, revalidateApplication } = await import('/@/libs/api')
      revalidateBuild('build-ok')
      revalidateApplication('app-runtime')
      await settle(ctx, () => expect(ctx.backend.requestsTo('getBuild').length).toBeGreaterThanOrEqual(2))
      await idle(ctx, 100)
      expect(ctx.backend.requestsTo('getBuildLog')).toHaveLength(1)
      expect(scrolls).not.toHaveBeenCalled()
    })
  })

  it('app page: an unchanged app revalidation does not refetch the metrics', async () => {
    await scenario('metrics on app revalidation', '/apps/app-runtime', async (ctx) => {
      await settle(ctx, () => expect(ctx.backend.requestsTo('getApplicationMetrics').length).toBeGreaterThan(0))
      await idle(ctx, 100)
      const metrics = ctx.backend.requestsTo('getApplicationMetrics').length
      const apps = ctx.backend.requestsTo('getApplication').length
      const { revalidateApplication } = await import('/@/libs/api')
      revalidateApplication('app-runtime')
      await settle(ctx, () => expect(ctx.backend.requestsTo('getApplication').length).toBeGreaterThan(apps))
      await idle(ctx, 100)
      expect(ctx.backend.requestsTo('getApplicationMetrics')).toHaveLength(metrics)
    })
  })

  describe('settings forms keep unsaved input across revalidation', () => {
    it('app general: a name typed after saving survives the follow-up refetches', async () => {
      await scenario('app rename twice', '/apps/app-runtime/settings', async (ctx) => {
        await settle(ctx, () => field(ctx, 'Application Name'))
        await idle(ctx)
        type(ctx, field(ctx, 'Application Name'), 'first-name')
        await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
        click(ctx, button(ctx, 'Save'))
        await settle(ctx, () => expect(ctx.backend.requestsTo('updateApplication')).toHaveLength(1))
        await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(true))

        type(ctx, field(ctx, 'Application Name'), 'second-name')
        // the save schedules a second refetch 1s later
        await idle(ctx, 1200)
        expect(field(ctx, 'Application Name').value).toBe('second-name')
        expect(button(ctx, 'Save').disabled).toBe(false)
      })
    })

    it('app urls: a newly added row survives an app revalidation', async () => {
      await scenario('add website then revalidate', '/apps/app-runtime/settings/urls', async (ctx) => {
        await waitForText(ctx, '8080')
        await idle(ctx)
        const rows = () => ctx.screen.querySelectorAll(`input[name$='"subdomain"]']`).length
        const before = rows()
        click(ctx, button(ctx, 'Add More'))
        await settle(ctx, () => expect(rows()).toBe(before + 1))

        const { revalidateApplication } = await import('/@/libs/api')
        await revalidateApplication('app-runtime')
        await idle(ctx, 200)
        expect(rows()).toBe(before + 1)
      })
    })

    it('app env vars: a draft survives retryAll', async () => {
      await scenario('env var draft', '/apps/app-runtime/settings/envVars', async (ctx) => {
        await waitForText(ctx, 'FOO', 'bar')
        await idle(ctx)
        const value = () => ctx.screen.querySelector<HTMLInputElement>(`input[name='["variables",0,"value"]']`)
        await settle(ctx, () => expect(value()?.value).toBe('bar'))
        type(ctx, value() as HTMLInputElement, 'draft')

        const { retryAll } = await import('/@/libs/session')
        await retryAll()
        await idle(ctx, 200)
        expect(value()?.value).toBe('draft')
      })
    })

    it('repo general: a name draft survives a repository revalidation', async () => {
      await scenario('repo name draft', '/repos/repo-github/settings', async (ctx) => {
        await settle(ctx, () => field(ctx, 'Repository Name'))
        await idle(ctx)
        type(ctx, field(ctx, 'Repository Name'), 'draft-repo')

        const { revalidateRepository } = await import('/@/libs/api')
        await revalidateRepository('repo-github')
        await idle(ctx, 200)
        expect(field(ctx, 'Repository Name').value).toBe('draft-repo')
      })
    })
  })

  describe('settings panels send only their own fields', () => {
    it('app general: saving a rename does not restore a website removed on the server', async () => {
      await scenario('rename keeps websites', '/apps/app-runtime/settings', async (ctx) => {
        await settle(ctx, () => field(ctx, 'Application Name'))
        await idle(ctx)
        // removed elsewhere while this panel is open
        const app = ctx.backend.fixtures.applications.find((a) => a.id === 'app-runtime')
        if (!app) throw new Error('no app-runtime fixture')
        app.websites = []

        type(ctx, field(ctx, 'Application Name'), 'renamed-app')
        await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
        click(ctx, button(ctx, 'Save'))
        await settle(ctx, () => expect(ctx.backend.requestsTo('updateApplication')).toHaveLength(1))
        const [req] = ctx.backend.requestsTo('updateApplication')
        expect(req).toMatchObject({ id: 'app-runtime', name: 'renamed-app' })
        expect(req.websites).toBeUndefined()
        expect(req.portPublications).toBeUndefined()
        expect(req.config).toBeUndefined()
        expect(req.ownerIds).toBeUndefined()
        expect(app.websites).toEqual([])
      })
    })

    it('repo general: renaming an SSH repository sends only its name', async () => {
      await scenario('rename ssh repo', '/repos/repo-gitea/settings', async (ctx) => {
        await settle(ctx, () => field(ctx, 'Repository Name'))
        await idle(ctx)
        type(ctx, field(ctx, 'Repository Name'), 'renamed-gitea')
        await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
        click(ctx, button(ctx, 'Save'))
        await settle(ctx, () => expect(ctx.backend.requestsTo('updateRepository')).toHaveLength(1))
        const [req] = ctx.backend.requestsTo('updateRepository')
        expect(req).toMatchObject({ id: 'repo-gitea', name: 'renamed-gitea' })
        expect(req.auth).toBeUndefined()
        expect(req.url).toBeUndefined()
      })
    })

    it('repo general: a BASIC-auth repository can be renamed', async () => {
      await scenario(
        'rename basic repo',
        '/repos/repo-github/settings',
        async (ctx) => {
          await settle(ctx, () => field(ctx, 'Repository Name'))
          await idle(ctx)
          type(ctx, field(ctx, 'Repository Name'), 'renamed-basic')
          await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
          click(ctx, button(ctx, 'Save'))
          // the hidden credential fields are not part of this panel, so they neither block nor join the request
          await settle(ctx, () => expect(ctx.backend.requestsTo('updateRepository')).toHaveLength(1))
          const [req] = ctx.backend.requestsTo('updateRepository')
          expect(req).toMatchObject({ id: 'repo-github', name: 'renamed-basic' })
          expect(req.auth).toBeUndefined()
        },
        {
          setup: (b) => {
            const repo = b.fixtures.repositories.find((r) => r.id === 'repo-github')
            if (repo) repo.authMethod = Repository_AuthMethod.BASIC
          },
        },
      )
    })

    it('repo authorization: changing only the URL keeps the SSH key', async () => {
      await scenario('repo url only', '/repos/repo-gitea/settings/authorization', async (ctx) => {
        await settle(ctx, () => field(ctx, 'Repository URL'))
        await idle(ctx)
        type(ctx, field(ctx, 'Repository URL'), 'git@git.trap.jp:traP/moved.git')
        await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
        click(ctx, button(ctx, 'Save'))
        await settle(ctx, () => expect(ctx.backend.requestsTo('updateRepository')).toHaveLength(1))
        const [req] = ctx.backend.requestsTo('updateRepository')
        expect(req).toMatchObject({ id: 'repo-gitea', url: 'git@git.trap.jp:traP/moved.git' })
        expect(req.auth).toBeUndefined()
        expect(req.name).toBeUndefined()
      })
    })
  })

  it('app build settings: the database controls follow the app when moving between apps', async () => {
    await scenario('runtime db across apps', '/apps/app-runtime/settings/build', async (ctx) => {
      // app-runtime uses MariaDB
      await waitForText(ctx, 'MariaDB')
      await idle(ctx)
      // broken-app is a runtime app without a database; navigate client-side so the form is reused
      const a = document.createElement('a')
      a.href = '/apps/app-error/settings/build'
      ctx.screen.append(a)
      click(ctx, a)
      await settle(ctx, () => expect(window.location.pathname).toBe('/apps/app-error/settings/build'))
      await waitForText(ctx, 'broken-app')
      await waitForNoText(ctx, 'MariaDB')
    })
  })

  describe('error boundaries', () => {
    // Settings: the key list is its own section, so its failure must not take down the page, header or nav.
    it('a failing section shows its own error and Retry, leaving the page and header usable', async () => {
      await scenario(
        'section error',
        '/settings',
        async (ctx) => {
          // the failure is shown in place of the section, with the section heading kept
          await waitForText(ctx, 'SSH Public Keys', 'Failed to load', 'keys backend is down')
          // everything around it stays rendered and usable
          expect(ctx.screen.textContent).not.toContain('An error has occurred')
          expect(ctx.screen.textContent).toContain('admin-user')
          expect(link(ctx, '/apps')).toBeTruthy()
          click(ctx, link(ctx, '/apps'))
          await waitForText(ctx, 'runtime-app')
        },
        {
          setup: (b) => b.failMethod('getUserKeys', 'keys backend is down'),
          expectedErrors: ['[section: SSH Public Keys]'],
        },
      )
    })

    it('Retry recovers the section once the backend is fixed', async () => {
      await scenario(
        'section retry',
        '/settings',
        async (ctx) => {
          await waitForText(ctx, 'Failed to load')
          const before = ctx.backend.requestsTo('getUserKeys').length

          // still down: Retry re-asks the backend and the section fails again, in place
          click(ctx, button(ctx, 'Retry'))
          await settle(ctx, () => expect(ctx.backend.requestsTo('getUserKeys').length).toBeGreaterThan(before))
          await waitForText(ctx, 'Failed to load')

          ctx.backend.recover('getUserKeys')
          click(ctx, button(ctx, 'Retry'))
          await waitForText(ctx, 'laptop')
          await waitForNoText(ctx, 'Failed to load')
        },
        {
          setup: (b) => b.failMethod('getUserKeys', 'keys backend is down'),
          expectedErrors: ['[section: SSH Public Keys]'],
        },
      )
    })

    it('app page: a failing Usage section does not take down the other sections or the tabs', async () => {
      await scenario(
        'usage error',
        '/apps/app-runtime',
        async (ctx) => {
          await waitForText(ctx, 'Usage', 'Failed to load', 'metrics are down')
          // siblings and the app navigation are intact
          await waitForText(ctx, 'Branch Resolution', 'feat: add feature A', 'previous log line')
          expect(ctx.screen.textContent).not.toContain('An error has occurred')
          expect(link(ctx, '/apps/app-runtime/builds')).toBeTruthy()

          ctx.backend.recover('getAvailableMetrics')
          click(ctx, button(ctx, 'Retry'))
          await waitForNoText(ctx, 'Failed to load')
        },
        {
          setup: (b) => b.failMethod('getAvailableMetrics', 'metrics are down'),
          expectedErrors: ['[section: Usage]'],
        },
      )
    })

    it('the failure of the header user info stays in the header', async () => {
      await scenario(
        'header error',
        '/settings',
        async (ctx) => {
          await waitForText(ctx, 'Retry')
          // the page itself is fine
          await waitForText(ctx, 'laptop')
          expect(ctx.screen.textContent).not.toContain('An error has occurred')
          ctx.backend.recover('getMe')
          click(ctx, buttonContaining(ctx, 'Retry'))
          await waitForText(ctx, 'admin-user')
          await waitForNoText(ctx, 'Retry')
        },
        { setup: (b) => b.failMethod('getMe', 'me is down'), expectedErrors: ['[section: User Info]'] },
      )
    })

    it('a page whose own data fails shows its error inside the page, then Retry recovers it', async () => {
      await scenario(
        'page error',
        '/repos/repo-github',
        async (ctx) => {
          await waitForText(ctx, 'An error has occurred', 'repo is down')
          // the header (and its nav) stay
          expect(link(ctx, '/apps')).toBeTruthy()
          ctx.backend.recover('getRepository')
          click(ctx, button(ctx, 'Retry'))
          await waitForText(ctx, 'github-repo', 'runtime-app')
          await waitForNoText(ctx, 'An error has occurred')
        },
        {
          expectErrorView: true,
          setup: (b) => b.failMethod('getRepository', 'repo is down'),
          expectedErrors: ['[page]'],
        },
      )
    })
  })
})
