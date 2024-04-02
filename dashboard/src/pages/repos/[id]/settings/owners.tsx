import { Show } from 'solid-js'
import toast from 'solid-toast'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '/@/components/layouts/DataTable'
import OwnerList from '/@/components/templates/OwnerList'
import { client, handleAPIError } from '/@/libs/api'
import { userFromId, users } from '/@/libs/useAllUsers'
import { useRepositoryData } from '/@/routes'

export default () => {
  const { repo, refetchRepo, hasPermission } = useRepositoryData()
  const loaded = () => !!(repo() && users())

  const handleAddOwner = async (user: User) => {
    const newOwnerIds = repo()?.ownerIds.concat(user.id)
    try {
      await client.updateRepository({
        id: repo()?.id,
        ownerIds: { ownerIds: newOwnerIds },
      })
      toast.success('リポジトリオーナーを追加しました')
      refetchRepo()
    } catch (e) {
      handleAPIError(e, 'リポジトリオーナーの追加に失敗しました')
    }
  }

  const handleDeleteOwner = async (user: User) => {
    const newOwnerIds = repo()?.ownerIds.filter((id) => id !== user.id)
    try {
      await client.updateRepository({
        id: repo()?.id,
        ownerIds: { ownerIds: newOwnerIds },
      })
      toast.success('リポジトリオーナーを削除しました')
      refetchRepo()
    } catch (e) {
      handleAPIError(e, 'リポジトリオーナーの削除に失敗しました')
    }
  }

  return (
    <DataTable.Container>
      <DataTable.Title>Owners</DataTable.Title>
      <DataTable.SubTitle>オーナーはリポジトリ設定の変更が可能になります</DataTable.SubTitle>
      <Show when={loaded()}>
        <OwnerList
          owners={repo()!.ownerIds.map(userFromId)}
          users={users()!}
          handleAddOwner={handleAddOwner}
          handleDeleteOwner={handleDeleteOwner}
          hasPermission={hasPermission()}
        />
      </Show>
    </DataTable.Container>
  )
}
