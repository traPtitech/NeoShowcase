import { type Component, Show } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '/@/components/styled-components'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import { isOwner } from '/@/libs/api'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { me } from '/@/libs/session'

const Container = styled('div', 'flex h-19 w-full items-center gap-8 bg-ui-primary p-4 pl-5')

const TitleContainer = styled('div', 'flex w-full items-center gap-2 overflow-hidden')

const RepositoryRowSkeleton: Component = () => {
  return (
    <Container class="pointer-events-none">
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
  const canEdit = () => isOwner(me(), props.repository)

  return (
    <Show when={props.repository} fallback={<RepositoryRowSkeleton />}>
      <Container>
        <TitleContainer>
          {originToIcon(repositoryURLToOrigin(props.repository!.url), 24)}
          <a href={`/repos/${props.repository?.id}`} class="overflow-hidden">
            <div class="h4-bold text-text-black truncate">{props.repository!.name}</div>
          </a>
          <div class="caption-regular text-text-grey whitespace-nowrap">{`${props.appCount} apps`}</div>
        </TitleContainer>
        <Show when={canEdit()}>
          <div class="shrink-0">
            <a href={`/apps/new?repositoryID=${props.repository!.id}`}>
              <Button
                variants="border"
                size="medium"
                tooltip={{ content: 'このリポジトリからアプリケーションを作成します' }}
              >
                Add New App
              </Button>
            </a>
          </div>
        </Show>
      </Container>
    </Show>
  )
}
