import { style } from '@macaron-css/core'
import { styled } from '@macaron-css/solid'
import { createVirtualizer } from '@tanstack/solid-virtual'
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

    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
  },
})
const bordered = style({
  selectors: {
    '&:not(:last-child)': {
      borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
    },
  }
})
const UserRowContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
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
    height: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
    alignItems: 'center',
    justifyContent: 'center',

    border: `1px solid ${colorVars.semantic.ui.border}`,
    borderRadius: '8px',
    background: colorVars.semantic.ui.primary,
    color: colorVars.semantic.text.black,
    ...textVars.h4.medium,
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

  const [containerRef, setContainerRef] = createSignal<HTMLDivElement | null>(null)
  const virtualizer = createMemo(() =>
    createVirtualizer({
      count: filteredUsers().length,
      getScrollElement: containerRef,
      estimateSize: () => 64,
    }),
  )
  const items = () => virtualizer().getVirtualItems()

  return (
    <AddOwnersContainer>
      <TextField
        placeholder="Search UserID"
        leftIcon={<MaterialSymbols>search</MaterialSymbols>}
        value={searchUserQuery()}
        onInput={(e) => setSearchUserQuery(e.currentTarget.value)}
      />
      <Show
        when={filteredUsers().length !== 0}
        fallback={
          <UserPlaceholder>
            <MaterialSymbols displaySize={80}>search</MaterialSymbols>
            No Users Found
          </UserPlaceholder>
        }
      >
        <UsersContainer ref={setContainerRef}>
          <div
            style={{
              width: '100%',
              height: `${virtualizer().getTotalSize()}px`,
              position: 'relative',
            }}
          >
            <For each={items() ?? []}>
              {(vRow) => (
                <div
                  data-index={vRow.index}
                  style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '100%',
                    height: `${vRow.size}px`,
                    transform: `translateY(${vRow.start}px)`,
                  }}
                  class={bordered}
                >
                  <UserRowContainer>
                    <UserAvatar user={filteredUsers()[vRow.index]} size={32} />
                    <UserName>{filteredUsers()[vRow.index].name}</UserName>
                    <Button variants="ghost" size="small" onClick={() => props.addOwner(filteredUsers()[vRow.index])}>
                      Add
                    </Button>
                  </UserRowContainer>
                </div>
              )}
            </For>
          </div>
        </UsersContainer>
      </Show>
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
