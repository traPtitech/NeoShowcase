import { type Component, createMemo, createSignal, For, Show } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '/@/components/styled-components'
import { Button } from '/@/components/UI/Button'
import UserAvatar from '/@/components/UI/UserAvater'
import useModal from '/@/libs/useModal'
import ModalDeleteConfirm from '../UI/ModalDeleteConfirm'
import { TextField } from '../UI/TextField'

const UserPlaceholder = styled(
  'div',
  'h4-medium flex h-full w-full flex-col items-center justify-center gap-6 rounded-lg border border-ui-border bg-ui-primary px-5 py-4 text-text-black',
)

const UsersContainer = styled('div', 'h-auto max-h-full w-full overflow-y-auto rounded-lg border border-ui-border')

const UserRowContainer = styled('div', 'flex w-full items-center gap-2 px-5 py-4')

const UserName = styled('div', 'w-full text-medium text-text-black')

const AddOwners: Component<{
  nonOwners: User[]
  addOwner: (user: User) => void
}> = (props) => {
  const [searchUserQuery, setSearchUserQuery] = createSignal('')
  const filteredUsers = createMemo(() => {
    if (searchUserQuery() === '') {
      return props.nonOwners
    }
    return props.nonOwners.filter(({ name }) => name.toLowerCase().includes(searchUserQuery().toLowerCase()))
  })

  return (
    <div class="flex h-full max-h-full w-full flex-col gap-4">
      <TextField
        placeholder="Search UserID"
        leftIcon={<div class="i-material-symbols:search shrink-0 text-2xl/6" />}
        value={searchUserQuery()}
        onInput={(e) => setSearchUserQuery(e.currentTarget.value)}
      />
      <Show
        when={filteredUsers().length !== 0}
        fallback={
          <UserPlaceholder>
            <div class="i-material-symbols:search text-20/20 shrink-0" />
            No Users Found
          </UserPlaceholder>
        }
      >
        <div class="border-ui-border min-h-0 flex-1 overflow-y-auto rounded-lg border">
          <For each={filteredUsers()} keyed={(user) => user.id}>
            {(user) => (
              <div class="border-ui-border w-full [&:not(:last-child)]:border-b">
                <UserRowContainer>
                  <UserAvatar user={user()} size={32} />
                  <UserName>{user().name}</UserName>
                  <Button variants="ghost" size="small" onClick={() => props.addOwner(user())}>
                    Add
                  </Button>
                </UserRowContainer>
              </div>
            )}
          </For>
        </div>
      </Show>
    </div>
  )
}

const OwnerRow: Component<{
  user: User
  deleteOwner?: (user: User) => void
}> = (props) => {
  const { Modal: DeleteUserModal, open: openDeleteUserModal, close: closeDeleteUserModal } = useModal()

  return (
    <>
      <div class="border-ui-border [&:not(:last-child)]:border-b">
        <UserRowContainer>
          <UserAvatar user={props.user} size={32} />
          <UserName>{props.user.name}</UserName>
          <Show when={props.deleteOwner !== undefined}>
            <Button variants="textError" size="small" onClick={openDeleteUserModal}>
              Delete
            </Button>
          </Show>
        </UserRowContainer>
      </div>
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
  const filteredOwners = createMemo(() => {
    if (searchUserQuery() === '') {
      return props.owners
    }
    return props.owners.filter(({ name }) => name.toLowerCase().includes(searchUserQuery().toLowerCase()))
  })

  const nonOwners = createMemo(() => props.users.filter((u) => !props.owners.some((o) => o.id === u.id)))
  const { Modal: AddUserModal, open: openAddUserModal } = useModal({
    showCloseButton: true,
  })

  return (
    <>
      <div class="flex items-center gap-4">
        <TextField
          placeholder="Search UserID"
          leftIcon={<div class="i-material-symbols:search shrink-0 text-2xl/6" />}
          value={searchUserQuery()}
          onInput={(e) => setSearchUserQuery(e.currentTarget.value)}
        />
        <Button
          variants="primary"
          size="medium"
          leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
          onClick={openAddUserModal}
          disabled={!props.hasPermission}
          tooltip={{
            content: !props.hasPermission ? '設定を変更するにはオーナーになる必要があります' : undefined,
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
      </div>
      <UsersContainer>
        <For each={filteredOwners()} keyed={(owner) => owner.id}>
          {(owner) => (
            <OwnerRow user={owner()} deleteOwner={props.hasPermission ? props.handleDeleteOwner : undefined} />
          )}
        </For>
        <Show when={filteredOwners().length === 0}>
          <UserPlaceholder>No Owners Found</UserPlaceholder>
        </Show>
      </UsersContainer>
    </>
  )
}

export default OwnerList
