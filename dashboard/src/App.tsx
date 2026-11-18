import { Title } from '@solidjs/meta'
import { Errored } from 'solid-js'
import ErrorView from '/@/components/layouts/ErrorView'
import { PageBoundary } from '/@/components/layouts/PageBoundary'
import { WithHeader } from '/@/components/layouts/WithHeader'
import { Toaster } from '/@/libs/toast'
import { Router } from '/@/router'

const App = () => (
  <>
    <Title>NeoShowcase</Title>
    <Toaster />
    {/* Last resort: only what the header and page boundaries cannot catch (the shell itself) ends up here. */}
    <Errored
      fallback={(error) => {
        console.error('[root]', error())
        return <ErrorView error={error()} />
      }}
    >
      <Router>
        {(props) => (
          <WithHeader>
            <PageBoundary>{props.children}</PageBoundary>
          </WithHeader>
        )}
      </Router>
    </Errored>
  </>
)

export default App
