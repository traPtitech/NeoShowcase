import type { JSX } from '@solidjs/web'
import { type Component, omit } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'

export interface UserAvatarProps extends JSX.HTMLAttributes<HTMLImageElement> {
  user: User
  size?: number
}

const UserAvatar: Component<UserAvatarProps> = (props) => {
  const imgProps = omit(props, 'user', 'size')
  return (
    <img
      class="aspect-square h-auto rounded-full"
      src={props.user.avatarUrl}
      style={{
        width: props.size ? `${props.size}px` : '100%',
      }}
      width={props.size}
      height={props.size}
      loading="lazy"
      {...imgProps}
      alt={props.user.name}
    />
  )
}

export default UserAvatar
