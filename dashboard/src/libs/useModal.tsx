import { Dialog, type PolymorphicCallbackProps } from '@kobalte/core'
import type {
  DialogContentOptions,
  DialogContentRenderProps,
  DialogDescriptionOptions,
  DialogDescriptionRenderProps,
} from '@kobalte/core/dialog'
import { keyframes, style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { type ComponentProps, type ParentComponent, Show, createSignal, mergeProps } from 'solid-js'
import { colorVars, textVars } from '/@/theme'
import { MaterialSymbols } from '../components/UI/MaterialSymbols'

const overlayShow = keyframes({
  from: {
    opacity: 0,
  },
  to: {
    opacity: 1,
  },
})
const overlayHide = keyframes({
  from: {
    opacity: 1,
  },
  to: {
    opacity: 0,
  },
})
const overlayStyle = style({
  position: 'fixed',
  inset: 0,
  background: colorVars.primitive.blackAlpha[600],
  animation: `${overlayHide} 0.2s`,
  selectors: {
    '&[data-expanded]': {
      animation: `${overlayShow} 0.2s`,
    },
  },
})
const DialogPositioner = styled('div', {
  base: {
    position: 'fixed',
    inset: 0,
    padding: '32px',
    display: 'grid',
    placeItems: 'center',
  },
})
const contentShow = keyframes({
  from: {
    opacity: 0,
    transform: 'scale(0.95)',
  },
  to: {
    opacity: 1,
    transform: 'scale(1)',
  },
})
const contentHide = keyframes({
  from: {
    opacity: 1,
    transform: 'scale(1)',
  },
  to: {
    opacity: 0,
    transform: 'scale(0.95)',
  },
})
const Content = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    maxWidth: '568px',
    maxHeight: '100%',
    display: 'flex',
    flexDirection: 'column',
    background: colorVars.semantic.ui.primary,
    borderRadius: '12px',
    opacity: 1,
    overflow: 'hidden',

    animation: `${contentHide} 0.3s`,
    selectors: {
      '&[data-expanded]': {
        animation: `${contentShow} 0.3s`,
      },
    },
  },
  variants: {
    fit: {
      true: {
        height: 'auto',
      },
      false: {
        height: '100%',
      },
    },
  },
  defaultVariants: {
    fit: true,
  },
})
const DialogHeader = styled('div', {
  base: {
    position: 'relative',
    width: '100%',
    height: '72px',
    padding: '8px 32px',
    flexShrink: 0,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
})
const titleStyle = style({
  color: colorVars.semantic.text.black,
  ...textVars.h2.medium,
})
const Description = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    maxHeight: '100%',
    display: 'flex',
    overflowY: 'hidden',
    padding: '24px 32px',
    selectors: {
      [`${DialogHeader.selector({})}~&`]: {
        borderTop: `2px solid ${colorVars.semantic.ui.border}`,
      },
    },
  },
  variants: {
    fit: {
      true: {
        height: 'auto',
      },
      false: {
        height: '100%',
      },
    },
  },
  defaultVariants: {
    fit: true,
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
    selectors: {
      [`${Description.selector({})}~&`]: {
        borderTop: `2px solid ${colorVars.semantic.ui.border}`,
      },
    },
  },
})
const closeButtonStyle = style({
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
})

const useModal = (options?: {
  showCloseButton?: boolean
  closeOnClickOutside?: boolean
}) => {
  const defaultOptions = {
    showCloseButton: false,
    closeOnClickOutside: true,
  }
  const mergedProps = mergeProps(defaultOptions, options)
  const [isOpen, setIsOpen] = createSignal(false)
  // モーダルを開くときはopen()を呼ぶ
  const open = () => setIsOpen(true)
  // モーダルを閉じるときはclose()を呼ぶ
  const close = () => setIsOpen(false)

  const Container: ParentComponent<{
    fit?: boolean
  }> = (props) => {
    return (
      <Dialog.Root open={isOpen()}>
        <Dialog.Portal>
          <Dialog.Overlay class={overlayStyle} />
          <DialogPositioner>
            <Dialog.Content
              onEscapeKeyDown={close}
              onPointerDownOutside={mergedProps.closeOnClickOutside ? close : undefined}
              as={(
                asProps: PolymorphicCallbackProps<
                  ComponentProps<typeof Content>,
                  DialogContentOptions,
                  DialogContentRenderProps
                >,
              ) => (
                <Content fit={props.fit} {...asProps}>
                  {props.children}
                </Content>
              )}
            />
          </DialogPositioner>
        </Dialog.Portal>
      </Dialog.Root>
    )
  }

  const Header: ParentComponent = (props) => {
    return (
      <DialogHeader>
        <Dialog.Title class="whitespace-nowrap text-bold text-text-black">{props.children}</Dialog.Title>
        <Show when={mergedProps.showCloseButton}>
          <Dialog.CloseButton class={closeButtonStyle} onClick={close}>
            <MaterialSymbols>close</MaterialSymbols>
          </Dialog.CloseButton>
        </Show>
      </DialogHeader>
    )
  }

  const Body: ParentComponent<{
    fit?: boolean
  }> = (props) => {
    return (
      <Dialog.Description
        as={(
          asProps: PolymorphicCallbackProps<
            ComponentProps<typeof Description>,
            DialogDescriptionOptions,
            DialogDescriptionRenderProps
          >,
        ) => (
          <Description fit={props.fit} {...asProps}>
            {props.children}
          </Description>
        )}
      />
    )
  }

  return {
    Modal: {
      Container,
      Header,
      Body,
      Footer: ModalFooter,
    },
    open,
    close,
    isOpen,
  }
}

export default useModal
