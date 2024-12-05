import { timestampDate } from '@bufbuild/protobuf/wkt'
import { A } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import type { Application, Build } from '/@/api/neoshowcase/protobuf/gateway_pb'
import Badge from '/@/components/UI/Badge'
import { ToolTip } from '/@/components/UI/ToolTip'
import type { CommitsMap } from '/@/libs/api'
import { diffHuman, shortSha } from '/@/libs/format'
import { BuildStatusIcon } from './BuildStatusIcon'

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
    const diff = diffHuman(timestampDate(c.commitDate))
    const localeString = timestampDate(c.commitDate).toLocaleString()
    return (
      <>
        {c.authorName}
        <span>, </span>
        <ToolTip props={{ content: localeString }}>
          <span>{diff()}</span>
        </ToolTip>
        <span>, </span>
        {shortSha(c.hash)}
      </>
    )
  }

  return (
    <A href={`/apps/${props.build.applicationId}/builds/${props.build.id}`}>
      <div class="w-full cursor-pointer bg-ui-primary p-4 pl-5 hover:bg-color-overlay-ui-primary-to-black-alpha-50">
        <div class="flex w-full items-center gap-2">
          <BuildStatusIcon state={props.build.status} />
          <div class="h4-regular w-auto truncate text-text-black">{commitHeadline()}</div>
          <Show when={props.isCurrent}>
            <ToolTip props={{ content: 'このビルドがデプロイされています' }}>
              <Badge variant="success">Current</Badge>
            </ToolTip>
          </Show>
          <div class="grow-1" />
          <Show when={props.build.queuedAt}>
            {(nonNullQueuedAt) => {
              const diff = diffHuman(timestampDate(nonNullQueuedAt()))
              const localeString = timestampDate(nonNullQueuedAt()).toLocaleString()
              return (
                <ToolTip props={{ content: localeString }}>
                  <div class="caption-regular shrink-0 text-text-grey">{diff()}</div>
                </ToolTip>
              )
            }}
          </Show>
        </div>
        <div class="caption-regular flex w-full items-center gap-1 pl-8 text-text-grey">
          <div class="w-fit truncate">{commitDetails()}</div>
          <div class="grow-1" />
          <Show when={props.app}>
            <div class="ml-auto flex w-fit items-center truncate text-right">
              <div class="i-material-symbols:deployed-code-outline shrink-0 text-xl/5" />
              {props.app!.name}
            </div>
          </Show>
        </div>
      </div>
    </A>
  )
}
