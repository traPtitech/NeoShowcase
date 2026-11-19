import { timestampDate } from '@bufbuild/protobuf/wkt'
import { type Component, Show } from 'solid-js'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { ToolTip } from '/@/components/UI/ToolTip'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import { diffHuman } from '/@/libs/format'
import { Nav } from '../Nav'

export const AppNav: Component<{
  app: Application
  repository: Repository
}> = (props) => {
  const edited = (
    <div class="ml-auto flex w-fit items-center gap-1 truncate text-right">
      <div class="i-material-symbols:edit-outline shrink-0" />
      <Show when={props.app.updatedAt} keyed>
        {(updatedAt) => {
          const date = timestampDate(updatedAt)
          const diff = diffHuman(date)
          return (
            <ToolTip content={`App Last Edited: ${date.toLocaleString()}`}>
              <div>{diff()}</div>
            </ToolTip>
          )
        }}
      </Show>
    </div>
  )

  return (
    <Nav
      title={props.app.name}
      icon={<div class="i-material-symbols:deployed-code-outline text-10/10 shrink-0" />}
      action={edited}
    >
      <div class="text-regular text-text-black mt-1 flex w-full items-center gap-2 whitespace-nowrap">
        created from
        <a href={`/repos/${props.repository.id}`} class="overflow-hidden">
          <div class="flex w-full items-center gap-1 overflow-x-hidden">
            {originToIcon(repositoryURLToOrigin(props.repository.url), 20)}
            <div class="w-full truncate">{props.repository.name}</div>
          </div>
        </a>
      </div>
    </Nav>
  )
}
