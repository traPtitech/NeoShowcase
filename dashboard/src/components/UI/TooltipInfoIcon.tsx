import { colorVars } from '/@/theme'
import { Component } from 'solid-js'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip, TooltipProps } from './ToolTip'

export const TooltipInfoIcon: Component<TooltipProps> = (props) => {
  return (
    <ToolTip {...props}>
      <MaterialSymbols opticalSize={20} color={colorVars.semantic.text.link}>
        help
      </MaterialSymbols>
    </ToolTip>
  )
}
