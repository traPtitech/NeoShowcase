import { createForm, Field, Form, reset } from '@formisch/solid/v2'
import { Title } from '@solidjs/meta'
import { type Component, createMemo, For, Show } from 'solid-js'
import * as v from 'valibot'
import { styled } from '/@/components/styled-components'
import { Button } from '/@/components/UI/Button'
import { client, getUserKeys, handleAPIError, revalidateUserKeys } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import type { DeleteUserKeyRequest, UserKey } from '../api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '../components/layouts/DataTable'
import { MainViewContainer } from '../components/layouts/MainView'
import { SectionBoundary } from '../components/layouts/SectionBoundary'
import { WithNav } from '../components/layouts/WithNav'
import { List } from '../components/templates/List'
import { Nav } from '../components/templates/Nav'
import ModalDeleteConfirm from '../components/UI/ModalDeleteConfirm'
import { TextField } from '../components/UI/TextField'
import { dateHuman } from '../libs/format'
import useModal from '../libs/useModal'

const SshKeyName = styled('div', 'h4-bold truncate text-text-black')
const SshKeyRowValue = styled('div', 'truncate text-text-black text-text-medium')
const SshKeyAddedAt = styled('div', 'truncate text-text-grey text-text-regular')

const AddedOn: Component<{ key: UserKey }> = (props) => (
  <Show
    when={props.key.createdAt && props.key.createdAt.seconds !== 0n ? props.key.createdAt : undefined}
    fallback={<SshKeyAddedAt>Added on ----</SshKeyAddedAt>}
  >
    {(createdAt) => <SshKeyAddedAt>Added on {dateHuman(createdAt())}</SshKeyAddedAt>}
  </Show>
)

const SshKeyRow: Component<{ key: UserKey }> = (props) => {
  const { Modal, open, close } = useModal()
  const handleDeleteKey = async (keyID: DeleteUserKeyRequest['keyId']) => {
    try {
      await client.deleteUserKey({ keyId: keyID })
      toast.success('公開鍵を削除しました')
      revalidateUserKeys()
    } catch (e) {
      handleAPIError(e, '公開鍵の削除に失敗しました')
    }
  }

  return (
    <>
      <div class="bg-ui-primary last:border-b-none flex w-full items-center gap-2 overflow-x-hidden px-5 py-4">
        <div class="flex w-full flex-col overflow-x-hidden">
          <SshKeyName>{props.key.name === '' ? '(Name not set)' : props.key.name}</SshKeyName>
          <SshKeyRowValue>{props.key.publicKey}</SshKeyRowValue>
          <AddedOn key={props.key} />
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
              <AddedOn key={props.key} />
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

const userKeyRequestSchema = v.object({
  name: v.pipe(v.string(), v.nonEmpty('Enter Name')),
  publicKey: v.pipe(v.string(), v.nonEmpty('Enter SSH Public Key')),
})

export default () => {
  const keys = createMemo(() => getUserKeys())
  const { Modal: AddNewKeyModal, open: newKeyOpen, close: newKeyClose } = useModal()

  const form = createForm({
    schema: userKeyRequestSchema,
    initialInput: { name: '', publicKey: '' },
  })

  const handleSubmit = async (values: v.InferOutput<typeof userKeyRequestSchema>) => {
    try {
      await client.createUserKey(values)
      toast.success('公開鍵を追加しました')
      newKeyClose()
      reset(form)
      revalidateUserKeys()
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
          <MainViewContainer>
            <SectionBoundary title="SSH Public Keys" bare>
              <DataTable.Container>
                <div class="flex w-full items-end justify-between">
                  <DataTable.Titles>
                    <DataTable.Title>SSH Public Keys</DataTable.Title>
                    <DataTable.SubTitle>
                      SSH鍵はruntimeアプリケーションのコンテナにssh接続するときに使います
                    </DataTable.SubTitle>
                  </DataTable.Titles>
                  <Show when={keys().length > 0}>
                    <AddNewSSHKeyButton />
                  </Show>
                </div>
                <Show
                  when={keys().length > 0}
                  fallback={
                    <List.Container>
                      <List.PlaceHolder>
                        <div class="i-material-symbols:key-off-outline text-20/20 shrink-0" />
                        No Keys Registered
                        <AddNewSSHKeyButton />
                      </List.PlaceHolder>
                    </List.Container>
                  }
                >
                  <List.Container>
                    <For each={keys()} keyed={(key) => key.id}>
                      {(key) => <SshKeyRow key={key()} />}
                    </For>
                  </List.Container>
                </Show>
              </DataTable.Container>
            </SectionBoundary>
          </MainViewContainer>
        </WithNav.Body>
      </WithNav.Container>
      <AddNewKeyModal.Container>
        <Form of={form} onSubmit={handleSubmit}>
          <AddNewKeyModal.Header>Add New SSH Key</AddNewKeyModal.Header>
          <AddNewKeyModal.Body>
            <div class="flex w-full flex-col gap-2">
              <Field of={form} path={['name']}>
                {(field) => (
                  <TextField label="Name" required {...field.props} value={field.input} error={field.errors?.[0]} />
                )}
              </Field>
              <Field of={form} path={['publicKey']}>
                {(field) => (
                  <TextField
                    label="Key"
                    required
                    multiline
                    placeholder="ssh-ed25519 AAA..."
                    {...field.props}
                    value={field.input}
                    error={field.errors?.[0]}
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
              disabled={!form.isValid || form.isSubmitting}
              loading={form.isSubmitting}
            >
              Add
            </Button>
          </AddNewKeyModal.Footer>
        </Form>
      </AddNewKeyModal.Container>
    </>
  )
}
