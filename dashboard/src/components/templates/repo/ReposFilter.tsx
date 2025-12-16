import { Checkbox, DropdownMenu } from '@kobalte/core'
import { type Component, For, type Setter, Show } from 'solid-js'
import { CheckBoxIcon } from '/@/components/UI/CheckBoxIcon'
import { originToIcon, type RepositoryOrigin } from '/@/libs/application'
import { clsx } from '/@/libs/clsx'
import { allOrigins } from '/@/pages/apps'

// TODO: AppsFilter と共通するスタイルが多いので共通化する

const selectItemStyle = clsx(
  'flex h-11 w-full cursor-pointer flex-nowrap items-center gap-2 whitespace-nowrap rounded-lg border-none bg-inherit p-2 text-bold text-text-black',
  'hover:bg-transparency-primary-hover data-[highlighted]:bg-transparency-primary-hover',
  '!data-[disabled]:bg-text-disabled !data-[disabled]:text-text-black data-[disabled]:cursor-not-allowed',
)

const ReposFilter: Component<{
  origin: RepositoryOrigin[]
  setOrigin: Setter<RepositoryOrigin[]>
}> = (props) => {
  const filtered = () => props.origin.length !== allOrigins.length

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        class={clsx(
          'flex cursor-pointer rounded bg-inherit p-2 text-text-black',
          'hover:bg-transparency-primary-hover',
          'active:bg-transparency-primary-selected active:text-primary-main',
        )}
      >
        <div class="relative size-6">
          <div class="i-material-symbols:tune shrink-0 text-2xl/6" />
          <Show when={filtered()}>
            <div class="absolute -top-0.5 -right-0.5 size-2 rounded bg-primary-main outline outline-1 outline-ui-background" />
          </Show>
        </div>
        <DropdownMenu.Icon class="size-6 transition-transform duration-200 data-[expanded]:rotate-180">
          <div class="i-material-symbols:expand-more shrink-0 text-2xl/6" />
        </DropdownMenu.Icon>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content class="z-1 flex origin-[--kb-menu-content-transform-origin] animate-duration-200 animate-name-wipe-hide-up gap-2 rounded-md bg-ui-primary p-4 shadow-default ease-in-out data-[expanded]:animate-name-wipe-show-down">
          <div class="flex flex-col gap-2 text-bold text-text-black">
            Origin
            <div class="flex w-full flex-col">
              <For each={allOrigins}>
                {(s) => (
                  <Checkbox.Root
                    checked={props.origin.includes(s.value)}
                    onChange={(selected) => {
                      if (selected) {
                        props.setOrigin([...props.origin, s.value])
                      } else {
                        props.setOrigin(props.origin.filter((v) => v !== s.value))
                      }
                    }}
                  >
                    <Checkbox.Input />
                    <Checkbox.Label class={selectItemStyle}>
                      <Checkbox.Indicator forceMount class="size-6">
                        <CheckBoxIcon checked={props.origin.includes(s.value)} />
                      </Checkbox.Indicator>
                      {originToIcon(s.value)}
                      {s.label}
                    </Checkbox.Label>
                  </Checkbox.Root>
                )}
              </For>
            </div>
          </div>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  )
}

export default ReposFilter
