import { type Accessor, createContext, createMemo, type ParentComponent, useContext } from 'solid-js'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import {
  type CommitsMap,
  getRepository,
  getRepositoryApps,
  getRepositoryCommits,
  isOwner,
  revalidateRepository,
  revalidateRepositoryApps,
} from '/@/libs/api'
import { me } from '/@/libs/session'

/** Data shared by every page under `/repos/:id`. All accessors are async: read them under `<Loading>`. */
export type RepositoryData = {
  repo: Accessor<Repository>
  apps: Accessor<Application[]>
  commits: Accessor<CommitsMap>
  hasPermission: Accessor<boolean>
  refetch: () => void
}

const RepositoryContext = createContext<RepositoryData>()

export const RepositoryDataProvider: ParentComponent<{ id: string }> = (props) => {
  const repo = createMemo(() => getRepository(props.id))
  const apps = createMemo(() => getRepositoryApps(props.id))
  const commits = createMemo(() => getRepositoryCommits(apps().map((a) => a.commit)))
  const hasPermission = createMemo(() => isOwner(me(), repo()))
  const refetch = () => {
    revalidateRepository(props.id)
    revalidateRepositoryApps(props.id)
  }
  return <RepositoryContext value={{ repo, apps, commits, hasPermission, refetch }}>{props.children}</RepositoryContext>
}

export const useRepositoryData = () => useContext(RepositoryContext)
