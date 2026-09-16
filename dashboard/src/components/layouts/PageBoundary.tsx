import { ErrorBoundary, type ParentComponent } from 'solid-js'
import ErrorView from './ErrorView'

/**
 * PageBoundary catches the failure of a page's main data, and of any part of the page without its own
 * SectionBoundary. The header and the navigation around the page stay rendered.
 */
export const PageBoundary: ParentComponent<{
  /**
   * Refetches the data sources that live outside this boundary, run before resetting it.
   * Omit it only when everything this boundary reads is created inside it.
   */
  onRetry?: () => unknown
}> = (props) => (
  <ErrorBoundary
    fallback={(err, reset) => {
      console.error('[page]', err)
      return (
        <ErrorView
          error={err}
          onRetry={async () => {
            await props.onRetry?.()
            reset()
          }}
        />
      )
    }}
  >
    {props.children}
  </ErrorBoundary>
)
