import { styled } from '@macaron-css/solid'
import { JSX, Component, splitProps } from 'solid-js'
import { User } from '/@/api/neoshowcase/protobuf/gateway_pb'

const UserAvatarImg = styled('img', {
  base: {
    height: 'auto',
    aspectRatio: '1',
    borderRadius: '50%',
  },
})

export interface UserAvatarProps extends JSX.HTMLAttributes<HTMLImageElement> {
  user: User
  size?: number
}

const UserAvatar: Component<UserAvatarProps> = (props) => {
  const [addedProps, originalImgProps] = splitProps(props, ['user', 'size'])
  return (
    <UserAvatarImg
      src={addedProps.user.avatarUrl}
      style={{
        width: addedProps.size ? `${addedProps.size}px` : '100%',
      }}
      {...originalImgProps}
    />
  )
}

export default UserAvatar
