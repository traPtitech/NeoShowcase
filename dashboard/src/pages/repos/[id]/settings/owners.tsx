import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '/@/components/layouts/DataTable'
import OwnerList from '/@/components/templates/OwnerList'
import { useRepositoryData } from '/@/features/repository/data'
import { client, handleAPIError } from '/@/libs/api'
import { users, usersFromIds } from '/@/libs/session'
import { toast } from '/@/libs/toast'

export default () => {
  const { repo, refetch, hasPermission } = useRepositoryData()

  const updateOwners = async (ownerIds: string[], success: string, failure: string) => {
    try {
      await client.updateRepository({ id: repo().id, ownerIds: { ownerIds } })
      toast.success(success)
      refetch()
    } catch (e) {
      handleAPIError(e, failure)
    }
  }

  const handleAddOwner = (user: User) =>
    updateOwners(
      repo().ownerIds.concat(user.id),
      'リポジトリオーナーを追加しました',
      'リポジトリオーナーの追加に失敗しました',
    )

  const handleDeleteOwner = (user: User) =>
    updateOwners(
      repo().ownerIds.filter((id) => id !== user.id),
      'リポジトリオーナーを削除しました',
      'リポジトリオーナーの削除に失敗しました',
    )

  return (
    <DataTable.Container>
      <DataTable.Title>Owners</DataTable.Title>
      <DataTable.SubTitle>オーナーはリポジトリ設定の変更が可能になります</DataTable.SubTitle>
      <OwnerList
        owners={usersFromIds(repo().ownerIds)}
        users={users()}
        handleAddOwner={handleAddOwner}
        handleDeleteOwner={handleDeleteOwner}
        hasPermission={hasPermission()}
      />
    </DataTable.Container>
  )
}
