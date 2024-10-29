import { defineConfig, presetIcons, presetUno, transformerVariantGroup } from 'unocss'
import { parseColor } from '@unocss/preset-mini/utils'

export default defineConfig({
  presets: [presetUno(), presetIcons()],
  rules: [
    [
      /^overflow-wrap-(normal|break-word|anywhere|inherit|initial|revert|unset)$/,
      ([, p]) => ({ 'overflow-wrap': p }),
      {
        autocomplete: 'overflow-wrap-(normal|break-word|anywhere|inherit|initial|revert|unset)',
      },
    ],
    [
      /^bg-color-overlay-(.*)-to-(.*)$/,
      ([, base, overlay], { theme }) => {
        const baseColor = parseColor(base, theme)?.color
        const overlayColor = parseColor(overlay, theme)?.color
        if (!baseColor || !overlayColor) return {}
        return {
          background: `linear-gradient(0deg, ${overlayColor} 0%, ${overlayColor} 100%), ${baseColor}`,
        }
      },
      { autocomplete: 'bg-color-overlay-$colors-to-$colors' },
    ],
    [
      /^scrollbar-gutter-both$/,
      () => ({
        'scrollbar-gutter': 'stable both-edges',
      }),
      {
        autocomplete: 'scrollbar-gutter-(auto|stable|inherit|initial|revert|revert-layer|unset|both)',
      },
    ],
    [
      /^scrollbar-gutter-(auto|stable|inherit|initial|revert|revert-layer|unset)$/,
      ([, p]) => ({
        'scrollbar-gutter': p,
      }),
    ],
    [
      /^leading-(.+)$/,
      ([, p]) => ({
        'line-height': Number(p) / 4,
      }),
      {
        autocomplete: 'leading-<num>',
      },
    ],
    [
      /^shadow-default$/,
      ([, p]) => ({
        'box-shadow': '0 0 20px 0 rgba(0, 0, 0, .1)',
      }),
      {
        autocomplete: 'shadow-default',
      },
    ],
    [
      /^animate-wipe-(show|hide)-(up|right|down|left)$/,
      ([, type, direction]) => {
        return {
          'animation-name': `keyframe-wipe-${type}-${direction}`,
        }
      },
      {
        autocomplete: 'animate-wipe-(show|hide)-(up|right|down|left)',
      },
    ],
    [
      /^animate-duration-(.+)$/,
      ([, p]) => ({
        'animation-duration': `${p}ms`,
      }),
      {
        autocomplete: 'animate-duration-(75|100|150|200|300|500|700|1000)',
      },
    ],
    [
      /^animate-(ease|linear|ease-in|ease-out|ease-in-out)$/,
      ([, p]) => ({
        'animation-timing-function': p,
      }),
      {
        autocomplete: 'animate-timing-(ease|linear|ease-in|ease-out|ease-in-out)',
      },
    ],
  ],
  shortcuts: {
    'h1-regular': 'text-7 leading-6',
    'h1-medium': 'text-7 font-medium leading-6',
    'h1-bold': 'text-7 font-bold leading-6',
    'h2-regular': 'text-6 leading-6',
    'h2-medium': 'text-6 font-medium leading-6',
    'h2-bold': 'text-6 font-bold leading-6',
    'h3-regular': 'text-5 leading-6',
    'h3-medium': 'text-5 font-medium leading-6',
    'h3-bold': 'text-5 font-bold leading-6',
    'h4-regular': 'text-4.5 leading-6',
    'h4-medium': 'text-4.5 font-medium leading-6',
    'h4-bold': 'text-4.5 font-bold leading-6',
    'text-regular': 'text-4 leading-6',
    'text-medium': 'text-4 font-medium leading-6',
    'text-bold': 'text-4 font-bold leading-6',
    'caption-regular': 'text-3.5 leading-6',
    'caption-medium': 'text-3.5 font-medium leading-6',
    'caption-bold': 'text-3.5 font-bold leading-6',
  },
  theme: {
    breakpoints: {
      md: '768px',
      lg: '1024px',
    },
    colors: {
      primary: {
        white: '#FFFFFF',
        main: '#005BAC',
      },
      accent: {
        error: '#F25151',
        warn: '#F1B61E',
        success: '#20BD77',
      },
      transparency: {
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
  },
  transformers: [transformerVariantGroup()],
})
