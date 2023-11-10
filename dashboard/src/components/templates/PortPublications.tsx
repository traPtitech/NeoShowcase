import { PortPublication, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { pickRandom, randIntN } from '/@/libs/random'
import useModal from '/@/libs/useModal'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { For } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { systemInfo } from '../../libs/api'
import { Button } from '../UI/Button'
import { MaterialSymbols } from '../UI/MaterialSymbols'
import { TextInput } from '../UI/TextInput'
import { SelectItem, SingleSelect } from './Select'

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
    alignItems: 'center',
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
        gridTemplateColumns: 'auto auto',
      },
    },
  },
})

const protocolItems: SelectItem<PortPublicationProtocol>[] = [
  { value: PortPublicationProtocol.TCP, title: 'TCP' },
  { value: PortPublicationProtocol.UDP, title: 'UDP' },
]

const protoToName: Record<PortPublicationProtocol, string> = {
  [PortPublicationProtocol.TCP]: 'TCP',
  [PortPublicationProtocol.UDP]: 'UDP',
}

interface PortPublicationProps {
  port: PlainMessage<PortPublication>
  setPort: <T extends keyof PlainMessage<PortPublication>>(
    valueName: T,
    value: PlainMessage<PortPublication>[T],
  ) => void
  deletePort: () => void
}

const PortSetting = (props: PortPublicationProps) => {
  const { Modal } = useModal()

  return (
    <PortRow>
      <PortVisualContainer variant="wrapper">
        <PortVisualContainer variant="from">
          <TextInput
            required
            placeholder="39000"
            type="number"
            value={props.port.internetPort || ''}
            onChange={(e) => props.setPort('internetPort', +e.target.value)}
            tooltip={{
              props: {
                content: 'インターネット側ポート',
              },
            }}
          />
          <span>/</span>
          <SingleSelect
            items={protocolItems}
            selected={props.port.protocol}
            setSelected={(proto) => {
              console.log(`setting ${proto}, type: ${typeof proto}`)
              props.setPort('protocol', proto)
            }}
          />
        </PortVisualContainer>
        <PortVisualContainer variant="to">
          <span> → </span>
          <TextInput
            required
            placeholder="8080"
            type="number"
            value={props.port.applicationPort || ''}
            onChange={(e) => props.setPort('applicationPort', +e.target.value)}
            tooltip={{
              props: {
                content: 'アプリ側ポート',
              },
            }}
          />
          <span>/{protoToName[props.port.protocol]}</span>
        </PortVisualContainer>
      </PortVisualContainer>
      <Button onclick={props.deletePort} variants="textError" size="medium" type="button">
        Delete
      </Button>
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

interface PortPublicationSettingsProps {
  ports: PlainMessage<PortPublication>[]
  setPorts: SetStoreFunction<PlainMessage<PortPublication>[]>
}

export const PortPublicationSettings = (props: PortPublicationSettingsProps) => {
  return (
    <PortsContainer>
      <For each={props.ports}>
        {(port, i) => (
          <PortSetting
            port={port}
            setPort={(valueName, value) => props.setPorts(i(), valueName, value)}
            deletePort={() => props.setPorts((current) => [...current.slice(0, i()), ...current.slice(i() + 1)])}
          />
        )}
      </For>
      <Button
        onclick={() => props.setPorts([...props.ports, newPort()])}
        variants="border"
        size="small"
        type="button"
        leftIcon={<MaterialSymbols opticalSize={20}>add</MaterialSymbols>}
      >
        Add More
      </Button>
    </PortsContainer>
  )
}
