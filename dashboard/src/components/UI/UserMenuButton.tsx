import { DropdownMenu } from '@kobalte/core/dropdown-menu'
import { type Component, For, Loading } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/session'
import { Button } from './Button'
import UserAvatar from './UserAvater'

const linkNameToMaterialIcon = (name: string): string => {
  // Manually assign icons to some known external link names
  const lowerName = name.toLowerCase()
  switch (lowerName) {
    case 'wiki':
    case 'help':
      return 'shrink-0 i-material-symbols:help-outline'
    case 'phpmyadmin':
    case 'adminer':
    case 'db admin':
      return 'shrink-0 i-material-symbols:database-outline'
  }
  if (lowerName.includes('mysql') || lowerName.includes('mongo')) {
    return 'shrink-0 i-material-symbols:database-outline'
  }
  return 'shrink-0 i-material-symbols:open-in-new'
}

export const UserMenuButton: Component<{
  user: User
}> = (props) => {
  return (
    <DropdownMenu placement="top-end">
      <DropdownMenu.Trigger
        class={[
          'w-fit-content relative flex h-11 cursor-pointer items-center gap-2 rounded-lg border-none bg-inherit px-2',
          'hover:bg-transparency-primary-hover active:bg-transparency-primary-selected',
        ]}
      >
        <UserAvatar user={props.user} size={32} />
        <span class="text-bold text-text-black max-md:hidden">{props.user.name}</span>
        <DropdownMenu.Icon class="data-[expanded]:rotate-180deg size-6 transition-transform duration-200">
          <div class="i-material-symbols:arrow-drop-down shrink-0 text-2xl/6" />
        </DropdownMenu.Icon>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          class={[
            'transform-origin-[var(--kb-menu-content-transform-origin)] bg-ui-primary shadow-default z-1 flex flex-col rounded-md p-1.5',
            'animate-duration-200 animate-ease-in-out animate-name-wipe-hide-up data-[expanded]:animate-name-wipe-show-down',
          ]}
        >
          <DropdownMenu.Item>
            <a href="/settings">
              <Button
                variants="text"
                size="medium"
                leftIcon={<div class="i-material-symbols:settings-outline shrink-0 text-2xl/6" />}
                full
              >
                Settings
              </Button>
            </a>
          </DropdownMenu.Item>
          <Loading>
            <For each={systemInfo().additionalLinks}>
              {(link) => (
                <DropdownMenu.Item>
                  <a href={link.url} target="_blank" rel="noopener noreferrer">
                    <Button
                      variants="text"
                      size="medium"
                      leftIcon={<div class={[linkNameToMaterialIcon(link.name), 'text-2xl/6']} />}
                      full
                    >
                      {link.name}
                    </Button>
                  </a>
                </DropdownMenu.Item>
              )}
            </For>
            <DropdownMenu.Item>
              <div class="text-regular text-text-grey flex flex-col px-4 py-1">
                <span>NeoShowcase</span>
                <span>
                  {systemInfo().version} ({systemInfo().revision})
                </span>
              </div>
            </DropdownMenu.Item>
          </Loading>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu>
  )
}
