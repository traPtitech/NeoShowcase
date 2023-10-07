import { Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { user } from '/@/libs/api'
import { colorOverlay } from '/@/libs/colorOverlay'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A, useNavigate } from '@solidjs/router'
import { Component, Show, createSignal } from 'solid-js'
import { Button } from '../UI/Button'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '76px',
    padding: '16px 16px 16px 20px',
    display: 'flex',
    alignItems: 'center',
    gap: '4px',

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
  },
})
const RepositoryName = styled('div', {
  base: {
    color: colorVars.semantic.text.black,
    ...textVars.h4.bold,
  },
})
const AppCount = styled('div', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})

export interface Props {
  repository: Repository
  appCount: number
}

export const RepositoryRow: Component<Props> = (props) => {
  const navigator = useNavigate()
  const [showNewAppButton, setShowNewAppButton] = createSignal(false)
  const canEdit = () => user()?.admin || props.repository.ownerIds.includes(user()?.id)

  return (
    <Container onMouseEnter={() => setShowNewAppButton(true)} onMouseLeave={() => setShowNewAppButton(false)}>
      <TitleContainer>
        <A href={`/repos/${props.repository.id}`}>
          <RepositoryName>{props.repository.name}</RepositoryName>
        </A>
        <AppCount>{`${props.appCount} apps`}</AppCount>
      </TitleContainer>
      <Show when={canEdit() && showNewAppButton()}>
        <Button
          color="border"
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
      </Show>
    </Container>
  )
}
