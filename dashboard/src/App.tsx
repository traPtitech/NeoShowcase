import { MetaProvider, Title } from '@solidjs/meta'
import { type RouteSectionProps, revalidate } from '@solidjs/router'
import { type Component, ErrorBoundary } from 'solid-js'
import { Toaster } from 'solid-toast'
import { Routes } from '/@/routes'
import ErrorView from './components/layouts/ErrorView'
import { PageBoundary } from './components/layouts/PageBoundary'
import { WithHeader } from './components/layouts/WithHeader'

const Root: Component<RouteSectionProps> = (props) => {
  return (
    <MetaProvider>
      <Title>NeoShowcase</Title>
      <Toaster
        toastOptions={{
          duration: 10000,
          position: 'bottom-left',
        }}
      />
      <ErrorBoundary
        fallback={(err) => {
          console.error('[root]', err)
          return <ErrorView error={err} />
        }}
      >
        <WithHeader>
          <PageBoundary onRetry={() => revalidate()}>{props.children}</PageBoundary>
        </WithHeader>
      </ErrorBoundary>
    </MetaProvider>
  )
}

const App = () => <Routes root={Root} />

export default App
