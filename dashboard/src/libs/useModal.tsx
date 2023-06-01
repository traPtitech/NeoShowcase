import { ParentComponent, Show, createSignal, onCleanup, onMount } from 'solid-js'
import { Portal } from 'solid-js/web'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'

const ModalWrapper = styled('div', {
  base: {
    position: 'relative',
    background: vars.bg.white1,
    borderRadius: '4px',
    padding: '24px',
    opacity: 1,
    minWidth: '400px',
  },
})
const ModalBackground = styled('div', {
  base: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    background: 'rgba(0, 0, 0, 0.5)',
    display: 'grid',
    placeItems: 'center',
  },
})

const useModal = (mount?: Node) => {
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

  const Modal: ParentComponent = (props) => {
    return (
      <Show when={isOpen()}>
        <Portal mount={mount ? mount : document.body}>
          <ModalBackground onClick={close}>
            <ModalWrapper onClick={(e) => e.stopPropagation()}>{props.children}</ModalWrapper>
          </ModalBackground>
        </Portal>
      </Show>
    )
  }

  return {
    Modal,
    open,
    close,
    isOpen,
  }
}

export default useModal
