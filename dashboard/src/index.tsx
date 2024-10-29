/* @refresh reload */
import { render } from 'solid-js/web'

import App from './App'
import '@unocss/reset/tailwind-compat.css'
import './global.css'
import './animate.css'
import 'virtual:uno.css'

render(() => <App />, document.getElementById('root') as HTMLElement)
