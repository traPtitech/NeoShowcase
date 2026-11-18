import { type Component, Show } from 'solid-js'
import { Button } from '../UI/Button'

const ErrorView: Component<{
  error: unknown
  /** Re-runs the failed subtree (e.g. `Errored`'s `reset`). Shows a "Retry" button when given. */
  reset?: () => void
}> = (props) => {
  return (
    <div class="flex w-full flex-col items-center justify-center gap-4">
      <div class="i-material-symbols:error text-16/16 text-accent-error shrink-0" />
      <h2 class="h2-bold text-accent-error">An error has occurred</h2>
      <Show when={props.error instanceof Error}>
        <p class="caption-medium text-text-grey">{(props.error as Error).message}</p>
      </Show>
      <div class="flex flex-col gap-2">
        <Show when={props.reset}>
          {(reset) => (
            <Button
              onClick={() => reset()()}
              size="medium"
              variants="border"
              leftIcon={<div class="i-material-symbols:refresh shrink-0 text-2xl/6" />}
            >
              Retry
            </Button>
          )}
        </Show>
        <a href="/">
          <Button
            size="medium"
            variants="border"
            leftIcon={<div class="i-material-symbols:arrow-back shrink-0 text-2xl/6" />}
          >
            Back to Home
          </Button>
        </a>
        <Button
          onClick={() => window.location.reload()}
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
