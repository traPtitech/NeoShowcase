import { Tooltip } from '@kobalte/core/tooltip'
import type { JSX } from '@solidjs/web'
import { children, createMemo, type ParentComponent, Show } from 'solid-js'

export type TooltipProps = {
  /** Tooltip body. When `undefined`, no tooltip is shown. */
  content?: JSX.Element
  /** Disable the tooltip while keeping the trigger rendered. */
  disabled?: boolean
  /** Text alignment inside the tooltip. @default "center" */
  style?: 'left' | 'center'
  placement?: 'top' | 'bottom' | 'left' | 'right'
  /** Extra classes for the wrapper element that acts as the hover target. */
  class?: string
}

/**
 * Wraps `children` in an inline hover target and shows `content` in a floating tooltip.
 * Works on disabled buttons too, because the hover target is the wrapper, not the button.
 */
export const ToolTip: ParentComponent<TooltipProps> = (props) => {
  // Resolve children once so the same nodes are reused when the tooltip toggles
  // (an input inside keeps its state instead of being re-created).
  const resolved = children(() => props.children)
  // Read in a tracked scope: `content`/`disabled` often depend on async data (e.g. permissions), and
  // Kobalte 2 (alpha) evaluates its root props untracked on construction, which throws while pending.
  const enabled = createMemo(() => props.content !== undefined && !props.disabled)
  const wrapperClass = () => ['inline-flex max-w-full', props.class]

  return (
    <Show when={enabled()} fallback={<span class={wrapperClass()}>{resolved()}</span>}>
      <Tooltip placement={props.placement ?? 'top'} openDelay={150} closeDelay={0} gutter={6}>
        <Tooltip.Trigger as="span" class={wrapperClass()}>
          {resolved()}
        </Tooltip.Trigger>
        <Tooltip.Portal>
          <Tooltip.Content
            class={[
              'text-3.5 text-text-white shadow-default z-50 flex max-w-250 flex-col rounded bg-gray-800 px-2.5 py-1.5 leading-5',
              'animate-duration-150 animate-name-fade-hide data-[expanded]:animate-name-fade-show origin-[--kb-tooltip-content-transform-origin]',
              (props.style ?? 'center') === 'center' ? 'items-center text-center' : 'items-start text-left',
            ]}
          >
            <Tooltip.Arrow size={10} />
            {props.content}
          </Tooltip.Content>
        </Tooltip.Portal>
      </Tooltip>
    </Show>
  )
}
