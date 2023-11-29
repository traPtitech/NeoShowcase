import { styled } from '@macaron-css/solid'
import { Component, Show } from 'solid-js'
import { writeToClipboard } from '/@/libs/clipboard'
import { colorVars } from '/@/theme'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip } from './ToolTip'

const Container = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    minHeight: 'calc(1lh + 8px)',
    marginTop: '4px',
    whiteSpace: 'pre-wrap',
    overflowX: 'auto',
    padding: '4px 8px',
    fontSize: '16px',
    lineHeight: '1.5',
    fontFamily: 'Menlo, Monaco, Consolas, Courier New, monospace !important',
    background: colorVars.semantic.ui.secondary,
    borderRadius: '4px',
    color: colorVars.semantic.text.black,
  },
  variants: {
    copyable: {
      true: {
        paddingRight: '40px',
      },
    },
  },
})
const CopyButton = styled('button', {
  base: {
    position: 'absolute',
    width: '24px',
    height: '24px',
    top: '4px',
    right: '8px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    border: `solid 1px ${colorVars.semantic.ui.border}`,
    borderRadius: '4px',

    cursor: 'pointer',
    color: colorVars.semantic.text.black,
    background: 'none',
    lineHeight: 1,

    selectors: {
      '&:hover': {
        background: colorVars.primitive.blackAlpha[200],
      },
      '&:active': {
        background: colorVars.primitive.blackAlpha[300],
      },
    },
  },
})

const Code: Component<{
  value: string
  copyable?: boolean
}> = (props) => {
  const handleCopy = () => {
    writeToClipboard(props.value)
  }

  return (
    <Container copyable={props.copyable}>
      {props.value}
      <Show when={props.copyable}>
        <ToolTip
          props={{
            content: 'copy to clipboard',
          }}
        >
          <CopyButton onClick={handleCopy} type="button">
            <MaterialSymbols opticalSize={20}>content_copy</MaterialSymbols>
          </CopyButton>
        </ToolTip>
      </Show>
    </Container>
  )
}

export default Code
