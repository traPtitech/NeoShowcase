import { type Component, For } from 'solid-js'
import { allOrigins, originToIcon, type RepositoryOrigin } from '/@/libs/application'
import { FilterCheckbox, FilterMenu, toggled } from '../app/FilterMenu'

const ReposFilter: Component<{
  origin: RepositoryOrigin[]
  setOrigin: (origin: RepositoryOrigin[]) => void
}> = (props) => (
  <FilterMenu filtered={props.origin.length !== allOrigins.length} contentClass="flex gap-2">
    <div class="text-bold text-text-black flex flex-col gap-2">
      Origin
      <div class="flex w-full flex-col">
        <For each={allOrigins}>
          {(s) => (
            <FilterCheckbox
              checked={props.origin.includes(s.value)}
              onChange={(on) => props.setOrigin(toggled(props.origin, s.value, on))}
            >
              {originToIcon(s.value)}
              {s.label}
            </FilterCheckbox>
          )}
        </For>
      </div>
    </div>
  </FilterMenu>
)

export default ReposFilter
