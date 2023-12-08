import { MetaProvider, Title } from '@solidjs/meta'
import { RouteSectionProps } from '@solidjs/router'
import { type Component, ErrorBoundary } from 'solid-js'
import { Toaster } from 'solid-toast'
import { Routes } from '/@/routes'
import ErrorView from './components/layouts/ErrorView'
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
      <WithHeader>
        <ErrorBoundary fallback={(props) => <ErrorView {...props} />}>{props.children}</ErrorBoundary>
      </WithHeader>
    </MetaProvider>
  )
}

const App = () => <Routes root={Root} />

export default App
