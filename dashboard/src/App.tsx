import { Router } from '@solidjs/router'
import type { Component } from 'solid-js'
import Routes from '/@/routes'

const App: Component = () => {
  return (
    <>
      <Router>
        <Routes />
      </Router>
    </>
  )
}

export default App
