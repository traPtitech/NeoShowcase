import { CreateRepositoryAuth, Repository, Repository_AuthMethod } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { AuthForm, RepositoryAuthSettings } from '/@/components/templates/RepositoryAuthSettings'
import { client, handleAPIError } from '/@/libs/api'
import { useRepositoryData } from '/@/routes'
import { PlainMessage } from '@bufbuild/protobuf'
import { SubmitHandler, createForm, getValue, reset } from '@modular-forms/solid'
import { Component, Show } from 'solid-js'
import toast from 'solid-toast'

const mapAuthMethod = (authMethod: Repository_AuthMethod): PlainMessage<CreateRepositoryAuth>['auth']['case'] => {
  switch (authMethod) {
    case Repository_AuthMethod.NONE:
      return 'none'
    case Repository_AuthMethod.BASIC:
      return 'basic'
    case Repository_AuthMethod.SSH:
      return 'ssh'
  }
}

const AuthConfig: Component<{
  repo: Repository
  refetchRepo: () => void
}> = (props) => {
  const [authForm, Auth] = createForm<AuthForm>({
    initialValues: {
      url: props.repo.url,
      case: mapAuthMethod(props.repo.authMethod),
      auth: {
        basic: {
          username: '',
          password: '',
        },
        ssh: {
          keyId: '',
        },
      },
    },
  })
  const authValue = (): PlainMessage<CreateRepositoryAuth>['auth'] => {
    const authMethod = getValue(authForm, 'case')
    switch (authMethod) {
      case 'none':
        return {
          case: 'none',
          value: '',
        }
      case 'basic':
        return {
          case: 'basic',
          value: {
            username: getValue(authForm, 'auth.basic.username'),
            password: getValue(authForm, 'auth.basic.password'),
          },
        }
      case 'ssh':
        return {
          case: 'ssh',
          value: {
            keyId: getValue(authForm, 'auth.ssh.keyId'),
          },
        }
    }
    throw new Error('unreachable')
  }

  const handleSubmit: SubmitHandler<AuthForm> = async (values) => {
    try {
      await client.updateRepository({
        id: props.repo.id,
        url: values.url,
        auth: {
          auth: authValue(),
        },
      })
      toast.success('リポジトリの設定を更新しました')
      props.refetchRepo()
    } catch (e) {
      handleAPIError(e, 'リポジトリの設定の更新に失敗しました')
    }
  }

  const discardChanges = () => {
    reset(authForm)
  }

  const AuthSetting = RepositoryAuthSettings({
    formStore: authForm,
    Form: Auth,
  })

  return (
    <Auth.Form onSubmit={handleSubmit}>
      <FormBox.Container>
        <FormBox.Forms>
          <AuthSetting.Url />
          <AuthSetting.AuthMethod />
          <AuthSetting.AuthConfig />
        </FormBox.Forms>
        <FormBox.Actions>
          <Show when={authForm.dirty && !authForm.submitting}>
            <Button color="borderError" size="small" onClick={discardChanges} type="button">
              Discard Changes
            </Button>
          </Show>
          <Button
            color="primary"
            size="small"
            type="submit"
            disabled={authForm.invalid || !authForm.dirty || authForm.submitting}
          >
            Save
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
    </Auth.Form>
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
