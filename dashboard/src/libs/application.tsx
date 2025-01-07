import { timestampDate } from '@bufbuild/protobuf/wkt'
import { AiFillGithub } from 'solid-icons/ai'
import { RiDevelopmentGitRepositoryLine } from 'solid-icons/ri'
import { SiGitea } from 'solid-icons/si'
import type { JSXElement } from 'solid-js'
import {
  type Application,
  Application_ContainerState,
  BuildStatus,
  type CreateWebsiteRequest,
  DeployType,
  PortPublicationProtocol,
  type Repository,
  type Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'

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

const autoShutdownEnabled = (app: Application): boolean => {
  switch (app.config?.buildConfig?.case) {
    case 'runtimeBuildpack':
    case 'runtimeCmd':
    case 'runtimeDockerfile':
      return app.config.buildConfig.value.runtimeConfig?.autoShutdown?.enabled ?? false
  }
  return false
}

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

export const originToIcon = (origin: RepositoryOrigin, size = 20): JSXElement => {
  switch (origin) {
    case 'GitHub':
      return <AiFillGithub size={size} class="text-text-black" />
    case 'Gitea':
      return <SiGitea size={size} class="text-text-black" />
    case 'Others':
      return <RiDevelopmentGitRepositoryLine size={size} class="text-text-black" />
  }
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

const newestAppDate = (apps: Application[]): number =>
  Math.max(0, ...apps.map((a) => (a.updatedAt ? timestampDate(a.updatedAt).getTime() : 0)))
const compareRepoWithApp =
  (sort: 'asc' | 'desc') =>
  (a: RepoWithApp, b: RepoWithApp): number => {
    // Sort by apps updated at
    if (a.apps.length > 0 && b.apps.length > 0) {
      if (sort === 'asc') {
        return newestAppDate(a.apps) - newestAppDate(b.apps)
      }
      return newestAppDate(b.apps) - newestAppDate(a.apps)
    }
    // Bring up repositories with 1 or more apps at top
    if ((a.apps.length > 0 && b.apps.length === 0) || (a.apps.length === 0 && b.apps.length > 0)) {
      return b.apps.length - a.apps.length
    }
    // Fallback to sort by repository id
    return a.repo.id.localeCompare(b.repo.id)
  }

export interface RepoWithApp {
  repo: Repository
  apps: Application[]
}

export const useApplicationsFilter = (
  repos: Repository[],
  apps: Application[],
  statuses: ApplicationState[],
  origins: RepositoryOrigin[],
  includeNoApp: boolean,
  sort: 'asc' | 'desc',
): RepoWithApp[] => {
  const filteredReposByOrigin = () => {
    return repos.filter((r) => origins.includes(repositoryURLToOrigin(r.url))) ?? []
  }
  const filteredApps = () => {
    return apps.filter((a) => statuses.includes(applicationState(a))) ?? []
  }

  const appsMap = {} as Record<string, Application[]>
  for (const app of filteredApps()) {
    if (!appsMap[app.repositoryId]) appsMap[app.repositoryId] = []
    appsMap[app.repositoryId].push(app)
  }
  const res = filteredReposByOrigin().reduce<RepoWithApp[]>((acc, repo) => {
    if (!includeNoApp && !appsMap[repo.id]) return acc
    acc.push({ repo, apps: appsMap[repo.id] || [] })
    return acc
  }, [])
  res.sort(compareRepoWithApp(sort))
  return res
}
