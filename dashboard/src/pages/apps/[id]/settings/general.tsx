import { useNavigate } from '@solidjs/router'
import { type Component, Show } from 'solid-js'
import toast from 'solid-toast'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import ModalDeleteConfirm from '/@/components/UI/ModalDeleteConfirm'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import GeneralConfigForm from '/@/features/application/components/form/GeneralConfigForm'
import { ApplicationFormProvider } from '/@/features/application/provider/applicationFormProvider'
import { client, handleAPIError } from '/@/libs/api'
import useModal from '/@/libs/useModal'
import { useApplicationData } from '/@/routes'

const DeleteApp: Component<{
  app: Application
  repo: Repository
  hasPermission: boolean
}> = (props) => {
  const { Modal, open, close } = useModal()
  const navigate = useNavigate()

  const deleteApplication = async () => {
    try {
      await client.deleteApplication({ id: props.app.id })
      toast.success('アプリケーションを削除しました')
      close()
      navigate('/apps')
    } catch (e) {
      handleAPIError(e, 'アプリケーションの削除に失敗しました')
    }
  }

  return (
    <>
      <FormBox.Container>
        <FormBox.Forms>
          <FormItem title="Delete Application">
            <div class="caption-regular text-text-grey">このアプリケーションを削除します。</div>
          </FormItem>
        </FormBox.Forms>
        <FormBox.Actions>
          <Button
            variants="primaryError"
            size="small"
            onClick={open}
            type="button"
            disabled={!props.hasPermission}
            tooltip={{
              props: {
                content: !props.hasPermission
                  ? 'アプリケーションを削除するにはオーナーになる必要があります'
                  : undefined,
              },
            }}
          >
            Delete Application
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Application</Modal.Header>
        <Modal.Body>
          <ModalDeleteConfirm>
            <MaterialSymbols>deployed_code</MaterialSymbols>
            {props.app.name}
          </ModalDeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button variants="text" size="medium" onClick={close} type="button">
            No, Cancel
          </Button>
          <Button variants="primaryError" size="medium" onClick={deleteApplication} type="button">
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </>
  )
}

export default () => {
  const { app, refetch, repo, hasPermission } = useApplicationData()
  const loaded = () => !!(app() && repo())

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <ApplicationFormProvider>
        <Show when={loaded()}>
          <GeneralConfigForm app={app()!} repo={repo()!} hasPermission={hasPermission()} refetchApp={refetch} />
          <DeleteApp app={app()!} repo={repo()!} hasPermission={hasPermission()} />
        </Show>
      </ApplicationFormProvider>
    </DataTable.Container>
  )
}
