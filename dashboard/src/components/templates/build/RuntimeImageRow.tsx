import type { Component } from 'solid-js'
import type { RuntimeImage } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { List } from '../List'

const humanReadableSize = (size: bigint): string => {
  const s = Number(size)
  if (size < 1024) {
    return `${size} B`
  }
  if (size < 1024 * 1024) {
    return `${(s / 1024).toFixed(2)} KB`
  }
  if (size < 1024 * 1024 * 1024) {
    return `${(s / 1024 / 1024).toFixed(2)} MB`
  }
  return `${(s / 1024 / 1024 / 1024).toFixed(2)} GB`
}

export interface Props {
  image: RuntimeImage
}

export const RuntimeImageRow: Component<Props> = (props) => {
  return (
    <List.Row>
      <List.RowContent>
        <List.RowTitle>Image Size</List.RowTitle>
        <List.RowData>{humanReadableSize(props.image.size)}</List.RowData>
      </List.RowContent>
    </List.Row>
  )
}
