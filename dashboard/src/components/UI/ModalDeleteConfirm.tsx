import { styled } from '@macaron-css/solid'
import type { ParentComponent } from 'solid-js'
import { colorVars, textVars } from '/@/theme'

const DeleteConfirm = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    borderRadius: '8px',
    background: colorVars.semantic.ui.secondary,
    overflowWrap: 'anywhere',
    color: colorVars.semantic.text.black,
    ...textVars.h3.regular,
  },
})
const ModalDeleteConfirm: ParentComponent = (props) => {
  return <DeleteConfirm>{props.children}</DeleteConfirm>
}

export default ModalDeleteConfirm
