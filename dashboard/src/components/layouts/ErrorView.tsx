import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'

const ErrorView: Component<{
  error: unknown
}> = (props) => {
  const handleReload = () => {
    window.location.reload()
  }

  return (
    <div class="flex w-full flex-col items-center justify-center gap-4">
      <MaterialSymbols fill displaySize={64} class="text-accent-error">
        error
      </MaterialSymbols>
      <h2 class="h2-bold text-accent-error">An error has occurred</h2>
      <Show when={props.error instanceof Error}>
        <p class="caption-medium text-text-grey">{(props.error as Error).message}</p>
      </Show>
      <div class="flex flex-col gap-2">
        <A href="/">
          <Button size="medium" variants="border" leftIcon={<MaterialSymbols>arrow_back</MaterialSymbols>}>
            Back to Home
          </Button>
        </A>
        <Button
          onClick={handleReload}
          size="medium"
          variants="border"
          leftIcon={<MaterialSymbols>refresh</MaterialSymbols>}
        >
          Reload Page
        </Button>
      </div>
    </div>
  )
}

export default ErrorView
