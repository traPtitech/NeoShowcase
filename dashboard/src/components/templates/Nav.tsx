import { Button } from '/@/components/UI/Button'
import { textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A, BeforeLeaveEventArgs, useBeforeLeave } from '@solidjs/router'
import { Component, JSX, Show, createSignal, onMount } from 'solid-js'
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
  action?: JSX.Element
}

const [prevPath, setPrevPath] = createSignal<string | undefined>(undefined)

export const Nav: Component<Props> = (props) => {
  const [backTo, setBackTo] = createSignal<string | undefined>(undefined)

  useBeforeLeave((e: BeforeLeaveEventArgs) => {
    setPrevPath(e.from.pathname)
  })
  onMount(() => {
    setBackTo(prevPath())
  })

  const backToTitle = () => {
    const reg = new RegExp(/\/\w+\/?/)
    const startsWith = backTo()?.match(reg)?.[0]
    switch (startsWith) {
      case '/apps':
        return 'Apps'
      case '/apps/':
        return 'App'
      case '/repos/':
        return 'Repository'
      default:
        return undefined
    }
  }

  return (
    <Container>
      <Show when={backTo()} fallback={<div />}>
        {(nonNullBackTo) => (
          <A href={nonNullBackTo()}>
            <Button variants="text" size="medium" leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}>
              {backToTitle()}
            </Button>
          </A>
        )}
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
