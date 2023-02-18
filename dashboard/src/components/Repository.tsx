import { JSXElement } from 'solid-js'
import {
  application,
  appsCount,
  container,
  header,
  headerLeft,
  addBranchButton,
  repoName,
  appName,
  appDetail,
  appFooter,
  appFooterRight,
  applicationNotLast,
} from '/@/components/Repository.css'
import { AiFillGithub, AiFillGitlab } from 'solid-icons/ai'
import { vars } from '/@/theme.css'
import { SiGitea } from 'solid-icons/si'
import { Status, StatusIcon } from '/@/components/StatusIcon'

export type Provider = 'GitHub' | 'GitLab' | 'Gitea'

export interface Application {
  name: string
  status: Status
  lastCommit: string
  lastCommitDate: string
  url: string
  updateDate: string
}

export interface Props {
  name: string
  provider: Provider
  apps: Application[]
}

const providerToIcon = (provider: Provider): JSXElement => {
  switch (provider) {
    case 'GitHub':
      return <AiFillGithub size={20} color={vars.text.black1} />
    case 'GitLab':
      return <AiFillGitlab size={20} color='#FC6D26' />
    case 'Gitea':
      return <SiGitea size={20} color={vars.text.black1} />
  }
}

export const Repository = ({ name, provider, apps }: Props): JSXElement => {
  return (
    <div class={container}>
      <div class={header}>
        <div class={headerLeft}>
          {providerToIcon(provider)}
          <div class={repoName}>{name}</div>
          <div class={appsCount}>
            {apps.length} {apps.length === 1 ? 'app' : 'apps'}
          </div>
        </div>
        <div class={addBranchButton}>
          <div>Add&nbsp;branch</div>
        </div>
      </div>
      {apps.map((app, i) => (
        <div class={i === apps.length - 1 ? application : applicationNotLast}>
          <StatusIcon status={app.status} />
          <div class={appDetail}>
            <div class={appName}>{app.name}</div>
            <div class={appFooter}>
              <div>
                {app.lastCommit}・{app.lastCommitDate}
              </div>
              <div class={appFooterRight}>
                <div>{app.url}</div>
                <div>{app.updateDate}</div>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
