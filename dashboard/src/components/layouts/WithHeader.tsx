import type { ParentComponent } from 'solid-js'
import { Header } from '../templates/Header'

export const WithHeader: ParentComponent = (props) => {
  return (
    <div class="grid h-full w-full grid-cols-1 grid-rows-[auto_1fr]">
      <Header />
      <div class="h-full w-full overflow-y-auto">{props.children}</div>
    </div>
  )
}
