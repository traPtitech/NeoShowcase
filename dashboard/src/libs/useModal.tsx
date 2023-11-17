import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { ParentComponent, Show, createSignal, onCleanup, onMount } from 'solid-js'
import { Portal } from 'solid-js/web'
import { MaterialSymbols } from '../components/UI/MaterialSymbols'

const ModalBackground = styled('div', {
  base: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    padding: '32px',
    background: colorVars.primitive.blackAlpha[600],
    display: 'grid',
    placeItems: 'center',
  },
})
const ModalWrapper = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    maxWidth: '568px',
    height: 'auto',
    maxHeight: '100%',
    display: 'flex',
    flexDirection: 'column',
    background: colorVars.semantic.ui.primary,
    borderRadius: '12px',
    opacity: 1,
    overflow: 'hidden',
  },
})
const ModalHeader = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    height: '72px',
    padding: '8px 32px',
    flexShrink: 0,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',

    color: colorVars.semantic.text.black,
    ...textVars.h2.medium,

    selectors: {
      '&:not(:last-child)': {
        borderBottom: `2px solid ${colorVars.semantic.ui.border}`,
      },
    },
  },
})
const ModalBody = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    maxHeight: '100%',
    display: 'flex',
    overflowY: 'hidden',
    padding: '24px 32px',
    selectors: {
      '&:not(:last-child)': {
        borderBottom: `2px solid ${colorVars.semantic.ui.border}`,
      },
    },
  },
})
const ModalFooter = styled('div', {
  base: {
    width: '100%',
    height: '72px',
    padding: '8px 32px',
    flexShrink: 0,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-end',
    alignItems: 'center',
    gap: '8px',
  },
})
const CloseButton = styled('button', {
  base: {
    position: 'absolute',
    width: '24px',
    height: '24px',
    top: '24px',
    right: '24px',
    padding: '0',
    background: 'none',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',

    color: colorVars.semantic.text.black,
    selectors: {
      '&:hover': {
        background: colorVars.semantic.transparent.primaryHover,
      },
      '&:active': {
        color: colorVars.semantic.primary.main,
        background: colorVars.semantic.transparent.primarySelected,
      },
    },
  },
})

const useModal = (options?: {
  mount?: Node
  size?: 'small' | 'medium'
  showCloseButton?: boolean
}) => {
  const [isOpen, setIsOpen] = createSignal(false)
  // モーダルを開くときはopen()を呼ぶ
  const open = () => setIsOpen(true)
  // モーダルを閉じるときはclose()を呼ぶ
  const close = () => setIsOpen(false)

  // ESCキーでモーダルを閉じる
  const closeOnEsc = (e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      close()
    }
  }

  onMount(() => {
    document.addEventListener('keydown', closeOnEsc)
  })
  onCleanup(() => {
    document.removeEventListener('keydown', closeOnEsc)
  })

  const Container: ParentComponent = (props) => {
    return (
      <Show when={isOpen()}>
        <Portal mount={options?.mount ? options?.mount : document.body}>
          <ModalBackground onClick={close}>
            <ModalWrapper onClick={(e) => e.stopPropagation()}>{props.children}</ModalWrapper>
          </ModalBackground>
        </Portal>
      </Show>
    )
  }
  const Header: ParentComponent = (props) => {
    return (
      <ModalHeader>
        {props.children}
        <Show when={options?.showCloseButton}>
          <CloseButton onClick={close}>
            <MaterialSymbols>close</MaterialSymbols>
          </CloseButton>
        </Show>
      </ModalHeader>
    )
  }

  return {
    Modal: {
      Container,
      Header,
      Body: ModalBody,
      Footer: ModalFooter,
    },
    open,
    close,
    isOpen,
  }
}

export default useModal
