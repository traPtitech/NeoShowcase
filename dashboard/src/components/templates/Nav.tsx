import { Button } from '/@/components/UI/Button'
import { textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, JSX, Show } from 'solid-js'
import { MaterialSymbols } from '../UI/MaterialSymbols'

const Container = styled('div', {
  base: {
    width: '100%',
    overflowX: 'hidden',
    padding: '32px 32px 32px 32px',
    paddingRight: 'max(calc(50% - 500px), 32px)',
    display: 'flex',
    gap: '8px',

    '@media': {
      'screen and (max-width: 768px)': {
        padding: '32px 16px',
      },
    },
  },
})
const TitleStickyContainer = styled('div', {
  base: {
    width: '100%',
    overflowX: 'clip',
  },
})
const TitleContainer = styled('div', {
  base: {
    position: 'sticky',
    width: '100%',
    maxWidth: '1000px',
    height: '44px',
    left: 'calc(75% - 250px)',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    overflowX: 'hidden',
  },
})
const Title = styled('h1', {
  base: {
    width: '100%',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    ...textVars.h1.medium,
  },
})
export interface Props {
  title: string
  icon?: JSX.Element
  backToTitle?: string
  action?: JSX.Element
}

export const Nav: Component<Props> = (props) => {
  return (
    <Container>
      <Show when={props.backToTitle} fallback={<div />}>
        <Button
          variants="text"
          size="medium"
          onClick={() => {
            window.history.back()
          }}
          leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}
        >
          {props.backToTitle}
        </Button>
      </Show>
      <TitleStickyContainer>
        <TitleContainer>
          <Show when={props.icon}>{props.icon}</Show>
          <Title>{props.title}</Title>
          <Show when={props.action}>{props.action}</Show>
        </TitleContainer>
      </TitleStickyContainer>
    </Container>
  )
}
