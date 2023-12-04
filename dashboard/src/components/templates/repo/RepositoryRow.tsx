import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import { user } from '/@/libs/api'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { colorVars, textVars } from '/@/theme'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '76px',
    padding: '16px 16px 16px 20px',
    display: 'flex',
    alignItems: 'center',
    gap: '32px',

    background: colorVars.semantic.ui.primary,
  },
})
const TitleContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    overflow: 'hidden',
  },
})
const RepositoryName = styled('div', {
  base: {
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.black,
    ...textVars.h4.bold,
  },
})
const AppCount = styled('div', {
  base: {
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const AddNewAppButtonContainer = styled('div', {
  base: {
    flexShrink: 0,
  },
})

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
          {providerToIcon(repositoryURLToProvider(props.repository!.url), 24)}
          <A
            href={`/repos/${props.repository!.id}`}
            style={{
              overflow: 'hidden',
            }}
          >
            <RepositoryName>{props.repository!.name}</RepositoryName>
          </A>
          <AppCount>{`${props.appCount} apps`}</AppCount>
        </TitleContainer>
        <Show when={canEdit()}>
          <AddNewAppButtonContainer>
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
          </AddNewAppButtonContainer>
        </Show>
      </Container>
    </Show>
  )
}
