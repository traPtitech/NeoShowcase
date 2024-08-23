import { styled } from '@macaron-css/solid'
import { Field, FieldArray, Form, type SubmitHandler, insert, reset, setValues } from '@modular-forms/solid'
import { type Component, For, Show, createEffect, onMount, untrack } from 'solid-js'
import toast from 'solid-toast'
import { type Application, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { MaterialSymbols } from '/@/components/UI/MaterialSymbols'
import FormBox from '/@/components/layouts/FormBox'
import { client, handleAPIError, systemInfo } from '/@/libs/api'
import { pickRandom, randIntN } from '/@/libs/random'
import { colorVars } from '/@/theme'
import { useApplicationForm } from '../../provider/applicationFormProvider'
import {
  type CreateOrUpdateApplicationInput,
  handleSubmitUpdateApplicationForm,
  updateApplicationFormInitialValues,
} from '../../schema/applicationSchema'
import type { PortPublicationInput } from '../../schema/portPublicationSchema'
import PortField from './portForwarding/PortField'

const PortsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '16px',
  },
})
const FallbackText = styled('div', {
  color: colorVars.semantic.text.black,
})

const suggestPort = (protocol: PortPublicationProtocol): number => {
  const available = systemInfo()?.ports.filter((a) => a.protocol === protocol) || []
  if (available.length === 0) return 0
  const range = pickRandom(available)
  return randIntN(range.endPort + 1 - range.startPort) + range.startPort
}

const newPort = (): PortPublicationInput => {
  return {
    internetPort: suggestPort(PortPublicationProtocol.TCP),
    applicationPort: 0,
    protocol: `${PortPublicationProtocol.TCP}`,
  }
}

type Props = {
  app: Application
  refetchApp: () => Promise<void>
  hasPermission: boolean
}

const PortForwardingForm: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  // `reset` doesn't work on first render when the Field not rendered
  // see: https://github.com/fabian-hiller/modular-forms/issues/157#issuecomment-1848567069
  onMount(() => {
    setValues(formStore, updateApplicationFormInitialValues(props.app))
  })

  // reset forms when props.app changed
  createEffect(() => {
    reset(
      untrack(() => formStore),
      {
        initialValues: updateApplicationFormInitialValues(props.app),
      },
    )
  })

  const handleAdd = () => {
    insert(formStore, 'form.portPublications', {
      value: newPort(),
    })
  }

  const handleSubmit: SubmitHandler<CreateOrUpdateApplicationInput> = (values) =>
    handleSubmitUpdateApplicationForm(values, async (output) => {
      try {
        await client.updateApplication(output)
        toast.success('ポート公開設定を更新しました')
        props.refetchApp()
      } catch (e) {
        handleAPIError(e, 'ポート公開設定の更新に失敗しました')
      }
    })

  return (
    <Form of={formStore} onSubmit={handleSubmit}>
      <Field of={formStore} name="type">
        {() => null}
      </Field>
      <Field of={formStore} name="form.id">
        {() => null}
      </Field>
      <FormBox.Container>
        <FormBox.Forms>
          <PortsContainer>
            <FieldArray of={formStore} name="form.portPublications">
              {(fieldArray) => (
                <For each={fieldArray.items} fallback={<FallbackText>ポート公開が設定されていません</FallbackText>}>
                  {(_, index) => <PortField index={index()} hasPermission={props.hasPermission} />}
                </For>
              )}
            </FieldArray>
            <Button
              onclick={handleAdd}
              variants="border"
              size="small"
              type="button"
              leftIcon={<MaterialSymbols opticalSize={20}>add</MaterialSymbols>}
            >
              Add Port Forwarding
            </Button>
          </PortsContainer>
        </FormBox.Forms>
        <FormBox.Actions>
          <Button
            variants="primary"
            size="small"
            type="submit"
            disabled={formStore.invalid || !formStore.dirty || formStore.submitting || !props.hasPermission}
            loading={formStore.submitting}
            tooltip={{
              props: {
                content: !props.hasPermission
                  ? '設定を変更するにはアプリケーションのオーナーになる必要があります'
                  : undefined,
              },
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
