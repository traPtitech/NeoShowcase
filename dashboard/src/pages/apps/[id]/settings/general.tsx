import { Application, Repository, UpdateApplicationRequest } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextInput } from '/@/components/UI/TextInput'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { ComboBox } from '/@/components/templates/Select'
import { client, handleAPIError } from '/@/libs/api'
import { providerToIcon, repositoryURLToProvider } from '/@/libs/application'
import { useBranchesSuggestion } from '/@/libs/branchesSuggestion'
import useModal from '/@/libs/useModal'
import { useApplicationData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { useNavigate } from '@solidjs/router'
import { Component, Show, createEffect } from 'solid-js'
import { createStore } from 'solid-js/store'
import toast from 'solid-toast'

const GeneralConfig: Component<{
  app: Application
  refetchApp: () => void
  repo: Repository
}> = (props) => {
  let formRef: HTMLFormElement

  const [updateReq, setUpdateReq] = createStore<PlainMessage<UpdateApplicationRequest>>({
    id: props.app.id,
    name: props.app.name,
    repositoryId: props.app.repositoryId,
    refName: props.app.refName,
  })
  const discardChanges = () => {
    setUpdateReq({
      name: props.app.name,
      repositoryId: props.app.repositoryId,
      refName: props.app.refName,
    })
  }
  const configChanged = () =>
    props.app.name !== updateReq.name ||
    props.app.repositoryId !== updateReq.repositoryId ||
    props.app.refName !== updateReq.refName
  const branchesSuggestion = useBranchesSuggestion(
    () => props.repo.id,
    () => updateReq.refName,
  )

  const saveChanges = async () => {
    try {
      // validate form
      if (!formRef.reportValidity()) {
        return
      }
      await client.updateApplication(updateReq)
      toast.success('アプリケーション設定を更新しました')
      props.refetchApp()
    } catch (e) {
      handleAPIError(e, 'アプリケーション設定の更新に失敗しました')
    }
  }

  return (
    <>
      <FormBox.Container ref={formRef}>
        <FormBox.Forms>
          <FormItem title="Application Name" required>
            <TextInput
              required
              value={updateReq.name}
              onInput={(e) => {
                setUpdateReq('name', e.target.value)
              }}
            />
          </FormItem>
          <FormItem title="Repository ID" required>
            <TextInput
              required
              value={updateReq.repositoryId}
              onInput={(e) => {
                setUpdateReq('repositoryId', e.target.value)
              }}
            />
          </FormItem>
          <FormItem title="Branch" required>
            <ComboBox
              required
              value={updateReq.refName}
              onInput={(e) => setUpdateReq('refName', e.target.value)}
              items={branchesSuggestion().map((branch) => ({
                title: branch,
                value: branch,
              }))}
              setSelected={(branch) => {
                setUpdateReq('refName', branch)
              }}
            />
          </FormItem>
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
    </>
  )
}

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
  const loaded = () => !!app()
  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <Show when={loaded()}>
        <GeneralConfig app={app()} refetchApp={refetchApp} repo={repo()} />
        <DeleteApp app={app()} repo={repo()} />
      </Show>
    </DataTable.Container>
  )
}
