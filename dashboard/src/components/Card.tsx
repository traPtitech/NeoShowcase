import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'

export const CardsContainer = styled('div', {
  base: {
    marginTop: '24px',
    display: 'flex',
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: '40px',
  },
})

export const Card = styled('div', {
  base: {
    minWidth: '320px',
    borderRadius: '4px',
    border: `1px solid ${vars.bg.white4}`,
    background: vars.bg.white1,
    padding: '24px 36px',

    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
  },
})

export const CardTitle = styled('div', {
  base: {
    fontSize: '24px',
    fontWeight: 600,
  },
})

export const CardItems = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
})

export const CardItem = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
  },
})

export const CardItemTitle = styled('div', {
  base: {
    fontSize: '16px',
    color: vars.text.black3,
  },
})

export const CardItemContent = styled('div', {
  base: {
    marginLeft: 'auto',
    fontSize: '16px',
    color: vars.text.black1,

    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '4px',
  },
})
