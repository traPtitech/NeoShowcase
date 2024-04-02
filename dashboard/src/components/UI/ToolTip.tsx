import { styled } from '@macaron-css/solid'
import { type FlowComponent, type JSX, children, mergeProps, onMount, splitProps } from 'solid-js'
import { type TippyOptions, tippy } from 'solid-tippy'
import type { Props } from 'tippy.js'
import 'tippy.js/animations/shift-away-subtle.css'
import 'tippy.js/dist/tippy.css'

const TooltipContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
  },
  variants: {
    align: {
      left: {
        alignItems: 'flex-start',
      },
      center: {
        alignItems: 'center',
      },
    },
  },
})

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
            content: (
              <TooltipContainer align={addedProps.style}>{propsWithDefaults.props?.content}</TooltipContainer>
            ) as Element,
          },
        }))
      }
    }
  })

  return <>{c()}</>
}
