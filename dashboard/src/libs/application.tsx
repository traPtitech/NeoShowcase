import { Provider } from '/@/components/RepositoryRow'
import {
  Application,
  ApplicationConfig,
  Build_BuildStatus,
  DeployType,
  PortPublicationProtocol,
  Repository_AuthMethod,
  Website,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { JSXElement } from 'solid-js'
import { AiFillGithub, AiFillGitlab } from 'solid-icons/ai'
import { SiGitea } from 'solid-icons/si'
import { vars } from '/@/theme'
import { AuthMethod } from '../components/RepositoryAuthSettings'

export const buildTypeStr: Record<ApplicationConfig['buildConfig']['case'], string> = {
  runtimeBuildpack: 'Runtime (Buildpack)',
  runtimeCmd: 'Runtime (command)',
  runtimeDockerfile: 'Runtime (Dockerfile)',
  staticCmd: 'Static (command)',
  staticDockerfile: 'Static (Dockerfile)',
}

export const buildStatusStr: Record<Build_BuildStatus, string> = {
  [Build_BuildStatus.QUEUED]: 'Queued',
  [Build_BuildStatus.BUILDING]: 'Building',
  [Build_BuildStatus.SUCCEEDED]: 'Succeeded',
  [Build_BuildStatus.FAILED]: 'Failed',
  [Build_BuildStatus.CANCELLED]: 'Cancelled',
  [Build_BuildStatus.SKIPPED]: 'Skipped',
}

export enum ApplicationState {
  Idle = 'Idle',
  Deploying = 'Deploying',
  Running = 'Running',
  Static = 'Static',
}

export const applicationState = (app: Application): ApplicationState => {
  if (!app.running) {
    return ApplicationState.Idle
  }
  if (app.wantCommit !== app.currentCommit) {
    return ApplicationState.Deploying
  }
  if (app.deployType === DeployType.RUNTIME) {
    return ApplicationState.Running
  } else {
    return ApplicationState.Static
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
      return <AiFillGitlab size={size} color='#FC6D26' />
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
  const repositoryName = lastSegment?.replace(/\.git$/, '') ?? ''
  return repositoryName
}

export const authMethodMap: Record<Repository_AuthMethod, AuthMethod> = {
  [Repository_AuthMethod.NONE]: 'none',
  [Repository_AuthMethod.BASIC]: 'basic',
  [Repository_AuthMethod.SSH]: 'ssh',
}

export const portPublicationProtocolMap: Record<PortPublicationProtocol, string> = {
  [PortPublicationProtocol.TCP]: 'TCP',
  [PortPublicationProtocol.UDP]: 'UDP',
}
