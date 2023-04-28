import { Provider } from '/@/components/RepositoryRow'
import { Application, BuildConfig, DeployType, Website } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { JSXElement } from 'solid-js'
import { AiFillGithub, AiFillGitlab } from 'solid-icons/ai'
import { SiGitea } from 'solid-icons/si'
import { vars } from '/@/theme.css'

export const buildTypeStr: Record<BuildConfig['buildConfig']['case'], string> = {
  runtimeBuildpack: 'Runtime (Buildpack)',
  runtimeCmd: 'Runtime (command)',
  runtimeDockerfile: 'Runtime (Dockerfile)',
  staticCmd: 'Static (command)',
  staticDockerfile: 'Static (Dockerfile)',
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
  url = url.toLowerCase()
  if (url.includes('github')) return 'GitHub'
  if (url.includes('gitlab')) return 'GitLab'
  if (url.includes('gitea')) return 'Gitea'
  if (url.includes('git.trap.jp')) return 'Gitea'
  return 'GitHub' // fallback?
}

export const providerToIcon = (provider: Provider, size: number = 20): JSXElement => {
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
