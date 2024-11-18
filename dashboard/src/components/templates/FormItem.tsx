import { type JSX, type ParentComponent, Show } from 'solid-js'
import { styled } from '/@/components/styled-components'
import type { TooltipProps } from '../UI/ToolTip'
import { TooltipInfoIcon } from '../UI/TooltipInfoIcon'

export const TitleContainer = styled('div', 'flex w-full items-center gap-2')

export const RequiredMark = styled('div', 'text-accent-error text-bold')

interface Props {
  title: string | JSX.Element
  required?: boolean
  helpText?: string
  tooltip?: TooltipProps
}

export const FormItem: ParentComponent<Props> = (props) => {
  return (
    <div class="flex w-full flex-col gap-2">
      <TitleContainer>
        <div class="whitespace-nowrap text-bold text-text-black">{props.title}</div>
        <Show when={props.required}>
          <RequiredMark>*</RequiredMark>
        </Show>
        <Show when={props.tooltip}>
          <TooltipInfoIcon {...props.tooltip} />
        </Show>
        <Show when={props.helpText !== ''}>
          <div class="caption-regular text-text-grey">{props.helpText}</div>
        </Show>
      </TitleContainer>
      {props.children}
    </div>
  )
}
