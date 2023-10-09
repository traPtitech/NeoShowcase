import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
import { colorVars, textVars } from '/@/theme'
import { styled } from '@macaron-css/solid'
import { Component, For, Show, createMemo, createResource } from 'solid-js'
import toast from 'solid-toast'
import { DeleteUserKeyRequest, UserKey } from '../api/neoshowcase/protobuf/gateway_pb'
import { MaterialSymbols } from '../components/UI/MaterialSymbols'
import { Textarea } from '../components/UI/Textarea'
import { DataTable } from '../components/layouts/DataTable'
import { MainViewContainer } from '../components/layouts/MainView'
import { WithNav } from '../components/layouts/WithNav'
import { FormItem } from '../components/templates/FormItem'
import { List } from '../components/templates/List'
import { Nav } from '../components/templates/Nav'
import useModal from '../libs/useModal'

const PlaceHolder = styled('div', {
  base: {
    width: '100%',
    height: '400px',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
    alignItems: 'center',
    justifyContent: 'center',

    color: colorVars.semantic.text.black,
    ...textVars.h4.medium,
  },
})
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

    borderBottom: `1px solid ${colorVars.semantic.ui.border}`,
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
const DeleteConfirm = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'column',

    overflowWrap: 'break-word',
    borderRadius: '8px',
    background: colorVars.semantic.ui.secondary,
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
        <Button color="textError" size="medium" onClick={open}>
          Delete
        </Button>
      </SshKeyRowContainer>
      <Modal.Container>
        <Modal.Header>Delete SSH Key</Modal.Header>
        <Modal.Body>
          <DeleteConfirm>{props.key.publicKey}</DeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button color="text" size="medium" onClick={close}>
            No, Cancel
          </Button>
          <Button color="primaryError" size="medium" onClick={() => handleDeleteKey(props.key.id)}>
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

export default () => {
  const [userKeys, { refetch: refetchKeys }] = createResource(() => client.getUserKeys({}))
  const { Modal: AddNewKeyModal, open: newKeyOpen, close: newKeyClose } = useModal()

  let textareaRef: HTMLTextAreaElement
  const handleAddNewKey = () => {
    try {
      client.createUserKey({ publicKey: textareaRef.value })
      toast.success('公開鍵を追加しました')
      newKeyClose()
      refetchKeys()
    } catch (e) {
      handleAPIError(e, '公開鍵の追加に失敗しました')
    }
  }

  const AddNewSSHKeyButton = () => (
    <Button color="primary" size="medium" leftIcon={<MaterialSymbols>add</MaterialSymbols>} onClick={newKeyOpen}>
      Add New SSH Key
    </Button>
  )

  const showPlaceHolder = createMemo(() => userKeys()?.keys.length === 0)

  return (
    <>
      <WithNav.Container>
        <WithNav.Navs>
          <Nav title="Settings" />
        </WithNav.Navs>
        <WithNav.Body>
          <Show when={userKeys()}>
            <MainViewContainer>
              <DataTable.Container>
                <TitleContainer>
                  <DataTable.Titles>
                    <DataTable.Title>SSH Public Keys</DataTable.Title>
                    <DataTable.SubTitle>
                      SSH鍵はruntimeアプリケーションのコンテナにssh接続するときに使います
                    </DataTable.SubTitle>
                  </DataTable.Titles>
                  <Show when={!showPlaceHolder()}>
                    <AddNewSSHKeyButton />
                  </Show>
                </TitleContainer>
                <Show when={showPlaceHolder()} fallback={<SshKeys keys={userKeys()?.keys} refetchKeys={refetchKeys} />}>
                  <List.Container>
                    <PlaceHolder>
                      <MaterialSymbols displaySize={80}>key_off</MaterialSymbols>
                      No Keys Registered
                      <AddNewSSHKeyButton />
                    </PlaceHolder>
                  </List.Container>
                </Show>
              </DataTable.Container>
            </MainViewContainer>
          </Show>
        </WithNav.Body>
      </WithNav.Container>
      <AddNewKeyModal.Container>
        <AddNewKeyModal.Header>Add New SSH Key</AddNewKeyModal.Header>
        <AddNewKeyModal.Body>
          <FormItem title="Key" required>
            <Textarea placeholder="ssh-ed25519 AAA..." ref={textareaRef} />
          </FormItem>
        </AddNewKeyModal.Body>
        <AddNewKeyModal.Footer>
          <Button color="text" size="medium" onClick={newKeyClose}>
            Cancel
          </Button>
          <Button color="primary" size="medium" onClick={() => handleAddNewKey()}>
            Add
          </Button>
        </AddNewKeyModal.Footer>
      </AddNewKeyModal.Container>
    </>
  )
}
