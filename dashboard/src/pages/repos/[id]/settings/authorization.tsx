import {
  CreateRepositoryAuth,
  Repository,
  Repository_AuthMethod,
  UpdateRepositoryRequest,
} from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { RepositoryAuthSettings } from '/@/components/templates/RepositoryAuthSettings'
import { client, handleAPIError } from '/@/libs/api'
import { useRepositoryData } from '/@/routes'
import { PlainMessage } from '@bufbuild/protobuf'
import { Component, Show } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

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
        <RepositoryAuthSettings
          url={updateReq.url}
          setUrl={(v) => setUpdateReq('url', v)}
          authConfig={authConfig}
          setAuthConfig={setAuthConfig}
        />
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
    <DataTable.Container>
      <DataTable.Title>Authorization</DataTable.Title>
      <Show when={loaded()}>
        <AuthConfig repo={repo()} refetchRepo={refetchRepo} />
      </Show>
    </DataTable.Container>
  )
}
