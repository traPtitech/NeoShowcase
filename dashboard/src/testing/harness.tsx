import { captureArtifact } from '@solidjs/diagnostics'
import type { DiagnosticEvent } from 'solid-js'
import { afterEach, expect, vi } from 'vitest'
import * as backend from '/@/testing/mockBackend'

/**
 * Helpers to render the real <App/> against the in-memory backend (src/testing/mockBackend).
 * Every scenario gets fresh modules (router, query cache, session singletons) and fails on
 * Solid diagnostics, console errors/warnings and error views.
 */

export type Ctx = {
  /** The in-memory backend (fixtures + recorded requests). */
  backend: typeof backend
  /** `document.body`: dialogs, tooltips and toasts are portalled out of the app root. */
  screen: HTMLElement
  flush: () => void
}

let cleanup: (() => void) | undefined
afterEach(() => {
  cleanup?.()
  cleanup = undefined
  document.body.innerHTML = ''
})

const mount = async (url: string, setup?: (backend: Ctx['backend']) => void): Promise<Ctx> => {
  vi.resetModules()
  window.history.replaceState(null, '', url)
  backend.resetMockBackend()
  setup?.(backend)
  sessionStorage.clear()
  localStorage.clear()
  const { render } = await import('@solidjs/web')
  const { flush } = await import('solid-js')
  const { default: App } = await import('/@/App')

  const root = document.createElement('div')
  document.body.append(root)
  cleanup = render(() => <App />, root)
  return { backend, screen: document.body, flush }
}

// Kobalte 2 (alpha) reads some props once while constructing its components (e.g. `others.placement` in
// Tooltip, `props.value` in TextField). Solid reports these as STRICT_READ_UNTRACKED inside the library's own
// components, which the app cannot fix. Everything else is a finding.
const kobalteReadOnConstruct =
  /^\[STRICT_READ_UNTRACKED\].*(in <(K[A-Z][\w.]*|Tooltip|TextFieldInputBase|ComboboxBase)>|at node_modules\/.*@kobalte)/s

const formatEvent = (e: DiagnosticEvent) => `[${e.code}] ${e.message}\n  in ${e.ownerPath?.join(' > ')}\n  at ${e.loc}`

/** Runs a scenario, then fails on Solid diagnostics, console errors/warnings and error views. */
export const scenario = async (
  name: string,
  url: string,
  body: (ctx: Ctx) => Promise<void>,
  {
    expectErrorView = false,
    setup,
    expectedErrors = [],
  }: {
    expectErrorView?: boolean
    /** Runs before the app renders, e.g. to make a backend method fail from the first request. */
    setup?: (backend: Ctx['backend']) => void
    /** Substrings of console errors the scenario is expected to log (boundaries log what they catch). */
    expectedErrors?: string[]
  } = {},
) => {
  const logs: string[] = []
  const spies = (['error', 'warn'] as const).map((level) =>
    vi.spyOn(console, level).mockImplementation((...args) => {
      const text = args.map(String).join(' ')
      // Solid diagnostics are checked through the artifact below
      if (/^\[[A-Z_]+\]/.test(text)) return
      if (level === 'error' && expectedErrors.some((part) => text.includes(part))) return
      logs.push(`console.${level}: ${text}`)
    }),
  )
  try {
    const { artifact } = await captureArtifact(
      async () => {
        const ctx = await mount(url, setup)
        await body(ctx)
        if (!expectErrorView) expect(ctx.screen.textContent).not.toContain('An error has occurred')
      },
      { scenario: name, attribution: false },
    )
    expect(logs).toEqual([])
    const findings = artifact.diagnostics.map(formatEvent).filter((text) => !kobalteReadOnConstruct.test(text))
    expect(findings).toEqual([])
  } finally {
    for (const s of spies) s.mockRestore()
  }
}

// --- DOM helpers (no @testing-library: plain DOM + Solid's flush)

const pageText = (ctx: Ctx) => {
  // inputs keep their value outside of textContent
  const fieldValues = [...ctx.screen.querySelectorAll('input, textarea')].map((el) => (el as HTMLInputElement).value)
  return [ctx.screen.textContent, ...fieldValues].join('\n')
}

/** Polls (flushing Solid's queue) until `check` stops throwing. */
// Generous default: the first visit to a lazy route transforms its modules on demand, which is slow on a cold run.
export const settle = (ctx: Ctx, check: () => void, timeout = 10_000) =>
  vi.waitFor(
    () => {
      ctx.flush()
      check()
    },
    { timeout, interval: 20 },
  )

export const waitForText = (ctx: Ctx, ...texts: string[]) =>
  settle(ctx, () => {
    const content = pageText(ctx)
    for (const t of texts) expect(content).toContain(t)
  })

export const waitForNoText = (ctx: Ctx, text: string) => settle(ctx, () => expect(pageText(ctx)).not.toContain(text))

/** Lets trailing async work (polling kick-off, commit lookups, effects) settle. */
export const idle = async (ctx: Ctx, ms = 50) => {
  await new Promise((r) => setTimeout(r, ms))
  ctx.flush()
}

/** The first element matching `selector` whose own text is exactly `text`. */
export const find = (ctx: Ctx, selector: string, text: string): HTMLElement => {
  const hit = [...ctx.screen.querySelectorAll<HTMLElement>(selector)].find((el) => el.textContent?.trim() === text)
  if (!hit) throw new Error(`no <${selector}> with text "${text}"`)
  return hit
}
export const link = (ctx: Ctx, href: string) => {
  const el = ctx.screen.querySelector<HTMLElement>(`a[href="${href}"]`)
  if (!el) throw new Error(`no link to ${href}`)
  return el
}
export const button = (ctx: Ctx, text: string) => find(ctx, 'button', text) as HTMLButtonElement

/** The first button whose text contains `text`. */
export const buttonContaining = (ctx: Ctx, text: string) => {
  const hit = [...ctx.screen.querySelectorAll<HTMLElement>('button')].find((b) => b.textContent?.includes(text))
  if (!hit) throw new Error(`no <button> containing "${text}"`)
  return hit
}

export const click = (ctx: Ctx, el: HTMLElement) => {
  el.click()
  ctx.flush()
}

/** The input/textarea of the labelled field. */
export const field = (ctx: Ctx, label: string): HTMLInputElement | HTMLTextAreaElement => {
  const lbl = [...ctx.screen.querySelectorAll('label')].find((l) => l.textContent?.replaceAll('*', '').trim() === label)
  if (!lbl) throw new Error(`no label "${label}" in:\n${pageText(ctx).slice(0, 1500)}`)
  const id = lbl.getAttribute('for')
  const el = id ? document.getElementById(id) : lbl.parentElement?.querySelector('input, textarea')
  if (!el) throw new Error(`no field for label "${label}"`)
  return el as HTMLInputElement
}

export const type = (ctx: Ctx, el: HTMLInputElement | HTMLTextAreaElement, value: string) => {
  el.focus()
  el.value = value
  el.dispatchEvent(new Event('input', { bubbles: true }))
  ctx.flush()
}

/** Moves to `url` like the browser's Back/Forward does (the router follows `popstate`). */
export const go = (ctx: Ctx, url: string) => {
  window.history.pushState(null, '', url)
  window.dispatchEvent(new PopStateEvent('popstate'))
  ctx.flush()
}

export const dialog = (ctx: Ctx) => ctx.screen.querySelector<HTMLElement>('[role=dialog]')
