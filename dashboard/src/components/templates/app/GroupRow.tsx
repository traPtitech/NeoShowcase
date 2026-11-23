import type { Component } from 'solid-js'
import type { ApplicationGroup } from '/@/api/neoshowcase/protobuf/gateway_pb'

export const GroupRow: Component<{
  group: ApplicationGroup
  appCount: number
}> = (props) => (
  <div class="bg-ui-primary flex h-19 w-full items-center gap-2 overflow-hidden p-4 pl-5">
    <div class="i-material-symbols:folder-outline text-text-black shrink-0 text-2xl/6" />
    <a href={`/groups/${props.group.id}`} class="overflow-hidden">
      <div class="h4-bold text-text-black truncate">{props.group.name}</div>
    </a>
    <div class="caption-regular text-text-grey whitespace-nowrap">{`${props.appCount} apps`}</div>
  </div>
)
