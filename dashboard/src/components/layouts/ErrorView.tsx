import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import { Button } from '../UI/Button'

const ErrorView: Component<{
  error: unknown
  /** Retries what failed. When omitted, only the navigation and reload buttons are shown. */
  onRetry?: () => void
}> = (props) => {
  const handleReload = () => {
    window.location.reload()
  }

  return (
    <div class="flex w-full flex-col items-center justify-center gap-4">
      <div class="i-material-symbols:error shrink-0 text-16/16 text-accent-error" />
      <h2 class="h2-bold text-accent-error">An error has occurred</h2>
      <Show when={props.error instanceof Error}>
        <p class="caption-medium text-text-grey">{(props.error as Error).message}</p>
      </Show>
      <div class="flex flex-col gap-2">
        <Show when={props.onRetry}>
          <Button
            onClick={() => props.onRetry?.()}
            size="medium"
            variants="border"
            leftIcon={<div class="i-material-symbols:refresh shrink-0 text-2xl/6" />}
          >
            Retry
          </Button>
        </Show>
        <A href="/">
          <Button
            size="medium"
            variants="border"
            leftIcon={<div class="i-material-symbols:arrow-back shrink-0 text-2xl/6" />}
          >
            Back to Home
          </Button>
        </A>
        <Button
          onClick={handleReload}
          size="medium"
          variants="border"
          leftIcon={<div class="i-material-symbols:refresh shrink-0 text-2xl/6" />}
        >
          Reload Page
        </Button>
      </div>
    </div>
  )
}

export default ErrorView
