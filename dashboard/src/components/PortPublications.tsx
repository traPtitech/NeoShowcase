import { PortPublication, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { FormButton, FormSettings, FormSettingsButton, SettingsContainer } from '/@/components/AppsNew'
import { Button } from '/@/components/Button'
import { InputBar, InputLabel } from '/@/components/Input'
import { Radio, RadioItem } from '/@/components/Radio'
import { Select, SelectItem } from '/@/components/Select'
import { pickRandom, randIntN } from '/@/libs/random'
import { PlainMessage } from '@bufbuild/protobuf'
import { styled } from '@macaron-css/solid'
import { AiOutlinePlusCircle } from 'solid-icons/ai'
import { FaRegularTrashCan } from 'solid-icons/fa'
import { For } from 'solid-js'
import { SetStoreFunction } from 'solid-js/store'
import { systemInfo } from '../libs/api'
import { portPublicationProtocolMap } from '../libs/application'
import { vars } from '../theme'

const AvailablePortContainer = styled('div', {
  base: {
    fontSize: '14px',
    color: vars.text.black2,
    padding: '8px',
  },
})

const AvailableDomainUl = styled('ul', {
  base: {
    margin: '8px 0',
  },
})

const PortVisualContainer = styled('div', {
  base: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: '2px',
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
  return (
    <FormSettings>
      <PortVisualContainer>
        <InputBar
          placeholder="39000"
          type="number"
          value={props.port.internetPort || ''}
          onChange={(e) => props.setPort('internetPort', +e.target.value)}
          width="tiny"
          tooltip="インターネット側ポート"
        />
        <span>/</span>
        <Select
          items={protocolItems}
          selected={props.port.protocol}
          onSelect={(proto) => {
            props.setPort('protocol', proto)
          }}
        />
        <span> → </span>
        <InputBar
          placeholder="8080"
          type="number"
          value={props.port.applicationPort || ''}
          onChange={(e) => props.setPort('applicationPort', +e.target.value)}
          width="tiny"
          tooltip="アプリ側ポート"
        />
        <span>/{protoToName[props.port.protocol]}</span>
      </PortVisualContainer>
      <FormSettingsButton>
        <Button onclick={props.deletePort} color="black1" size="large" width="auto" type="button">
          <FaRegularTrashCan />
          <span> この設定を削除</span>
        </Button>
      </FormSettingsButton>
    </FormSettings>
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
    <SettingsContainer>
      <AvailablePortContainer>
        使用可能なポート
        <AvailableDomainUl>
          <For each={systemInfo()?.ports || []}>
            {(port) => (
              <li>
                {port.startPort}/{portPublicationProtocolMap[port.protocol]}~{port.endPort}/
                {portPublicationProtocolMap[port.protocol]}
              </li>
            )}
          </For>
        </AvailableDomainUl>
      </AvailablePortContainer>
      <For each={props.ports}>
        {(port, i) => (
          <PortSetting
            port={port}
            setPort={(valueName, value) => props.setPorts(i(), valueName, value)}
            deletePort={() => props.setPorts((current) => [...current.slice(0, i()), ...current.slice(i() + 1)])}
          />
        )}
      </For>

      <FormButton>
        <Button
          onclick={() => props.setPorts([...props.ports, newPort()])}
          color="black1"
          size="large"
          width="auto"
          type="button"
        >
          <AiOutlinePlusCircle />
          <span> 設定を追加</span>
        </Button>
      </FormButton>
    </SettingsContainer>
  )
}
