/* @refresh reload */
import { render } from 'solid-js/web'
import './global-style'

import App from './App'

render(() => <App />, document.getElementById('root') as HTMLElement)
