import { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { applicationState, getWebsiteURL } from '/@/libs/application'
import { DiffHuman, shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { AppStatusIcon } from '../UI/AppStatusIcon'

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
    textAlign: 'right',
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

    background: colorVars.primitive.blackAlpha[200],
    color: colorVars.semantic.text.black,
    ...textVars.caption.regular,
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
          <AppStatusIcon state={applicationState(props.app)} />
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
          <LastCommitName>{shortSha(props.app.commit)}</LastCommitName>
          <Show when={props.app.websites.length > 0}>
            <UrlContainer>{getWebsiteURL(props.app.websites[0])}</UrlContainer>
            <Show when={props.app.websites.length > 1}>
              <UrlCount>{`+${props.app.websites.length - 1}`}</UrlCount>
            </Show>
          </Show>
        </MetaContainer>
      </Container>
    </A>
  )
}
