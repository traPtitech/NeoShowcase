import { Application, User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import { TextInput } from '/@/components/UI/TextInput'
import UserAvatar from '/@/components/UserAvatar'
import { DataTable } from '/@/components/layouts/DataTable'
import { client, handleAPIError } from '/@/libs/api'
import { userFromId, users } from '/@/libs/useAllUsers'
import useModal from '/@/libs/useModal'
import { useApplicationData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import Fuse from 'fuse.js'
import { Component, For, Show, createMemo, createSignal } from 'solid-js'
import toast from 'solid-toast'

const SearchUserRow = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '16px',
  },
})
const AddOwnersContainer = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    maxHeight: '100%',
    display: 'grid',
    gridTemplateRows: 'auto 1fr',
    gap: '16px',
  },
})
const UsersContainer = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    maxHeight: '100%',
    overflowY: 'auto',
    scrollbarGutter: 'stable',

    display: 'flex',
    flexDirection: 'column',

    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
  },
})
const UserRowContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

    selectors: {
      '&:not(:last-child)': {
        borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
      },
    },
  },
})
const UserName = styled('div', {
  base: {
    width: '100%',
    color: colorVars.semantic.text.black,
    ...textVars.text.medium,
  },
})
const UserPlaceholder = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: colorVars.semantic.text.grey,
    ...textVars.text.medium,
  },
})
const DeleteConfirm = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    borderRadius: '8px',
    background: colorVars.semantic.ui.secondary,
  },
})

const AddOwners: Component<{
  nonOwners: User[]
  addOwner: (user: User) => void
}> = (props) => {
  const [searchUserQuery, setSearchUserQuery] = createSignal('')
  const fuse = createMemo(
    () =>
      new Fuse(props.nonOwners, {
        keys: ['name'],
      }),
  )
  const filteredUsers = createMemo(() => {
    if (searchUserQuery() === '') {
      return props.nonOwners
    } else {
      return fuse()
        .search(searchUserQuery())
        .map((result) => result.item)
    }
  })

  return (
    <AddOwnersContainer>
      <TextInput
        placeholder="Search UserID"
        leftIcon={<MaterialSymbols>search</MaterialSymbols>}
        value={searchUserQuery()}
        onInput={(e) => setSearchUserQuery(e.target.value)}
      />
      <UsersContainer>
        <For each={filteredUsers()}>
          {(user) => (
            <UserRowContainer>
              <UserAvatar user={user} size={32} />
              <UserName>{user.name}</UserName>
              <Button color="ghost" size="small" onClick={() => props.addOwner(user)}>
                Add
              </Button>
            </UserRowContainer>
          )}
        </For>
        <Show when={filteredUsers().length === 0}>
          <UserPlaceholder>No Users Found</UserPlaceholder>
        </Show>
      </UsersContainer>
    </AddOwnersContainer>
  )
}

const OwnerRow: Component<{ user: User; deleteOwner: (user: User) => void }> = (props) => {
  const { Modal: DeleteUserModal, open: openDeleteUserModal, close: closeDeleteUserModal } = useModal()

  return (
    <>
      <UserRowContainer>
        <UserAvatar user={props.user} size={32} />
        <UserName>{props.user.name}</UserName>
        <Button color="textError" size="small" onClick={openDeleteUserModal}>
          Delete
        </Button>
      </UserRowContainer>
      <DeleteUserModal.Container>
        <DeleteUserModal.Header>Delete Owner</DeleteUserModal.Header>
        <DeleteUserModal.Body>
          <DeleteConfirm>
            <UserAvatar user={props.user} size={32} />
            <UserName>{props.user.name}</UserName>
          </DeleteConfirm>
        </DeleteUserModal.Body>
        <DeleteUserModal.Footer>
          <Button color="text" size="medium" onClick={closeDeleteUserModal} type="button">
            No, Cancel
          </Button>
          <Button
            color="primaryError"
            size="medium"
            onClick={() => {
              props.deleteOwner(props.user)
              closeDeleteUserModal()
            }}
            type="button"
          >
            Yes, Delete
          </Button>
        </DeleteUserModal.Footer>
      </DeleteUserModal.Container>
    </>
  )
}

const OwnerConfig: Component<{
  app: Application
  users: User[]
  refetchApp: () => void
}> = (props) => {
  const [searchUserQuery, setSearchUserQuery] = createSignal('')
  const owners = () => props.app.ownerIds.map(userFromId)
  const fuse = createMemo(
    () =>
      new Fuse(owners(), {
        keys: ['name'],
      }),
  )
  const filteredOwners = createMemo(() => {
    if (searchUserQuery() === '') {
      return owners()
    } else {
      return fuse()
        .search(searchUserQuery())
        .map((result) => result.item)
    }
  })
  const handleDeleteOwner = async (user: User) => {
    const newOwnerIds = props.app.ownerIds.filter((id) => id !== user.id)
    try {
      await client.updateApplication({
        id: props.app.id,
        ownerIds: { ownerIds: newOwnerIds },
      })
      toast.success('アプリケーションオーナーを削除しました')
      props.refetchApp()
    } catch (e) {
      handleAPIError(e, 'アプリケーションオーナーの削除に失敗しました')
    }
  }

  const nonOwners = createMemo(() => props.users.filter((u) => !owners().some((o) => o.id === u.id)))
  const { Modal: AddUserModal, open: openAddUserModal } = useModal({
    showCloseButton: true,
  })
  const handleAddOwner = async (user: User) => {
    const newOwnerIds = props.app.ownerIds.concat(user.id)
    try {
      await client.updateApplication({
        id: props.app.id,
        ownerIds: { ownerIds: newOwnerIds },
      })
      toast.success('アプリケーションオーナーを追加しました')
      props.refetchApp()
    } catch (e) {
      handleAPIError(e, 'アプリケーションオーナーの追加に失敗しました')
    }
  }

  return (
    <>
      <SearchUserRow>
        <TextInput
          placeholder="Search UserID"
          leftIcon={<MaterialSymbols>search</MaterialSymbols>}
          value={searchUserQuery()}
          onInput={(e) => setSearchUserQuery(e.target.value)}
        />
        <Button
          color="primary"
          size="medium"
          leftIcon={<MaterialSymbols>add</MaterialSymbols>}
          onClick={openAddUserModal}
        >
          Add Owners
        </Button>
        <AddUserModal.Container>
          <AddUserModal.Header>Add Owner</AddUserModal.Header>
          <AddUserModal.Body>
            <AddOwners addOwner={handleAddOwner} nonOwners={nonOwners()} />
          </AddUserModal.Body>
        </AddUserModal.Container>
      </SearchUserRow>
      <UsersContainer>
        <For each={filteredOwners()}>{(owner) => <OwnerRow user={owner} deleteOwner={handleDeleteOwner} />}</For>
        <Show when={filteredOwners().length === 0}>
          <UserPlaceholder>No Owners Found</UserPlaceholder>
        </Show>
      </UsersContainer>
    </>
  )
}

export default () => {
  const { app, refetchApp } = useApplicationData()
  const loaded = () => !!(app() && users())

  return (
    <DataTable.Container>
      <DataTable.Title>Owner</DataTable.Title>
      <DataTable.SubTitle>
        オーナーはアプリ設定の変更, アプリログ/メトリクスの閲覧, 環境変数の閲覧, ビルドログの閲覧が可能になります
      </DataTable.SubTitle>
      <Show when={loaded()}>
        <OwnerConfig app={app()} users={users()} refetchApp={refetchApp} />
      </Show>
    </DataTable.Container>
  )
}
