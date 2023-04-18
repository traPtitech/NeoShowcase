import { useParams } from '@solidjs/router'
import { createResource } from 'solid-js'
import { client } from '/@/libs/api'
import { Header } from '/@/components/Header'
import { container, contentContainer } from '/@/pages/apps.css'
import { applicationState, getWebsiteURL, providerToIcon } from '/@/libs/application'
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
import { Application_ContainerState, DeployType } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DiffHuman } from '/@/libs/format'
import { url } from '/@/theme.css'

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
      </div>
    </div>
  )
}
