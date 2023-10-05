import CopyIcon from '/@/assets/icons/24/content_copy.svg'
import { writeToClipboard } from '/@/libs/clipboard'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component } from 'solid-js'

const Container = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'stretch',

    borderRadius: '8px',
    border: 'none',
    outline: `1px solid ${colorVars.semantic.ui.border}`,
    overflow: 'hidden',
  },
})
const Input = styled('input', {
  base: {
    width: '100%',
    height: '48px',
    padding: '10px 16px',
    background: colorVars.semantic.ui.primary,
    borderRadius: '0',
    border: 'none',
    outline: 'none',

    color: colorVars.semantic.text.black,
    ...textVars.text.regular,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})
const CopyButton = styled('button', {
  base: {
    width: '48px',
    flexShrink: 0,
    borderRadius: '0',
    border: 'none',
    borderLeft: `1px solid ${colorVars.semantic.ui.border}`,

    cursor: 'pointer',
    color: colorVars.semantic.text.black,
    background: colorVars.primitive.blackAlpha[100],

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

export const CopyableInput: Component<{
  value: string
}> = (props) => {
  const handleCopy = () => writeToClipboard(props.value)

  return (
    <Container>
      <Input value={props.value} readOnly />
      <CopyButton onClick={handleCopy} type="button">
        <CopyIcon />
      </CopyButton>
    </Container>
  )
}
