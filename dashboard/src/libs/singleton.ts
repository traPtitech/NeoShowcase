import { createRoot, runWithOwner } from 'solid-js'

/**
 * Creates app-lifetime reactive state (signals, memos, timers) that is never disposed.
 *
 * In Solid 2 a `createRoot` is owned by whatever owner is current and is disposed with it. A module may first be
 * evaluated while some component is rendering (e.g. a lazily loaded page), so a bare `createRoot` at module scope
 * could be torn down together with that component. Running it with no owner detaches it for good.
 */
export const createSingleton = <T>(fn: () => T): T => runWithOwner(null, () => createRoot(fn)) as T
