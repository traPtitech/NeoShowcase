import {
  CreateRepositoryAuth,
  Repository,
  Repository_AuthMethod,
  UpdateRepositoryRequest,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import FormBox from '/@/components/layouts/FormBox'
import { RepositoryAuthSettings } from '/@/components/templates/RepositoryAuthSettings'
import { client, handleAPIError } from '/@/libs/api'
import { useRepositoryData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Component, Show, onMount } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

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
  refetchRepo: () => void
}> = (props) => {
  let formRef: HTMLFormElement

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

  const discardChanges = () => {
    setUpdateReq({
      id: props.repo.id,
      url: props.repo.url,
    })
    setAuthConfig({ auth: mapAuthMethod(props.repo.authMethod) })
  }
  const saveChanges = async () => {
    try {
      // validate form
      if (!formRef.reportValidity()) {
        return
      }
      await client.updateRepository({ ...updateReq, auth: authConfig })
      toast.success('リポジトリの設定を更新しました')
      props.refetchRepo()
    } catch (e) {
      handleAPIError(e, 'リポジトリの設定の更新に失敗しました')
    }
  }

  return (
    <FormBox.Container ref={formRef}>
      <FormBox.Forms>
        <ItemsContainer>
          <RepositoryAuthSettings
            url={updateReq.url}
            setUrl={(v) => setUpdateReq('url', v)}
            authConfig={authConfig}
            setAuthConfig={setAuthConfig}
          />
        </ItemsContainer>
      </FormBox.Forms>
      <FormBox.Actions>
        <Button color="borderError" size="small" onClick={discardChanges} type="button">
          Discard Changes
        </Button>
        <Button color="primary" size="small" onClick={saveChanges} type="button">
          Save
        </Button>
      </FormBox.Actions>
    </FormBox.Container>
  )
}

export default () => {
  const { repo, refetchRepo } = useRepositoryData()
  const loaded = () => !!repo()

  return (
    <Container>
      Authorization
      <Show when={loaded()}>
        <AuthConfig repo={repo()} refetchRepo={refetchRepo} />
      </Show>
    </Container>
  )
}
