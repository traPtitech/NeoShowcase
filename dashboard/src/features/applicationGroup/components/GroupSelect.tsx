import { type Component, createMemo, Show } from 'solid-js'
import type { ApplicationGroup } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { type SelectOption, SingleSelect } from '/@/components/templates/Select'
import type { TooltipProps } from '/@/components/UI/ToolTip'
import { getMyApplicationGroups } from '/@/libs/api'

// The option value of "no group". Kobalte's Select treats the empty string as no selection and shows nothing.
const noGroup = '__none__'

type Props = {
  name?: string
  label?: string
  info?: TooltipProps
  /** Group ID, or undefined for a standalone application. */
  value: string | undefined
  setValue: (groupId: string | undefined) => void
  /** The group the application is in now. Listed even when the user does not own it. */
  current?: ApplicationGroup
  readOnly?: boolean
}

/** Selects one of the groups the user owns, or no group. Read it under `<Loading>`. */
const GroupSelect: Component<Props> = (props) => {
  const options = createMemo(async (): Promise<SelectOption<string>[]> => {
    const current = props.current
    const mine = await getMyApplicationGroups()
    const listed = current && !mine.some((g) => g.id === current.id) ? [current, ...mine] : mine
    return [{ value: noGroup, label: 'なし (単体のアプリ)' }, ...listed.map((g) => ({ value: g.id, label: g.name }))]
  })

  // Resolve the options in a tracked scope first, and re-create the select when they change:
  // Kobalte reads its root props untracked, on construction and in its own effects.
  return (
    <Show when={options()} keyed>
      {(opts) => (
        <SingleSelect
          name={props.name}
          label={props.label}
          info={props.info}
          options={opts}
          value={props.value || noGroup}
          setValue={(v) => props.setValue(v === noGroup ? undefined : v)}
          readOnly={props.readOnly}
        />
      )}
    </Show>
  )
}

export default GroupSelect
