import Routes from '/@/routes'
import { Router } from '@solidjs/router'
import type { Component } from 'solid-js'
import { Toaster } from 'solid-toast'

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
        <Routes />
      </Router>
    </>
  )
}

export default App
