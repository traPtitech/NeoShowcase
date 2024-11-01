import type { Component } from 'solid-js'
import { ToolTip, type TooltipProps } from './ToolTip'

export const TooltipInfoIcon: Component<TooltipProps> = (props) => {
  return (
    <ToolTip {...props}>
      <div class="i-material-symbols:help-outline text-text-black text-xl/5" />
    </ToolTip>
  )
}
