import { AnyMessage, Message, PlainMessage } from '@bufbuild/protobuf'

export const getPlainMessage = <T extends Message<T> = AnyMessage>(message: Message<T>): PlainMessage<T> => {
  return Object.keys(message).reduce((acc, key) => {
    const value = message[key as keyof Message<T>]
    if (Array.isArray(value)) {
      return {
        ...acc,
        [key]: value.map(getPlainMessage),
      }
    } else if (value instanceof Message) {
      return {
        ...acc,
        [key]: getPlainMessage(value),
      }
    } else if (isCaseValueObject(value)) {
      return {
        ...acc,
        [key]: {
          case: value.case,
          value: value.value instanceof Message ? getPlainMessage(value.value) : value.value,
        },
      }
    } else {
      return {
        ...acc,
        [key]: value,
      }
    }
  }, {} as PlainMessage<T>)
}
const isCaseValueObject = (obj: unknown): obj is { case: string; value: unknown } => {
  return typeof obj === 'object' && obj !== null && 'case' in obj && 'value' in obj
}
