import { createEffect, createSignal, type Signal } from 'solid-js'

/**
 * A signal persisted to `sessionStorage` / `localStorage` as JSON.
 * Storage access is wrapped because it can throw (private mode, quota, disabled storage).
 */
const createStorageSignal = <T>(storage: () => Storage, key: string, init: T): Signal<T> => {
  let initial = init
  try {
    const stored = storage().getItem(key)
    if (stored !== null) initial = JSON.parse(stored) as T
  } catch {
    // ignore broken / unavailable storage
  }
  // biome-ignore lint/complexity/noBannedTypes: mirrors createSignal's own overload (non-function values)
  const [state, setState] = createSignal<T>(initial as Exclude<T, Function>)
  createEffect(state, (value) => {
    try {
      storage().setItem(key, JSON.stringify(value))
    } catch {
      // ignore
    }
  })
  return [state, setState] as Signal<T>
}

export const createSessionSignal = <T>(key: string, init: T): Signal<T> =>
  createStorageSignal(() => sessionStorage, key, init)
export const createLocalSignal = <T>(key: string, init: T): Signal<T> =>
  createStorageSignal(() => localStorage, key, init)
