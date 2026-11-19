import { useNavigate } from '@solidjs/router'

/** `/` has no page of its own; send users to the app list. */
export default () => {
  useNavigate()('/apps', { replace: true })
  return null
}
