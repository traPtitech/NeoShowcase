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
import { StatusIcon } from '/@/components/StatusIcon'
import { Application, Repository } from '/@/api/neoshowcase/protobuf/apiserver_pb'
import { applicationState, providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { durationHuman, shortSha } from '/@/libs/format'
import { A } from '@solidjs/router'

export type Provider = 'GitHub' | 'GitLab' | 'Gitea'

export interface Props {
  repo: Repository
  apps: Application[]
}

export const RepositoryRow = ({ repo, apps }: Props): JSXElement => {
  const provider = repositoryURLToProvider(repo.url)
  return (
    <div class={container}>
      <div class={header}>
        <div class={headerLeft}>
          {providerToIcon(provider)}
          <div class={repoName}>{repo.name}</div>
          <div class={appsCount}>
            {apps.length} {apps.length === 1 ? 'app' : 'apps'}
          </div>
        </div>
        <div class={addBranchButton}>
          <div>Add&nbsp;branch</div>
        </div>
      </div>
      {apps.map((app, i) => (
        <A href={`/apps/${app.id}`}>
          <div className={i === apps.length - 1 ? application : applicationNotLast}>
            <StatusIcon state={applicationState(app)} />
            <div className={appDetail}>
              <div className={appName}>{app.name}</div>
              <div className={appFooter}>
                <div>{shortSha(app.currentCommit)}</div>
                <div className={appFooterRight}>
                  <div>{app.websites[0]?.fqdn || ''}</div>
                  <div>{durationHuman(3 * 60 * 1000) /* TODO: use updatedAt */}</div>
                </div>
              </div>
            </div>
          </div>
        </A>
      ))}
    </div>
  )
}
