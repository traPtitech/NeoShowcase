import type { Component } from 'solid-js'
import type { RuntimeImage } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { formatBytes } from '/@/libs/format'
import { List } from '../List'

export interface Props {
  image: RuntimeImage
}

export const RuntimeImageRow: Component<Props> = (props) => {
  return (
    <List.Row>
      <List.RowContent>
        <List.RowTitle>Size</List.RowTitle>
        <List.RowData>{formatBytes(Number(props.image.size))}</List.RowData>
      </List.RowContent>
    </List.Row>
  )
}
