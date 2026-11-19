import type { VoidComponent } from 'solid-js'
import { ToolTip } from '/@/components/UI/ToolTip'

const JumpButton: VoidComponent<{ href: string; tooltip?: string }> = (props) => (
  <ToolTip content={props.tooltip}>
    <a href={props.href}>
      <div
        class={[
          'text-text-black flex size-6 shrink-0 cursor-pointer items-center justify-center rounded-md border-none bg-inherit',
          'hover:bg-transparency-primary-hover',
          'active:bg-transparency-primary-selected active:text-primary-main data-[active]:bg-transparency-primary-selected data-[active]:text-primary-main',
          'disabled:!bg-text-disabled disabled:!text-text-black disabled:cursor-not-allowed disabled:!border-none',
        ]}
      >
        <div class="i-material-symbols:arrow-outward shrink-0 text-xl/5" />
      </div>
    </a>
  </ToolTip>
)

export default JumpButton
