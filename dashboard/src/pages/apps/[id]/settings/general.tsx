import { useNavigate } from '@solidjs/router'
import type { Component } from 'solid-js'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { Button } from '/@/components/UI/Button'
import ModalDeleteConfirm from '/@/components/UI/ModalDeleteConfirm'
import GeneralConfigForm from '/@/features/application/components/form/GeneralConfigForm'
import { useApplicationData } from '/@/features/application/data'
import { client, handleAPIError } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import useModal from '/@/libs/useModal'

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
              content: !props.hasPermission ? 'アプリケーションを削除するにはオーナーになる必要があります' : undefined,
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
            <div class="i-material-symbols:deployed-code-outline shrink-0 text-2xl/6" />
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

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <GeneralConfigForm app={app()} repo={repo()} hasPermission={hasPermission()} refetch={refetch} />
      <DeleteApp app={app()} repo={repo()} hasPermission={hasPermission()} />
    </DataTable.Container>
  )
}
