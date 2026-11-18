import { type Component, Errored, Loading, type ParentComponent, Show } from 'solid-js'
import { retryAll } from '/@/libs/session'
import SectionSkeleton from '../templates/SectionSkeleton'
import { Button } from '../UI/Button'
import Skeleton from '../UI/Skeleton'
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
    <div class="border-ui-border bg-ui-primary flex w-full items-center gap-4 rounded-lg border px-5 py-4">
      <div class="i-material-symbols:error text-accent-error shrink-0 text-2xl/6" />
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
    leftIcon={<div class="i-material-symbols:error text-accent-error shrink-0 text-2xl/6" />}
    tooltip={{
      style: 'left',
      content: (
        <>
          <div>Failed to load {props.title}</div>
          <Show when={errorMessage(props.error)}>{(message) => <div>{message()}</div>}</Show>
        </>
      ),
    }}
  >
    Retry
  </Button>
)

/**
 * SectionBoundary confines the failure of one independently fetched piece of data to the block that shows it,
 * leaving sibling blocks rendered. On failure it replaces the block with its heading, the message and a Retry
 * button, so the page keeps its shape. While the block is loading for the first time it shows `SectionSkeleton`
 * (or a header-sized one for 'compact'), with the heading already in place.
 *
 * Retry goes through `retryAll`, so placing a boundary never involves naming the resources beneath it.
 */
export const SectionBoundary: ParentComponent<{
  /** Name of the data this boundary reads. The 'section' variant renders it as the section heading. */
  title: string
  /**
   * How the boundary presents itself. 'section' (the default) owns a heading and a full-width body.
   * 'compact' fits a slot the size of a single control, such as the header, and moves the detail into a tooltip.
   */
  variant?: 'section' | 'compact'
  /**
   * 'section' only: the children render their own heading (e.g. one with actions that depend on the data), so
   * the boundary adds none while healthy. The failure view still shows `title`.
   */
  bare?: boolean
}> = (props) => (
  <Errored
    fallback={(err, reset) => {
      console.error(`[section: ${props.title}]`, err())
      const onRetry = async () => {
        await retryAll()
        reset()
      }
      return props.variant === 'compact' ? (
        <CompactErrorView title={props.title} error={err()} onRetry={onRetry} />
      ) : (
        <SectionErrorView title={props.title} error={err()} onRetry={onRetry} />
      )
    }}
  >
    {props.variant === 'compact' ? (
      // The header slot is sized by its content, so the placeholder cannot take a percentage width.
      <Loading fallback={<Skeleton width={120} height={32} />}>{props.children}</Loading>
    ) : props.bare ? (
      <Loading fallback={<SectionSkeleton />}>{props.children}</Loading>
    ) : (
      <DataTable.Container>
        <DataTable.Title>{props.title}</DataTable.Title>
        <Loading fallback={<SectionSkeleton />}>{props.children}</Loading>
      </DataTable.Container>
    )}
  </Errored>
)
