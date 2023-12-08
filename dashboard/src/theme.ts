import { createGlobalTheme } from '@macaron-css/core'

export const colorVars = createGlobalTheme(':root', {
  semantic: {
    primary: {
      white: '#FFFFFF',
      main: '#005BAC',
    },
    accent: {
      error: '#F25151',
      warn: '#F1B61E',
      success: '#20BD77',
    },
    transparent: {
      primaryHover: 'rgba(0, 91, 172, 0.06)',
      primarySelected: 'rgba(0, 91, 172, 0.10)',
      successHover: 'rgba(32, 189, 119, 0.06)',
      successSelected: 'rgba(32, 189, 119, 0.10)',
      warnHover: 'rgba(241, 182, 30, 0.06)',
      warnSelected: 'rgba(241, 182, 30, 0.10)',
      errorHover: 'rgba(242, 81, 81, 0.06)',
      errorSelected: 'rgba(242, 81, 81, 0.10)',
    },
    text: {
      black: '#2F3438',
      white: '#FFFFFF',
      grey: '#606A71',
      link: '#005BAC',
      disabled: '#B9BEC1',
    },
    ui: {
      border: '#CED6DB',
      background: '#F9F9F9',
      primary: '#FFFFFF',
      secondary: '#F0F2F5',
      tertiary: '#E2E5E9',
    },
  },
  primitive: {
    blackAlpha: {
      50: 'rgba(0, 0, 0, 0.04)',
      100: 'rgba(0, 0, 0, 0.06)',
      200: 'rgba(0, 0, 0, 0.08)',
      300: 'rgba(0, 0, 0, 0.16)',
      600: 'rgba(0, 0, 0, 0.48)',
    },
    gray: {
      700: '#2D3748',
      800: '#1A202C',
      900: '#171923',
    },
    blue: {
      500: '#3182CE',
    },
  },
})

export const textVars = createGlobalTheme(':root', {
  h1: {
    regular: {
      fontSize: '1.75rem',
      fontWeight: '400',
      lineHeight: '1.5',
    },
    medium: {
      fontSize: '1.75rem',
      fontWeight: '500',
      lineHeight: '1.5',
    },
    bold: {
      fontSize: '1.75rem',
      fontWeight: '700',
      lineHeight: '1.5',
    },
  },
  h2: {
    regular: {
      fontSize: '1.5rem',
      fontWeight: '400',
      lineHeight: '1.5',
    },
    medium: {
      fontSize: '1.5rem',
      fontWeight: '500',
      lineHeight: '1.5',
    },
    bold: {
      fontSize: '1.5rem',
      fontWeight: '700',
      lineHeight: '1.5',
    },
  },
  h3: {
    regular: {
      fontSize: '1.25rem',
      fontWeight: '400',
      lineHeight: '1.5',
    },
    medium: {
      fontSize: '1.25rem',
      fontWeight: '500',
      lineHeight: '1.5',
    },
    bold: {
      fontSize: '1.25rem',
      fontWeight: '700',
      lineHeight: '1.5',
    },
  },
  h4: {
    regular: {
      fontSize: '1.125rem',
      fontWeight: '400',
      lineHeight: '1.5',
    },
    medium: {
      fontSize: '1.125rem',
      fontWeight: '500',
      lineHeight: '1.5',
    },
    bold: {
      fontSize: '1.125rem',
      fontWeight: '700',
      lineHeight: '1.5',
    },
  },
  text: {
    regular: {
      fontSize: '1rem',
      fontWeight: '400',
      lineHeight: '1.5',
    },
    medium: {
      fontSize: '1rem',
      fontWeight: '500',
      lineHeight: '1.5',
    },
    bold: {
      fontSize: '1rem',
      fontWeight: '700',
      lineHeight: '1.5',
    },
  },
  caption: {
    regular: {
      fontSize: '0.875rem',
      fontWeight: '400',
      lineHeight: '1.5',
    },
    medium: {
      fontSize: '0.875rem',
      fontWeight: '500',
      lineHeight: '1.5',
    },
    bold: {
      fontSize: '0.875rem',
      fontWeight: '700',
      lineHeight: '1.5',
    },
  },
})

export const media = {
  mobile: 'screen and (max-width: 768px)',
}
