import { createStore } from 'solid-js/store'

// ストアを作成するためのユーティリティ関数
export const storify = <T extends {},>(class_: T) => {
  const [store, _] = createStore(class_)
  return store
}
