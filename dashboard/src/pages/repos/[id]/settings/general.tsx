import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { SubmitHandler, createForm, required, reset } from '@modular-forms/solid'
import { useNavigate } from '@solidjs/router'
import { Component, Show, createEffect } from 'solid-js'
import toast from 'solid-toast'
import { Application, Repository, UpdateRepositoryRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import ModalDeleteConfirm from '/@/components/UI/ModalDeleteConfirm'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { client, handleAPIError } from '/@/libs/api'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import useModal from '/@/libs/useModal'
import { useRepositoryData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { TextField } from '../../../../components/UI/TextField'

type GeneralForm = Required<Pick<PlainMessage<UpdateRepositoryRequest>, 'name'>>

const NameConfig: Component<{
  repo: Repository
  refetchRepo: () => void
  hasPermission: boolean
}> = (props) => {
  const [generalForm, General] = createForm<GeneralForm>({
    initialValues: {
      name: props.repo.name,
    },
  })

  createEffect(() => {
    reset(generalForm, 'name', {
      initialValue: props.repo.name,
    })
  })

  const handleSubmit: SubmitHandler<GeneralForm> = async (values) => {
    try {
      await client.updateRepository({
        id: props.repo.id,
        name: values.name,
      })
      toast.success('Project名を更新しました')
      props.refetchRepo()
    } catch (e) {
      handleAPIError(e, 'Project名の更新に失敗しました')
    }
  }
  const discardChanges = () => {
    reset(generalForm)
  }

  return (
    <General.Form onSubmit={handleSubmit}>
      <FormBox.Container>
        <FormBox.Forms>
          <General.Field name="name" validate={[required('Enter Project Name')]}>
            {(field, fieldProps) => (
              <TextField
                label="Project Name"
                required
                {...fieldProps}
                value={field.value ?? ''}
                error={field.error}
                readOnly={!props.hasPermission}
              />
            )}
          </General.Field>
        </FormBox.Forms>
        <FormBox.Actions>
          <Show when={generalForm.dirty && !generalForm.submitting}>
            <Button variants="borderError" size="small" onClick={discardChanges} type="button">
              Discard Changes
            </Button>
          </Show>
          <Button
            variants="primary"
            size="small"
            type="submit"
            disabled={generalForm.invalid || !generalForm.dirty || generalForm.submitting || !props.hasPermission}
            loading={generalForm.submitting}
            tooltip={{
              props: {
                content: !props.hasPermission
                  ? '設定を変更するにはリポジトリのオーナーになる必要があります'
                  : undefined,
              },
            }}
          >
            Save
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
    </General.Form>
  )
}

const DeleteProjectNotice = styled('div', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})

const DeleteProject: Component<{
  repo: Repository
  apps: Application[]
  hasPermission: boolean
}> = (props) => {
  const { Modal, open, close } = useModal()
  const navigate = useNavigate()

  const deleteRepository = async () => {
    try {
      await client.deleteRepository({ repositoryId: props.repo.id })
      toast.success('Projectを削除しました')
      close()
      navigate('/apps')
    } catch (e) {
      handleAPIError(e, 'Projectの削除に失敗しました')
    }
  }
  const canDeleteRepository = () => props.apps.length === 0

  return (
    <>
      <FormBox.Container>
        <FormBox.Forms>
          <FormItem title="Delete Project">
            <DeleteProjectNotice>
              Projectを削除するには、このプロジェクト内のすべてのAppを削除する必要があります。
            </DeleteProjectNotice>
          </FormItem>
        </FormBox.Forms>
        <FormBox.Actions>
          <Button
            variants="primaryError"
            size="small"
            onClick={open}
            type="button"
            disabled={!canDeleteRepository() || !props.hasPermission}
            tooltip={{
              props: {
                content: !props.hasPermission
                  ? 'Projectを削除するにはリポジトリのオーナーになる必要があります'
                  : !canDeleteRepository()
                  ? 'Project内にAppが存在するため削除できません'
                  : undefined,
              },
            }}
          >
            Delete Project
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Repository</Modal.Header>
        <Modal.Body>
          <ModalDeleteConfirm>
            {providerToIcon(repositoryURLToProvider(props.repo.url), 24)}
            {props.repo.name}
          </ModalDeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button variants="text" size="medium" onClick={close} type="button">
            No, Cancel
          </Button>
          <Button variants="primaryError" size="medium" onClick={deleteRepository} type="button">
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </>
  )
}

export default () => {
  const { repo, refetchRepo, apps, hasPermission } = useRepositoryData()
  const loaded = () => !!(repo() && apps())

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <Show when={loaded()}>
        <NameConfig repo={repo()!} refetchRepo={refetchRepo} hasPermission={hasPermission()} />
        <DeleteProject repo={repo()!} apps={apps()!} hasPermission={hasPermission()} />
      </Show>
    </DataTable.Container>
  )
}
