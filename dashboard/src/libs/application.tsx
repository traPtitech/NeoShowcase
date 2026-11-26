import { timestampDate } from '@bufbuild/protobuf/wkt'
import {
  type Application,
  Application_ContainerState,
  BuildStatus,
  type CreateWebsiteRequest,
  DeployType,
  PortPublicationProtocol,
  type Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import type { SelectOption } from '/@/components/templates/Select'

export const buildStatusStr: Record<BuildStatus, string> = {
  [BuildStatus.QUEUED]: 'Queued',
  [BuildStatus.BUILDING]: 'Building',
  [BuildStatus.SUCCEEDED]: 'Succeeded',
  [BuildStatus.FAILED]: 'Failed',
  [BuildStatus.CANCELLED]: 'Cancelled',
  [BuildStatus.SKIPPED]: 'Skipped',
}

export enum ApplicationState {
  Idle = 'Idle',
  Deploying = 'Deploying',
  Running = 'Running',
  Sleeping = 'Sleeping',
  Serving = 'Serving',
  Error = 'Error',
}

const autoShutdownEnabled = (app: Application): boolean => app.config?.autoShutdown?.enabled ?? false

export const deploymentState = (app: Application): ApplicationState => {
  // App is not running
  if (!app.running) {
    return ApplicationState.Idle
  }
  if (app.currentBuild === '') {
    // First build may still be running
    return ApplicationState.Idle
  }
  if (app.deployType === DeployType.RUNTIME) {
    switch (app.container) {
      case Application_ContainerState.MISSING:
        // Has auto shutdown enabled, and the container is missing - app is sleeping, and will start on HTTP access
        if (autoShutdownEnabled(app)) {
          return ApplicationState.Sleeping
        }
        return ApplicationState.Deploying
      case Application_ContainerState.STARTING:
        return ApplicationState.Deploying
      case Application_ContainerState.RUNNING:
        return ApplicationState.Running
      case Application_ContainerState.RESTARTING:
      case Application_ContainerState.EXITED:
        return ApplicationState.Idle
      case Application_ContainerState.ERRORED:
      case Application_ContainerState.UNKNOWN:
        return ApplicationState.Error
    }
  }
  return ApplicationState.Serving
}

export const errorCommit = '0'.repeat(40)

export const applicationState = (app: Application): ApplicationState => {
  if (!app.running) {
    return ApplicationState.Idle
  }
  if (app.commit === errorCommit) {
    return ApplicationState.Error
  }
  switch (app.latestBuildStatus) {
    case BuildStatus.QUEUED:
      return ApplicationState.Deploying
    case BuildStatus.BUILDING:
      return ApplicationState.Deploying
    case BuildStatus.SUCCEEDED:
      return deploymentState(app)
    case BuildStatus.FAILED:
      return ApplicationState.Error
    case BuildStatus.CANCELLED:
      return deploymentState(app)
    case BuildStatus.SKIPPED:
      return deploymentState(app)
    case undefined:
      return ApplicationState.Error
  }
}

export type RepositoryOrigin = 'GitHub' | 'Gitea' | 'Others'

export const repositoryURLToOrigin = (url: string): RepositoryOrigin => {
  const normalizedURL = url.toLowerCase()
  if (normalizedURL.includes('github')) return 'GitHub'
  if (normalizedURL.includes('gitea')) return 'Gitea'
  if (normalizedURL.includes('git.trap.jp')) return 'Gitea'
  return 'Others'
}

const originIconClass: Record<RepositoryOrigin, string> = {
  GitHub: 'i-mdi:github',
  Gitea: 'i-simple-icons:gitea',
  Others: 'i-mdi:source-repository',
}

export const originToIcon = (origin: RepositoryOrigin, size = 20) => (
  <div class={['text-text-black shrink-0', originIconClass[origin]]} style={{ 'font-size': `${size}px` }} />
)

export const allStatuses: SelectOption<ApplicationState>[] = [
  { label: 'Idle', value: ApplicationState.Idle },
  { label: 'Deploying', value: ApplicationState.Deploying },
  { label: 'Running', value: ApplicationState.Running },
  { label: 'Sleeping', value: ApplicationState.Sleeping },
  { label: 'Serving', value: ApplicationState.Serving },
  { label: 'Error', value: ApplicationState.Error },
]

export const allOrigins: SelectOption<RepositoryOrigin>[] = [
  { label: 'GitHub', value: 'GitHub' },
  { label: 'Gitea', value: 'Gitea' },
  { label: 'Others', value: 'Others' },
]

export const sortItems: { [k in 'desc' | 'asc']: SelectOption<k> } = {
  desc: { value: 'desc', label: 'Newest' },
  asc: { value: 'asc', label: 'Oldest' },
}

export const getWebsiteURL = (website: Website | CreateWebsiteRequest): string => {
  const scheme = website.https ? 'https' : 'http'
  return `${scheme}://${website.fqdn}${website.pathPrefix}`
}

export const websiteWarnings = (subdomain: string | undefined, isHTTPS: boolean | undefined): string[] => {
  const warnings = []
  if (subdomain?.includes('_')) {
    warnings.push('アンダースコア「_」を含むホスト名は非推奨です。将来非対応になる可能性があります。')
  }
  const labels = subdomain?.split('.')
  if (isHTTPS && labels && labels.length >= 2) {
    warnings.push('このホスト名では専用の証明書が取得されます。可能な限りホストのラベル数は少なくしてください。')
  }
  return warnings
}

export const extractRepositoryNameFromURL = (url: string): string => {
  const segments = url.split('/')
  const lastSegment = segments.pop() || segments.pop() // 末尾のスラッシュを除去
  return lastSegment?.replace(/\.git$/, '') ?? ''
}

export const portPublicationProtocolMap: Record<PortPublicationProtocol, string> = {
  [PortPublicationProtocol.TCP]: 'TCP',
  [PortPublicationProtocol.UDP]: 'UDP',
}

const appDate = (app: Application): number => (app.updatedAt ? timestampDate(app.updatedAt).getTime() : 0)
const newestAppDate = (apps: Application[]): number => Math.max(0, ...apps.map(appDate))

/**
 * An entry of the app list: the apps of a group a user created, or standalone apps that are next to each other
 * in the list (`groupId` undefined). `apps` is never empty.
 */
export interface AppListEntry {
  groupId?: string
  apps: Application[]
}

/**
 * Lists apps as entries: the apps of the same group share one entry, and so do standalone apps that are next to
 * each other. Groups and standalone apps are ordered by the newest update of their apps, and the apps in a group
 * by their update, both by `sort`.
 */
export const groupApplications = (apps: Application[], sort: 'asc' | 'desc'): AppListEntry[] => {
  const direction = sort === 'asc' ? 1 : -1
  const items: AppListEntry[] = []
  const groups = new Map<string, AppListEntry>()
  for (const app of apps) {
    if (app.groupId === undefined) {
      items.push({ apps: [app] })
      continue
    }
    const group = groups.get(app.groupId)
    if (group) {
      group.apps.push(app)
    } else {
      const entry = { groupId: app.groupId, apps: [app] }
      groups.set(app.groupId, entry)
      items.push(entry)
    }
  }
  for (const item of items) {
    item.apps.sort((a, b) => direction * (appDate(a) - appDate(b)) || a.id.localeCompare(b.id))
  }
  items.sort(
    (a, b) => direction * (newestAppDate(a.apps) - newestAppDate(b.apps)) || a.apps[0].id.localeCompare(b.apps[0].id),
  )
  return items.reduce<AppListEntry[]>((entries, item) => {
    const last = entries.at(-1)
    if (item.groupId === undefined && last !== undefined && last.groupId === undefined) {
      last.apps.push(...item.apps)
    } else {
      entries.push(item)
    }
    return entries
  }, [])
}
