import { ErrorBoundary, type ParentComponent } from 'solid-js'
import { retryAll } from '/@/libs/api'
import ErrorView from './ErrorView'

/**
 * PageBoundary catches the failure of a page's main data, and of any part of the page without its own
 * SectionBoundary. The header and the navigation around the page stay rendered.
 */
export const PageBoundary: ParentComponent = (props) => (
  <ErrorBoundary
    fallback={(err, reset) => {
      console.error('[page]', err)
      return (
        <ErrorView
          error={err}
          onRetry={async () => {
            await retryAll()
            reset()
          }}
        />
      )
    }}
  >
    {props.children}
  </ErrorBoundary>
)
