import { styled } from '@macaron-css/solid'
import Fuse from 'fuse.js'
import { FlowComponent, For, JSX, Show, createMemo, createSignal } from 'solid-js'
import { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { vars } from '/@/theme'
import { InputBar } from '/@/components/Input'
import UserAvatar from '/@/components/UserAvatar'

const UserSearchContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
  },
})
const UserRow = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: '16px',
  },
})
const UserContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    gap: '8px',
    alignItems: 'center',
  },
})
const UserName = styled('div', {
  base: {
    fontSize: '16px',
    color: vars.text.black1,
  },
})
const UsersList = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
    overflowY: 'auto',
    maxHeight: '400px',
  },
})
const NoUsersFound = styled('div', {
  base: {
    color: vars.text.black2,
    textAlign: 'center',
  },
})

export const UserSearch: FlowComponent<
  {
    users: User[]
  },
  (user: User) => JSX.Element
> = (props) => {
  const [userSearchQuery, setUserSearchQuery] = createSignal('')

  // users()の更新時にFuseインスタンスを再生成する
  const fuse = createMemo(
    () =>
      new Fuse(props.users, {
        keys: ['name'],
      }),
  )
  // userSearchQuery()の更新時に検索を実行する
  const userSearchResults = createMemo(() => {
    // 検索クエリが空の場合は全ユーザーを表示する
    if (userSearchQuery() === '') {
      return props.users
    } else {
      return fuse()
        .search(userSearchQuery())
        .map((result) => result.item)
    }
  })

  return (
    <UserSearchContainer>
      <InputBar
        type='text'
        value={userSearchQuery()}
        onInput={(e) => setUserSearchQuery(e.target.value)}
        placeholder='search users...'
      />
      <UsersList>
        <For each={userSearchResults()}>
          {(user) => (
            <UserRow>
              <UserContainer>
                <UserAvatar user={user} size={32} />
                <UserName>{user.name}</UserName>
              </UserContainer>
              {props.children(user)}
            </UserRow>
          )}
        </For>
        <Show when={userSearchResults().length === 0}>
          <NoUsersFound>No Users Found</NoUsersFound>
        </Show>
      </UsersList>
    </UserSearchContainer>
  )
}
