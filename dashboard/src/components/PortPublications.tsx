import { PortPublication, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { InputBar, InputLabel } from '/@/components/Input'
import { Radio, RadioItem } from '/@/components/Radio'
import { Button } from '/@/components/Button'
import { SetStoreFunction } from 'solid-js/store'
import { FormButton, FormSettings, FormSettingsButton, SettingsContainer } from '/@/components/AppsNew'
import { For } from 'solid-js'
import { styled } from '@macaron-css/solid'
import { vars } from '../theme'
import { availablePorts } from '../libs/api'
import { portPublicationProtocolMap } from '../libs/application'
import { PlainMessage } from '@bufbuild/protobuf'

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
  portPublication: PlainMessage<PortPublication>
  setPortPublication: <T extends keyof PlainMessage<PortPublication>>(
    valueName: T,
    value: PlainMessage<PortPublication>[T],
  ) => void
  deletePortPublication: () => void
}

const PortPublications = (props: PortPublicationProps) => {
  return (
    <FormSettings>
      <div>
        <InputLabel>Protocol</InputLabel>
        <Radio
          items={protocolItems}
          selected={props.portPublication.protocol}
          setSelected={(proto) => props.setPortPublication('protocol', proto)}
        />
      </div>
      <div>
        <InputLabel>Internet Port</InputLabel>
        <InputBar
          placeholder='39000'
          type='number'
          onChange={(e) => props.setPortPublication('internetPort', +e.target.value)}
        />
        <AvailablePortContainer>
          使用可能なポート
          <AvailableDomainUl>
            <For each={availablePorts()?.availablePorts}>
              {(port) => (
                <li>
                  {port.startPort}/{portPublicationProtocolMap[port.protocol]}~{port.endPort}/
                  {portPublicationProtocolMap[port.protocol]}
                </li>
              )}
            </For>
          </AvailableDomainUl>
        </AvailablePortContainer>
      </div>
      <div>
        <InputLabel>Application Port</InputLabel>
        <InputBar
          placeholder='8080'
          type='number'
          onChange={(e) => props.setPortPublication('applicationPort', +e.target.value)}
        />
      </div>

      <FormSettingsButton>
        <Button onclick={props.deletePortPublication} color='black1' size='large' width='auto' type='button'>
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

const newPort = (): PlainMessage<PortPublication> => ({
  internetPort: 0,
  applicationPort: 0,
  protocol: PortPublicationProtocol.TCP,
})

interface PortPublicationSettingsProps {
  ports: PlainMessage<PortPublication>[]
  setPorts: SetStoreFunction<PlainMessage<PortPublication>[]>
}

export const PortPublicationSettings = (props: PortPublicationSettingsProps) => {
  return (
    <SettingsContainer>
      <For each={props.ports}>
        {(portPublication, i) => (
          <PortPublications
            portPublication={portPublication}
            setPortPublication={(valueName, value) => {
              props.setPorts(i(), valueName, value)
            }}
            deletePortPublication={() =>
              props.setPorts((current) => [...current.slice(0, i()), ...current.slice(i() + 1)])
            }
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
