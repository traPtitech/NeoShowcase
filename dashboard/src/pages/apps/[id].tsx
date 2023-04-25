import { A, useParams } from '@solidjs/router'
import { createResource } from 'solid-js'
import { client } from '/@/libs/api'
import { Header } from '/@/components/Header'
import { container, contentContainer } from '/@/pages/apps.css'
import { applicationState, buildTypeStr, getWebsiteURL, providerToIcon } from '/@/libs/application'
import {
  appTitle,
  appTitleContainer,
  card,
  cardItem,
  cardItemContent,
  cardItems,
  cardItemTitle,
  cardTitle,
  centerInline,
} from '/@/pages/apps/[id].css'
import { StatusIcon } from '/@/components/StatusIcon'
import { titleCase } from '/@/libs/casing'
import { Application_ContainerState, BuildConfig, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DiffHuman, shortSha } from '/@/libs/format'
import { url } from '/@/theme.css'

interface BuildConfigInfoProps {
  config: BuildConfig
}

const BuildConfigInfo = (props: BuildConfigInfoProps) => {
  const c = props.config.buildConfig
  switch (c.case) {
    case 'runtimeCmd':
      return (
        <>
          <div class={cardItem}>
            <div class={cardItemTitle}>Base Image</div>
            <div class={cardItemContent}>{c.value.baseImage || 'Scratch'}</div>
          </div>
          {c.value.buildCmd && (
            <div class={cardItem}>
              <div class={cardItemTitle}>Build Command{c.value.buildCmdShell && ' (Shell)'}</div>
              <div class={cardItemContent}>{c.value.buildCmd}</div>
            </div>
          )}
        </>
      )
    case 'runtimeDockerfile':
      return (
        <>
          <div class={cardItem}>
            <div class={cardItemTitle}>Dockerfile</div>
            <div class={cardItemContent}>{c.value.dockerfileName}</div>
          </div>
        </>
      )
    case 'staticCmd':
      return (
        <>
          <div class={cardItem}>
            <div class={cardItemTitle}>Base Image</div>
            <div class={cardItemContent}>{c.value.baseImage || 'Scratch'}</div>
          </div>
          {c.value.buildCmd && (
            <div class={cardItem}>
              <div class={cardItemTitle}>Build Command{c.value.buildCmdShell && ' (Shell)'}</div>
              <div class={cardItemContent}>{c.value.buildCmd}</div>
            </div>
          )}
          <div class={cardItem}>
            <div class={cardItemTitle}>Artifact Path</div>
            <div class={cardItemContent}>{c.value.artifactPath}</div>
          </div>
        </>
      )
    case 'staticDockerfile':
      return (
        <>
          <div class={cardItem}>
            <div class={cardItemTitle}>Dockerfile</div>
            <div class={cardItemContent}>{c.value.dockerfileName}</div>
          </div>
          <div class={cardItem}>
            <div class={cardItemTitle}>Artifact Path</div>
            <div class={cardItemContent}>{c.value.artifactPath}</div>
          </div>
        </>
      )
  }
}

export default () => {
  const params = useParams()
  const [app] = createResource(
    () => params.id,
    (id) => client.getApplication({ id }),
  )
  const [repo] = createResource(
    () => app()?.repositoryId,
    (id) => client.getRepository({ repositoryId: id }),
  )

  return (
    <div class={container}>
      <Header />
      <div class={appTitleContainer}>
        <div class={centerInline}>{providerToIcon('GitHub', 36)}</div>
        <div class={appTitle}>
          <div>{repo()?.name}</div>
          <div>/</div>
          <div>{app()?.name}</div>
        </div>
      </div>
      <div class={contentContainer}>
        <div class={card}>
          <div class={cardTitle}>Overall</div>
          <div class={cardItems}>
            <div class={cardItem}>
              <div class={cardItemTitle}>状態</div>
              <div class={cardItemContent}>
                {app() && <StatusIcon state={applicationState(app())} size={24} />}
                {app() && applicationState(app())}
              </div>
            </div>
            {app() && app().deployType === DeployType.RUNTIME && (
              <div class={cardItem}>
                <div class={cardItemTitle}>コンテナの状態</div>
                <div class={cardItemContent}>{app() && titleCase(Application_ContainerState[app().container])}</div>
              </div>
            )}
            <div class={cardItem}>
              <div class={cardItemTitle}>起動時刻</div>
              <div class={cardItemContent}>
                {app()?.running && <DiffHuman target={app().updatedAt.toDate()} />}
                {app() && !app().running && '-'}
              </div>
            </div>
            <div class={cardItem}>
              <div class={cardItemTitle}>作成日</div>
              <div class={cardItemContent}>{app() && <DiffHuman target={app().createdAt.toDate()} />}</div>
            </div>
            {app() && app().websites.length > 0 && (
              <div class={cardItem}>
                <div class={cardItemTitle}>URLs</div>
              </div>
            )}
            {app() && app().websites.length > 0 && (
              <div class={cardItem}>
                <div class={cardItemContent}>
                  {app()?.websites.map((website) => (
                    <a class={url} href={getWebsiteURL(website)} target='_blank' rel="noreferrer">
                      {getWebsiteURL(website)}
                    </a>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
        <div class={card}>
          <div class={cardTitle}>Info</div>
          <div class={cardItems}>
            <div class={cardItem}>
              <div class={cardItemTitle}>ID</div>
              <div class={cardItemContent}>{app()?.id}</div>
            </div>
            <div class={cardItem}>
              <div class={cardItemTitle}>Name</div>
              <div class={cardItemContent}>{app()?.name}</div>
            </div>
            <A class={cardItem} href={`/repos/${repo()?.id}`}>
              <div class={cardItemTitle}>Repository</div>
              <div class={cardItemContent}>
                <div class={centerInline}>{providerToIcon('GitHub', 20)}</div>
                {repo()?.name}
              </div>
            </A>
            <div class={cardItem}>
              <div class={cardItemTitle}>Git ref (short)</div>
              <div class={cardItemContent}>{app()?.refName}</div>
            </div>
            <div class={cardItem}>
              <div class={cardItemTitle}>Deploy type</div>
              <div class={cardItemContent}>{app() && titleCase(DeployType[app().deployType])}</div>
            </div>
            <div class={cardItem}>
              <div class={cardItemTitle}>Commit</div>
              <div class={cardItemContent}>
                {app() && app().currentCommit !== app().wantCommit && (
                  <div>
                    {shortSha(app().currentCommit)} → {shortSha(app().wantCommit)} (Deploying)
                  </div>
                )}
                {app() && app().currentCommit === app().wantCommit && <div>{shortSha(app().currentCommit)}</div>}
              </div>
            </div>
            <div class={cardItem}>
              <div class={cardItemTitle}>Use MariaDB</div>
              <div class={cardItemContent}>{app() && `${app().config.useMariadb}`}</div>
            </div>
            <div class={cardItem}>
              <div class={cardItemTitle}>Use MongoDB</div>
              <div class={cardItemContent}>{app() && `${app().config.useMongodb}`}</div>
            </div>
          </div>
        </div>
        <div class={card}>
          <div class={cardTitle}>Build Config</div>
          <div class={cardItems}>
            <div class={cardItem}>
              <div class={cardItemTitle}>Build Type</div>
              <div class={cardItemContent}>{app() && buildTypeStr[app().config.buildConfig.buildConfig.case]}</div>
            </div>
            {app()?.config.buildConfig && <BuildConfigInfo config={app().config.buildConfig} />}
            {app()?.config.entrypoint && (
              <div class={cardItem}>
                <div class={cardItemTitle}>Entrypoint</div>
                <div class={cardItemContent}>{app()?.config.entrypoint}</div>
              </div>
            )}
            {app()?.config.command && (
              <div class={cardItem}>
                <div class={cardItemTitle}>Command</div>
                <div class={cardItemContent}>{app()?.config.command}</div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
