import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import type { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import { styled } from '/@/components/styled-components'
import { user } from '/@/libs/api'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'

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
  const canEdit = () => user()?.admin || (user.latest !== undefined && props.repository?.ownerIds.includes(user()?.id))

  return (
    <Show when={props.repository} fallback={<RepositoryRowSkeleton />}>
      <Container>
        <TitleContainer>
          {originToIcon(repositoryURLToOrigin(props.repository!.url), 24)}
          <A href={`/repos/${props.repository?.id}`} class="overflow-hidden">
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
