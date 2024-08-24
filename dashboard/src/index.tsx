/* @refresh reload */
import { render } from 'solid-js/web'

import App from './App'
import './global-style'
import 'virtual:uno.css'

render(() => <App />, document.getElementById('root') as HTMLElement)
