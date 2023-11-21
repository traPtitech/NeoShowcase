import { Button } from '/@/components/UI/Button'
import { media, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { JSX, ParentComponent, Show } from 'solid-js'
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
      [media.mobile]: {
        padding: '32px 16px',
      },
    },
  },
})
const BackToTitle = styled('div', {
  base: {
    '@media': {
      [media.mobile]: {
        display: 'none',
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
    height: 'auto',
    left: 'calc(75% - 250px)',
    overflowX: 'hidden',
  },
})
const Titles = styled('div', {
  base: {
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
  backTo?: string
  backToTitle?: string
  icon?: JSX.Element
  action?: JSX.Element
}

export const Nav: ParentComponent<Props> = (props) => {
  return (
    <Container>
      <Show when={props.backTo} fallback={<div />}>
        {(nonNullBackTo) => (
          <A href={nonNullBackTo()}>
            <Button variants="text" size="medium" leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}>
              <BackToTitle>{props.backToTitle}</BackToTitle>
            </Button>
          </A>
        )}
      </Show>
      <TitleStickyContainer>
        <TitleContainer>
          <Titles>
            <Show when={props.icon}>{props.icon}</Show>
            <Title>{props.title}</Title>
            <Show when={props.action}>{props.action}</Show>
          </Titles>
          {props.children}
        </TitleContainer>
      </TitleStickyContainer>
    </Container>
  )
}
