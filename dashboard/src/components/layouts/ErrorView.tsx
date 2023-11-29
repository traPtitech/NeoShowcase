import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { colorVars, textVars } from '/@/theme'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    gap: '1rem',
  },
})
const Title = styled('h2', {
  base: {
    color: colorVars.semantic.accent.error,
    ...textVars.h2.bold,
  },
})
const Message = styled('p', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.caption.medium,
  },
})

const ErrorView: Component<{
  error: unknown
}> = (props) => {
  return (
    <Container>
      <MaterialSymbols fill displaySize={64} color={colorVars.semantic.accent.error}>
        error
      </MaterialSymbols>
      <Title>An error has occurred</Title>
      <Show when={props.error instanceof Error}>
        <Message>{(props.error as Error).message}</Message>
      </Show>
      <A href="/">
        <Button size="medium" variants="border" leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}>
          Back to Home
        </Button>
      </A>
    </Container>
  )
}

export default ErrorView
