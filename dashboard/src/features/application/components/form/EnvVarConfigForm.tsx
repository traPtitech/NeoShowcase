import { createForm, Form, reset } from '@formisch/solid/v2'
import { type Component, createEffect, Show, untrack } from 'solid-js'
import type { ApplicationEnvVars } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import { envVarSchema, envVarsMessageToSchema } from '../../schema/envVarSchema'
import EnvVarConfigField from './config/EnvVarConfigField'

type Props = {
  appId: string
  envVars: ApplicationEnvVars
  refetch: () => void
}

const EnvVarConfigForm: Component<Props> = (props) => {
  const form = createForm({
    schema: envVarSchema,
    initialInput: untrack(() => envVarsMessageToSchema(props.envVars)),
  })

  const discardChanges = () => reset(form, { initialInput: envVarsMessageToSchema(props.envVars) })

  // サーバー側のデータが更新されたらフォームを巻き戻す
  createEffect(
    () => props.envVars,
    (envVars) => reset(form, { initialInput: envVarsMessageToSchema(envVars) }),
    { defer: true },
  )

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        const oldVars = new Map(
          props.envVars.variables.filter((envVar) => !envVar.system).map((envVar) => [envVar.key, envVar.value]),
        )
        const newVars = new Map(
          output.variables
            .filter((envVar) => !envVar.system && envVar.key !== '')
            .map((envVar) => [envVar.key, envVar.value]),
        )
        const addedKeys = [...newVars.keys()].filter((key) => !oldVars.has(key))
        const deletedKeys = [...oldVars.keys()].filter((key) => !newVars.has(key))
        const updatedKeys = [...oldVars.keys()].filter(
          (key) => newVars.has(key) && newVars.get(key) !== oldVars.get(key),
        )

        try {
          await Promise.all([
            ...[...addedKeys, ...updatedKeys].map((key) =>
              client.setEnvVar({ applicationId: props.appId, key, value: newVars.get(key) }),
            ),
            ...deletedKeys.map((key) => client.deleteEnvVar({ applicationId: props.appId, key })),
          ])
          toast.success('環境変数を更新しました')
          props.refetch()
          // 非同期でビルドが開始されるので1秒程度待ってから再度リロード
          setTimeout(props.refetch, 1000)
        } catch (e) {
          handleAPIError(e, '環境変数の更新に失敗しました')
        }
      }}
    >
      <FormBox.Container>
        <FormBox.Forms>
          <EnvVarConfigField form={form} />
        </FormBox.Forms>
        <FormBox.Actions>
          <Show when={form.isDirty && !form.isSubmitting}>
            <Button variants="ghost" size="small" type="button" onClick={discardChanges}>
              Discard Changes
            </Button>
          </Show>
          <Button
            variants="primary"
            size="small"
            type="submit"
            disabled={!form.isDirty || form.isSubmitting}
            loading={form.isSubmitting}
          >
            Save
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
    </Form>
  )
}

export default EnvVarConfigForm
