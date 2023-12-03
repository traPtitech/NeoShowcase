import { MetaProvider, Title } from '@solidjs/meta'
import { Router } from '@solidjs/router'
import { type Component, ErrorBoundary } from 'solid-js'
import { Toaster } from 'solid-toast'
import Routes from '/@/routes'
import ErrorView from './components/layouts/ErrorView'
import { WithHeader } from './components/layouts/WithHeader'

const App: Component = () => {
  return (
    <MetaProvider>
      <Title>NeoShowcase</Title>
      <Toaster
        toastOptions={{
          duration: 10000,
          position: 'bottom-left',
        }}
      />
      <Router>
        <WithHeader>
          <ErrorBoundary fallback={(props) => <ErrorView {...props} />}>
            <Routes />
          </ErrorBoundary>
        </WithHeader>
      </Router>
    </MetaProvider>
  )
}

export default App
