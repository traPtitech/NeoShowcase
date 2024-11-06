import { Field, getValue, remove } from '@modular-forms/solid'
import { type Component, type ComponentProps, type ParentComponent, Show, splitProps } from 'solid-js'
import { PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { styled } from '/@/components/styled-components'
import { type SelectOption, SingleSelect } from '/@/components/templates/Select'
import { clsx } from '/@/libs/clsx'
import { useApplicationForm } from '../../../provider/applicationFormProvider'

const PortVisualContainer: ParentComponent<ComponentProps<'div'> & { variant: 'from' | 'to' | 'wrapper' }> = (
  props,
) => {
  const [_, extraProps] = splitProps(props, ['variant', 'class'])
  return (
    <div
      class={clsx(
        'items-start gap-2',
        props.variant === 'from' &&
          'grid w-full grow-1 basis-[calc(60%-4px)] grid-cols-[minmax(calc(8ch+32px),1fr)_auto_minmax(calc(4ch+60px),1fr)]',
        props.variant === 'to' &&
          'grid w-full grow-1 basis-[calc(40%-4px)] grid-cols-[auto_minmax(calc(8ch+32px),1fr)_auto]',
        props.variant === 'wrapper' && 'flex w-full flex-wrap',
        props.class,
      )}
      {...extraProps}
    />
  )
}
const PortItem = styled('div', 'flex h-12 items-center')

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
    <div class="flex w-full items-center gap-6">
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
    </div>
  )
}

export default PortField
