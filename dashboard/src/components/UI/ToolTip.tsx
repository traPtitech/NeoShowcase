import { children, type FlowComponent, type JSX, mergeProps, onMount, splitProps } from 'solid-js'
import { type TippyOptions, tippy } from 'solid-tippy'
import type { Props } from 'tippy.js'
import 'tippy.js/animations/shift-away-subtle.css'
import 'tippy.js/dist/tippy.css'
import { clsx } from '/@/libs/clsx'

export type TooltipProps = Omit<TippyOptions, 'props'> & {
  props?: Partial<
    Omit<Props, 'content'> & {
      content?: JSX.Element
    }
  >
} & {
  /**
   * @default "center"
   */
  style?: 'left' | 'center'
}

export const ToolTip: FlowComponent<TooltipProps> = (props) => {
  const defaultOptions: TooltipProps = {
    style: 'center',
    hidden: true,
    props: {
      allowHTML: true,
      maxWidth: 1000,
      animation: 'shift-away-subtle',
    },
    disabled: props.props?.content === undefined,
  }
  const propsWithDefaults = mergeProps(defaultOptions, props)
  const [addedProps, tippyProps] = splitProps(propsWithDefaults, ['style', 'children'])
  const c = children(() => props.children)

  onMount(() => {
    for (const child of c.toArray()) {
      if (child instanceof Element) {
        tippy(child, () => ({
          ...tippyProps,
          props: {
            ...tippyProps.props,
            content: (
              <div
                class={clsx(
                  'flex flex-col',
                  addedProps.style === 'center' && 'items-center',
                  addedProps.style === 'left' && 'items-start',
                )}
              >
                {propsWithDefaults.props?.content}
              </div>
            ) as Element,
          },
        }))
      }
    }
  })

  return <>{c()}</>
}
