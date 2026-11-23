import { createForm, Field, Form, reset } from '@formisch/solid/v2'
import { createEffect, Loading, Show, untrack } from 'solid-js'
import * as v from 'valibot'
import type { Application } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '/@/components/layouts/DataTable'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import Skeleton from '/@/components/UI/Skeleton'
import { useApplicationData } from '/@/features/application/data'
import GroupSelect from '/@/features/applicationGroup/components/GroupSelect'
import { client, handleAPIError, revalidateApplicationGroupApps, revalidateMyApplicationGroups } from '/@/libs/api'
import { toast } from '/@/libs/toast'

const moveSchema = v.object({ groupId: v.optional(v.string()) })

const toInput = (app: Application) => ({ groupId: app.groupId })

export default () => {
  const { app, group: currentGroup, refetch, hasPermission } = useApplicationData()

  const form = createForm({ schema: moveSchema, initialInput: untrack(() => toInput(app())) })
  // サーバー側のデータが更新されたらフォームを巻き戻す
  createEffect(
    () => app(),
    (app) => reset(form, { initialInput: toInput(app) }),
    { defer: true },
  )

  return (
    <DataTable.Container>
      <DataTable.Title>Group</DataTable.Title>
      <DataTable.SubTitle>
        同じグループのアプリ同士は URL を入れ子にでき、グループのオーナーがグループ内のすべてのアプリを操作できます。
        別のグループへ移すには、移動元と移動先の両方のオーナーである必要があります。
      </DataTable.SubTitle>
      <Form
        of={form}
        onSubmit={async (output) => {
          const from = app().groupId
          try {
            await client.moveApplication({
              applicationId: app().id,
              groupId: output.groupId,
            })
            toast.success('アプリケーションのグループを変更しました')
            refetch()
            revalidateMyApplicationGroups()
            if (from) revalidateApplicationGroupApps(from)
            if (output.groupId) revalidateApplicationGroupApps(output.groupId)
          } catch (e) {
            handleAPIError(e, 'アプリケーションのグループの変更に失敗しました')
          }
        }}
      >
        <FormBox.Container>
          <FormBox.Forms>
            <Loading fallback={<Skeleton height={48} />}>
              <Field of={form} path={['groupId']}>
                {(field) => (
                  <GroupSelect
                    name={field.props.name}
                    label="Group"
                    value={field.input}
                    setValue={(groupId) => field.onInput(groupId)}
                    current={currentGroup()}
                    readOnly={!hasPermission()}
                  />
                )}
              </Field>
            </Loading>
            <Show when={currentGroup()}>
              {(group) => (
                <a href={`/groups/${group().id}`} class="text-regular text-primary-main w-fit">
                  {group().name} のグループページへ
                </a>
              )}
            </Show>
          </FormBox.Forms>
          <FormBox.Actions>
            <Button
              variants="primary"
              size="small"
              type="submit"
              disabled={!form.isDirty || form.isSubmitting || !hasPermission()}
              loading={form.isSubmitting}
              tooltip={{
                content: !hasPermission() ? 'グループを変更するにはオーナーになる必要があります' : undefined,
              }}
            >
              Save
            </Button>
          </FormBox.Actions>
        </FormBox.Container>
      </Form>
    </DataTable.Container>
  )
}
