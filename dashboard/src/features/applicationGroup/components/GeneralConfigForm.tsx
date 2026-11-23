import { createForm, Field, Form, reset } from '@formisch/solid/v2'
import { type Component, createEffect, untrack } from 'solid-js'
import type { ApplicationGroup } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { client, handleAPIError } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import { type ApplicationGroupInput, applicationGroupSchema } from '../schema/applicationGroupSchema'

type Props = {
  group: ApplicationGroup
  hasPermission: boolean
  refetch: () => void
}

const toInput = (group: ApplicationGroup): ApplicationGroupInput => ({ name: group.name })

const GeneralConfigForm: Component<Props> = (props) => {
  const form = createForm({
    schema: applicationGroupSchema,
    initialInput: untrack(() => toInput(props.group)),
  })

  // サーバー側のデータが更新されたらフォームを巻き戻す
  createEffect(
    () => props.group,
    (group) => reset(form, { initialInput: toInput(group) }),
    { defer: true },
  )

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        try {
          await client.updateApplicationGroup({ id: props.group.id, name: output.name })
          toast.success('グループ名を更新しました')
          props.refetch()
        } catch (e) {
          handleAPIError(e, 'グループ名の更新に失敗しました')
        }
      }}
    >
      <FormBox.Container>
        <FormBox.Forms>
          <Field of={form} path={['name']}>
            {(field) => (
              <TextField
                label="Group Name"
                required
                {...field.props}
                value={field.input ?? ''}
                error={field.errors?.[0]}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
        </FormBox.Forms>
        <FormBox.Actions>
          <Button
            variants="primary"
            size="small"
            type="submit"
            disabled={!form.isDirty || form.isSubmitting || !props.hasPermission}
            loading={form.isSubmitting}
            tooltip={{
              content: !props.hasPermission ? '設定を変更するにはグループのオーナーになる必要があります' : undefined,
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
