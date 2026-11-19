import type { JSX } from '@solidjs/web'
import { Errored, Loading, type ParentComponent } from 'solid-js'
import { retryAll } from '/@/libs/session'
import ErrorView from './ErrorView'

/**
 * PageBoundary catches the failure of a page's main data, and of any part of the page without its own
 * SectionBoundary, and holds the page's first-load placeholder. The header and the navigation around the page
 * stay rendered. Retry re-asks the data layer (`retryAll`) before resetting.
 */
export const PageBoundary: ParentComponent<{ fallback?: JSX.Element }> = (props) => (
  <Errored
    fallback={(err, reset) => {
      console.error('[page]', err())
      return (
        <ErrorView
          error={err()}
          reset={async () => {
            await retryAll()
            reset()
          }}
        />
      )
    }}
  >
    <Loading fallback={props.fallback}>{props.children}</Loading>
  </Errored>
)
