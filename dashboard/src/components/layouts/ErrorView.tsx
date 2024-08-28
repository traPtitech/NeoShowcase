import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { colorVars } from '/@/theme'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'

const Container = styled('div', 'flex w-full flex-col items-center justify-center gap-4')

const Title = styled('h2', 'h2-bold text-accent-error')

const Message = styled('p', 'caption-medium text-text-grey')

const ButtonsContainer = styled('div', 'flex flex-col gap-2')

const ErrorView: Component<{
  error: unknown
}> = (props) => {
  const handleReload = () => {
    window.location.reload()
  }

  return (
    <Container>
      <MaterialSymbols fill displaySize={64} color={colorVars.semantic.accent.error}>
        error
      </MaterialSymbols>
      <Title>An error has occurred</Title>
      <Show when={props.error instanceof Error}>
        <Message>{(props.error as Error).message}</Message>
      </Show>
      <ButtonsContainer>
        <A href="/">
          <Button size="medium" variants="border" leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}>
            Back to Home
          </Button>
        </A>
        <Button
          onClick={handleReload}
          size="medium"
          variants="border"
          leftIcon={<MaterialSymbols>refresh</MaterialSymbols>}
        >
          Reload Page
        </Button>
      </ButtonsContainer>
    </Container>
  )
}

export default ErrorView
