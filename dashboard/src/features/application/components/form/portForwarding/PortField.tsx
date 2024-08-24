import { styled } from '@macaron-css/solid'
import { Field, getValue, remove } from '@modular-forms/solid'
import { type Component, Show } from 'solid-js'
import { PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { type SelectOption, SingleSelect } from '/@/components/templates/Select'
import { useApplicationForm } from '../../../provider/applicationFormProvider'

const PortRow = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    gap: '24px',
  },
})
const PortVisualContainer = styled('div', {
  base: {
    alignItems: 'flex-start',
    gap: '8px',
  },
  variants: {
    variant: {
      from: {
        width: '100%',
        flexBasis: 'calc(60% - 4px)',
        display: 'grid',
        flexGrow: '1',
        gridTemplateColumns: 'minmax(calc(8ch + 32px), 1fr) auto minmax(calc(4ch + 60px), 1fr)',
      },
      to: {
        width: '100%',
        flexBasis: 'calc(40% - 4px)',
        display: 'grid',
        flexGrow: '1',
        gridTemplateColumns: 'auto minmax(calc(8ch + 32px), 1fr) auto',
      },
      wrapper: {
        width: '100%',
        display: 'flex',
        flexWrap: 'wrap',
      },
    },
  },
})
const PortItem = styled('div', {
  base: {
    height: '48px',
    display: 'flex',
    alignItems: 'center',
  },
})

const protocolItems: SelectOption<`${PortPublicationProtocol}`>[] = [
  { value: `${PortPublicationProtocol.TCP}`, label: 'TCP' },
  { value: `${PortPublicationProtocol.UDP}`, label: 'UDP' },
]

const protoToName: Record<`${PortPublicationProtocol}`, string> = {
  [PortPublicationProtocol.TCP]: 'TCP',
  [PortPublicationProtocol.UDP]: 'UDP',
}

type Props = {
  index: number
  hasPermission: boolean
}

const PortField: Component<Props> = (props) => {
  const { formStore } = useApplicationForm()

  const handleDelete = (index: number) => {
    remove(formStore, 'form.portPublications', { at: index })
  }

  return (
    <PortRow>
      <PortVisualContainer variant="wrapper">
        <PortVisualContainer variant="from">
          <Field of={formStore} name={`form.portPublications.${props.index}.internetPort`} type="number">
            {(field, fieldProps) => (
              <TextField
                type="number"
                placeholder="39000"
                tooltip={{
                  props: {
                    content: 'インターネット側ポート',
                  },
                }}
                {...fieldProps}
                value={field.value?.toString() ?? ''}
                error={field.error}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
          <PortItem>/</PortItem>
          <Field of={formStore} name={`form.portPublications.${props.index}.protocol`}>
            {(field, fieldProps) => (
              <SingleSelect
                {...fieldProps}
                options={protocolItems}
                value={field.value}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
        </PortVisualContainer>
        <PortVisualContainer variant="to">
          <PortItem> → </PortItem>
          <Field of={formStore} name={`form.portPublications.${props.index}.applicationPort`} type="number">
            {(field, fieldProps) => (
              <TextField
                type="number"
                placeholder="8080"
                tooltip={{
                  props: {
                    content: 'アプリ側ポート',
                  },
                }}
                {...fieldProps}
                value={field.value?.toString() ?? ''}
                error={field.error}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
          <Show when={getValue(formStore, `form.portPublications.${props.index}.protocol`)}>
            {(protocol) => <PortItem>/{protoToName[protocol()]}</PortItem>}
          </Show>
        </PortVisualContainer>
      </PortVisualContainer>
      <Show when={props.hasPermission}>
        <Button onclick={() => handleDelete(props.index)} variants="textError" size="medium" type="button">
          Delete
        </Button>
      </Show>
    </PortRow>
  )
}

export default PortField
