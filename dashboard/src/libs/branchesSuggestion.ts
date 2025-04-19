import { createMemo, createResource } from 'solid-js'
import toast from 'solid-toast'
import { client } from '/@/libs/api'

export const useBranchesSuggestion = (repoID: () => string, current: () => string): (() => string[]) => {
  const [refs] = createResource(
    () => repoID(),
    (id) => client.getRepositoryRefs({ repositoryId: id }),
  )

  const branches = createMemo(() => {
    if (refs.state === 'ready') {
      const branches = refs()
        .refs.map((r) => r.refName)
        .filter((b) => !b.startsWith('refs/'))
      const normal = branches?.filter((b) => !b.includes('/'))
      const long = branches?.filter((b) => b.includes('/'))
      return [normal, long]
    }
    return [[], []]
  })

  return createMemo(() => {
    const query = current()
    const normal = branches()[0]
    const long = branches()[1]
    if (!query) return normal.concat(long)

    const p0 = normal.filter((r) => r.toLowerCase().includes(query.toLowerCase()))
    const p1 = long.filter((r) => r.toLowerCase().includes(query.toLowerCase()))
    return p0.concat(p1)
  })
}

export const useBranches = (repoID: () => string): (() => string[]) => {
  const [refs] = createResource(
    () => repoID(),
    (id) =>
      client.getRepositoryRefs({ repositoryId: id }).catch((err) => {
        // ブランチ取得の失敗時は502エラーが返ってくるのでcatchで処理する
        console.trace(err)
        toast.error('ブランチの取得に失敗しました。リポジトリへのアクセス権を確認してください。')
        return { refs: [] }
      }),
  )

  return createMemo(() => {
    return (
      refs()
        ?.refs.map((r) => r.refName)
        .filter((b) => !b.startsWith('refs/')) ?? []
    )
  })
}
