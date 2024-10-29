import type { Component } from 'solid-js'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip, type TooltipProps } from './ToolTip'

export const TooltipInfoIcon: Component<TooltipProps> = (props) => {
  return (
    <ToolTip {...props}>
      <MaterialSymbols opticalSize={20} class="text-text-black">
        help
      </MaterialSymbols>
    </ToolTip>
  )
}
