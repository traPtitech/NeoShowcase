import { timestampDate } from '@bufbuild/protobuf/wkt'
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
    if (!c?.commitDate) return '<no info>'
    const diff = diffHuman(timestampDate(c.commitDate))
    const localeString = timestampDate(c.commitDate).toLocaleString()
    return (
      <>
        {c.authorName}
        <span>, </span>
        <ToolTip content={localeString}>
          <span>{diff()}</span>
        </ToolTip>
        <span>, </span>
        {shortSha(c.hash)}
      </>
    )
  }

  return (
    <a href={`/apps/${props.build.applicationId}/builds/${props.build.id}`}>
      <div class="bg-ui-primary hover:bg-color-overlay-ui-primary-to-black-alpha-50 w-full cursor-pointer p-4 pl-5">
        <div class="flex w-full items-center gap-2">
          <BuildStatusIcon state={props.build.status} />
          <div class="h4-regular text-text-black w-auto truncate">{commitHeadline()}</div>
          <Show when={props.isCurrent}>
            <ToolTip content="このビルドがデプロイされています">
              <Badge variant="success">Current</Badge>
            </ToolTip>
          </Show>
          <div class="grow-1" />
          <Show when={props.build.queuedAt} keyed>
            {(queuedAt) => {
              const date = timestampDate(queuedAt)
              const diff = diffHuman(date)
              return (
                <ToolTip content={date.toLocaleString()}>
                  <div class="caption-regular text-text-grey shrink-0">{diff()}</div>
                </ToolTip>
              )
            }}
          </Show>
        </div>
        <div class="caption-regular text-text-grey flex w-full items-center gap-1 pl-8">
          <div class="w-fit truncate">{commitDetails()}</div>
          <div class="grow-1" />
          <Show when={props.app} keyed>
            {(app) => (
              <div class="ml-auto flex w-fit items-center truncate text-right">
                <div class="i-material-symbols:deployed-code-outline shrink-0 text-xl/5" />
                {app.name}
              </div>
            )}
          </Show>
        </div>
      </div>
    </a>
  )
}
