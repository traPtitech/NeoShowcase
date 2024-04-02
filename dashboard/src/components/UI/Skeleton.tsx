import { Skeleton as KSkeleton } from '@kobalte/core'
import { keyframes, style } from '@macaron-css/core'
import { type Component, mergeProps } from 'solid-js'

const skeletonAnimation = keyframes({
  from: {
    transform: 'translateX(-100%)',
  },
  to: {
    transform: 'translateX(100%)',
  },
})
const skeletonClass = style({
  position: 'relative',
  width: 'auto',
  height: 'auto',
  flexShrink: '0',
  opacity: '0.2',

  selectors: {
    "&[data-visible='true']": {
      overflow: 'hidden',
    },
    "&[data-visible='true']::after": {
      position: 'absolute',
      content: '""',
      inset: '0',
      backgroundColor: 'currentcolor',
      backgroundImage: 'linear-gradient(90deg, transparent, #fff4, transparent)',
    },
    "&[data-visible='true']::before": {
      position: 'absolute',
      content: '""',
      inset: '0',
      backgroundColor: 'currentcolor',
    },
    "&[data-animate='true']::after": {
      animation: `${skeletonAnimation} 1.5s linear infinite`,
    },
  },
})

type SkeletonProps = Parameters<typeof KSkeleton.Root>[0]

const defaultProps: SkeletonProps = {
  radius: 999,
  width: -1, // for `width: auto`
}

const Skeleton: Component<SkeletonProps> = (props) => {
  const mergedProps = mergeProps(defaultProps, props)

  return (
    <KSkeleton.Root {...mergedProps} class={skeletonClass}>
      {props.children}
    </KSkeleton.Root>
  )
}

export default Skeleton
