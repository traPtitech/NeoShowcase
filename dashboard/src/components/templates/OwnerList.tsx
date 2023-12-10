import { styled } from '@macaron-css/solid'
import Fuse from 'fuse.js'
import { Component, For, Show, createMemo, createSignal } from 'solid-js'
import { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import UserAvatar from '/@/components/UI/UserAvater'
import useModal from '/@/libs/useModal'
import { colorVars, textVars } from '/@/theme'
import ModalDeleteConfirm from '../UI/ModalDeleteConfirm'
import { TextField } from '../UI/TextField'

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
    height: '100%',
    maxHeight: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
  },
})
const UsersContainer = styled('div', {
  base: {
    width: '100%',
    height: 'auto',
    maxHeight: '100%',
    overflowY: 'auto',
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
    }
    return fuse()
      .search(searchUserQuery())
      .map((result) => result.item)
  })

  return (
    <AddOwnersContainer>
      <TextField
        placeholder="Search UserID"
        leftIcon={<MaterialSymbols>search</MaterialSymbols>}
        value={searchUserQuery()}
        onInput={(e) => setSearchUserQuery(e.currentTarget.value)}
      />
      <UsersContainer>
        <For each={filteredUsers()}>
          {(user) => (
            <UserRowContainer>
              <UserAvatar user={user} size={32} />
              <UserName>{user.name}</UserName>
              <Button variants="ghost" size="small" onClick={() => props.addOwner(user)}>
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

const OwnerRow: Component<{
  user: User
  deleteOwner?: (user: User) => void
}> = (props) => {
  const { Modal: DeleteUserModal, open: openDeleteUserModal, close: closeDeleteUserModal } = useModal()

  return (
    <>
      <UserRowContainer>
        <UserAvatar user={props.user} size={32} />
        <UserName>{props.user.name}</UserName>
        <Show when={props.deleteOwner !== undefined}>
          <Button variants="textError" size="small" onClick={openDeleteUserModal}>
            Delete
          </Button>
        </Show>
      </UserRowContainer>
      <DeleteUserModal.Container>
        <DeleteUserModal.Header>Delete Owner</DeleteUserModal.Header>
        <DeleteUserModal.Body>
          <ModalDeleteConfirm>
            <UserAvatar user={props.user} size={32} />
            <UserName>{props.user.name}</UserName>
          </ModalDeleteConfirm>
        </DeleteUserModal.Body>
        <DeleteUserModal.Footer>
          <Button variants="text" size="medium" onClick={closeDeleteUserModal} type="button">
            No, Cancel
          </Button>
          <Button
            variants="primaryError"
            size="medium"
            onClick={() => {
              props.deleteOwner?.(props.user)
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

const OwnerList: Component<{
  owners: User[]
  users: User[]
  handleAddOwner: (user: User) => Promise<void>
  handleDeleteOwner: (user: User) => Promise<void>
  hasPermission: boolean
}> = (props) => {
  const [searchUserQuery, setSearchUserQuery] = createSignal('')
  const fuse = createMemo(
    () =>
      new Fuse(props.owners, {
        keys: ['name'],
      }),
  )
  const filteredOwners = createMemo(() => {
    if (searchUserQuery() === '') {
      return props.owners
    }
    return fuse()
      .search(searchUserQuery())
      .map((result) => result.item)
  })

  const nonOwners = createMemo(() => props.users.filter((u) => !props.owners.some((o) => o.id === u.id)))
  const { Modal: AddUserModal, open: openAddUserModal } = useModal({
    showCloseButton: true,
  })

  return (
    <>
      <SearchUserRow>
        <TextField
          placeholder="Search UserID"
          leftIcon={<MaterialSymbols>search</MaterialSymbols>}
          value={searchUserQuery()}
          onInput={(e) => setSearchUserQuery(e.currentTarget.value)}
        />
        <Button
          variants="primary"
          size="medium"
          leftIcon={<MaterialSymbols>add</MaterialSymbols>}
          onClick={openAddUserModal}
          disabled={!props.hasPermission}
          tooltip={{
            props: {
              content: !props.hasPermission ? '設定を変更するにはオーナーになる必要があります' : undefined,
            },
          }}
        >
          Add Owners
        </Button>
        <AddUserModal.Container fit={false}>
          <AddUserModal.Header>Add Owner</AddUserModal.Header>
          <AddUserModal.Body fit={false}>
            <AddOwners addOwner={props.handleAddOwner} nonOwners={nonOwners()} />
          </AddUserModal.Body>
        </AddUserModal.Container>
      </SearchUserRow>
      <UsersContainer>
        <For each={filteredOwners()}>
          {(owner) => <OwnerRow user={owner} deleteOwner={props.hasPermission ? props.handleDeleteOwner : undefined} />}
        </For>
        <Show when={filteredOwners().length === 0}>
          <UserPlaceholder>No Owners Found</UserPlaceholder>
        </Show>
      </UsersContainer>
    </>
  )
}

export default OwnerList
