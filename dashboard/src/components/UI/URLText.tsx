import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component } from 'solid-js'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip } from './ToolTip'

const StyledAnchor = styled('a', {
  base: {
    color: colorVars.semantic.text.link,
    ...textVars.text.regular,
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
    <ToolTip
      props={{
        content: props.href,
      }}
      disabled={props.text === props.href}
    >
      <StyledAnchor href={props.href} target="_blank" rel="noreferrer">
        <ContentContainer>
          {props.text}
          <MaterialSymbols opticalSize={20}>open_in_new</MaterialSymbols>
        </ContentContainer>
      </StyledAnchor>
    </ToolTip>
  )
}
