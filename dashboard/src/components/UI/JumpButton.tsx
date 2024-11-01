import { A } from '@solidjs/router'
import type { VoidComponent } from 'solid-js'
import { ToolTip } from '/@/components/UI/ToolTip'
import { clsx } from '/@/libs/clsx'

const JumpButton: VoidComponent<{ href: string; tooltip?: string }> = (props) => (
  <ToolTip props={{ content: props.tooltip }} disabled={!props.tooltip}>
    <A href={props.href}>
      <div
        class={clsx(
          'flex size-6 shrink-0 cursor-pointer items-center justify-center rounded-md border-none bg-inherit text-text-black',
          'hover:bg-transparency-primary-hover',
          'active:bg-transparency-primary-selected active:text-primary-main data-[active]:bg-transparency-primary-selected data-[active]:text-primary-main',
          '!disabled:border-none !disabled:bg-text-disabled !disabled:text-text-black disabled:cursor-not-allowed',
        )}
      >
        <div class="i-material-symbols:arrow-outward shrink-0 text-xl/5" />
      </div>
    </A>
  </ToolTip>
)

export default JumpButton
