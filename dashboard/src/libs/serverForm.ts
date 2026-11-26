import { type FormSchema, type FormStore, getInput, isDirty, type RequiredPath, reset } from '@formisch/solid/v2'
import { createEffect, untrack } from 'solid-js'
import type * as v from 'valibot'

type Path = readonly (string | number)[]
type Group = Record<string, unknown>

export const isDeepEqual = (a: unknown, b: unknown): boolean => {
  if (Object.is(a, b)) return true
  if (typeof a !== 'object' || typeof b !== 'object' || a === null || b === null) return false
  if (Array.isArray(a) !== Array.isArray(b)) return false
  const aKeys = Object.keys(a).filter((k) => (a as Group)[k] !== undefined)
  const bKeys = Object.keys(b).filter((k) => (b as Group)[k] !== undefined)
  if (aKeys.length !== bKeys.length) return false
  return aKeys.every((k) => isDeepEqual((a as Group)[k], (b as Group)[k]))
}

const at = (value: unknown, path: Path): Group =>
  (path.reduce<unknown>((v, key) => (v as Group | undefined)?.[key], value) ?? {}) as Group

/**
 * Keeps a settings form in sync with the server data it edits, without throwing away what the user is typing.
 *
 * The object at `root` is synced per top-level key ("group"): when the server value of a group changes,
 * an untouched group takes the new value and an edited group keeps the user's draft (with the new server
 * value as its clean baseline). Revalidations that bring back the same values change nothing.
 * When `key` (the resource id) changes, the whole form is re-initialized.
 *
 * Capture `snapshot()` when the form is submitted and pass it to `markSaved` after a successful save: the
 * submitted input becomes the clean state, so input typed while the request was in flight stays dirty.
 */
export const syncFormWithServer = <TSchema extends FormSchema>(
  form: FormStore<TSchema>,
  options: {
    key: () => string
    input: () => v.InferInput<TSchema>
    root?: Path
  },
) => {
  const root = options.root ?? []
  // Formisch's path types cannot follow the dynamic group paths used here
  const path = (key: string) => [...root, key] as unknown as RequiredPath
  const anyForm = form as unknown as FormStore<FormSchema>

  let currentKey = untrack(options.key)
  /** The server value each group was last synced with. */
  let base: Group = { ...at(untrack(options.input), root) }

  /** Makes `value` the clean state of the group, keeping its current input. */
  const rebase = (key: string, value: unknown) => {
    const current = getInput(anyForm, { path: path(key) } as never)
    if (isDeepEqual(current, value)) reset(anyForm, { path: path(key), initialInput: value } as never)
    else
      reset(anyForm, {
        path: path(key),
        initialInput: value,
        keepInput: true,
        keepTouched: true,
        keepEdited: true,
        keepErrors: true,
      } as never)
  }

  /** Applies a new server value: re-initializes on a new resource, otherwise syncs the changed groups. */
  const apply = (key: string, input: v.InferInput<TSchema>) => {
    if (key !== currentKey) {
      currentKey = key
      base = { ...at(input, root) }
      reset(form, { initialInput: input } as never)
      return
    }
    const next = at(input, root)
    for (const group of Object.keys(next)) {
      if (isDeepEqual(base[group], next[group])) continue
      base[group] = next[group]
      if (isDirty(anyForm, { path: path(group) } as never)) rebase(group, next[group])
      else reset(anyForm, { path: path(group), initialInput: next[group] } as never)
    }
  }

  createEffect(
    () => ({ key: options.key(), input: options.input() }),
    // apply reads the form's own state, which must not become (or be reported as) a dependency
    ({ key, input }) => untrack(() => apply(key, input)),
    { defer: true },
  )

  return {
    /** The current input, to be captured at submit time and passed to `markSaved`. */
    snapshot: () => ({ key: currentKey, input: getInput(form) as v.InferInput<TSchema> }),
    markSaved: (submitted: { key: string; input: v.InferInput<TSchema> }) => {
      // the form has moved on to another resource meanwhile
      if (submitted.key !== currentKey) return
      const saved = at(submitted.input, root)
      for (const group of Object.keys(saved)) {
        base[group] = saved[group]
        rebase(group, saved[group])
      }
    },
    /** Drops every draft and goes back to the latest server values. */
    discard: () => {
      for (const group of Object.keys(base)) reset(anyForm, { path: path(group), initialInput: base[group] } as never)
    },
  }
}
