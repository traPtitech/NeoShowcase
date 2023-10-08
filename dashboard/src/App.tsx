import Routes from '/@/routes'
import { Router } from '@solidjs/router'
import type { Component } from 'solid-js'
import { Toaster } from 'solid-toast'
import { WithHeader } from './components/layouts/WithHeader'

const App: Component = () => {
  return (
    <>
      <Toaster
        toastOptions={{
          duration: 10000,
          position: 'bottom-right',
        }}
      />
      <Router>
        <WithHeader>
          <Routes />
        </WithHeader>
      </Router>
    </>
  )
}

export default App
