import { useNavigate } from '@solidjs/router'
import type { Component } from 'solid-js'
import toast from 'solid-toast'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import ModalDeleteConfirm from '/@/components/UI/ModalDeleteConfirm'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { client, handleAPIError } from '/@/libs/api'
import { originToIcon, repositoryURLToOrigin } from '/@/libs/application'
import useModal from '/@/libs/useModal'

type Props = {
  repo: Repository
  apps: Application[]
  hasPermission: boolean
}

const DeleteForm: Component<Props> = (props) => {
  const { Modal, open, close } = useModal()
  const navigate = useNavigate()

  const canDeleteRepository = () => props.apps.length === 0

  const deleteRepository = async () => {
    try {
      await client.deleteRepository({ repositoryId: props.repo.id })
      toast.success('リポジトリを削除しました')
      close()
      navigate('/apps')
    } catch (e) {
      handleAPIError(e, 'リポジトリの削除に失敗しました')
    }
  }

  return (
    <>
      <FormBox.Container>
        <FormBox.Forms>
          <FormItem title="Delete Repository">
            <div class="caption-regular text-text-grey">
              リポジトリを削除するには、このリポジトリ内のすべてのアプリケーションを削除する必要があります。
            </div>
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
                  ? 'リポジトリを削除するにはオーナーになる必要があります'
                  : !canDeleteRepository()
                    ? 'リポジトリ内にアプリケーションが存在するため削除できません'
                    : undefined,
              },
            }}
          >
            Delete Repository
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Repository</Modal.Header>
        <Modal.Body>
          <ModalDeleteConfirm>
            {originToIcon(repositoryURLToOrigin(props.repo.url), 24)}
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

export default DeleteForm
