import { createForm, FieldArray, Form, insert } from '@formisch/solid/v2'
import { type Component, For, untrack } from 'solid-js'
import { type Application, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import FormBox from '/@/components/layouts/FormBox'
import { Button } from '/@/components/UI/Button'
import { client, handleAPIError } from '/@/libs/api'
import { pickRandom, randIntN } from '/@/libs/random'
import { syncFormWithServer } from '/@/libs/serverForm'
import { systemInfo } from '/@/libs/session'
import { toast } from '/@/libs/toast'
import {
  createOrUpdateApplicationSchema,
  getInitialValueOfUpdateAppForm,
  toUpdateApplicationRequest,
} from '../../schema/applicationSchema'
import type { PortPublicationInput } from '../../schema/portPublicationSchema'
import PortField from './portForwarding/PortField'

const suggestPort = (protocol: PortPublicationProtocol): number => {
  const available = systemInfo().ports.filter((a) => a.protocol === protocol)
  if (available.length === 0) return 0
  const range = pickRandom(available)
  return randIntN(range.endPort + 1 - range.startPort) + range.startPort
}

const newPort = (): PortPublicationInput => ({
  internetPort: suggestPort(PortPublicationProtocol.TCP),
  applicationPort: 0,
  protocol: `${PortPublicationProtocol.TCP}`,
})

type Props = {
  app: Application
  hasPermission: boolean
  refetch: () => void
}

const PortForwardingForm: Component<Props> = (props) => {
  // このパネルが編集するフィールドだけを持たせ、他のパネルの設定を検証・送信しない
  const toInput = (app: Application) => getInitialValueOfUpdateAppForm(app, ['portPublications'] as const)
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
          toast.success('ポート公開設定を更新しました')
          props.refetch()
        } catch (e) {
          handleAPIError(e, 'ポート公開設定の更新に失敗しました')
        }
      }}
    >
      <FormBox.Container>
        <FormBox.Forms>
          <div class="flex w-full flex-col items-center gap-4">
            <FieldArray of={form} path={['form', 'portPublications']}>
              {(fieldArray) => (
                <For
                  each={fieldArray.items}
                  fallback={<div class="text-text-black">ポート公開が設定されていません</div>}
                >
                  {(_, index) => <PortField form={form} index={index()} hasPermission={props.hasPermission} />}
                </For>
              )}
            </FieldArray>
            <Button
              onClick={() => insert(form, { path: ['form', 'portPublications'], initialInput: newPort() })}
              variants="border"
              size="small"
              type="button"
              leftIcon={<div class="i-material-symbols:add shrink-0 text-xl/5" />}
            >
              Add Port Forwarding
            </Button>
          </div>
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

export default PortForwardingForm
