import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import type { Application, Build } from '/@/api/neoshowcase/protobuf/gateway_pb'
import Badge from '/@/components/UI/Badge'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { ToolTip } from '/@/components/UI/ToolTip'
import type { CommitsMap } from '/@/libs/api'
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

const leftFit = style({
  width: 'fit-content',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const rightFit = style({
  width: 'fit-content',
  marginLeft: 'auto',
  textAlign: 'right',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const center = style({
  display: 'flex',
  alignItems: 'center',
})

export interface Props {
  build: Build
  app?: Application
  commits?: CommitsMap
  isCurrent: boolean
}

export const BuildRow: Component<Props> = (props) => {
  const commit = () => props.commits?.[props.build.commit]
  const commitHeadline = () => {
    const c = commit()
    if (!c) return `Build at ${shortSha(props.build.commit)}`
    return c.message.split('\n')[0]
  }
  const commitDetails = () => {
    const c = commit()
    if (!c || !c.commitDate) return '<no info>'
    const diff = diffHuman(c.commitDate.toDate())
    const localeString = c.commitDate.toDate().toLocaleString()
    return (
      <>
        {c.authorName}
        <span>, </span>
        <ToolTip props={{ content: localeString }}>
          <span>{diff}</span>
        </ToolTip>
        <span>, </span>
        {shortSha(c.hash)}
      </>
    )
  }

  return (
    <A href={`/apps/${props.build.applicationId}/builds/${props.build.id}`}>
      <Container>
        <TitleContainer>
          <BuildStatusIcon state={props.build.status} />
          <BuildName>{commitHeadline()}</BuildName>
          <Show when={props.isCurrent}>
            <ToolTip props={{ content: 'このビルドがデプロイされています' }}>
              <Badge variant="success">Current</Badge>
            </ToolTip>
          </Show>
          <Spacer />
          <Show when={props.build.queuedAt}>
            {(nonNullQueuedAt) => {
              const diff = diffHuman(nonNullQueuedAt().toDate())
              const localeString = nonNullQueuedAt().toDate().toString()
              return (
                <ToolTip props={{ content: localeString }}>
                  <UpdatedAt>{diff}</UpdatedAt>
                </ToolTip>
              )
            }}
          </Show>
        </TitleContainer>
        <MetaContainer>
          <div class={leftFit}>{commitDetails()}</div>
          <Spacer />
          <Show when={props.app}>
            <div class={`${rightFit} ${center}`}>
              <MaterialSymbols displaySize={20}>deployed_code</MaterialSymbols>
              {props.app!.name}
            </div>
          </Show>
        </MetaContainer>
      </Container>
    </A>
  )
}
