import { createForm, Form, reset } from '@formisch/solid/v2'
import { type Component, createEffect, untrack } from 'solid-js'
import type { Application, Repository } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
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
  const form = createForm({
    schema: createOrUpdateApplicationSchema,
    initialInput: untrack(() => getInitialValueOfUpdateAppForm(props.app)),
  })

  // サーバー側のデータが更新されたらフォームを巻き戻す
  createEffect(
    () => getInitialValueOfUpdateAppForm(props.app),
    (initialInput) => reset(form, { initialInput }),
    { defer: true },
  )

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        try {
          await client.updateApplication(toUpdateApplicationRequest(output))
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
