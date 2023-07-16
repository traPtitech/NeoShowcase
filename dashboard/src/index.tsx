/* @refresh reload */
import { render } from 'solid-js/web'
import { TippyOptions } from 'solid-tippy'
import App from './App'
import './global-style'

render(() => <App />, document.getElementById('root') as HTMLElement)
