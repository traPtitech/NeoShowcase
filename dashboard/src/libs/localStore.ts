import { createEffect, createSignal } from 'solid-js'
import { createStore } from 'solid-js/store'

// https://www.solidjs.com/examples/todos

export const createLocalSignal = <T>(name: string, init: T): ReturnType<typeof createSignal<T>> => {
  const localState = localStorage.getItem(name)
  const [state, setState] = createSignal<T>(localState ? JSON.parse(localState) : init)
  createEffect(() => localStorage.setItem(name, JSON.stringify(state())))
  return [state, setState]
}

export const createLocalStore = <T extends {}>(name: string, init: T): ReturnType<typeof createStore<T>> => {
  const localState = localStorage.getItem(name)
  const [state, setState] = createStore<T>(localState ? JSON.parse(localState) : init)
  createEffect(() => localStorage.setItem(name, JSON.stringify(state)))
  return [state, setState]
}
