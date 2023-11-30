import { styled } from '@macaron-css/solid'

const SuspenseContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',

    opacity: 1,
    transition: 'opacity 0.2s ease-in-out',
  },
  variants: {
    isPending: {
      true: {
        opacity: 0.5,
        pointerEvents: 'none',
      },
    },
  },
})

export default SuspenseContainer
