import { A } from '@solidjs/router'
import { BiRegularPencil } from 'solid-icons/bi'
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
      <BiRegularPencil />
      <Show when={props.app.updatedAt}>
        {(nonNullUpdatedAt) => {
          const diff = diffHuman(nonNullUpdatedAt().toDate())
          const localeString = nonNullUpdatedAt().toDate().toLocaleString()
          return (
            <ToolTip props={{ content: `App Last Edited: ${localeString}` }}>
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
      icon={<span class="i-material-symbols:deployed-code-outline text-10/10" />}
      action={edited}
    >
      <div class="mt-1 flex w-full items-center gap-2 whitespace-nowrap text-regular text-text-black">
        created from
        <A href={`/repos/${props.repository.id}`} class="overflow-hidden">
          <div class="flex w-full items-center gap-1 overflow-x-hidden">
            {originToIcon(repositoryURLToOrigin(props.repository.url), 20)}
            <div class="w-full truncate">{props.repository.name}</div>
          </div>
        </A>
      </div>
    </Nav>
  )
}
