import { Header } from '/@/components/Header'
import { createResource, createSignal, For, JSX, Show } from 'solid-js'
import { client } from '/@/libs/api'
import { CreateUserKeyRequest, DeleteUserKeyRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { Container } from '/@/libs/layout'
import { Button } from '/@/components/Button'
import { createStore } from 'solid-js/store'
import { InputBar, InputLabel } from '/@/components/Input'
import toast from 'solid-toast'
import { ConnectError } from '@bufbuild/connect'
import { useParams } from '@solidjs/router'

// copy from /pages/apps AppsTitle component
const PageTitle = styled('div', {
  base: {
    marginTop: '48px',
    fontSize: '32px',
    fontWeight: 'bold',
    color: vars.text.black1,
  },
})

// copy from /pages/apps
// and delete unnecessary styles
const ContentContainer = styled('div', {
  base: {
    marginTop: '24px',
  },
})

const SidebarTitle = styled('div', {
  base: {
    fontSize: '24px',
    fontWeight: 500,
    color: vars.text.black1,
  },
})

const MainContentContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
})

const UserKeysContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',

    background: vars.bg.white3,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

const PublicKey = styled('div', {
  base: {
    background: vars.bg.white1,
    border: `1px solid ${vars.bg.white5}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

const PublicKeyContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '6px',

    background: vars.bg.white2,
    border: `1px solid ${vars.bg.white5}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

const CreateKeyContainer = styled('div', {
  base: {
    background: vars.bg.white3,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

const CreatingKeyContainer = styled('form', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',

    background: vars.bg.white2,
    border: `1px solid ${vars.bg.white5}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

export const FormButton = styled('div', {
  base: {
    marginLeft: '4px',
  },
})

export default () => {
  const params = useParams()
  const [userKeys, { refetch: refetchApp }] = createResource(() => client.getUserKeys({}))

  const [createKeyToggle, setCreateKeyToggle] = createSignal(false)
  const [createKey, setCreateKey] = createStore(new CreateUserKeyRequest())

  const [deleteKeyID, setDeleteKeyID] = createStore(new DeleteUserKeyRequest())

  let formContainer: HTMLFormElement
  const createKeyRequest: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
    // prevent default form submit (reload page)
    e.preventDefault()

    // validate form
    if (!formContainer.reportValidity()) {
      return
    }

    try {
      const res = await client.createUserKey(createKey)
      toast.success('User Key を登録しました')
      refetchApp()
    } catch (e) {
      console.error(e)
      // gRPCエラー
      if (e instanceof ConnectError) {
        toast.error('User Key の登録に失敗しました\n' + e.message)
      }
    }
  }

  const deleteKeyRequest = (keyID: string) => {
    try {
      const a = new DeleteUserKeyRequest()
      a.keyId = keyID
      client.deleteUserKey(a)
      toast.success('User Key を削除しました')
      refetchApp()
    } catch (e) {
      console.error(e)
      // gRPCエラー
      if (e instanceof ConnectError) {
        toast.error('User Key の削除に失敗しました\n' + e.message)
      }
    }
  }

  return (
    <Container>
      <Header />
      <PageTitle>User Settings</PageTitle>
      <ContentContainer>
        <MainContentContainer>
          <UserKeysContainer>
            <SidebarTitle>登録済みSSH公開鍵</SidebarTitle>
            <For each={!userKeys() || userKeys().keys}>
              {(key) => (
                <PublicKeyContainer>
                  <div>
                    <PublicKey>{key.publicKey}</PublicKey>
                  </div>
                  <FormButton>
                    <Button
                      color='black1'
                      size='large'
                      width='auto'
                      onclick={(e) => {
                        setDeleteKeyID('keyId', key.id)
                        deleteKeyRequest(key.id)
                      }}
                      type='submit'
                    >
                      削除
                    </Button>
                  </FormButton>
                </PublicKeyContainer>
              )}
            </For>
          </UserKeysContainer>

          <CreateKeyContainer>
            <Show
              when={createKeyToggle()}
              fallback={
                <Button color='black1' size='large' width='auto' onClick={() => setCreateKeyToggle(true)}>
                  SSH公開鍵の追加
                </Button>
              }
            >
              <CreatingKeyContainer ref={formContainer}>
                <InputLabel>SSH公開鍵の追加</InputLabel>
                <InputBar
                  placeholder='my-app'
                  value={createKey.publicKey}
                  onInput={(e) => setCreateKey('publicKey', e.target.value)}
                  required
                />
                <Button color='black1' size='large' width='auto' onclick={createKeyRequest} type='submit'>
                  + SSH公開鍵の追加
                </Button>
              </CreatingKeyContainer>
            </Show>
          </CreateKeyContainer>
        </MainContentContainer>
      </ContentContainer>
    </Container>
  )
}