import { createForm, Field, Form } from '@formisch/solid/v2'
import { useNavigate } from '@solidjs/router'
import type { Component } from 'solid-js'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { client, handleAPIError, revalidateMyApplicationGroups } from '/@/libs/api'
import { toast } from '/@/libs/toast'
import { applicationGroupSchema } from '../schema/applicationGroupSchema'

const CreateForm: Component = () => {
  const navigate = useNavigate()
  const form = createForm({ schema: applicationGroupSchema, initialInput: { name: '' } })

  return (
    <Form
      of={form}
      onSubmit={async (output) => {
        try {
          const group = await client.createApplicationGroup({ name: output.name })
          toast.success('グループを作成しました')
          revalidateMyApplicationGroups()
          navigate(`/groups/${group.id}`)
        } catch (e) {
          handleAPIError(e, 'グループの作成に失敗しました')
        }
      }}
    >
      <div class="flex w-full flex-col items-end gap-10">
        <div class="bg-ui-primary mx-auto flex w-full flex-col gap-6 rounded-lg px-6 py-5">
          <div class="caption-regular text-text-grey">
            グループにまとめたアプリ同士は URL を入れ子にでき (例: example.trap.show/ と example.trap.show/api)、
            グループのオーナーはグループ内のすべてのアプリを操作できます。
          </div>
          <Field of={form} path={['name']}>
            {(field) => (
              <TextField
                label="Group Name"
                required
                {...field.props}
                value={field.input ?? ''}
                error={field.errors?.[0]}
              />
            )}
          </Field>
        </div>
        <Button
          variants="primary"
          size="medium"
          rightIcon={<div class="i-material-symbols:arrow-forward shrink-0 text-2xl/6" />}
          type="submit"
          disabled={form.isSubmitting}
          loading={form.isSubmitting}
        >
          Create
        </Button>
      </div>
    </Form>
  )
}

export default CreateForm
