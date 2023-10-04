import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, ParentComponent, Show } from 'solid-js'
import { InfoTooltip } from '../InfoTooltip'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
})
const TitleContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '2px',
  },
})
const Title = styled('div', {
  base: {
    color: colorVars.semantic.text.black,
    ...textVars.text.bold,
  },
})
const Required = styled('div', {
  base: {
    color: colorVars.semantic.accent.error,
    ...textVars.text.bold,
  },
})
const HelpText = styled('div', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})

interface Props {
  title: string
  required?: boolean
  helpText?: string
  tooltip?: string | string[]
}

export const FormItem: ParentComponent<Props> = (props) => {
  return (
    <Container>
      <TitleContainer>
        <Title>{props.title}</Title>
        <Show when={props.required}>
          <Required>*</Required>
        </Show>
        <Show when={props.tooltip}>
          <InfoTooltip tooltip={props.tooltip} />
        </Show>
        <Show when={props.helpText}>
          <HelpText>{props.helpText}</HelpText>
        </Show>
      </TitleContainer>
      {props.children}
    </Container>
  )
}
