import type { PartialMessage } from '@bufbuild/protobuf'
import * as v from 'valibot'
import { type PortPublication, PortPublicationProtocol } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { systemInfo } from '/@/libs/api'

const isValidPort = (port?: number, protocol?: PortPublicationProtocol): boolean => {
  if (port === undefined) return false
  const available = systemInfo()?.ports.filter((a) => a.protocol === protocol) || []
  if (available.length === 0) return false
  return available.some((range) => port >= range.startPort && port <= range.endPort)
}

// KobalteのSelectではstringしか扱えないためform内では文字列として持つ
const protocolSchema = v.pipe(
  v.union([v.literal(`${PortPublicationProtocol.TCP}`), v.literal(`${PortPublicationProtocol.UDP}`)]),
  v.transform((input): PortPublicationProtocol => {
    switch (input) {
      case `${PortPublicationProtocol.TCP}`: {
        return PortPublicationProtocol.TCP
      }
      case `${PortPublicationProtocol.UDP}`: {
        return PortPublicationProtocol.UDP
      }
      default: {
        const _unreachable: never = input
        throw new Error('unknown PortPublicationProtocol')
      }
    }
  }),
)

export const portPublicationSchema = v.pipe(
  v.object({
    internetPort: v.pipe(v.number(), v.integer()),
    applicationPort: v.pipe(v.number(), v.integer()),
    protocol: protocolSchema,
  }),
  v.transform((input): PartialMessage<PortPublication> => input),
  v.forward(
    v.partialCheck(
      [['internetPort'], ['protocol']],
      (input) => isValidPort(input.internetPort, input.protocol),
      'Please enter the available port',
    ),
    ['internetPort'],
  ),
  v.forward(
    v.partialCheck(
      [['applicationPort'], ['protocol']],
      (input) => isValidPort(input.applicationPort, input.protocol),
      'Please enter the available port',
    ),
    ['applicationPort'],
  ),
)

export type PortPublicationInput = v.InferInput<typeof portPublicationSchema>

export const portPublicationMessageToSchema = (portPublication: PortPublication): PortPublicationInput => {
  return {
    applicationPort: portPublication.applicationPort,
    internetPort: portPublication.internetPort,
    protocol: `${portPublication.protocol}`,
  }
}
