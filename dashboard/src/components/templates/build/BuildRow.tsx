import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { Build } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { diffHuman, shortSha } from '/@/libs/format'
import { colorVars, textVars } from '/@/theme'
import Badge from '../../UI/Badge'
import { ToolTip } from '../../UI/ToolTip'
import { BuildStatusIcon } from './BuildStatusIcon'

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
const CommitName = styled('div', {
  base: {
    width: 'fit-content',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})

export interface Props {
  appId: string
  build: Build
  showAppId?: boolean
  isDeployed?: boolean
}

export const BuildRow: Component<Props> = (props) => {
  return (
    <A href={`/apps/${props.appId}/builds/${props.build.id}`}>
      <Container>
        <TitleContainer>
          <BuildStatusIcon state={props.build.status} />
          <BuildName>
            Build {props.build.id}
            {props.showAppId && ` (App ${props.appId})`}
          </BuildName>
          <Show when={props.isDeployed}>
            <Badge variant="success">Deployed</Badge>
          </Show>
          <Spacer />
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
        </TitleContainer>
        <MetaContainer>
          <CommitName>{shortSha(props.build.commit)}</CommitName>
        </MetaContainer>
      </Container>
    </A>
  )
}
