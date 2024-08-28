import { Skeleton as KSkeleton } from '@kobalte/core'
import { type Component, mergeProps } from 'solid-js'
import { clsx } from '/@/libs/clsx'
import styles from './Skeleton.module.css'

type SkeletonProps = Parameters<typeof KSkeleton.Root>[0]

const defaultProps: SkeletonProps = {
  radius: 999,
  width: -1, // for `width: auto`
}

const Skeleton: Component<SkeletonProps> = (props) => {
  const mergedProps = mergeProps(defaultProps, props)

  return (
    <KSkeleton.Root
      {...mergedProps}
      class={clsx(
        'relative h-auto w-auto shrink-0 opacity-20',
        'data-[visible]:overflow-hidden',
        "data-[visible]:after:absolute data-[visible]:after:inset-0 data-[visible]:after:bg-[currentColor] data-[visible]:after:bg-[linear-gradient(90deg,transparent,#fff4,transparent)] data-[visible]:after:content-['']",
        "data-[visible]:before:absolute data-[visible]:before:inset-0 data-[visible]:before:bg-[currentColor] data-[visible]:before:content-['']",
        styles['skeleton-animation'],
      )}
    >
      {props.children}
    </KSkeleton.Root>
  )
}

export default Skeleton
