import { PortPublication, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { pickRandom, randIntN } from '/@/libs/random'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { Field, FieldArray, FormStore, custom, getValue, insert, remove, setValue } from '@modular-forms/solid'
import { For, Show } from 'solid-js'
import { systemInfo } from '../../libs/api'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { TextField } from '../UI/TextField'
import { SelectOption, SingleSelect } from './Select'

const PortsContainer = styled('div', {
  base: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '16px',
  },
})
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

const protocolItems: SelectOption<PortPublicationProtocol>[] = [
  { value: PortPublicationProtocol.TCP, label: 'TCP' },
  { value: PortPublicationProtocol.UDP, label: 'UDP' },
]

const protoToName: Record<PortPublicationProtocol, string> = {
  [PortPublicationProtocol.TCP]: 'TCP',
  [PortPublicationProtocol.UDP]: 'UDP',
}

interface PortPublicationProps {
  formStore: FormStore<PortSettingsStore, undefined>
  name: `ports.${number}`
  deletePort: () => void
  hasPermission: boolean
}

const isValidPort = (port?: number, proto?: PortPublicationProtocol): boolean => {
  if (port === undefined) return false
  const available = systemInfo()?.ports.filter((a) => a.protocol === proto) || []
  if (available.length === 0) return false
  return available.some((range) => port >= range.startPort && port <= range.endPort)
}

const PortSetting = (props: PortPublicationProps) => {
  return (
    <PortRow>
      <PortVisualContainer variant="wrapper">
        <PortVisualContainer variant="from">
          <Field
            of={props.formStore}
            name={`${props.name}.internetPort`}
            type="number"
            validate={custom(
              (port) => isValidPort(port, getValue(props.formStore, `${props.name}.protocol`)),
              'Please enter the available port',
            )}
          >
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
          <Field of={props.formStore} name={`${props.name}.protocol`} type="number">
            {(field, fieldProps) => (
              <SingleSelect
                {...fieldProps}
                options={protocolItems}
                value={field.value}
                setValue={(value) => {
                  setValue(props.formStore, `${props.name}.protocol`, value)
                }}
                readOnly={!props.hasPermission}
              />
            )}
          </Field>
        </PortVisualContainer>
        <PortVisualContainer variant="to">
          <PortItem> → </PortItem>
          <Field of={props.formStore} name={`${props.name}.applicationPort`} type="number">
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
          <PortItem>/{protoToName[getValue(props.formStore, `${props.name}.protocol`)]}</PortItem>
        </PortVisualContainer>
      </PortVisualContainer>
      <Show when={props.hasPermission}>
        <Button onclick={props.deletePort} variants="textError" size="medium" type="button">
          Delete
        </Button>
      </Show>
    </PortRow>
  )
}

const suggestPort = (proto: PortPublicationProtocol): number => {
  const available = systemInfo()?.ports.filter((a) => a.protocol === proto) || []
  if (available.length === 0) return 0
  const range = pickRandom(available)
  return randIntN(range.endPort + 1 - range.startPort) + range.startPort
}

const newPort = (): PlainMessage<PortPublication> => {
  return {
    internetPort: suggestPort(PortPublicationProtocol.TCP),
    applicationPort: 0,
    protocol: PortPublicationProtocol.TCP,
  }
}

export type PortSettingsStore = {
  ports: PlainMessage<PortPublication>[]
}
interface PortPublicationSettingsProps {
  formStore: FormStore<PortSettingsStore, undefined>
  hasPermission: boolean
}

export const PortPublicationSettings = (props: PortPublicationSettingsProps) => {
  return (
    <PortsContainer>
      <FieldArray of={props.formStore} name="ports">
        {(fieldArray) => (
          <For each={fieldArray.items} fallback={'No Port Forwarding Configured'}>
            {(_, index) => (
              <PortSetting
                formStore={props.formStore}
                name={`${fieldArray.name}.${index()}`}
                deletePort={() => remove(props.formStore, 'ports', { at: index() })}
                hasPermission={props.hasPermission}
              />
            )}
          </For>
        )}
      </FieldArray>
      <Show when={props.hasPermission}>
        <Button
          onclick={() =>
            insert(props.formStore, 'ports', {
              value: newPort(),
            })
          }
          variants="border"
          size="small"
          type="button"
          leftIcon={<MaterialSymbols opticalSize={20}>add</MaterialSymbols>}
        >
          Add Port Forwarding
        </Button>
      </Show>
    </PortsContainer>
  )
}
