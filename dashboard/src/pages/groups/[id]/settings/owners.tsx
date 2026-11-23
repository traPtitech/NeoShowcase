import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '/@/components/layouts/DataTable'
import OwnerList from '/@/components/templates/OwnerList'
import { useApplicationGroupData } from '/@/features/applicationGroup/data'
import { client, handleAPIError } from '/@/libs/api'
import { users, usersFromIds } from '/@/libs/session'
import { toast } from '/@/libs/toast'

export default () => {
  const { group, refetch, hasPermission } = useApplicationGroupData()

  const updateOwners = async (ownerIds: string[], success: string, failure: string) => {
    try {
      await client.updateApplicationGroup({ id: group().id, ownerIds: { ownerIds } })
      toast.success(success)
      refetch()
    } catch (e) {
      handleAPIError(e, failure)
    }
  }

  const handleAddOwner = (user: User) =>
    updateOwners(
      group().ownerIds.concat(user.id),
      'グループオーナーを追加しました',
      'グループオーナーの追加に失敗しました',
    )

  const handleDeleteOwner = (user: User) =>
    updateOwners(
      group().ownerIds.filter((id) => id !== user.id),
      'グループオーナーを削除しました',
      'グループオーナーの削除に失敗しました',
    )

  return (
    <DataTable.Container>
      <DataTable.Title>Owners</DataTable.Title>
      <DataTable.SubTitle>
        オーナーはグループ設定の変更と, グループ内のすべてのアプリの操作 (設定の変更,
        ログ/メトリクス/環境変数/ビルドログの閲覧) が可能になります
      </DataTable.SubTitle>
      <OwnerList
        owners={usersFromIds(group().ownerIds)}
        users={users()}
        handleAddOwner={handleAddOwner}
        handleDeleteOwner={handleDeleteOwner}
        hasPermission={hasPermission()}
      />
    </DataTable.Container>
  )
}
