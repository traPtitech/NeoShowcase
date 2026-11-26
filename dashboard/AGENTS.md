# dashboard — conventions

SolidJS **2.0** SPA (client-only). Solid 1.x habits are wrong here; when in doubt read
`node_modules/solid-js/CHEATSHEET.md` (especially the "What changed from 1.x" list at the bottom).

## Stack

| Concern | Library | Import from |
| --- | --- | --- |
| UI runtime | solid-js 2, @solidjs/web | `solid-js`, `@solidjs/web` (`JSX`, `Portal`, `dynamic`, `render`) |
| Routing | @solidjs/router 2 (`createRouter`, config objects) | `/@/router` (instance), `@solidjs/router` (hooks, `query`) |
| Headless UI | @kobalte/core 2 (alpha) | `@kobalte/core/<component>` (subpath only) |
| Forms | Formisch | `@formisch/solid/v2` (**the `/v2` entry**, not the root) |
| Toast | Kobalte toast, wrapped | `/@/libs/toast` only — never `@kobalte/core/toast` in features |
| Tooltip | Kobalte tooltip, wrapped | `/@/components/UI/ToolTip` |
| Icons | UnoCSS preset-icons | `<div class="i-material-symbols:xxx" />`, `i-mdi:*`, `i-simple-icons:*` |
| Charts | chart.js (no wrapper lib) | `chart.js` |
| API | Connect-RPC client | `/@/libs/api` |

Removed and must not come back: `solid-toast`, `solid-tippy`/`tippy.js`, `solid-icons`, `solid-chartjs`,
`@modular-forms/solid`, `async-lock`, `vite-plugin-solid(-svg)`, `solid-js/web`, `solid-js/store`, `/@/libs/clsx`.
Biome's `noRestrictedImports` enforces this.

## Solid 2 rules that bite

- `createEffect(compute, apply)` — two functions. Writes go in `apply`, event handlers, or `onSettled`.
  Writing a signal in a component body / memo / effect compute throws in dev.
- An effect whose compute reads query-backed data (directly or via props, e.g. `() => props.appID` from `app().id`)
  re-runs whenever a revalidation settles, **even if the value is equal**. When the effect resets state or (re)starts
  work, key it on a memo of the primitive: `const id = createMemo(() => props.appID); createEffect(id, ...)`.
- Async work started for a resource (`await` in a handler or effect) must check it is still for the current visit
  before writing signals: route components are reused across param changes, so give each run an `AbortController`
  that the effect cleanup aborts, rather than comparing ids.
- `onMount` → `onSettled(() => { ...; return cleanup })`. Component-level timers/listeners live there.
- `createResource` → `createMemo(() => fetch...())` (async memo). Reading it before it resolves suspends the
  nearest `<Loading>`. `.latest`/`.loading`/`.state` do not exist: use `<Loading>` for first load,
  `isPending(() => x())` for "a change is in flight", `refresh(x)` to re-run.
- `Suspense` → `Loading`, `ErrorBoundary` → `Errored` (`fallback={(err, reset) => ...}`, `err` is an accessor).
- `useTransition`/`startTransition` are gone. Navigation is already a transition; show staleness with
  `isPending`.
- `splitProps` → `omit(props, 'a', 'b')`; `mergeProps` → `merge(...)` (`undefined` overrides!).
- `<For>` default (keyed): `(item, i) => ...` where `item` is a value, `i` an accessor.
  That keys rows by object identity, and every refetch returns new objects, so rows remount (focus, open dialogs
  and inputs inside them are lost). Lists of API entities: `<For each={xs()} keyed={(x) => x.id}>` — then `item`
  is an accessor (`x()`), updated in place when a row with the same id gets new data.
  `<For keyed={false}>`: `item` is an accessor, `i` a number. `<Show>` function children get an accessor.
  Never read accessors at the top of these callbacks — read them in JSX.
- `Context.Provider` → `<Context value={...}>`. Default-less `useContext` throws when missing — no wrapper hooks.
- Classes: `class={['base', props.class, { active: on() }]}` — no string concatenation, no `clsx`.
  `classList` is a compile error.
- Attributes are HTML-cased: `tabindex`, `readonly`, `srcset`, `for`. Events stay camelCase (`onClick`);
  lowercase `onclick` is an *attribute* in Solid 2.
- `use:` directives are gone → `ref={directive(opts)}`. Refs are unowned (no `onCleanup` inside a ref).
- `<Dynamic>` is deprecated → `const Tag = dynamic(() => tag)` hoisted out of JSX.
- Don't destructure props. Pass values, not accessors (`value={x()}`).

## Data

- Session-wide data: `/@/libs/session` → `me`, `systemInfo`, `availableMetrics`, `users`, `usersFromIds`.
  These are async memos (call them: `me().admin`).
- App-lifetime reactive state at module scope: wrap it in `createSingleton(() => ...)` (`/@/libs/singleton`),
  never a bare `createRoot` (it would be owned by, and disposed with, whatever component first imports the module).
- Cached API reads: `query` wrappers in `/@/libs/api` (`getApplication`, `getBuilds`, …) plus
  `revalidateXxx(id)`. Consume with `createMemo(() => getApplication(id))`.
- Route-scoped data: the `/apps/:id` layout provides `useApplicationData()` and `/repos/:id` provides
  `useRepositoryData()` (`/@/features/*/data.tsx`). Child pages read from context — they do not refetch.
- Mutations: call `client.xxx()` in an event handler, `toast.success(...)`, then `refetch()`/`revalidateXxx()`.
  Errors: `handleAPIError(e, 'メッセージ')`.
- Polling: `onSettled(() => { const t = setInterval(refetch, 10000); return () => clearInterval(t) })`.

## Routing

- Routes are config objects in `src/router.tsx`. Pages live in `src/pages/` (file names mirror the URL).
- Links are plain `<a href="/apps/...">`. The router intercepts them and sets `aria-current="page"` (exact)
  and `data-active` (prefix) — style active tabs with `aria-[current=page]:` / `data-[active]:` variants.
- Programmatic navigation: `useNavigate()`.

## Forms (Formisch)

- `const form = createForm({ schema, initialInput })` per form component, `<Form of={form} onSubmit={...}>`.
- Fields: `<Field of={form} path={['form', 'name']}>{(field) => <TextField {...field.props} value={field.input} error={field.errors?.[0]} />}</Field>`.
- Settings forms editing server data: `syncFormWithServer(form, { key, input, root })` (`/@/libs/serverForm`), not an
  effect that calls `reset` — every revalidation yields fresh objects and would wipe unsaved input. Capture
  `sync.snapshot()` at submit and call `sync.markSaved(snapshot)` after the request succeeds.
- A settings panel's initial input holds only the fields it edits (`getInitialValueOfUpdateAppForm(app, ['name'])`),
  so it neither validates nor sends other panels' settings (update requests treat absent fields as unchanged).
- State: `form.isSubmitting`, `form.isDirty`, `form.isValid` (reactive getters).

## Kobalte 2 (alpha) caveats

- Kobalte reads many root props **untracked on construction**. Never pass a value that may depend on a pending
  async memo (e.g. `hasPermission()`, `me()`) straight into a Kobalte root prop: resolve it first in a tracked
  scope (`<Show when={...}>` / a memo read in JSX) — see `components/UI/ToolTip.tsx`. Otherwise Solid throws
  `PENDING_ASYNC_UNTRACKED_READ` and the page falls into the error view.
- Its `STRICT_READ_UNTRACKED` dev warnings are known noise; the test harness filters only those.

## Tests

- `src/testing/mockBackend.ts` is an in-memory Connect transport with fixtures; `src/testing/harness.tsx` renders
  the real `App` at a URL and fails on console errors, error views and Solid diagnostics.
- `src/app.smoke.test.tsx` renders every route; `src/app.interaction.test.tsx` covers key flows.
  Add a route → add a smoke case. Fix a runtime bug → add a test that reproduces it.
- The harness calls `vi.resetModules()` per scenario: import app modules a scenario drives directly (e.g.
  `revalidateApplication`) with `await import(...)` after mount — a top-level import targets a stale instance.
- Race scenarios: `backend.holdNext(method)` delays the next unary call, `backend.holdMethod(method)` every call,
  until the returned `release()`.
- `backend.dropOutputStream(n)` makes the next `n` container log streams break after their lines.

## Commands

`pnpm typecheck` · `pnpm lint` · `pnpm test --run` · `pnpm build`
