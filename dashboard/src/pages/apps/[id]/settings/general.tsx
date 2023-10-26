import { Application, Repository, UpdateApplicationRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { GeneralConfig } from '/@/components/templates/GeneralConfig'
import { client, handleAPIError } from '/@/libs/api'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import useModal from '/@/libs/useModal'
import { useApplicationData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import { Component, Show } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const DeleteAppNotice = styled('div', {
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
const DeleteApp: Component<{
  app: Application
  repo: Repository
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
            <DeleteAppNotice>このアプリケーションを削除します。</DeleteAppNotice>
          </FormItem>
        </FormBox.Forms>
        <FormBox.Actions>
          <Button color="primaryError" size="small" onClick={open} type="button">
            Delete Application
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
      <Modal.Container>
        <Modal.Header>Delete Application</Modal.Header>
        <Modal.Body>
          <DeleteConfirm>
            {providerToIcon(repositoryURLToProvider(props.repo.url), 24)}
            {`${props.repo.name}/${props.app.name}`}
          </DeleteConfirm>
        </Modal.Body>
        <Modal.Footer>
          <Button color="text" size="medium" onClick={close} type="button">
            No, Cancel
          </Button>
          <Button color="primaryError" size="medium" onClick={deleteApplication} type="button">
            Yes, Delete
          </Button>
        </Modal.Footer>
      </Modal.Container>
    </>
  )
}

export default () => {
  const { app, refetchApp, repo } = useApplicationData()
  const loaded = () => !!(app() && repo())

  let formRef: HTMLFormElement
  const [updateReq, setUpdateReq] = createStore<PlainMessage<UpdateApplicationRequest>>({
    id: app()?.id,
    name: app()?.name,
    repositoryId: app()?.repositoryId,
    refName: app()?.refName,
  })
  const discardChanges = () => {
    setUpdateReq({
      name: app()?.name,
      repositoryId: app()?.repositoryId,
      refName: app()?.refName,
    })
  }
  const configChanged = () =>
    app()?.name !== updateReq.name ||
    app()?.repositoryId !== updateReq.repositoryId ||
    app()?.refName !== updateReq.refName

  const saveChanges = async () => {
    try {
      // validate form
      if (!formRef.reportValidity()) {
        return
      }
      await client.updateApplication(updateReq)
      toast.success('アプリケーション設定を更新しました')
      refetchApp()
    } catch (e) {
      handleAPIError(e, 'アプリケーション設定の更新に失敗しました')
    }
  }

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <Show when={loaded()}>
        <FormBox.Container ref={formRef}>
          <FormBox.Forms>
            <GeneralConfig repo={repo()} config={updateReq} setConfig={setUpdateReq} editBranchId />
          </FormBox.Forms>
          <FormBox.Actions>
            <Show when={configChanged()}>
              <Button color="borderError" size="small" onClick={discardChanges} type="button">
                Discard Changes
              </Button>
            </Show>
            <Button color="primary" size="small" onClick={saveChanges} type="button" disabled={!configChanged()}>
              Save
            </Button>
          </FormBox.Actions>
        </FormBox.Container>
        <DeleteApp app={app()} repo={repo()} />
      </Show>
    </DataTable.Container>
  )
}
