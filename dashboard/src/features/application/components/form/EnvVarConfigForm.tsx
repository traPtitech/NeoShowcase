import { Form, reset, type SubmitHandler, setValues } from '@modular-forms/solid'
import { type Component, createEffect, onMount, Show, untrack } from 'solid-js'
import toast from 'solid-toast'
import type { ApplicationEnvVars } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
import { useEnvVarConfigForm } from '../../provider/envVarConfigFormProvider'
import { type EnvVarInput, envVarsMessageToSchema, handleSubmitEnvVarForm } from '../../schema/envVarSchema'
import EnvVarConfigField from './config/EnvVarConfigField'

type Props = {
  appId: string
  refetch: () => void
  envVars: ApplicationEnvVars
}

const EnvVarConfigForm: Component<Props> = (props) => {
  const { formStore } = useEnvVarConfigForm()

  const discardChanges = () => {
    reset(
      untrack(() => formStore),
      {
        initialValues: envVarsMessageToSchema(props.envVars),
      },
    )
  }

  // `reset` doesn't work on first render when the Field not rendered
  // see: https://github.com/fabian-hiller/modular-forms/issues/157#issuecomment-1848567069
  onMount(() => {
    setValues(formStore, envVarsMessageToSchema(props.envVars))
  })

  // reset forms when props.envVars changed
  createEffect(() => {
    discardChanges()
  })

  const handleSubmit: SubmitHandler<EnvVarInput> = (values) =>
    handleSubmitEnvVarForm(values, async (input) => {
      const oldVars = new Map(
        props.envVars.variables.filter((envVar) => !envVar.system).map((envVar) => [envVar.key, envVar.value]),
      )
      const newVars = new Map(
        input.variables
          .filter((envVar) => !envVar.system && envVar.key !== '')
          .map((envVar) => [envVar.key, envVar.value]),
      )
      const addedKeys = Array.from(newVars.keys()).filter((key) => !oldVars.has(key))
      const deletedKeys = Array.from(oldVars.keys()).filter((key) => !newVars.has(key))
      const updatedKeys = Array.from(oldVars.keys()).filter(
        (key) => newVars.has(key) && newVars.get(key) !== oldVars.get(key),
      )

      const addEnvVarRequests = [...addedKeys, ...updatedKeys].map((key) => {
        return client.setEnvVar({
          applicationId: props.appId,
          key,
          value: newVars.get(key),
        })
      })
      const deleteEnvVarRequests = deletedKeys.map((key) => {
        return client.deleteEnvVar({
          applicationId: props.appId,
          key,
        })
      })
      try {
        await Promise.all([...addEnvVarRequests, ...deleteEnvVarRequests])
        toast.success('環境変数を更新しました')
        props.refetch()
        // 非同期でビルドが開始されるので1秒程度待ってから再度リロード
        setTimeout(props.refetch, 1000)
      } catch (e) {
        handleAPIError(e, '環境変数の更新に失敗しました')
      }
    })

  return (
    <Form of={formStore} onSubmit={handleSubmit} shouldActive={false}>
      <FormBox.Container>
        <FormBox.Forms>
          <EnvVarConfigField />
        </FormBox.Forms>
        <FormBox.Actions>
          <Show when={formStore.dirty && !formStore.submitting}>
            <Button variants="ghost" size="small" type="button" onClick={discardChanges}>
              Discard Changes
            </Button>
          </Show>
          <Button
            variants="primary"
            size="small"
            type="submit"
            disabled={formStore.invalid || !formStore.dirty || formStore.submitting}
            loading={formStore.submitting}
          >
            Save
          </Button>
        </FormBox.Actions>
      </FormBox.Container>
    </Form>
  )
}

export default EnvVarConfigForm
