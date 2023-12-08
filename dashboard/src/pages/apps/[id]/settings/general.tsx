import { styled } from '@macaron-css/solid'
import { SubmitHandler, createForm, reset } from '@modular-forms/solid'
import { useNavigate } from '@solidjs/router'
import { Component, Show, createEffect, on } from 'solid-js'
import toast from 'solid-toast'
import { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import ModalDeleteConfirm from '/@/components/UI/ModalDeleteConfirm'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { FormItem } from '/@/components/templates/FormItem'
import { AppGeneralConfig, AppGeneralForm } from '/@/components/templates/app/AppGeneralConfig'
import { client, handleAPIError } from '/@/libs/api'
import useModal from '/@/libs/useModal'
import { useApplicationData } from '/@/routes'
import { colorVars, textVars } from '/@/theme'

const DeleteAppNotice = styled('div', {
  base: {
    color: colorVars.semantic.text.grey,
    ...textVars.caption.regular,
  },
})
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
            <DeleteAppNotice>このアプリケーションを削除します。</DeleteAppNotice>
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
  const { app, refetchApp, repo, hasPermission } = useApplicationData()
  const loaded = () => !!(app() && repo())

  const [generalForm, General] = createForm<AppGeneralForm>({
    initialValues: {
      name: app()?.name,
      repositoryId: app()?.repositoryId,
      refName: app()?.refName,
    },
  })

  createEffect(
    on(app, (app) => {
      reset(generalForm, {
        initialValues: {
          name: app?.name,
          repositoryId: app?.repositoryId,
          refName: app?.refName,
        },
      })
    }),
  )

  const handleSubmit: SubmitHandler<AppGeneralForm> = async (values) => {
    try {
      await client.updateApplication({
        id: app()?.id,
        ...values,
      })
      toast.success('アプリケーション設定を更新しました')
      refetchApp()
    } catch (e) {
      handleAPIError(e, 'アプリケーション設定の更新に失敗しました')
    }
  }
  const discardChanges = () => {
    reset(generalForm)
  }

  return (
    <DataTable.Container>
      <DataTable.Title>General</DataTable.Title>
      <Show when={loaded()}>
        <General.Form onSubmit={handleSubmit}>
          <FormBox.Container>
            <FormBox.Forms>
              <AppGeneralConfig repo={repo()!} formStore={generalForm} editBranchId hasPermission={hasPermission()} />
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
                disabled={generalForm.invalid || !generalForm.dirty || generalForm.submitting || !hasPermission()}
                loading={generalForm.submitting}
                tooltip={{
                  props: {
                    content: !hasPermission()
                      ? '設定を変更するにはアプリケーションのオーナーになる必要があります'
                      : undefined,
                  },
                }}
              >
                Save
              </Button>
            </FormBox.Actions>
          </FormBox.Container>
        </General.Form>
        <DeleteApp app={app()!} repo={repo()!} hasPermission={hasPermission()} />
      </Show>
    </DataTable.Container>
  )
}
