import { type Component, Show } from 'solid-js'
import type { ApplicationGroup } from '/@/api/neoshowcase/protobuf/gateway_pb'

/** A row that links to the page of a group. The number of its apps is shown when `appCount` is given. */
export const GroupRow: Component<{
  group: ApplicationGroup
  appCount?: number
}> = (props) => (
  <a href={`/groups/${props.group.id}`}>
    <div class="bg-ui-primary hover:bg-color-overlay-ui-primary-to-black-alpha-50 flex h-19 w-full cursor-pointer items-center gap-2 overflow-hidden p-4 pl-5">
      <div class="i-material-symbols:folder-outline text-text-black shrink-0 text-2xl/6" />
      <div class="h4-bold text-text-black truncate">{props.group.name}</div>
      <Show when={props.appCount !== undefined}>
        <div class="caption-regular text-text-grey whitespace-nowrap">{`${props.appCount} apps`}</div>
      </Show>
    </div>
  </a>
)
