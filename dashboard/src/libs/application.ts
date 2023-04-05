import { Provider } from '/@/components/RepositoryRow'
import { Application, BuildType } from '/@/api/neoshowcase/protobuf/apiserver_pb'

export enum ApplicationState {
  Idle = 0,
  Deploying = 1,
  Running = 2,
  Static = 3,
}

export const applicationState = (app: Application): ApplicationState => {
  if (!app.running) {
    return ApplicationState.Idle
  }
  if (app.wantCommit !== app.currentCommit) {
    return ApplicationState.Deploying
  }
  if (app.buildType === BuildType.RUNTIME) {
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
