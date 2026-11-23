import { type Accessor, createContext, createMemo, type ParentComponent, useContext } from 'solid-js'
import type { Application, ApplicationGroup } from '/@/api/neoshowcase/protobuf/gateway_pb'
import {
  type CommitsMap,
  getApplicationGroup,
  getApplicationGroupApps,
  getRepositoryCommits,
  isOwner,
  revalidateApplicationGroup,
  revalidateApplicationGroupApps,
  revalidateMyApplicationGroups,
} from '/@/libs/api'
import { me } from '/@/libs/session'

/** Data shared by every page under `/groups/:id`. All accessors are async: read them under `<Loading>`. */
export type ApplicationGroupData = {
  group: Accessor<ApplicationGroup>
  apps: Accessor<Application[]>
  commits: Accessor<CommitsMap>
  hasPermission: Accessor<boolean>
  refetch: () => void
}

const ApplicationGroupContext = createContext<ApplicationGroupData>()

export const ApplicationGroupDataProvider: ParentComponent<{ id: string }> = (props) => {
  const group = createMemo(() => getApplicationGroup(props.id))
  const apps = createMemo(() => getApplicationGroupApps(props.id))
  const commits = createMemo(() => getRepositoryCommits(apps().map((a) => a.commit)))
  const hasPermission = createMemo(() => isOwner(me(), group()))
  const refetch = () => {
    revalidateApplicationGroup(props.id)
    revalidateApplicationGroupApps(props.id)
    revalidateMyApplicationGroups()
  }
  return (
    <ApplicationGroupContext value={{ group, apps, commits, hasPermission, refetch }}>
      {props.children}
    </ApplicationGroupContext>
  )
}

export const useApplicationGroupData = () => useContext(ApplicationGroupContext)
