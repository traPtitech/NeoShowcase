import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import Badge from '/@/components/UI/Badge'
import Skeleton from '/@/components/UI/Skeleton'
import { ToolTip } from '/@/components/UI/ToolTip'
import { applicationState, getWebsiteURL } from '/@/libs/application'
import { colorOverlay } from '/@/libs/colorOverlay'
import { diffHuman, shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { AppStatusIcon } from './AppStatusIcon'

const Container = styled('div', {
  base: {
    width: '100%',
    padding: '16px 16px 16px 20px',
    cursor: 'pointer',
    background: colorVars.semantic.ui.primary,

    selectors: {
      '&:hover': {
        background: colorOverlay(colorVars.semantic.ui.primary, colorVars.primitive.blackAlpha[50]),
      },
    },
  },
  variants: {
    dark: {
      true: {
        background: colorVars.semantic.ui.secondary,
        selectors: {
          '&:hover': {
            background: colorOverlay(colorVars.semantic.ui.secondary, colorVars.primitive.blackAlpha[50]),
          },
        },
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

const AppRowSkeleton: Component<{
  dark?: boolean
}> = (props) => {
  return (
    <Container dark={props.dark} style={{ 'pointer-events': 'none' }}>
      <TitleContainer>
        <Skeleton height={24} circle />
        <AppName>
          <Skeleton>App Name Placeholder</Skeleton>
        </AppName>
        <UpdatedAt>
          <Skeleton>1 day ago</Skeleton>
        </UpdatedAt>
      </TitleContainer>
      <MetaContainer>
        <Skeleton>0000000</Skeleton>
        <UrlContainer>
          <Skeleton>https://example.com</Skeleton>
        </UrlContainer>
      </MetaContainer>
    </Container>
  )
}

export interface Props {
  app?: Application
  dark?: boolean
}

export const AppRow: Component<Props> = (props) => {
  return (
    <Show when={props.app} fallback={<AppRowSkeleton dark={props.dark} />}>
      <A href={`/apps/${props.app!.id}`}>
        <Container dark={props.dark}>
          <TitleContainer>
            <AppStatusIcon state={applicationState(props.app!)} />
            <AppName>{props.app!.name}</AppName>
            <Show when={props.app!.updatedAt}>
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
            <LastCommitName>{shortSha(props.app!.commit)}</LastCommitName>
            <Show when={props.app!.websites.length > 0}>
              <UrlContainer>{getWebsiteURL(props.app!.websites[0])}</UrlContainer>
              <Show when={props.app!.websites.length > 1}>
                <Badge variant="text">{`+${props.app!.websites.length - 1}`}</Badge>
              </Show>
            </Show>
          </MetaContainer>
        </Container>
      </A>
    </Show>
  )
}
