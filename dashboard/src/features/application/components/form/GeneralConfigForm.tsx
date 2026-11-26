import { createForm, Form } from '@formisch/solid/v2'
import { type Component, untrack } from 'solid-js'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
import { syncFormWithServer } from '/@/libs/serverForm'
import { toast } from '/@/libs/toast'
import {
  createOrUpdateApplicationSchema,
  getInitialValueOfUpdateAppForm,
  toUpdateApplicationRequest,
} from '../../schema/applicationSchema'
import BranchField from './general/BranchField'
import NameField from './general/NameField'

type Props = {
  app: Application
  repo: Repository
  hasPermission: boolean
  refetch: () => void
}

const GeneralConfigForm: Component<Props> = (props) => {
  // このパネルが編集するフィールドだけを持たせ、他のパネルの設定を検証・送信しない
  const toInput = (app: Application) => getInitialValueOfUpdateAppForm(app, ['name', 'refName'] as const)
  const form = createForm({
    schema: createOrUpdateApplicationSchema,
    initialInput: untrack(() => toInput(props.app)),
  })
  const sync = syncFormWithServer(form, { key: () => props.app.id, input: () => toInput(props.app), root: ['form'] })

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        const submitted = sync.snapshot()
        try {
          await client.updateApplication(toUpdateApplicationRequest(output))
          sync.markSaved(submitted)
          toast.success('アプリケーション設定を更新しました')
          props.refetch()
          // 非同期でビルドが開始されるので1秒程度待ってから再度リロード
          setTimeout(props.refetch, 1000)
        } catch (e) {
          handleAPIError(e, 'アプリケーション設定の更新に失敗しました')
        }
      }}
    >
      <FormBox.Container>
        <FormBox.Forms>
          <NameField form={form} hasPermission={props.hasPermission} />
          <BranchField form={form} repo={props.repo} hasPermission={props.hasPermission} />
        </FormBox.Forms>
        <FormBox.Actions>
          <Button
            variants="primary"
            size="small"
            type="submit"
            disabled={!form.isDirty || form.isSubmitting || !props.hasPermission}
            loading={form.isSubmitting}
            tooltip={{
              content: !props.hasPermission
                ? '設定を変更するにはアプリケーションのオーナーになる必要があります'
                : undefined,
            }}
          >
            Save
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
    </Form>
  )
}

export default GeneralConfigForm
