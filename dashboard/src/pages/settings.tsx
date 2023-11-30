import { styled } from '@macaron-css/solid'
import { SubmitHandler, createForm, required } from '@modular-forms/solid'
import { Component, For, Show, createResource } from 'solid-js'
import toast from 'solid-toast'
import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { DeleteUserKeyRequest, UserKey } from '../api/neoshowcase/protobuf/gateway_pb'
import { MaterialSymbols } from '../components/UI/MaterialSymbols'
import ModalDeleteConfirm from '../components/UI/ModalDeleteConfirm'
import { TextField } from '../components/UI/TextField'
import { DataTable } from '../components/layouts/DataTable'
import { MainViewContainer } from '../components/layouts/MainView'
import { WithNav } from '../components/layouts/WithNav'
import { List } from '../components/templates/List'
import { Nav } from '../components/templates/Nav'
import useModal from '../libs/useModal'

const TitleContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'flex-end',
    justifyContent: 'space-between',
  },
})
const SshKeyRowContainer = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    overflowX: 'hidden',
    background: colorVars.semantic.ui.primary,

    selectors: {
      '&:last-child': {
        borderBottom: 'none',
      },
    },
  },
})
const SshKeyRowContent = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    overflowX: 'hidden',
  },
})
const SshKeyRowValue = styled('h3', {
  base: {
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    color: colorVars.semantic.text.black,
    ...textVars.text.medium,
  },
})
const FormContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
})

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
      <SshKeyRowContainer>
        <SshKeyRowContent>
          <SshKeyRowValue>{props.key.publicKey}</SshKeyRowValue>
        </SshKeyRowContent>
        <Button variants="textError" size="medium" onClick={open}>
          Delete
        </Button>
      </SshKeyRowContainer>
      <Modal.Container>
        <Modal.Header>Delete SSH Key</Modal.Header>
        <Modal.Body>
          <ModalDeleteConfirm>{props.key.publicKey}</ModalDeleteConfirm>
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

type AddNewKeyForm = {
  value: string
}

export default () => {
  const [userKeys, { refetch: refetchKeys }] = createResource(() => client.getUserKeys({}))
  const { Modal: AddNewKeyModal, open: newKeyOpen, close: newKeyClose } = useModal()

  const [key, Key] = createForm<AddNewKeyForm>({
    initialValues: {
      value: '',
    },
  })

  const handleSubmit: SubmitHandler<AddNewKeyForm> = async (values) => {
    try {
      await client.createUserKey({ publicKey: values.value })
      toast.success('公開鍵を追加しました')
      newKeyClose()
      refetchKeys()
    } catch (e) {
      handleAPIError(e, '公開鍵の追加に失敗しました')
    }
  }

  const AddNewSSHKeyButton = () => (
    <Button variants="primary" size="medium" leftIcon={<MaterialSymbols>add</MaterialSymbols>} onClick={newKeyOpen}>
      Add New SSH Key
    </Button>
  )

  return (
    <>
      <WithNav.Container>
        <WithNav.Navs>
          <Nav title="Settings" />
        </WithNav.Navs>
        <WithNav.Body>
          <Show when={userKeys.state === 'ready'}>
            <MainViewContainer>
              <DataTable.Container>
                <TitleContainer>
                  <DataTable.Titles>
                    <DataTable.Title>SSH Public Keys</DataTable.Title>
                    <DataTable.SubTitle>
                      SSH鍵はruntimeアプリケーションのコンテナにssh接続するときに使います
                    </DataTable.SubTitle>
                  </DataTable.Titles>
                  <Show when={userKeys()!.keys.length !== 0}>
                    <AddNewSSHKeyButton />
                  </Show>
                </TitleContainer>
                <Show
                  when={userKeys()!.keys.length > 0}
                  fallback={
                    <List.Container>
                      <List.PlaceHolder>
                        <MaterialSymbols displaySize={80}>key_off</MaterialSymbols>
                        No Keys Registered
                        <AddNewSSHKeyButton />
                      </List.PlaceHolder>
                    </List.Container>
                  }
                >
                  <SshKeys keys={userKeys()!.keys!} refetchKeys={refetchKeys} />
                </Show>
              </DataTable.Container>
            </MainViewContainer>
          </Show>
        </WithNav.Body>
      </WithNav.Container>
      <AddNewKeyModal.Container>
        <Key.Form onSubmit={handleSubmit}>
          <AddNewKeyModal.Header>Add New SSH Key</AddNewKeyModal.Header>
          <AddNewKeyModal.Body>
            <FormContainer>
              <Key.Field name="value" validate={[required('Enter SSH Public Key')]}>
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
              </Key.Field>
            </FormContainer>
          </AddNewKeyModal.Body>
          <AddNewKeyModal.Footer>
            <Button variants="text" size="medium" type="button" onClick={newKeyClose}>
              Cancel
            </Button>
            <Button variants="primary" size="medium" type="submit" disabled={key.invalid || key.submitting}>
              Add
            </Button>
          </AddNewKeyModal.Footer>
        </Key.Form>
      </AddNewKeyModal.Container>
    </>
  )
}
