import { Dialog } from '@kobalte/core/dialog'
import { useLocation } from '@solidjs/router'
import { createEffect, createMemo, createSignal, merge, type ParentComponent, Show } from 'solid-js'
import { styled } from '/@/components/styled-components'

const useModal = (options?: { showCloseButton?: boolean; closeOnClickOutside?: boolean }) => {
  const config = merge({ showCloseButton: false, closeOnClickOutside: true }, options)
  const [isOpen, setIsOpen] = createSignal(false)
  // モーダルを開くときはopen()を呼ぶ
  const open = () => setIsOpen(true)
  // モーダルを閉じるときはclose()を呼ぶ
  const close = () => {
    setIsOpen(false)
  }
  // 遷移したら閉じる (ルートのコンポーネントはパラメータが変わっても使い回されるため、
  // 開いたままだと別のリソースに対する操作になってしまう)
  const location = useLocation()
  const pathname = createMemo(() => location.pathname)
  createEffect(
    pathname,
    () => {
      close()
    },
    { defer: true },
  )

  const Container: ParentComponent<{
    fit?: boolean
  }> = (props) => {
    return (
      // ESC・外側クリック・閉じるボタンはすべて onOpenChange 経由で閉じる
      <Dialog open={isOpen()} onOpenChange={setIsOpen}>
        <Dialog.Portal>
          <Dialog.Overlay class="animate-duration-200 animate-name-fade-hide bg-black-alpha-600 data-[expanded]:animate-name-fade-show fixed inset-0" />
          <div class="fixed inset-0 grid place-items-center p-8">
            <Dialog.Content
              onPointerDownOutside={(e) => {
                if (!config.closeOnClickOutside) e.preventDefault()
              }}
              as="div"
              class={[
                'animate-duration-300 animate-name-pop-hide bg-ui-primary data-[expanded]:animate-name-pop-show relative flex max-h-full w-full max-w-142 flex-col overflow-hidden rounded-xl',
                (props.fit ?? true) ? 'h-auto' : 'h-full',
              ]}
            >
              {props.children}
            </Dialog.Content>
          </div>
        </Dialog.Portal>
      </Dialog>
    )
  }

  const Header: ParentComponent = (props) => {
    return (
      <div class="dialog-header relative flex h-18 w-full shrink-0 items-center px-8 py-2">
        <Dialog.Title class="h2-medium text-text-black">{props.children}</Dialog.Title>
        <Show when={config.showCloseButton}>
          <Dialog.CloseButton
            class={[
              'text-text-black absolute top-6 right-6 size-6 cursor-pointer rounded border-none bg-inherit p-0',
              'hover:bg-transparency-primary-hover active:bg-transparency-primary-selected active:text-primary-main',
            ]}
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
        class={[
          'description flex h-auto max-h-full w-full overflow-y-hidden px-8 py-6',
          '[.dialog-header~&]:border-ui-border [.dialog-header~&]:border-t-2',
          (props.fit ?? true) ? 'h-auto' : 'h-full',
        ]}
      >
        {props.children}
      </Dialog.Description>
    )
  }

  const Footer = styled(
    'div',
    'flex h-18 w-full items-center justify-end gap-2 px-8 py-2 [.description~&]:border-ui-border [.description~&]:border-t-2',
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
