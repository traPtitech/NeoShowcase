import { Dialog } from '@kobalte/core'
import { createSignal, mergeProps, type ParentComponent, Show } from 'solid-js'
import { styled } from '/@/components/styled-components'
import { clsx } from '/@/libs/clsx'

const useModal = (options?: { showCloseButton?: boolean; closeOnClickOutside?: boolean }) => {
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
          <Dialog.Overlay class="fixed inset-0 animate-duration-200 animate-name-fade-hide bg-black-alpha-600 data-[expanded]:animate-name-fade-show" />
          <div class="fixed inset-0 grid place-items-center p-8">
            <Dialog.Content
              onEscapeKeyDown={close}
              onPointerDownOutside={mergedProps.closeOnClickOutside ? close : undefined}
              as="div"
              class={clsx(
                'relative flex max-h-full w-full max-w-142 animate-duration-300 animate-name-pop-hide flex-col overflow-hidden rounded-xl bg-ui-primary data-[expanded]:animate-name-pop-show',
                (props.fit ?? true) ? 'h-auto' : 'h-full',
              )}
            >
              {props.children}
            </Dialog.Content>
          </div>
        </Dialog.Portal>
      </Dialog.Root>
    )
  }

  const Header: ParentComponent = (props) => {
    return (
      <div class="dialog-header relative flex h-18 w-full shrink-0 items-center px-8 py-2">
        <Dialog.Title class="h2-medium text-text-black">{props.children}</Dialog.Title>
        <Show when={mergedProps.showCloseButton}>
          <Dialog.CloseButton
            class={clsx(
              'absolute top-6 right-6 size-6 cursor-pointer rounded border-none bg-inherit p-0 text-text-black',
              'hover:bg-transparency-primary-hover active:bg-transparency-primary-selected active:text-primary-main',
            )}
            onClick={close}
          >
            <div class="i-material-symbols:close shrink-0 text-2xl/6" />
          </Dialog.CloseButton>
        </Show>
      </div>
    )
  }

  const Body: ParentComponent<{
    fit?: boolean
  }> = (props) => {
    return (
      <Dialog.Description
        as="div"
        class={clsx(
          'description flex h-auto max-h-full w-full overflow-y-hidden px-8 py-6',
          '[.dialog-header~&]:border-ui-border [.dialog-header~&]:border-t-2',
          (props.fit ?? true) ? 'h-auto' : 'h-full',
        )}
      >
        {props.children}
      </Dialog.Description>
    )
  }

  const Footer = styled(
    'div',
    clsx(
      'flex h-18 w-full items-center justify-end gap-2 px-8 py-2',
      '[.description~&]:border-ui-border [.description~&]:border-t-2',
    ),
  )

  return {
    Modal: {
      Container,
      Header,
      Body,
      Footer,
    },
    open,
    close,
    isOpen,
  }
}

export default useModal
