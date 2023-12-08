import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { Build } from '/@/api/neoshowcase/protobuf/gateway_pb'
import Badge from '/@/components/UI/Badge'
import { ToolTip } from '/@/components/UI/ToolTip'
import { colorOverlay } from '/@/libs/colorOverlay'
import { diffHuman, shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import { BuildStatusIcon } from './BuildStatusIcon'

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
})
const TitleContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
})
const BuildName = styled('div', {
  base: {
    width: 'auto',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
    color: colorVars.semantic.text.black,
    ...textVars.h4.regular,
  },
})
const Spacer = styled('div', {
  base: {
    flexGrow: 1,
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
const AppName = styled('div', {
  base: {
    width: 'fit-content',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})

export interface Props {
  build: Build
  appName?: string
  isCurrent: boolean
}

export const BuildRow: Component<Props> = (props) => {
  return (
    <A href={`/apps/${props.build.applicationId}/builds/${props.build.id}`}>
      <Container>
        <TitleContainer>
          <BuildStatusIcon state={props.build.status} />
          <BuildName>Build at {shortSha(props.build.commit)}</BuildName>
          <Show when={props.isCurrent}>
            <ToolTip props={{ content: 'このビルドがデプロイされています' }}>
              <Badge variant="success">Current</Badge>
            </ToolTip>
          </Show>
          <Spacer />
        </TitleContainer>
        <MetaContainer>
          <Show when={props.appName}>
            <AppName>{props.appName}・</AppName>
          </Show>
          <Show when={props.build.queuedAt}>
            {(nonNullQueuedAt) => {
              const { diff, localeString } = diffHuman(nonNullQueuedAt().toDate())
              return (
                <ToolTip props={{ content: localeString }}>
                  <UpdatedAt>{diff}</UpdatedAt>
                </ToolTip>
              )
            }}
          </Show>
        </MetaContainer>
      </Container>
    </A>
  )
}
