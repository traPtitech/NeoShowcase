import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { Component } from 'solid-js'
import { tippy as tippyDir } from 'solid-tippy'
import { VsLinkExternal } from 'solid-icons/vs'

// https://github.com/solidjs/solid/discussions/845
const tippy = tippyDir

const StyledAnchor = styled('a', {
  base: {
    color: vars.text.url,
  },
})

const ContentContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '4px',
    alignItems: 'center',
  },
})

export interface URLTextProps {
  href: string
  text: string
}

export const URLText: Component<URLTextProps> = (props) => {
  return (
    <div
      use:tippy={{
        props: { content: props.href, maxWidth: 1000 },
        disabled: props.text === props.href,
        hidden: true,
      }}
    >
      <StyledAnchor href={props.href} target='_blank' rel='noreferrer'>
        <ContentContainer>
          {props.text}
          <VsLinkExternal />
        </ContentContainer>
      </StyledAnchor>
    </div>
  )
}
