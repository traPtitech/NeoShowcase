import type { Component } from 'solid-js'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip } from './ToolTip'

export interface URLTextProps {
  href: string
  text: string
}

export const URLText: Component<URLTextProps> = (props) => {
  return (
    <ToolTip
      props={{
        content: props.href,
      }}
      disabled={props.text === props.href}
    >
      <a class="text-regular text-text-link" href={props.href} target="_blank" rel="noreferrer">
        <div class="flex items-center gap-1">
          {props.text}
          <MaterialSymbols opticalSize={20}>open_in_new</MaterialSymbols>
        </div>
      </a>
    </ToolTip>
  )
}
