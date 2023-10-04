import { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { applicationState, getWebsiteURL } from '/@/libs/application'
import { colorOverlay } from '/@/libs/colorOverlay'
import { DiffHuman } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { AppStatus } from '../UI/AppStatus'
import { Button } from '../UI/Button'

const Container = styled('div', {
  base: {
    width: '100%',
    padding: '16px 16px 16px 20px',
    cursor: 'pointer',

    selectors: {
      '&:hover': {
        background: colorVars.primitive.blackAlpha[50],
      },
    },
  },
})
const TitleContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
})
const AppName = styled('div', {
  base: {
    width: '100%',
    color: colorVars.semantic.text.black,
    ...textVars.h4.regular,
  },
})
const UpdatedAt = styled('div', {
  base: {
    flexShrink: 0,
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const MetaContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '4px',
    padding: '0 0 0 32px',

    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const LastCommitName = styled('div', {
  base: {
    width: '100%',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})
const UrlContainer = styled('div', {
  base: {
    width: '100%',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})
const UrlCount = styled('div', {
  base: {
    height: '20px',
    padding: '0 8px',
    borderRadius: '9999px',

    color: colorVars.semantic.text.black,
    background: colorVars.primitive.blackAlpha[200],
  },
})

export interface Props {
  app: Application
}

export const AppRow: Component<Props> = (props) => {
  return (
    <A href={`/apps/${props.app.id}`}>
      <Container>
        <TitleContainer>
          <AppStatus state={applicationState(props.app)} />
          <AppName>{props.app.name}</AppName>
          <Show when={props.app.updatedAt}>
            {(nonNullUpdatedAt) => (
              <UpdatedAt>
                <DiffHuman target={nonNullUpdatedAt().toDate()} />
              </UpdatedAt>
            )}
          </Show>
        </TitleContainer>
        <MetaContainer>
          <LastCommitName>LastCommitName</LastCommitName>
          <Show when={props.app.websites.length > 0}>
            <div>{getWebsiteURL(props.app.websites[0])}</div>
            <Show when={props.app.websites.length > 1}>
              <UrlCount>{`+${props.app.websites.length - 1}`}</UrlCount>
            </Show>
          </Show>
        </MetaContainer>
      </Container>
    </A>
  )
}
