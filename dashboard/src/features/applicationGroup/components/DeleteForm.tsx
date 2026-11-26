import { useNavigate } from '@solidjs/router'
import type { Component } from 'solid-js'
import type { Application, ApplicationGroup } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { Button } from '/@/components/UI/Button'
import ModalDeleteConfirm from '/@/components/UI/ModalDeleteConfirm'
import { client, handleAPIError, revalidateMyApplicationGroups } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import useModal from '/@/libs/useModal'

type Props = {
  group: ApplicationGroup
  apps: Application[]
  hasPermission: boolean
}

const DeleteForm: Component<Props> = (props) => {
  const { Modal, open, close } = useModal()
  const navigate = useNavigate()

  const canDeleteGroup = () => props.apps.length === 0

  const deleteGroup = async () => {
    try {
      await client.deleteApplicationGroup({ groupId: props.group.id })
      toast.success('グループを削除しました')
      close()
      revalidateMyApplicationGroups()
      navigate('/groups')
    } catch (e) {
      handleAPIError(e, 'グループの削除に失敗しました')
    }
  }

  return (
    <>
      <FormBox.Container>
        <FormBox.Forms>
          <FormItem title="Delete Group">
            <div class="caption-regular text-text-grey">
              グループを削除するには、グループ内のすべてのアプリケーションを別のグループへ移すか削除する必要があります。
            </div>
          </FormItem>
        </FormBox.Forms>
        <FormBox.Actions>
          <Button
            variants="primaryError"
            size="small"
            onClick={open}
            type="button"
            disabled={!canDeleteGroup() || !props.hasPermission}
            tooltip={{
              content: !props.hasPermission
                ? 'グループを削除するにはオーナーになる必要があります'
                : !canDeleteGroup()
                  ? 'グループ内にアプリケーションが存在するため削除できません'
                  : undefined,
            }}
          >
            Delete Group
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Group</Modal.Header>
        <Modal.Body>
          <ModalDeleteConfirm>
            <div class="i-material-symbols:folder-outline shrink-0 text-2xl/6" />
            {props.group.name}
          </ModalDeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button variants="text" size="medium" onClick={close} type="button">
            No, Cancel
          </Button>
          <Button variants="primaryError" size="medium" onClick={deleteGroup} type="button">
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </>
  )
}

export default DeleteForm
