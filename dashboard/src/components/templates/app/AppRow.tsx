import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { applicationState, getWebsiteURL } from '/@/libs/application'
import { diffHuman, shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import Badge from '../../UI/Badge'
import { ToolTip } from '../../UI/ToolTip'
import { AppStatusIcon } from './AppStatusIcon'

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
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
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
    width: 'fit-content',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})
const UrlContainer = styled('div', {
  base: {
    width: 'fit-content',
    marginLeft: 'auto',
    textAlign: 'right',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
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
            {(nonNullUpdatedAt) => {
              const { diff, localeString } = diffHuman(nonNullUpdatedAt().toDate())
              return (
                <ToolTip props={{ content: localeString }}>
                  <UpdatedAt>{diff}</UpdatedAt>
                </ToolTip>
              )
            }}
          </Show>
        </TitleContainer>
        <MetaContainer>
          <LastCommitName>{shortSha(props.app.commit)}</LastCommitName>
          <Show when={props.app.websites.length > 0}>
            <UrlContainer>{getWebsiteURL(props.app.websites[0])}</UrlContainer>
            <Show when={props.app.websites.length > 1}>
              <Badge variant="text">{`+${props.app.websites.length - 1}`}</Badge>
            </Show>
          </Show>
        </MetaContainer>
      </Container>
    </A>
  )
}
