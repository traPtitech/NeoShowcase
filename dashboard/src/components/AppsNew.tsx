import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'

export const FormButton = styled('div', {
  base: {
    marginLeft: '8px',
  },
})

export const FormTextBig = styled('div', {
  base: {
    fontSize: '20px',
    alignItems: 'center',
    fontWeight: 900,
    color: vars.text.black1,

    marginBottom: '4px',
  },
})

export const FormCheckBox = styled('div', {
  base: {
    background: vars.bg.white1,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    marginLeft: '4px',
    padding: '8px 12px',

    width: '320px',
  },
})

export const SettingsContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '10px',
  },
})

export const FormSettings = styled('div', {
  base: {
    background: vars.bg.white2,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    marginLeft: '4px',
    padding: '8px 12px',

    display: 'flex',
    flexDirection: 'column',
    gap: '12px',
  },
})

export const FormSettingsButton = styled('div', {
  base: {
    display: 'flex',
    gap: '8px',
    marginBottom: '4px',
  },
})
