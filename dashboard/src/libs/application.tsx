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
  // if app is not running or autoShutdown is enabled and container is missing, it's idle
  if (!app.running || (autoShutdownEnabled(app) && app.container === Application_ContainerState.MISSING)) {
    return ApplicationState.Idle
  }
  if (app.currentBuild === '') {
    // First build may still be running
    return ApplicationState.Idle
  }
  if (app.deployType === DeployType.RUNTIME) {
    switch (app.container) {
      case Application_ContainerState.MISSING:
      case Application_ContainerState.STARTING:
        return ApplicationState.Deploying
      case Application_ContainerState.RUNNING:
        return ApplicationState.Running
      case Application_ContainerState.RESTARTING:
      case Application_ContainerState.EXITED:
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
