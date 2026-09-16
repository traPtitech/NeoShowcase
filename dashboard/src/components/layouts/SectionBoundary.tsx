import { type Component, ErrorBoundary, type ParentComponent, Show } from 'solid-js'
import { retryAll } from '/@/libs/api'
import { Button } from '../UI/Button'
import { DataTable } from './DataTable'

type ErrorViewProps = {
  title: string
  error: unknown
  onRetry: () => void
}

const errorMessage = (error: unknown) => (error instanceof Error ? error.message : undefined)

const SectionErrorView: Component<ErrorViewProps> = (props) => (
  <DataTable.Container>
    <DataTable.Title>{props.title}</DataTable.Title>
    <div class="flex w-full items-center gap-4 rounded-lg border border-ui-border bg-ui-primary px-5 py-4">
      <div class="i-material-symbols:error shrink-0 text-2xl/6 text-accent-error" />
      <div class="flex w-full flex-col overflow-x-hidden">
        <span class="h4-medium text-accent-error">Failed to load</span>
        <Show when={errorMessage(props.error)}>
          {(message) => <span class="overflow-wrap-anywhere text-regular text-text-grey">{message()}</span>}
        </Show>
      </div>
      <Button
        onClick={props.onRetry}
        size="medium"
        variants="border"
        leftIcon={<div class="i-material-symbols:refresh shrink-0 text-2xl/6" />}
      >
        Retry
      </Button>
    </div>
  </DataTable.Container>
)

const CompactErrorView: Component<ErrorViewProps> = (props) => (
  <Button
    onClick={props.onRetry}
    size="medium"
    variants="text"
    leftIcon={<div class="i-material-symbols:error shrink-0 text-2xl/6 text-accent-error" />}
    tooltip={{
      style: 'left',
      props: {
        content: (
          <>
            <div>Failed to load {props.title}</div>
            <Show when={errorMessage(props.error)}>{(message) => <div>{message()}</div>}</Show>
          </>
        ),
      },
    }}
  >
    Retry
  </Button>
)

/**
 * SectionBoundary confines the failure of one independently fetched piece of data to the block that shows it,
 * leaving sibling blocks rendered. On failure it replaces the block with its heading, the message and a Retry
 * button, so the page keeps its shape.
 *
 * It catches failures only. The block places its own `Suspense` around the part that reads the data, which
 * keeps the heading out of the pending state and leaves the choice of where the placeholder goes at the call
 * site, where the size of the slot is known. Sections share `SectionSkeleton` as that placeholder.
 *
 * Retry goes through `retryAll`, so placing a boundary never involves naming the resources beneath it.
 */
export const SectionBoundary: ParentComponent<{
  /** Name of the data this boundary reads. The 'section' variant also renders it as the section heading. */
  title: string
  /**
   * How the boundary presents itself. 'section' (the default) owns a heading and a full-width body.
   * 'compact' fits a slot the size of a single control, such as the header, and moves the detail into a tooltip.
   */
  variant?: 'section' | 'compact'
}> = (props) => {
  const isCompact = () => props.variant === 'compact'

  return (
    <ErrorBoundary
      fallback={(err, reset) => {
        console.error(`[section: ${props.title}]`, err)
        const onRetry = async () => {
          await retryAll()
          reset()
        }
        return isCompact() ? (
          <CompactErrorView title={props.title} error={err} onRetry={onRetry} />
        ) : (
          <SectionErrorView title={props.title} error={err} onRetry={onRetry} />
        )
      }}
    >
      {props.children}
    </ErrorBoundary>
  )
}
