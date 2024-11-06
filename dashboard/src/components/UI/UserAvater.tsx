import { type Component, type JSX, splitProps } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'

export interface UserAvatarProps extends JSX.HTMLAttributes<HTMLImageElement> {
  user: User
  size?: number
}

const UserAvatar: Component<UserAvatarProps> = (props) => {
  const [addedProps, originalImgProps] = splitProps(props, ['user', 'size'])
  return (
    <img
      class="aspect-square h-auto rounded-full"
      src={addedProps.user.avatarUrl}
      style={{
        width: addedProps.size ? `${addedProps.size}px` : '100%',
      }}
      width={addedProps.size}
      height={addedProps.size}
      loading="lazy"
      {...originalImgProps}
      alt={addedProps.user.name}
    />
  )
}

export default UserAvatar
