import { Field, getInput, remove } from '@formisch/solid/v2'
import type { ComponentProps } from '@solidjs/web'
import { type Component, omit, type ParentComponent, Show } from 'solid-js'
import { PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { styled } from '/@/components/styled-components'
import { type SelectOption, SingleSelect } from '/@/components/templates/Select'
import { Button } from '/@/components/UI/Button'
import { TextField } from '/@/components/UI/TextField'
import { inputToNumber, numberToInput } from '/@/libs/schemaUtil'
import type { ApplicationFormStore } from '../../../schema/applicationSchema'

const PortVisualContainer: ParentComponent<ComponentProps<'div'> & { variant: 'from' | 'to' | 'wrapper' }> = (
  props,
) => {
  const extraProps = omit(props, 'variant', 'class')
  return (
    <div
      class={[
        'items-start gap-2',
        props.variant === 'from' &&
          'grid w-full grow-1 basis-[calc(60%-4px)] grid-cols-[minmax(calc(8ch+32px),1fr)_auto_minmax(calc(4ch+60px),1fr)]',
        props.variant === 'to' &&
          'grid w-full grow-1 basis-[calc(40%-4px)] grid-cols-[auto_minmax(calc(8ch+32px),1fr)_auto]',
        props.variant === 'wrapper' && 'flex w-full flex-wrap',
        props.class,
      ]}
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
  form: ApplicationFormStore
  index: number
  hasPermission: boolean
}

const PortField: Component<Props> = (props) => (
  <div class="flex w-full items-center gap-6">
    <PortVisualContainer variant="wrapper">
      <PortVisualContainer variant="from">
        <Field of={props.form} path={['form', 'portPublications', props.index, 'internetPort']}>
          {(field) => (
            <TextField
              type="number"
              placeholder="39000"
              tooltip={{ content: 'インターネット側ポート' }}
              {...field.props}
              onInput={(e) => field.onInput(inputToNumber(e.currentTarget.value))}
              value={numberToInput(field.input)}
              error={field.errors?.[0]}
              readOnly={!props.hasPermission}
            />
          )}
        </Field>
        <PortItem>/</PortItem>
        <Field of={props.form} path={['form', 'portPublications', props.index, 'protocol']}>
          {(field) => (
            <SingleSelect
              {...field.props}
              options={protocolItems}
              value={field.input}
              setValue={(v) => field.onInput(v)}
              readOnly={!props.hasPermission}
            />
          )}
        </Field>
      </PortVisualContainer>
      <PortVisualContainer variant="to">
        <PortItem> → </PortItem>
        <Field of={props.form} path={['form', 'portPublications', props.index, 'applicationPort']}>
          {(field) => (
            <TextField
              type="number"
              placeholder="8080"
              tooltip={{ content: 'アプリ側ポート' }}
              {...field.props}
              onInput={(e) => field.onInput(inputToNumber(e.currentTarget.value))}
              value={numberToInput(field.input)}
              error={field.errors?.[0]}
              readOnly={!props.hasPermission}
            />
          )}
        </Field>
        <Show when={getInput(props.form, { path: ['form', 'portPublications', props.index, 'protocol'] })}>
          {(protocol) => <PortItem>/{protoToName[protocol()]}</PortItem>}
        </Show>
      </PortVisualContainer>
    </PortVisualContainer>
    <Show when={props.hasPermission}>
      <Button
        onClick={() => remove(props.form, { path: ['form', 'portPublications'], at: props.index })}
        variants="textError"
        size="medium"
        type="button"
      >
        Delete
      </Button>
    </Show>
  </div>
)

export default PortField
