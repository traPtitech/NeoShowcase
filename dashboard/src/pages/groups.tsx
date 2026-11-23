import { Title } from '@solidjs/meta'
import { createMemo, For, isPending, Loading, Show } from 'solid-js'
import { GetApplicationGroupsRequest_Scope } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { MainViewContainer } from '/@/components/layouts/MainView'
import PendingContainer from '/@/components/layouts/PendingContainer'
import { WithNav } from '/@/components/layouts/WithNav'
import { List } from '/@/components/templates/List'
import { Nav } from '/@/components/templates/Nav'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import { TabRound } from '/@/components/UI/TabRound'
import { client } from '/@/libs/api'
import { createSessionSignal } from '/@/libs/storage'

const scopeItems = [
  { value: GetApplicationGroupsRequest_Scope.MINE, label: 'My Groups' },
  { value: GetApplicationGroupsRequest_Scope.ALL, label: 'All Groups' },
]

export default () => {
  const [scope, setScope] = createSessionSignal('groups-scope', GetApplicationGroupsRequest_Scope.MINE)
  const groups = createMemo(() =>
    client
      .getApplicationGroups({ scope: scope() })
      .then((res) => res.groups.toSorted((a, b) => a.name.localeCompare(b.name))),
  )

  return (
    <WithNav.Container>
      <Title>Groups - NeoShowcase</Title>
      <WithNav.Navs>
        <Nav title="Groups" />
        <WithNav.Tabs>
          <For each={scopeItems}>
            {(s) => (
              <TabRound state={s.value === scope() ? 'active' : 'default'} onClick={() => setScope(s.value)}>
                <div class="i-material-symbols:folder-outline shrink-0 text-2xl/6" />
                {s.label}
              </TabRound>
            )}
          </For>
          <a href="/groups/new" class="ml-auto">
            <Button
              variants="primary"
              size="medium"
              leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
            >
              Add New Group
            </Button>
          </a>
        </WithNav.Tabs>
      </WithNav.Navs>
      <WithNav.Body>
        <MainViewContainer>
          <Loading fallback={<Skeleton height={200} />}>
            <PendingContainer isPending={isPending(() => groups())}>
              <List.Container>
                <Show
                  when={groups().length > 0}
                  fallback={
                    <List.PlaceHolder>
                      <div class="i-material-symbols:folder-outline text-20/20 shrink-0" />
                      No Groups
                    </List.PlaceHolder>
                  }
                >
                  <For each={groups()}>
                    {(group) => (
                      <List.Row>
                        <div class="i-material-symbols:folder-outline shrink-0 text-2xl/6" />
                        <a href={`/groups/${group.id}`} class="h4-bold text-text-black truncate">
                          {group.name}
                        </a>
                      </List.Row>
                    )}
                  </For>
                </Show>
              </List.Container>
            </PendingContainer>
          </Loading>
        </MainViewContainer>
      </WithNav.Body>
    </WithNav.Container>
  )
}
