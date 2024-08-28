import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import { user } from '/@/libs/api'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { styled } from '/@/components/styled-components'

const Container = styled('div', 'w-full h-19 p-4 pl-5 flex items-center gap-8 bg-ui-primary')

const TitleContainer = styled('div', 'w-full flex items-center gap-2 overflow-hidden')

const RepositoryRowSkeleton: Component = () => {
  return (
    <Container style={{ 'pointer-events': 'none' }}>
      <TitleContainer>
        <Skeleton width={24} height={24} circle />
        <Skeleton>Repository Name Placeholder</Skeleton>
      </TitleContainer>
    </Container>
  )
}

export interface Props {
  repository?: Repository
  appCount?: number
}

export const RepositoryRow: Component<Props> = (props) => {
  const canEdit = () => user()?.admin || (user.latest !== undefined && props.repository?.ownerIds.includes(user()?.id))

  return (
    <Show when={props.repository} fallback={<RepositoryRowSkeleton />}>
      <Container>
        <TitleContainer>
          {originToIcon(repositoryURLToOrigin(props.repository!.url), 24)}
          <A
            href={`/repos/${props.repository?.id}`}
            style={{
              overflow: 'hidden',
            }}
          >
            <div class="h4-bold truncate text-text-black">{props.repository!.name}</div>
          </A>
          <div class="caption-regular whitespace-nowrap text-text-grey">{`${props.appCount} apps`}</div>
        </TitleContainer>
        <Show when={canEdit()}>
          <div class="shrink-0">
            <A href={`/apps/new?repositoryID=${props.repository!.id}`}>
              <Button
                variants="border"
                size="medium"
                tooltip={{
                  props: {
                    content: 'このリポジトリからアプリケーションを作成します',
                  },
                }}
              >
                Add New App
              </Button>
            </A>
          </div>
        </Show>
      </Container>
    </Show>
  )
}
