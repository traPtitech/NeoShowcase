import { styled } from '@macaron-css/solid'
import { AiOutlineInfoCircle } from 'solid-icons/ai'
import { Component, For, createMemo } from 'solid-js'
import { tippy as tippyDir } from 'solid-tippy'
import { Content } from 'tippy.js'

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
  },
  variants: {
    align: {
      left: {
        alignItems: 'flex-start',
      },
      center: {
        alignItems: 'center',
      },
    },
  },
})

export interface InfoTooltipProps {
  tooltip: string | string[]
  style?: 'bullets' | 'bullets-with-title' | 'left' | 'center'
}

export const InfoTooltip: Component<InfoTooltipProps> = (props) => {
  const content = createMemo((): Content => {
    if (typeof props.tooltip === 'string') return props.tooltip
    if (props.style === 'bullets-with-title') {
      return (
        <TooltipContainer align="left">
          {props.tooltip[0]}
          <ul>
            <For each={props.tooltip.slice(1)}>{(line) => <li>{line}</li>}</For>
          </ul>
        </TooltipContainer>
      ) as Element
    }
    if (props.style === 'bullets') {
      return (
        <TooltipContainer align="left">
          <ul>
            <For each={props.tooltip}>{(line) => <li>{line}</li>}</For>
          </ul>
        </TooltipContainer>
      ) as Element
    }
    return (
      <TooltipContainer align={props.style ?? 'center'}>
        <For each={props.tooltip}>{(line) => <span>{line}</span>}</For>
      </TooltipContainer>
    ) as Element
  })

  return (
    <div
      use:tippy={{
        props: { content: content(), animation: 'shift-away-subtle', allowHTML: true, maxWidth: 1000 },
        hidden: true,
      }}
    >
      <Container>
        <AiOutlineInfoCircle />
      </Container>
    </div>
  )
}
