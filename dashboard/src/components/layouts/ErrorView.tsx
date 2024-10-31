import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import { Button } from '../UI/Button'

const ErrorView: Component<{
  error: unknown
}> = (props) => {
  const handleReload = () => {
    window.location.reload()
  }

  return (
    <div class="flex w-full flex-col items-center justify-center gap-4">
      <span class="i-material-symbols:error text-16/16 text-accent-error" />
      <h2 class="h2-bold text-accent-error">An error has occurred</h2>
      <Show when={props.error instanceof Error}>
        <p class="caption-medium text-text-grey">{(props.error as Error).message}</p>
      </Show>
      <div class="flex flex-col gap-2">
        <A href="/">
          <Button size="medium" variants="border" leftIcon={<span class="i-material-symbols:arrow-back text-2xl/6" />}>
            Back to Home
          </Button>
        </A>
        <Button
          onClick={handleReload}
          size="medium"
          variants="border"
          leftIcon={<span class="i-material-symbols:refresh text-2xl/6" />}
        >
          Reload Page
        </Button>
      </div>
    </div>
  )
}

export default ErrorView
