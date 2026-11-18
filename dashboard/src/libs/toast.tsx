import { Toast, toaster } from '@kobalte/core/toast'
import { type JSX, Portal } from '@solidjs/web'
import type { Component } from 'solid-js'

/**
 * Thin wrapper around Kobalte's toaster.
 * Call sites must use this module instead of importing a toast library directly,
 * so the underlying implementation can be swapped without touching them.
 */

type Variant = 'success' | 'error' | 'loading'

const variantIcon: Record<Variant, string> = {
  success: 'i-material-symbols:check-circle text-accent-success',
  error: 'i-material-symbols:error text-accent-error',
  loading: 'i-material-symbols:progress-activity animate-spin text-text-grey',
}

const ToastItem: Component<{ toastId: number; variant: Variant; message: JSX.Element; persistent?: boolean }> = (
  props,
) => (
  <Toast
    toastId={props.toastId}
    persistent={props.persistent}
    class="animate-duration-200 bg-ui-primary text-regular text-text-black shadow-default data-[closed]:animate-name-fade-hide data-[opened]:animate-name-fade-show pointer-events-auto flex w-90 max-w-[calc(100vw-32px)] items-start gap-2 rounded-lg px-4 py-3"
  >
    <div class={['mt-0.5 shrink-0 text-xl/5', variantIcon[props.variant]]} />
    <Toast.Description class="overflow-wrap-anywhere min-w-0 grow whitespace-pre-wrap">
      {props.message}
    </Toast.Description>
    <Toast.CloseButton class="text-text-grey hover:bg-transparency-primary-hover grid size-5 shrink-0 cursor-pointer place-items-center rounded border-none bg-transparent p-0">
      <div class="i-material-symbols:close text-base" />
    </Toast.CloseButton>
  </Toast>
)

const show = (variant: Variant, message: JSX.Element) =>
  toaster.show((p) => <ToastItem toastId={p.toastId} variant={variant} message={message} />)

export const toast = {
  success: (message: JSX.Element) => show('success', message),
  error: (message: JSX.Element) => show('error', message),
  /** Shows a loading toast while `promise` is pending, then success / error. Resolves/rejects with the promise. */
  promise: <T,>(promise: Promise<T>, messages: { loading: JSX.Element; success: JSX.Element; error: JSX.Element }) => {
    toaster.promise(promise, (p) => (
      <ToastItem
        toastId={p.toastId}
        variant={p.state === 'pending' ? 'loading' : p.state === 'fulfilled' ? 'success' : 'error'}
        message={p.state === 'pending' ? messages.loading : p.state === 'fulfilled' ? messages.success : messages.error}
        persistent={p.state === 'pending'}
      />
    ))
    return promise
  },
  dismiss: (id: number) => toaster.dismiss(id),
}

/** Mount once at the app root. */
export const Toaster: Component = () => (
  <Portal>
    <Toast.Region duration={10000} swipeDirection="left" limit={5}>
      <Toast.List class="pointer-events-none fixed bottom-0 left-0 z-100 m-0 flex list-none flex-col gap-2 p-4" />
    </Toast.Region>
  </Portal>
)
