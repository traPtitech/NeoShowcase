import { type Component, Show } from 'solid-js'
import { writeToClipboard } from '/@/libs/clipboard'
import { ToolTip } from './ToolTip'

const Code: Component<{
  value: string
  copyable?: boolean
}> = (props) => {
  const handleCopy = () => {
    writeToClipboard(props.value)
  }

  return (
    <div
      class={[
        'bg-ui-secondary text-regular text-text-black relative mt-1 w-full min-w-[calc(1lh+8px)] overflow-x-auto rounded px-2 py-1 !font-[Menlo,Monaco,Consolas,Courier_New,monospace] whitespace-pre-wrap',
        { 'pr-10': props.copyable },
      ]}
    >
      {props.value}
      <Show when={props.copyable}>
        <ToolTip content="copy to clipboard">
          <button
            class="border-ui-border text-text-black hover:bg-black-alpha-200 active:bg-black-alpha-300 absolute top-1 right-2 grid size-6 cursor-pointer place-content-center rounded border bg-inherit leading-4"
            onClick={handleCopy}
            type="button"
          >
            <div class="i-material-symbols:content-copy-outline shrink-0 text-xl/5" />
          </button>
        </ToolTip>
      </Show>
    </div>
  )
}

export default Code
