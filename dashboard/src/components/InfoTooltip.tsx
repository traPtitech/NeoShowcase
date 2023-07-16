import { Component, createMemo, For } from 'solid-js'
import { styled } from '@macaron-css/solid'
import { AiOutlineInfoCircle } from 'solid-icons/ai'
import { tippy as tippyDir, TippyOptions } from 'solid-tippy'
import 'tippy.js/dist/tippy.css'
import 'tippy.js/animations/shift-away-subtle.css'
import { Content } from 'tippy.js'

declare module 'solid-js' {
  namespace JSX {
    interface Directives {
      tippy: TippyOptions
    }
  }
}

// https://github.com/solidjs/solid/discussions/845
const tippy = tippyDir

const Container = styled('div', {
  base: {
    position: 'relative',

    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',

    width: '20px',
    height: '20px',
  },
})

const TooltipContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
})

export interface InfoTooltipProps {
  tooltip: string | string[]
}

export const InfoTooltip: Component<InfoTooltipProps> = (props) => {
  const content = createMemo((): Content => {
    if (typeof props.tooltip === 'string') return props.tooltip
    return (
      <TooltipContainer>
        <For each={props.tooltip}>{(line) => <span>{line}</span>}</For>
      </TooltipContainer>
    ) as Element
  })

  return (
    <Container>
      <div use:tippy={{ props: { content: content(), animation: 'shift-away-subtle', allowHTML: true }, hidden: true }}>
        <AiOutlineInfoCircle />
      </div>
    </Container>
  )
}
