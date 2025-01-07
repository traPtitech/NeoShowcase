import { Checkbox, DropdownMenu, RadioGroup } from '@kobalte/core'
import { type Component, type ComponentProps, For, type Setter, Show } from 'solid-js'
import { CheckBoxIcon } from '/@/components/UI/CheckBoxIcon'
import { RadioIcon } from '/@/components/UI/RadioIcon'
import { styled } from '/@/components/styled-components'
import {
  type ApplicationState,
  type RepoWithApp,
  type RepositoryOrigin,
  applicationState,
  originToIcon,
  repositoryURLToOrigin,
} from '/@/libs/application'
import { clsx } from '/@/libs/clsx'
import { allOrigins, allStatuses, sortItems } from '/@/pages/apps'
import { AppStatusIcon } from './AppStatusIcon'

const ItemsContainer = styled('div', 'flex w-full flex-col')

const selectItemStyle = clsx(
  'flex h-11 w-full cursor-pointer flex-nowrap items-center gap-2 whitespace-nowrap rounded-lg border-none bg-inherit p-2 text-bold text-text-black',
  'hover:bg-transparency-primary-hover data-[highlighted]:bg-transparency-primary-hover',
  '!data-[disabled]:bg-text-disabled !data-[disabled]:text-text-black data-[disabled]:cursor-not-allowed',
)

const FilterItemContainer = styled('div', 'flex flex-col gap-2 text-bold text-text-black')

const AppsFilter: Component<{
  allRepoWithApps: RepoWithApp[]
  statuses: ApplicationState[]
  setStatues: Setter<ApplicationState[]>
  origin: RepositoryOrigin[]
  setOrigin: Setter<RepositoryOrigin[]>
  sort: keyof typeof sortItems
  setSort: Setter<keyof typeof sortItems>
  includeNoApp: boolean
  setIncludeNoApp: Setter<boolean>
}> = (props) => {
  const filtered = () =>
    props.statuses.length !== allStatuses.length || props.origin.length !== allOrigins.length || props.includeNoApp

  const appCountByStatus = (status: ApplicationState): number => {
    return props.allRepoWithApps
      .filter((repo) => props.origin.includes(repositoryURLToOrigin(repo.repo.url)))
      .flatMap((repo) => repo.apps.filter((app) => applicationState(app) === status)).length
  }

  const repoCountByOrigin = (origin: RepositoryOrigin): number => {
    return props.allRepoWithApps
      .filter((repo) => repositoryURLToOrigin(repo.repo.url) === origin)
      .filter(
        (repo) =>
          props.includeNoApp || repo.apps.filter((app) => props.statuses.includes(applicationState(app))).length > 0,
      ).length
  }

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
            <div class="-right-0.5 -top-0.5 absolute size-2 rounded bg-primary-main outline outline-1 outline-ui-background" />
          </Show>
        </div>
        <DropdownMenu.Icon class="size-6 transition-transform duration-200 data-[expanded]:rotate-180">
          <div class="i-material-symbols:expand-more shrink-0 text-2xl/6" />
        </DropdownMenu.Icon>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          class={clsx(
            'z-1 grid max-w-[--kb-popper-content-available-width] origin-[--kb-menu-content-transform-origin] grid-cols-[repeat(3,auto)] grid-rows-[1fr_auto] gap-2 overflow-x-auto rounded-md bg-ui-primary p-4 shadow-default',
            'animate-duration-200 animate-ease-in-out animate-name-wipe-hide-up data-[expanded]:animate-name-wipe-show-down',
          )}
          style={{
            'grid-template-areas': `
              "status provider sort"
              "status noapp noapp"
            `,
          }}
        >
          <FilterItemContainer class="grid-area-[status]">
            App Status
            <ItemsContainer>
              <For each={allStatuses}>
                {(s) => (
                  <Checkbox.Root
                    checked={props.statuses.includes(s.value)}
                    onChange={(selected) => {
                      if (selected) {
                        props.setStatues([...props.statuses, s.value])
                      } else {
                        props.setStatues(props.statuses.filter((v) => v !== s.value))
                      }
                    }}
                  >
                    <Checkbox.Input />
                    <Checkbox.Label class={selectItemStyle}>
                      <Checkbox.Indicator forceMount class="size-6">
                        <CheckBoxIcon checked={props.statuses.includes(s.value)} />
                      </Checkbox.Indicator>
                      <AppStatusIcon state={s.value} hideTooltip />
                      <span>
                        {s.label} ({appCountByStatus(s.value)})
                      </span>
                    </Checkbox.Label>
                  </Checkbox.Root>
                )}
              </For>
            </ItemsContainer>
          </FilterItemContainer>
          <FilterItemContainer class="grid-area-[provider]">
            Repo Origin
            <ItemsContainer>
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
                      <span>
                        {s.label} ({repoCountByOrigin(s.value)})
                      </span>
                    </Checkbox.Label>
                  </Checkbox.Root>
                )}
              </For>
            </ItemsContainer>
          </FilterItemContainer>
          <RadioGroup.Root
            onChange={props.setSort}
            as={(asProps: ComponentProps<typeof FilterItemContainer>) => (
              <FilterItemContainer class="grid-area-[sort]" {...asProps}>
                <RadioGroup.Label>Sort</RadioGroup.Label>
                <ItemsContainer>
                  <For each={Object.values(sortItems)}>
                    {(s) => (
                      <RadioGroup.Item value={s.value}>
                        <RadioGroup.ItemInput />
                        <RadioGroup.ItemLabel
                          class={clsx(
                            selectItemStyle,
                            'data-[selected]:bg-transparency-primary-selected data-[selected]:text-primary-main',
                          )}
                        >
                          <RadioGroup.ItemIndicator forceMount>
                            <RadioIcon selected={props.sort === s.value} />
                          </RadioGroup.ItemIndicator>
                          {s.label}
                        </RadioGroup.ItemLabel>
                        <RadioGroup.ItemDescription />
                      </RadioGroup.Item>
                    )}
                  </For>
                </ItemsContainer>
              </FilterItemContainer>
            )}
          />
          <FilterItemContainer class="grid-area-[noapp]">
            <Checkbox.Root checked={props.includeNoApp} onChange={props.setIncludeNoApp}>
              <Checkbox.Input />
              <Checkbox.Label class={selectItemStyle}>
                <Checkbox.Indicator forceMount class="size-6">
                  <CheckBoxIcon checked={props.includeNoApp} />
                </Checkbox.Indicator>
                アプリを持たないリポジトリを表示
              </Checkbox.Label>
            </Checkbox.Root>
          </FilterItemContainer>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  )
}

export default AppsFilter
