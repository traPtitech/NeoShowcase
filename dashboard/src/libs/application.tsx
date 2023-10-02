import {
  Application,
  Application_ContainerState,
  BuildStatus,
  DeployType,
  PortPublicationProtocol,
  Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Provider } from '/@/components/RepositoryRow'
import { vars } from '/@/theme'
import { AiFillGithub, AiFillGitlab } from 'solid-icons/ai'
import { SiGitea } from 'solid-icons/si'
import { JSXElement } from 'solid-js'
import { BuildConfigMethod } from '../components/BuildConfigs'

export const buildTypeStr: Record<BuildConfigMethod, string> = {
  runtimeBuildpack: 'Runtime (Buildpack)',
  runtimeCmd: 'Runtime (command)',
  runtimeDockerfile: 'Runtime (Dockerfile)',
  staticBuildpack: 'Static (Buildpack)',
  staticCmd: 'Static (command)',
  staticDockerfile: 'Static (Dockerfile)',
}

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
  Static = 'Static',
  Error = 'Error',
}

const useDeployState = (app: Application): ApplicationState => {
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
  } else {
    return ApplicationState.Static
  }
}

const errorCommit = '0'.repeat(40)

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
      return useDeployState(app)
    case BuildStatus.FAILED:
      return ApplicationState.Error
    case BuildStatus.CANCELLED:
      return useDeployState(app)
    case BuildStatus.SKIPPED:
      return useDeployState(app)
  }
}

export const repositoryURLToProvider = (url: string): Provider => {
  const normalizedURL = url.toLowerCase()
  if (normalizedURL.includes('github')) return 'GitHub'
  if (normalizedURL.includes('gitlab')) return 'GitLab'
  if (normalizedURL.includes('gitea')) return 'Gitea'
  if (normalizedURL.includes('git.trap.jp')) return 'Gitea'
  return 'GitHub' // fallback?
}

export const providerToIcon = (provider: Provider, size = 20): JSXElement => {
  switch (provider) {
    case 'GitHub':
      return <AiFillGithub size={size} color={vars.text.black1} />
    case 'GitLab':
      return <AiFillGitlab size={size} color="#FC6D26" />
    case 'Gitea':
      return <SiGitea size={size} color={vars.text.black1} />
  }
}

export const getWebsiteURL = (website: Website): string => {
  const scheme = website.https ? 'https' : 'http'
  return `${scheme}://${website.fqdn}${website.pathPrefix}`
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
