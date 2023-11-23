import { styled } from '@macaron-css/solid'
import { A, useNavigate } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { user } from '/@/libs/api'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { Button } from '../../UI/Button'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '76px',
    padding: '16px 16px 16px 20px',
    display: 'flex',
    alignItems: 'center',
    gap: '32px',

    background: colorVars.semantic.ui.primary,
    cursor: 'pointer',

    selectors: {
      '&:hover': {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.primitive.blackAlpha[50]),
      },
    },
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

export interface Props {
  repository: Repository
  appCount: number
}

export const RepositoryRow: Component<Props> = (props) => {
  const navigator = useNavigate()
  const canEdit = () => user()?.admin || props.repository.ownerIds.includes(user()?.id)

  return (
    <Container>
      <TitleContainer>
        {providerToIcon(repositoryURLToProvider(props.repository.url), 24)}
        <A
          href={`/repos/${props.repository.id}`}
          style={{
            overflow: 'hidden',
          }}
        >
          <RepositoryName>{props.repository.name}</RepositoryName>
        </A>
        <AppCount>{`${props.appCount} apps`}</AppCount>
      </TitleContainer>
      <Show when={canEdit()}>
        <AddNewAppButtonContainer>
          <Button
            variants="border"
            size="medium"
            onClick={() => {
              navigator(`/apps/new?repositoryID=${props.repository.id}`)
            }}
            tooltip={{
              props: {
                content: 'このリポジトリからアプリケーションを作成します',
              },
            }}
          >
            Add New App
          </Button>
        </AddNewAppButtonContainer>
      </Show>
    </Container>
  )
}
