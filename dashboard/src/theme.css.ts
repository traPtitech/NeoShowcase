import { createTheme } from '@vanilla-extract/css'

export const [themeClass, vars] = createTheme({
  color: {
    bg: {
      black1: '#35495E',
      white1: '#FFFFFF',
      white2: '#F2F2F2',
      white3: '#ECECEC',
      white4: '#E8E8E8',
    },
    text: {
      black1: '#2F2D2A',
      black2: '#35495E',
      black3: '#9CA3AF',
      black4: '#ADB5BC',
      white1: '#FAFAFA',
    },
    icon: {
      error: '#EB5E28',
      pending: '#FFCE4F',
      success1: '#41B883',
      success2: '#68B3C8',
    },
  },
})
