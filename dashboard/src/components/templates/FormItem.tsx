import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { type JSX, type ParentComponent, Show } from 'solid-js'
import { colorVars, textVars } from '/@/theme'
import type { TooltipProps } from '../UI/ToolTip'
import { TooltipInfoIcon } from '../UI/TooltipInfoIcon'

export const containerStyle = style({
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  gap: '8px',
})
export const TitleContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '2px',
  },
})
export const titleStyle = style({
  whiteSpace: 'nowrap',
  color: colorVars.semantic.text.black,
  ...textVars.text.bold,
})
export const RequiredMark = styled('div', {
  base: {
    color: colorVars.semantic.accent.error,
    ...textVars.text.bold,
  },
})
export const errorTextStyle = style({
  width: '100%',
  color: colorVars.semantic.accent.error,
  ...textVars.text.regular,
})
const HelpText = styled('div', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})

interface Props {
  title: string | JSX.Element
  required?: boolean
  helpText?: string
  tooltip?: TooltipProps
}

export const FormItem: ParentComponent<Props> = (props) => {
  return (
    <div class={containerStyle}>
      <TitleContainer>
        <div class={titleStyle}>{props.title}</div>
        <Show when={props.required}>
          <RequiredMark>*</RequiredMark>
        </Show>
        <Show when={props.tooltip}>
          <TooltipInfoIcon {...props.tooltip} />
        </Show>
        <Show when={props.helpText !== ''}>
          <HelpText>{props.helpText}</HelpText>
        </Show>
      </TitleContainer>
      {props.children}
    </div>
  )
}
