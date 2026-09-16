import type { Component } from 'solid-js'
import Skeleton from '../UI/Skeleton'

/**
 * SectionSkeleton is the single loading display shared by every section body, so that adding a section does
 * not mean designing a placeholder for it. It fills the width it is given and holds a fixed height.
 */
const SectionSkeleton: Component = () => <Skeleton width={Number.NaN} height={80} />

export default SectionSkeleton
