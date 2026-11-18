/* @refresh reload */
import { render } from '@solidjs/web'
import App from './App'
import '@unocss/reset/tailwind-compat.css'
import './global.css'
import './animate.css'
import 'virtual:uno.css'

render(() => <App />, document.getElementById('root')!)
