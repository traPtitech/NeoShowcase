import { Header } from '/@/components/Header'
import { createMemo, createResource, For } from 'solid-js'
import { RadioItem } from '/@/components/Radio'
import { client } from '/@/libs/api'
import { GetRepositoriesRequest_Scope } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '@macaron-css/solid'
import { vars } from '/@/theme'
import { Container } from '/@/libs/layout'

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

const UserKeysContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',

    background: vars.bg.white1,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})


const PublicKey = styled('div', {
  base: {
    background: vars.bg.white3,
    border: `1px solid ${vars.bg.white4}`,
    borderRadius: '4px',
    padding: '8px 12px',
  },
})

export default () => {

  const [userKeys] = createResource(() => client.getUserKeys({}))

  const userKeysMemo = createMemo(() => {
    if (!userKeys()) return
    return userKeys().keys
  })

  return (
    <Container>
      <Header />
      <PageTitle>User Settings</PageTitle>
      <ContentContainer>
        <UserKeysContainer>
          <SidebarTitle>User Keys</SidebarTitle>
          <For each={userKeysMemo()}>
            {(key, i) => (
              <div>
                {i()}
                <PublicKey>
                  {key.publicKey}
                </PublicKey>
              </div>
            )}
          </For>
        </UserKeysContainer>
      </ContentContainer>
    </Container>
  )
}

