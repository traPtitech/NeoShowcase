import { style } from '@vanilla-extract/css'
import { vars } from '/@/theme.css'

export const container = style({
  borderRadius: '4px',
  border: `1px solid ${vars.bg.white4}`,
})

export const header = style({
  height: '60px',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'space-between',

  padding: '0 20px',
  backgroundColor: vars.bg.white3,
})

export const headerLeft = style({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '8px',
  width: '100%',
})

export const repoName = style({
  fontSize: '16px',
  color: vars.text.black1,
})

export const appsCount = style({
  fontSize: '11px',
  color: vars.text.black3,
})

export const addBranchButton = style({
  display: 'flex',
  alignItems: 'center',

  padding: '8px 16px',
  borderRadius: '4px',
  backgroundColor: vars.bg.white5,

  fontSize: '12px',
  color: vars.text.black2,
})

const appBorder = style({
  borderWidth: '1px 0',
  borderStyle: 'solid',
  borderColor: vars.bg.white4,
})

export const application = style({
  height: '40px',
  display: 'grid',
  gridTemplateColumns: '20px 1fr',
  gap: '8px',
  padding: '12px 20px',

  backgroundColor: vars.bg.white1,
})

export const applicationNotLast = style([appBorder, application])

export const appDetail = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '4px',
})

export const appName = style({
  fontSize: '14px',
  color: vars.text.black1,
})

export const appFooter = style({
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'space-between',
  width: '100%',

  fontSize: '11px',
  color: vars.text.black3,
})

export const appFooterRight = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '48px',
})
