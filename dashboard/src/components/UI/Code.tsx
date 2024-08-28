import { type Component, Show } from 'solid-js'
import { writeToClipboard } from '/@/libs/clipboard'
import { MaterialSymbols } from './MaterialSymbols'
import { ToolTip } from './ToolTip'
import { clsx } from '/@/libs/clsx'

const Code: Component<{
  value: string
  copyable?: boolean
}> = (props) => {
  const handleCopy = () => {
    writeToClipboard(props.value)
  }

  return (
    <div
      class={clsx(
        '!font-[Menlo,Monaco,Consolas,Courier_New,monospace] relative mt-1 w-full min-w-[calc(1lh+8px)] overflow-x-auto whitespace-pre-wrap rounded bg-ui-secondary px-2 py-1 text-regular text-text-black',
        props.copyable && 'pr-10',
      )}
    >
      {props.value}
      <Show when={props.copyable}>
        <ToolTip
          props={{
            content: 'copy to clipboard',
          }}
        >
          <button
            class="absolute top-1 right-2 grid size-6 cursor-pointer place-content-center rounded border border-ui-border bg-none text-text-black leading-4 hover:bg-black-alpha-200 active:bg-black-alpha-300"
            onClick={handleCopy}
            type="button"
          >
            <MaterialSymbols opticalSize={20}>content_copy</MaterialSymbols>
          </button>
        </ToolTip>
      </Show>
    </div>
  )
}

export default Code
