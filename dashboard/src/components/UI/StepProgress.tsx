import { styled } from '@macaron-css/solid'
import { Component } from 'solid-js'
import { colorVars, textVars } from '/@/theme'
import { MaterialSymbols } from './MaterialSymbols'

const Steps = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    gap: '16px',
  },
})
const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
})
const Bar = styled('div', {
  base: {
    width: '100%',
    height: '4px',
    borderRadius: '2px',
  },
  variants: {
    state: {
      complete: {
        background: colorVars.semantic.primary.main,
      },
      current: {
        background: colorVars.semantic.primary.main,
      },
      incomplete: {
        background: colorVars.semantic.ui.tertiary,
      },
    },
  },
})
const Content = styled('div', {
  base: {
    width: '100%',
    padding: '6px 12px 8px 12px',
    display: 'flex',
    flexDirection: 'column',
    borderRadius: '4px',
  },
  variants: {
    state: {
      complete: {
        color: colorVars.semantic.primary.main,
      },
      current: {
        color: colorVars.semantic.primary.main,
        background: colorVars.semantic.transparent.primaryHover,
      },
      incomplete: {
        color: colorVars.semantic.text.grey,
      },
    },
  },
})
const Title = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '6px',
    ...textVars.h3.bold,
  },
})
const Description = styled('div', {
  base: {
    ...textVars.caption.regular,
  },
})

const StepProgress: Component<{
  title: string
  description: string
  state: 'complete' | 'current' | 'incomplete'
}> = (props) => {
  return (
    <Container>
      <Bar state={props.state} />
      <Content state={props.state}>
        <Title>
          {props.title}
          <MaterialSymbols>
            {props.state === 'complete' ? 'check_circle' : props.state === 'current' ? 'adjust' : 'circle'}
          </MaterialSymbols>
        </Title>
        <Description>{props.description}</Description>
      </Content>
    </Container>
  )
}

export const Progress = {
  Container: Steps,
  Step: StepProgress,
}
