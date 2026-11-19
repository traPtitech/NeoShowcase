import { Skeleton as KSkeleton } from '@kobalte/core/skeleton'
import { type Component, merge } from 'solid-js'
import styles from './Skeleton.module.css'

type SkeletonProps = Parameters<typeof KSkeleton>[0]

const defaultProps: SkeletonProps = {
  radius: 999,
  width: -1, // for `width: auto`
}

const Skeleton: Component<SkeletonProps> = (props) => (
  <KSkeleton
    {...merge(defaultProps, props)}
    class={[
      'relative h-auto w-auto shrink-0 !rounded-lg opacity-20',
      'data-[visible]:overflow-hidden',
      "data-[visible]:after:absolute data-[visible]:after:inset-0 data-[visible]:after:bg-[currentColor] data-[visible]:after:bg-[linear-gradient(90deg,transparent,#fff4,transparent)] data-[visible]:after:content-['']",
      "data-[visible]:before:absolute data-[visible]:before:inset-0 data-[visible]:before:bg-[currentColor] data-[visible]:before:content-['']",
      styles['skeleton-animation'],
    ]}
  >
    {props.children}
  </KSkeleton>
)

export default Skeleton
