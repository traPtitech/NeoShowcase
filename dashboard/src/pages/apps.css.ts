import { style } from '@vanilla-extract/css'
import { vars } from '/@/theme.css'

export const container = style({
  padding: '40px 72px',
})

export const appsTitle = style({
  marginTop: '48px',
  fontSize: '32px',
  fontWeight: 'bold',
  color: vars.text.black1,
})

export const contentContainer = style({
  marginTop: '24px',
  display: 'grid',
  gridTemplateColumns: '380px 1fr',
  gap: '40px',
})

export const sidebarContainer = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '22px',

  padding: '24px 40px',
  backgroundColor: vars.bg.white1,
  borderRadius: '4px',
  border: `1px solid ${vars.bg.white4}`,
})

export const sidebarSection = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '16px',
})

export const sidebarTitle = style({
  fontSize: '24px',
  fontWeight: 500,
  color: vars.text.black1,
})

export const sidebarOptions = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '12px',

  fontSize: '20px',
  color: vars.text.black1,
})

export const statusCheckboxContainer = style({
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'space-between',
  alignItems: 'center',
  width: '100%',
})

export const statusCheckboxContainerLeft = style({
  display: 'flex',
  flexDirection: 'row',
  gap: '8px',
  alignItems: 'center',
})

export const mainContentContainer = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '20px',
})

export const searchBarContainer = style({
  display: 'grid',
  gridTemplateColumns: '1fr 180px',
  gap: '20px',
  height: '44px',
})

export const searchBar = style({
  padding: '12px 20px',
  borderRadius: '4px',
  border: `1px solid ${vars.bg.white4}`,
  fontSize: '14px',

  '::placeholder': {
    color: vars.text.black3,
  },
})

export const createAppButton = style({
  display: 'flex',
  borderRadius: '4px',
  backgroundColor: vars.bg.black1,
})

export const createAppText = style({
  margin: 'auto',
  color: vars.text.white1,
  fontSize: '16px',
  fontWeight: 'bold',
})

export const repositoriesContainer = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '20px',
})
