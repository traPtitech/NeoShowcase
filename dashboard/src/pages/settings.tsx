import { Field, Form, type SubmitHandler, createFormStore, reset, valiForm } from '@modular-forms/solid'
import { Title } from '@solidjs/meta'
import { type Component, For, Show, createResource } from 'solid-js'
import toast from 'solid-toast'
import * as v from 'valibot'
import { Button } from '/@/components/UI/Button'
import { styled } from '/@/components/styled-components'
import { client, handleAPIError } from '/@/libs/api'
import type { DeleteUserKeyRequest, UserKey } from '../api/neoshowcase/protobuf/gateway_pb'
import ModalDeleteConfirm from '../components/UI/ModalDeleteConfirm'
import { TextField } from '../components/UI/TextField'
import { DataTable } from '../components/layouts/DataTable'
import { MainViewContainer } from '../components/layouts/MainView'
import { WithNav } from '../components/layouts/WithNav'
import { List } from '../components/templates/List'
import { Nav } from '../components/templates/Nav'
import { dateHuman } from '../libs/format'
import useModal from '../libs/useModal'

const SshKeyName = styled('div', 'h4-bold truncate text-text-black')
const SshKeyRowValue = styled('div', 'truncate text-text-black text-text-medium')
const SshKeyAddedAt = styled('div', 'truncate text-text-grey text-text-regular')

const SshKeyRow: Component<{ key: UserKey; refetchKeys: () => void }> = (props) => {
  const { Modal, open, close } = useModal()
  const handleDeleteKey = async (keyID: DeleteUserKeyRequest['keyId']) => {
    try {
      await client.deleteUserKey({ keyId: keyID })
      toast.success('公開鍵を削除しました')
      props.refetchKeys()
    } catch (e) {
      handleAPIError(e, '公開鍵の削除に失敗しました')
    }
  }

  return (
    <>
      <div class="flex w-full items-center gap-2 overflow-x-hidden bg-ui-primary px-5 py-4 last:border-b-none">
        <div class="flex w-full flex-col overflow-x-hidden">
          <SshKeyName>{props.key.name === '' ? '(Name not set)' : props.key.name}</SshKeyName>
          <SshKeyRowValue>{props.key.publicKey}</SshKeyRowValue>
          <Show
            when={props.key.createdAt && props.key.createdAt?.seconds !== 0n ? props.key.createdAt : false}
            fallback={<SshKeyAddedAt>Added on ----</SshKeyAddedAt>}
          >
            {(createdAt) => (
              <SshKeyAddedAt>Added on {dateHuman(createdAt())}</SshKeyAddedAt>
            )}
          </Show>
        </div>
        <Button variants="textError" size="medium" onClick={open}>
          Delete
        </Button>
      </div>
      <Modal.Container>
        <Modal.Header>Delete SSH Key</Modal.Header>
        <Modal.Body>
          <ModalDeleteConfirm>
            <div class="flex w-full flex-col">
              <SshKeyName>{props.key.name === '' ? '(Name not set)' : props.key.name}</SshKeyName>
              {props.key.publicKey}
              <Show
                when={ props.key.createdAt?.seconds !==0n ? props.key.createdAt : false}
                fallback={<SshKeyAddedAt>Added on ----</SshKeyAddedAt>}
              >
                {(createdAt) => (
                  <SshKeyAddedAt>
                    Added on {dateHuman(createdAt())}
                  </SshKeyAddedAt> 
                )}
              </Show>
            </div>
          </ModalDeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button variants="text" size="medium" onClick={close}>
            No, Cancel
          </Button>
          <Button variants="primaryError" size="medium" onClick={() => handleDeleteKey(props.key.id)}>
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </>
  )
}

const SshKeys: Component<{ keys: UserKey[]; refetchKeys: () => void }> = (props) => {
  return (
    <List.Container>
      <For each={props.keys}>{(key) => <SshKeyRow key={key} refetchKeys={props.refetchKeys} />}</For>
    </List.Container>
  )
}

const userKeyRequestSchema = v.object({
  name: v.pipe(v.string(), v.nonEmpty('Enter Name')),
  publicKey: v.pipe(v.string(), v.nonEmpty('Enter SSH Public Key')),
})
type UserKeyRequestInput = v.InferInput<typeof userKeyRequestSchema>

export default () => {
  const [userKeys, { refetch: refetchKeys }] = createResource(() => client.getUserKeys({}))
  const { Modal: AddNewKeyModal, open: newKeyOpen, close: newKeyClose } = useModal()

  const formStore = createFormStore<UserKeyRequestInput>({
    validate: valiForm(userKeyRequestSchema),
    initialValues: {
      name: '',
      publicKey: '',
    },
  })

  const handleSubmit: SubmitHandler<UserKeyRequestInput> = async (values) => {
    try {
      await client.createUserKey(values)
      toast.success('公開鍵を追加しました')
      newKeyClose()
      reset(formStore)
      refetchKeys()
    } catch (e) {
      handleAPIError(e, '公開鍵の追加に失敗しました')
    }
  }

  const AddNewSSHKeyButton = () => (
    <Button
      variants="primary"
      size="medium"
      leftIcon={<div class="i-material-symbols:add shrink-0 text-2xl/6" />}
      onClick={newKeyOpen}
    >
      Add New SSH Key
    </Button>
  )

  return (
    <>
      <Title>Settings - NeoShowcase</Title>
      <WithNav.Container>
        <WithNav.Navs>
          <Nav title="Settings" />
        </WithNav.Navs>
        <WithNav.Body>
          <Show when={userKeys()}>
            {(keysData) => (
              <MainViewContainer>
                <DataTable.Container>
                  <div class="flex w-full items-end justify-between">
                    <DataTable.Titles>
                      <DataTable.Title>SSH Public Keys</DataTable.Title>
                      <DataTable.SubTitle>
                        SSH鍵はruntimeアプリケーションのコンテナにssh接続するときに使います
                      </DataTable.SubTitle>
                    </DataTable.Titles>
                    <Show when={userKeys()?.keys.length !== 0}>
                      <AddNewSSHKeyButton />
                    </Show>
                  </div>
                  <Show
                    when={keysData().keys.length > 0}
                    fallback={
                      <List.Container>
                        <List.PlaceHolder>
                          <div class="i-material-symbols:key-off-outline shrink-0 text-20/20" />
                          No Keys Registered
                          <AddNewSSHKeyButton />
                        </List.PlaceHolder>
                      </List.Container>
                    }
                  >
                    <SshKeys keys={keysData().keys} refetchKeys={refetchKeys} />
                  </Show>
                </DataTable.Container>
              </MainViewContainer>
            )}
          </Show>
        </WithNav.Body>
      </WithNav.Container>
      <AddNewKeyModal.Container>
        <Form of={formStore} onSubmit={handleSubmit}>
          <AddNewKeyModal.Header>Add New SSH Key</AddNewKeyModal.Header>
          <AddNewKeyModal.Body>
            <div class="flex w-full flex-col gap-2">
              <Field of={formStore} name="name">
                {(field, fieldProps) => (
                  <TextField label="Name" required {...fieldProps} value={field.value} error={field.error} />
                )}
              </Field>
              <Field of={formStore} name="publicKey">
                {(field, fieldProps) => (
                  <TextField
                    label="Key"
                    required
                    multiline
                    placeholder="ssh-ed25519 AAA..."
                    {...fieldProps}
                    value={field.value}
                    error={field.error}
                  />
                )}
              </Field>
            </div>
          </AddNewKeyModal.Body>
          <AddNewKeyModal.Footer>
            <Button variants="text" size="medium" type="button" onClick={newKeyClose}>
              Cancel
            </Button>
            <Button
              variants="primary"
              size="medium"
              type="submit"
              disabled={formStore.invalid || formStore.submitting}
              loading={formStore.submitting}
            >
              Add
            </Button>
          </AddNewKeyModal.Footer>
        </Form>
      </AddNewKeyModal.Container>
    </>
  )
}
