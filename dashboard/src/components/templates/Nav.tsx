import { Button } from '/@/components/UI/Button'
import { textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, JSX, Show } from 'solid-js'
import { MaterialSymbols } from '../UI/MaterialSymbols'

const Container = styled('div', {
  base: {
    width: '100%',
    padding: '32px',
    display: 'grid',
    gridTemplateColumns: '104px 1fr 104px',
  },
})
const TitleContainer = styled('div', {
  base: {
    width: '100%',
    height: '44px',
    maxWidth: '1000px',
    margin: '0 auto',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
})
const Title = styled('h1', {
  base: {
    width: '100%',
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
          color="text"
          size="medium"
          onClick={() => {
            window.history.back()
          }}
          leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}
        >
          {props.backToTitle}
        </Button>
      </Show>
      <TitleContainer>
        <Show when={props.icon}>{props.icon}</Show>
        <Title>{props.title}</Title>
        <Show when={props.action}>{props.action}</Show>
      </TitleContainer>
      <div />
    </Container>
  )
}
