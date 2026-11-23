import { Code } from '@connectrpc/connect'
import { describe, expect, it } from 'vitest'
import {
  button,
  buttonContaining,
  click,
  field,
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

  it('app general settings: renaming an app in a group succeeds', async () => {
    await scenario('rename grouped app', '/apps/app-static/settings', async (ctx) => {
      await waitForText(ctx, 'Branch')
      await idle(ctx)
      await settle(ctx, () => field(ctx, 'Application Name'))
      type(ctx, field(ctx, 'Application Name'), 'renamed-static')
      await settle(ctx, () => expect(button(ctx, 'Save').disabled).toBe(false))
      click(ctx, button(ctx, 'Save'))

      await settle(ctx, () => expect(ctx.backend.requestsTo('updateApplication')).toHaveLength(1))
      // owners of an app in a group belong to the group, so settings forms must not send them
      expect(ctx.backend.requestsTo('updateApplication')[0].ownerIds).toBeUndefined()
      await waitForText(ctx, 'renamed-static')
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
