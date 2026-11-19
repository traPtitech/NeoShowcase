import { createForm, Form, reset } from '@formisch/solid/v2'
import { type Component, createEffect, untrack } from 'solid-js'
import type { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import {
  createOrUpdateApplicationSchema,
  getInitialValueOfUpdateAppForm,
  toUpdateApplicationRequest,
} from '../../schema/applicationSchema'
import WebsiteListField from './website/WebsiteListField'

type Props = {
  app: Application
  hasPermission: boolean
  refetch: () => void
}

const WebsitesConfigForm: Component<Props> = (props) => {
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

  const isRuntimeApp = () => {
    const configCase = props.app.config?.buildConfig.case
    return configCase === 'runtimeBuildpack' || configCase === 'runtimeDockerfile' || configCase === 'runtimeCmd'
  }

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        try {
          await client.updateApplication(toUpdateApplicationRequest(output))
          toast.success('ウェブサイト設定を更新しました')
          props.refetch()
          // 非同期でビルドが開始されるので1秒程度待ってから再度リロード
          setTimeout(props.refetch, 1000)
        } catch (e) {
          handleAPIError(e, 'ウェブサイト設定の更新に失敗しました')
        }
      }}
    >
      <div class="border-ui-border flex w-full flex-col gap-0.25 overflow-hidden rounded-lg border">
        <WebsiteListField form={form} isRuntimeApp={isRuntimeApp()} readonly={!props.hasPermission} />
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
      </div>
    </Form>
  )
}

export default WebsitesConfigForm
