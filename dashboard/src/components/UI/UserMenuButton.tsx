import { DropdownMenu } from '@kobalte/core'
import { A } from '@solidjs/router'
import { type Component, For } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'
import { clsx } from '/@/libs/clsx'
import { Button } from './Button'
import { MaterialSymbols } from './MaterialSymbols'
import UserAvatar from './UserAvater'

const linkNameToMaterialIcon = (name: string): string => {
  // Manually assign icons to some known external link names
  const lowerName = name.toLowerCase()
  switch (lowerName) {
    case 'wiki':
    case 'help':
      return 'help'
    case 'phpmyadmin':
    case 'adminer':
    case 'db admin':
      return 'database'
  }
  if (lowerName.includes('mysql') || lowerName.includes('mongo')) {
    return 'database'
  }
  return 'open_in_new'
}

export const UserMenuButton: Component<{
  user: User
}> = (props) => {
  return (
    <DropdownMenu.Root placement="top-end">
      <DropdownMenu.Trigger
        class={clsx(
          'relative flex h-11 w-fit-content cursor-pointer items-center gap-2 rounded-lg border-none bg-inherit px-2',
          'hover:bg-transparency-primary-hover active:bg-transparency-primary-selected',
        )}
      >
        <UserAvatar user={props.user} size={32} />
        <span class="text-bold text-text-black max-md:hidden">{props.user.name}</span>
        <DropdownMenu.Icon class="size-6 transition-transform duration-200 data-[expanded]:rotate-180deg">
          <MaterialSymbols>arrow_drop_down</MaterialSymbols>
        </DropdownMenu.Icon>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          class={clsx(
            'transform-origin-[var(--kb-menu-content-transform-origin)] z-1 flex flex-col rounded-md bg-ui-primary p-1.5 shadow-default',
            'animate-duration-200 animate-ease-in-out animate-wipe-hide-up data-[expanded]:animate-wipe-show-down',
          )}
        >
          <DropdownMenu.Item>
            <A href="/settings">
              <Button variants="text" size="medium" leftIcon={<MaterialSymbols>settings</MaterialSymbols>} full>
                Settings
              </Button>
            </A>
          </DropdownMenu.Item>
          <For each={systemInfo()?.additionalLinks}>
            {(link) => (
              <DropdownMenu.Item>
                <a href={link.url} target="_blank" rel="noopener noreferrer">
                  <Button
                    variants="text"
                    size="medium"
                    leftIcon={<MaterialSymbols>{linkNameToMaterialIcon(link.name)}</MaterialSymbols>}
                    full
                  >
                    {link.name}
                  </Button>
                </a>
              </DropdownMenu.Item>
            )}
          </For>
          <DropdownMenu.Item>
            <div class="flex flex-col px-4 py-1 text-regular text-text-grey">
              <span>NeoShowcase</span>
              <span>
                {systemInfo()?.version} ({systemInfo()?.revision})
              </span>
            </div>
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  )
}
