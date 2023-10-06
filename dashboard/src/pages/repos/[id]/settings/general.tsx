import { Application, Repository, UpdateRepositoryRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextInput } from '/@/components/UI/TextInput'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { client, handleAPIError } from '/@/libs/api'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import useModal from '/@/libs/useModal'
import { useRepositoryData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const NameConfig: Component<{
  repo: Repository
  refetchRepo: () => void
}> = (props) => {
  let formRef: HTMLFormElement

  const [updateReq, setUpdateReq] = createStore<PlainMessage<UpdateRepositoryRequest>>({
    id: props.repo.id,
    name: props.repo.name,
  })
  const discardChanges = () => {
    setUpdateReq({
      name: props.repo.name,
    })
  }
  const nameChanged = () => props.repo.name !== updateReq.name
  const saveChanges = async () => {
    try {
      // validate form
      if (!formRef.reportValidity()) {
        return
      }
      await client.updateRepository(updateReq)
      toast.success('Project名を更新しました')
      props.refetchRepo()
    } catch (e) {
      handleAPIError(e, 'Project名の更新に失敗しました')
    }
  }

  return (
    <FormBox.Container ref={formRef}>
      <FormBox.Forms>
        <FormItem title="Project Name" required>
          <TextInput
            required
            value={updateReq.name}
            onInput={(e) => {
              setUpdateReq('name', e.target.value)
            }}
          />
        </FormItem>
      </FormBox.Forms>
      <FormBox.Actions>
        <Show when={nameChanged()}>
          <Button color="borderError" size="small" onClick={discardChanges} type="button">
            Discard Changes
          </Button>
        </Show>
        <Button color="primary" size="small" onClick={saveChanges} type="button" disabled={!nameChanged()}>
          Save
        </Button>
      </FormBox.Actions>
    </FormBox.Container>
  )
}

const DeleteProjectNotice = styled('div', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
const DeleteConfirm = styled('div', {
  base: {
    width: '100%',
    padding: '16px 20px',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '8px',
    borderRadius: '8px',
    background: colorVars.semantic.ui.secondary,
    color: colorVars.semantic.text.black,
    ...textVars.h3.regular,
  },
})

const DeleteProject: Component<{
  repo: Repository
  apps: Application[]
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
          <Button color="primaryError" size="small" onClick={open} type="button" disabled={!canDeleteRepository()}>
            Delete Project
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Repository</Modal.Header>
        <Modal.Body>
          <DeleteConfirm>
            {providerToIcon(repositoryURLToProvider(props.repo.url), 24)}
            {props.repo.name}
          </DeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button color="text" size="medium" onClick={close} type="button">
            No, Cancel
          </Button>
          <Button color="primaryError" size="medium" onClick={deleteRepository} type="button">
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </>
  )
}

export default () => {
  const { repo, refetchRepo, apps } = useRepositoryData()
  const loaded = () => !!(repo() && apps())

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <Show when={loaded()}>
        <NameConfig repo={repo()} refetchRepo={refetchRepo} />
        <DeleteProject repo={repo()} apps={apps()} />
      </Show>
    </DataTable.Container>
  )
}
