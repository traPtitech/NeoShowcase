import {
  CreateRepositoryAuth,
  Repository,
  Repository_AuthMethod,
  UpdateRepositoryRequest,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextInput } from '/@/components/UI/TextInput'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { RepositoryAuthSettings } from '/@/components/templates/RepositoryAuthSettings'
import { useRepositoryData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, Show } from 'solid-js'
import { createStore } from 'solid-js/store'

const Container = styled('div', {
  base: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '16px',
    overflowY: 'auto',

    color: colorVars.semantic.text.black,
    ...textVars.h2.medium,
  },
})
const ItemsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
  },
})

const AuthConfig: Component<{
  repo: Repository
}> = (props) => {
  const [updateReq, setUpdateReq] = createStore<PlainMessage<UpdateRepositoryRequest>>({
    id: props.repo.id,
    url: props.repo.url,
  })
  const mapAuthMethod = (authMethod: Repository_AuthMethod): PlainMessage<CreateRepositoryAuth>['auth'] => {
    switch (authMethod) {
      case Repository_AuthMethod.NONE:
        return { case: 'none', value: {} }
      case Repository_AuthMethod.BASIC:
        return { case: 'basic', value: { username: '', password: '' } }
      case Repository_AuthMethod.SSH:
        return { case: 'ssh', value: { keyId: '' } }
    }
  }
  const [authConfig, setAuthConfig] = createStore<PlainMessage<CreateRepositoryAuth>>({
    auth: mapAuthMethod(props.repo.authMethod),
  })

  return (
    <FormBox.Container>
      <FormBox.Forms>
        <ItemsContainer>
          <FormItem title="Repository URL" required>
            <TextInput
              value={updateReq.url}
              onInput={(e) => setUpdateReq('url', (e.target as HTMLInputElement).value)}
            />
          </FormItem>
          <RepositoryAuthSettings authConfig={authConfig} setAuthConfig={setAuthConfig} />
        </ItemsContainer>
      </FormBox.Forms>
      <FormBox.Actions>
        <Button color="borderError" size="small">
          Discard Changes
        </Button>
        <Button color="primary" size="small">
          Save
        </Button>
      </FormBox.Actions>
    </FormBox.Container>
  )
}

export default () => {
  const { repo } = useRepositoryData()
  const loaded = () => !!repo()

  return (
    <div>
      <Show when={loaded()}>
        <AuthConfig repo={repo()} />
      </Show>
    </div>
  )
}
