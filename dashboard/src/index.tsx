/* @refresh reload */
import { render } from 'solid-js/web'
import './global.css'
import './font.css'

import App from './App'

render(() => <App />, document.getElementById('root') as HTMLElement)
