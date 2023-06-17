import { PortPublication, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { InputBar, InputLabel } from '/@/components/Input'
import { Radio, RadioItem } from '/@/components/Radio'
import { Button } from '/@/components/Button'
import { SetStoreFunction } from 'solid-js/store'
import { FormButton, FormSettings, FormSettingsButton, SettingsContainer } from '/@/components/AppsNew'
import { createEffect, For } from 'solid-js'
import { styled } from '@macaron-css/solid'
import { vars } from '../theme'
import { availablePorts } from '../libs/api'
import { portPublicationProtocolMap } from '../libs/application'
import { PlainMessage } from '@bufbuild/protobuf'
import { pickRandom, randIntN } from '/@/libs/random'

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

interface PortPublicationProps {
  port: PlainMessage<PortPublication>
  setPort: <T extends keyof PlainMessage<PortPublication>>(
    valueName: T,
    value: PlainMessage<PortPublication>[T],
  ) => void
  deletePort: () => void
}

const PortSetting = (props: PortPublicationProps) => {
  createEffect(() => {
    props.setPort('internetPort', suggestPort(props.port.protocol))
  })

  return (
    <FormSettings>
      <div>
        <InputLabel>Protocol</InputLabel>
        <Radio
          items={protocolItems}
          selected={props.port.protocol}
          setSelected={(proto) => props.setPort('protocol', proto)}
        />
      </div>
      <div>
        <InputLabel>Internet Port</InputLabel>
        <InputBar
          placeholder='39000'
          type='number'
          value={props.port.internetPort || ''}
          onChange={(e) => props.setPort('internetPort', +e.target.value)}
        />
      </div>
      <div>
        <InputLabel>Application Port</InputLabel>
        <InputBar
          placeholder='8080'
          type='number'
          value={props.port.applicationPort || ''}
          onChange={(e) => props.setPort('applicationPort', +e.target.value)}
        />
      </div>

      <FormSettingsButton>
        <Button onclick={props.deletePort} color='black1' size='large' width='auto' type='button'>
          Delete port publication
        </Button>
      </FormSettingsButton>
    </FormSettings>
  )
}

const protocolItems: RadioItem<PortPublicationProtocol>[] = [
  { value: PortPublicationProtocol.TCP, title: 'TCP' },
  { value: PortPublicationProtocol.UDP, title: 'UDP' },
]

const suggestPort = (proto: PortPublicationProtocol): number => {
  const available = availablePorts()?.availablePorts.filter((a) => a.protocol === proto) || []
  if (available.length === 0) return 0
  const range = pickRandom(available)
  return randIntN(range.endPort + 1 - range.startPort) + range.startPort
}

const newPort = (): PlainMessage<PortPublication> => {
  return {
    internetPort: 0,
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
          <For each={availablePorts()?.availablePorts || []}>
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
          color='black1'
          size='large'
          width='auto'
          type='button'
        >
          Add port publication
        </Button>
      </FormButton>
    </SettingsContainer>
  )
}
