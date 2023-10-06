import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'

export const SettingsContainer = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',

    color: colorVars.semantic.text.black,
    ...textVars.h2.medium,
  },
})
