import { RadioGroup } from '@kobalte/core/radio-group'
import { type Component, For } from 'solid-js'
import { RadioIcon } from '/@/components/UI/RadioIcon'
import {
  type ApplicationState,
  allOrigins,
  allStatuses,
  applicationState,
  originToIcon,
  type RepositoryOrigin,
  type RepoWithApp,
  repositoryURLToOrigin,
  sortItems,
} from '/@/libs/application'
import { AppStatusIcon } from './AppStatusIcon'
import { FilterCheckbox, FilterMenu, filterItemStyle, toggled } from './FilterMenu'

const filterGroupStyle = 'flex flex-col gap-2 text-bold text-text-black'

const AppsFilter: Component<{
  allRepoWithApps: RepoWithApp[]
  statuses: ApplicationState[]
  setStatues: (statuses: ApplicationState[]) => void
  origin: RepositoryOrigin[]
  setOrigin: (origin: RepositoryOrigin[]) => void
  sort: keyof typeof sortItems
  setSort: (sort: keyof typeof sortItems) => void
  includeNoApp: boolean
  setIncludeNoApp: (includeNoApp: boolean) => void
}> = (props) => {
  const filtered = () =>
    props.statuses.length !== allStatuses.length || props.origin.length !== allOrigins.length || props.includeNoApp

  const appCountByStatus = (status: ApplicationState): number =>
    props.allRepoWithApps
      .filter((repo) => props.origin.includes(repositoryURLToOrigin(repo.repo.url)))
      .flatMap((repo) => repo.apps.filter((app) => applicationState(app) === status)).length

  const repoCountByOrigin = (origin: RepositoryOrigin): number =>
    props.allRepoWithApps
      .filter((repo) => repositoryURLToOrigin(repo.repo.url) === origin)
      .filter(
        (repo) =>
          props.includeNoApp || repo.apps.filter((app) => props.statuses.includes(applicationState(app))).length > 0,
      ).length

  return (
    <FilterMenu
      filtered={filtered()}
      contentClass="grid max-w-[--kb-popper-content-available-width] grid-cols-[repeat(3,auto)] grid-rows-[1fr_auto] gap-2 overflow-x-auto"
      contentStyle={{
        'grid-template-areas': `
          "status provider sort"
          "status noapp noapp"
        `,
      }}
    >
      <div class={[filterGroupStyle, 'grid-area-[status]']}>
        App Status
        <div class="flex w-full flex-col">
          <For each={allStatuses}>
            {(s) => (
              <FilterCheckbox
                checked={props.statuses.includes(s.value)}
                onChange={(on) => props.setStatues(toggled(props.statuses, s.value, on))}
              >
                <AppStatusIcon state={s.value} hideTooltip />
                <span>
                  {s.label} ({appCountByStatus(s.value)})
                </span>
              </FilterCheckbox>
            )}
          </For>
        </div>
      </div>
      <div class={[filterGroupStyle, 'grid-area-[provider]']}>
        Repo Origin
        <div class="flex w-full flex-col">
          <For each={allOrigins}>
            {(s) => (
              <FilterCheckbox
                checked={props.origin.includes(s.value)}
                onChange={(on) => props.setOrigin(toggled(props.origin, s.value, on))}
              >
                {originToIcon(s.value)}
                <span>
                  {s.label} ({repoCountByOrigin(s.value)})
                </span>
              </FilterCheckbox>
            )}
          </For>
        </div>
      </div>
      <RadioGroup
        class={[filterGroupStyle, 'grid-area-[sort]']}
        value={props.sort}
        onChange={(v) => props.setSort(v as keyof typeof sortItems)}
      >
        <RadioGroup.Label>Sort</RadioGroup.Label>
        <div class="flex w-full flex-col">
          <For each={Object.values(sortItems)}>
            {(s) => (
              <RadioGroup.Item value={s.value}>
                <RadioGroup.ItemInput />
                <RadioGroup.ItemLabel
                  class={[
                    filterItemStyle,
                    'data-[checked]:bg-transparency-primary-selected data-[checked]:text-primary-main',
                  ]}
                >
                  <RadioGroup.ItemIndicator forceMount>
                    <RadioIcon selected={props.sort === s.value} />
                  </RadioGroup.ItemIndicator>
                  {s.label}
                </RadioGroup.ItemLabel>
              </RadioGroup.Item>
            )}
          </For>
        </div>
      </RadioGroup>
      <div class={[filterGroupStyle, 'grid-area-[noapp]']}>
        <FilterCheckbox checked={props.includeNoApp} onChange={props.setIncludeNoApp}>
          アプリを持たないリポジトリを表示
        </FilterCheckbox>
      </div>
    </FilterMenu>
  )
}

export default AppsFilter
