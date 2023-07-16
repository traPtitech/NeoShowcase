import { Header } from '/@/components/Header'
import { Component, createResource, createSignal, For, JSX, Show } from 'solid-js'
import { client, handleAPIError } from '/@/libs/api'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { Container, PageTitle } from '/@/libs/layout'
import { Button } from '/@/components/Button'
import { InputBar, InputLabel } from '/@/components/Input'
import toast from 'solid-toast'
import { style } from '@macaron-css/core'
import { InfoTooltip } from '/@/components/InfoTooltip'

// copy from /pages/apps
// and delete unnecessary styles
const ContentContainer = styled('div', {
  base: {
    marginTop: '24px',
  },
})

const SidebarTitle = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',

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

const CreatingKeyContainerClass = style({
  display: 'flex',
  flexDirection: 'column',
  gap: '16px',

  background: vars.bg.white2,
  border: `1px solid ${vars.bg.white5}`,
  borderRadius: '4px',
  padding: '8px 12px',
})
styled('form', {
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
  const [userKeys, { refetch: refetchKeys }] = createResource(() => client.getUserKeys({}))

  const [createKeyToggle, setCreateKeyToggle] = createSignal(false)

  const deleteKeyRequest = async (keyID: string) => {
    try {
      await client.deleteUserKey({ keyId: keyID })
      toast.success('公開鍵を削除しました')
      refetchKeys()
    } catch (e) {
      return handleAPIError(e, '公開鍵の削除に失敗しました')
    }
  }

  const CreatingKeyContainer: Component = () => {
    let formRef: HTMLFormElement
    const [input, setInput] = createSignal('')

    const createKeyRequest: JSX.EventHandler<HTMLButtonElement, MouseEvent> = async (e) => {
      // prevent default form submit (reload page)
      e.preventDefault()

      // validate form
      if (!formRef.reportValidity()) {
        return
      }

      try {
        await client.createUserKey({ publicKey: input() })
        toast.success('公開鍵を登録しました')
        setInput('')
        refetchKeys()
      } catch (e) {
        return handleAPIError(e, '公開鍵の登録に失敗しました')
      }
    }

    return (
      <form class={CreatingKeyContainerClass} ref={formRef}>
        <InputLabel>SSH公開鍵の追加</InputLabel>
        <InputBar
          placeholder='ssh-ed25519 AAA...'
          type='text'
          value={input()}
          onInput={(e) => setInput(e.target.value)}
          required
        />
        <Button color='black1' size='large' width='auto' onclick={createKeyRequest} type='submit'>
          + SSH公開鍵の追加
        </Button>
      </form>
    )
  }

  return (
    <Container>
      <Header />
      <PageTitle>User Settings</PageTitle>
      <ContentContainer>
        <MainContentContainer>
          <UserKeysContainer>
            <SidebarTitle>
              登録済みSSH公開鍵
              <InfoTooltip tooltip='アプリへSSH接続時に使用します' />
            </SidebarTitle>
            <Show when={userKeys()} fallback={<SidebarTitle>登録済みSSH公開鍵を読み込み中...</SidebarTitle>}>
              <For each={userKeys()?.keys}>
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
                        onclick={() => {
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
            </Show>
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
              <CreatingKeyContainer />
            </Show>
          </CreateKeyContainer>
        </MainContentContainer>
      </ContentContainer>
    </Container>
  )
}
