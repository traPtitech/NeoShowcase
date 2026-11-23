import { timestampDate } from '@bufbuild/protobuf/wkt'
import { type Accessor, createContext, createMemo, type ParentComponent, useContext } from 'solid-js'
import type { Application, ApplicationGroup, Build, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import {
  type CommitsMap,
  getApplication,
  getApplicationGroup,
  getBuilds,
  getRepository,
  getRepositoryCommits,
  isOwner,
  revalidateApplication,
  revalidateApplicationGroup,
  revalidateBuilds,
} from '/@/libs/api'
import { me } from '/@/libs/session'

/** Data shared by every page under `/apps/:id`. All accessors are async: read them under `<Loading>`. */
export type ApplicationData = {
  app: Accessor<Application>
  repo: Accessor<Repository>
  /** The group a user created that the application is in, or undefined for a standalone application. */
  group: Accessor<ApplicationGroup | undefined>
  /** Builds sorted by queued time, newest first. */
  builds: Accessor<Build[]>
  commits: Accessor<CommitsMap>
  hasPermission: Accessor<boolean>
  /** Re-fetch the application, its group and its builds (quiet: content stays on screen). */
  refetch: () => void
}

const byQueuedAtDesc = (a: Build, b: Build) =>
  (b.queuedAt ? timestampDate(b.queuedAt).getTime() : 0) - (a.queuedAt ? timestampDate(a.queuedAt).getTime() : 0)

const ApplicationContext = createContext<ApplicationData>()

export const ApplicationDataProvider: ParentComponent<{ id: string }> = (props) => {
  const app = createMemo(() => getApplication(props.id))
  const repo = createMemo(() => getRepository(app().repositoryId))
  const group = createMemo(() => {
    const groupId = app().groupId
    return groupId ? getApplicationGroup(groupId) : undefined
  })
  const builds = createMemo(async () => (await getBuilds(props.id)).toSorted(byQueuedAtDesc))
  const commits = createMemo(() => getRepositoryCommits([app().commit, ...builds().map((b) => b.commit)]))
  const hasPermission = createMemo(() => isOwner(me(), app()))
  const refetch = () => {
    revalidateApplication(props.id)
    const groupId = app().groupId
    if (groupId) revalidateApplicationGroup(groupId)
    revalidateBuilds(props.id)
  }
  return (
    <ApplicationContext value={{ app, repo, group, builds, commits, hasPermission, refetch }}>
      {props.children}
    </ApplicationContext>
  )
}

export const useApplicationData = () => useContext(ApplicationContext)
