import { JSXElement } from 'solid-js'
import {
  application,
  appsCount,
  container,
  header,
  headerLeft,
  addSelectButton,
  repoName,
  appName,
  appDetail,
  appFooter,
  appFooterRight,
  applicationNotLast, appDescription,
} from '/@/components/Repository.css'
import { StatusIcon } from '/@/components/StatusIcon'
import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { applicationState, providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { DiffHuman, shortSha } from '/@/libs/format'
import { A } from '@solidjs/router'

export type Provider = 'GitHub' | 'GitLab' | 'Gitea'

export interface Props {
  repo: Repository
  apps: Application[]
  func: Function
}

export const RepositoryRow = ({ repo, apps, func }: Props): JSXElement => {
  const provider = repositoryURLToProvider(repo.url)
  const sortedApps = apps.sort((a, b) => {
    if (a.updatedAt.toDate() < b.updatedAt.toDate()) {
        return 1
    }
    if (a.updatedAt.toDate() > b.updatedAt.toDate()) {
        return -1
    }
    return 0
  })
  return (
    <div class={container}>
      <div class={header}>
        <div class={headerLeft}>
          {providerToIcon(provider)}
          <div class={repoName}>{repo.name}</div>
          <div class={appDescription}>
            {/*{apps.length} {apps.length === 1 ? 'app' : 'apps'}*/}
            {apps.length && ` ${sortedApps[0].refName}/${sortedApps[0].currentCommit.slice(0, 7)}・`}
            {apps.length && <DiffHuman target={sortedApps[0].updatedAt.toDate()} />}
          </div>
        </div>
        <div class={addSelectButton}>
          <div onclick={void func}>Select</div>
        </div>
      </div>
      {/*{apps.map((app, i) => (*/}
      {/*  <div class={i === apps.length - 1 ? application : applicationNotLast}>*/}
      {/*    <StatusIcon state={applicationState(app)} />*/}
      {/*    <div class={appDetail}>*/}
      {/*      <div class={appName}>{app.name}</div>*/}
      {/*      <div class={appFooter}>*/}
      {/*        <div>{shortSha(app.currentCommit)}</div>*/}
      {/*        <div class={appFooterRight}>*/}
      {/*          <div>{app.websites[0]?.fqdn || ''}</div>*/}
      {/*          <DiffHuman target={app.updatedAt.toDate()} />*/}
      {/*        </div>*/}
      {/*      </div>*/}
      {/*    </div>*/}
      {/*  </div>*/}
      {/*))}*/}
    </div>
  )
}
