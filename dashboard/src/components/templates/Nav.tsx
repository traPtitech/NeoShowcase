import { A } from '@solidjs/router'
import { type JSX, type ParentComponent, Show } from 'solid-js'
import { Button } from '/@/components/UI/Button'

export interface Props {
  title: string
  backTo?: string
  backToTitle?: string
  icon?: JSX.Element
  action?: JSX.Element
}

export const Nav: ParentComponent<Props> = (props) => {
  return (
    <div class="flex w-full gap-2 overflow-x-hidden p-8 pr-[max(calc(50%-500px),32px)] max-md:px-4 max-md:py-8">
      <Show when={props.backTo} fallback={<div />}>
        {(nonNullBackTo) => (
          <A href={nonNullBackTo()}>
            <Button variants="text" size="medium" leftIcon={<div class="i-material-symbols:arrow-back text-2xl/6" />}>
              <div class="max-md:hidden">{props.backToTitle}</div>
            </Button>
          </A>
        )}
      </Show>
      <div class="w-full overflow-x-clip">
        <div class="sticky left-[calc(75%-250px)] h-auto w-full max-w-250 overflow-x-hidden">
          <div class="flex items-center gap-2 overflow-x-hidden">
            <Show when={props.icon}>{props.icon}</Show>
            <div class="h1-medium truncate">{props.title}</div>
            <Show when={props.action}>{props.action}</Show>
          </div>
          {props.children}
        </div>
      </div>
    </div>
  )
}
