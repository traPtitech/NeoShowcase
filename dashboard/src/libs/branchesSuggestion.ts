import { createMemo, createResource } from 'solid-js'
import { client } from '/@/libs/api'
import Fuse from 'fuse.js'

export const useBranchesSuggestion = (repoID: () => string, current: () => string): (() => string[]) => {
  const [refs] = createResource(
    () => repoID(),
    (id) => client.getRepositoryRefs({ repositoryId: id }),
  )

  const branches = createMemo(() => {
    if (!refs()) return
    const branches = refs()
      .refs.map((r) => r.refName)
      .filter((b) => !b.startsWith('refs/'))
    const normal = branches.filter((b) => !b.includes('/'))
    const long = branches.filter((b) => b.includes('/'))
    return [normal, long]
  })
  const branchesFuse = createMemo(() => {
    if (!branches()) return
    const [normal, long] = branches()
    return [new Fuse(normal), new Fuse(long)]
  })

  return createMemo(() => {
    const query = current()

    if (!branchesFuse()) return

    if (!query) return branches()[0].concat(branches()[1])

    const p0 = branchesFuse()[0]
      .search(query)
      .map((r) => r.item)
    const p1 = branchesFuse()[1]
      .search(query)
      .map((r) => r.item)
    return p0.concat(p1)
  })
}
