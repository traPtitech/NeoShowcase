import { Component } from 'solid-js'
import { styled } from '@macaron-css/solid'
import { AiOutlineInfoCircle } from 'solid-icons/ai'
import { tippy as tippyDir, TippyOptions } from 'solid-tippy'
import 'tippy.js/dist/tippy.css'

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

export interface InfoTooltipProps {
  tooltip: string
}

export const InfoTooltip: Component<InfoTooltipProps> = (props) => {
  return (
    <Container>
      <div use:tippy={{ props: { content: props.tooltip }, hidden: true }}>
        <AiOutlineInfoCircle />
      </div>
    </Container>
  )
}
