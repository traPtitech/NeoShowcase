import { type Component, Show } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '/@/components/styled-components'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import { isOwner } from '/@/libs/api'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { me } from '/@/libs/session'

const Container = styled('div', 'flex h-19 w-full items-center gap-8 bg-ui-primary')

const containerPadding = 'p-4 pl-5'

const TitleContainer = styled('div', 'flex w-full items-center gap-2 overflow-hidden')

const RepositoryRowSkeleton: Component = () => {
  return (
    <Container class="pointer-events-none">
      <TitleContainer class={containerPadding}>
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
      <Container class="hover:bg-color-overlay-ui-primary-to-black-alpha-50">
        {/* "Add New App" is a link of its own, so the link to the repository covers the rest of the row instead of the whole row. */}
        <a href={`/repos/${props.repository?.id}`} class={['min-w-0 flex-1 self-stretch', containerPadding]}>
          <TitleContainer class="h-full">
            {originToIcon(repositoryURLToOrigin(props.repository!.url), 24)}
            <div class="h4-bold text-text-black truncate">{props.repository!.name}</div>
            <div class="caption-regular text-text-grey whitespace-nowrap">{`${props.appCount} apps`}</div>
          </TitleContainer>
        </a>
        <Show when={canEdit()}>
          <div class="shrink-0 pr-4">
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
