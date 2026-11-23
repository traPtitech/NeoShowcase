import { createMemo, Show } from 'solid-js'
import type { User } from '/@/api/neoshowcase/protobuf/gateway_pb'
import { DataTable } from '/@/components/layouts/DataTable'
import OwnerList, { ReadOnlyOwnerList } from '/@/components/templates/OwnerList'
import { useApplicationData } from '/@/features/application/data'
import { client, getApplicationGroup, handleAPIError } from '/@/libs/api'
import { users, usersFromIds } from '/@/libs/session'
import { toast } from '/@/libs/toast'

export default () => {
  const { app, refetch, hasPermission } = useApplicationData()

  const updateOwners = async (ownerIds: string[], success: string, failure: string) => {
    try {
      await client.updateApplication({ id: app().id, ownerIds: { ownerIds } })
      toast.success(success)
      refetch()
    } catch (e) {
      handleAPIError(e, failure)
    }
  }

  const handleAddOwner = (user: User) =>
    updateOwners(
      app().ownerIds.concat(user.id),
      'アプリケーションオーナーを追加しました',
      'アプリケーションオーナーの追加に失敗しました',
    )

  const handleDeleteOwner = (user: User) =>
    updateOwners(
      app().ownerIds.filter((id) => id !== user.id),
      'アプリケーションオーナーを削除しました',
      'アプリケーションオーナーの削除に失敗しました',
    )

  const group = createMemo(() => {
    const groupId = app().groupId
    return groupId ? getApplicationGroup(groupId) : undefined
  })

  return (
    <DataTable.Container>
      <DataTable.Title>Owners</DataTable.Title>
      <DataTable.SubTitle>
        オーナーはアプリ設定の変更, アプリログ/メトリクスの閲覧, 環境変数の閲覧, ビルドログの閲覧が可能になります
      </DataTable.SubTitle>
      <Show
        when={group()}
        fallback={
          <OwnerList
            owners={usersFromIds(app().ownerIds)}
            users={users()}
            handleAddOwner={handleAddOwner}
            handleDeleteOwner={handleDeleteOwner}
            hasPermission={hasPermission()}
          />
        }
      >
        {(group) => (
          <ReadOnlyOwnerList
            owners={usersFromIds(group().ownerIds)}
            notice={
              <>
                このアプリはグループ {group().name} に属しています。オーナーは
                <a href={`/groups/${group().id}/settings/owners`} class="text-primary-main">
                  グループのオーナー設定
                </a>
                から編集してください
              </>
            }
          />
        )}
      </Show>
    </DataTable.Container>
  )
}
