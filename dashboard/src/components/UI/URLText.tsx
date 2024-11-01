import type { Component } from 'solid-js'
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
          <div class="i-material-symbols:open-in-new text-xl/5" />
        </div>
      </a>
    </ToolTip>
  )
}
