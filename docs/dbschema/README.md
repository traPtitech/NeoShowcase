# neoshowcase

## Tables

| Name | Columns | Comment | Type |
| ---- | ------- | ------- | ---- |
| [applications](applications.md) | 12 | アプリケーションテーブル | BASE TABLE |
| [application_config](application_config.md) | 13 | アプリケーション詳細設定テーブル | BASE TABLE |
| [application_owners](application_owners.md) | 2 | アプリケーション所有者テーブル | BASE TABLE |
| [artifacts](artifacts.md) | 6 | 静的ファイル生成物テーブル | BASE TABLE |
| [builds](builds.md) | 10 | ビルドテーブル | BASE TABLE |
| [environments](environments.md) | 4 | 環境変数テーブル | BASE TABLE |
| [port_publications](port_publications.md) | 4 | 公開ポートテーブル | BASE TABLE |
| [repositories](repositories.md) | 3 | Gitリポジトリテーブル | BASE TABLE |
| [repository_auth](repository_auth.md) | 5 | Gitリポジトリ認証情報テーブル | BASE TABLE |
| [repository_commits](repository_commits.md) | 9 | コミットメタ情報テーブル | BASE TABLE |
| [repository_owners](repository_owners.md) | 2 | リポジトリ所有者テーブル | BASE TABLE |
| [runtime_images](runtime_images.md) | 3 | ランタイムイメージテーブル | BASE TABLE |
| [users](users.md) | 3 | ユーザーテーブル | BASE TABLE |
| [user_keys](user_keys.md) | 5 | ユーザーSSHキーテーブル | BASE TABLE |
| [websites](websites.md) | 9 | Webサイトテーブル | BASE TABLE |

## Relations

```mermaid
erDiagram

"applications" }o--|| "repositories" : "FOREIGN KEY (repository_id) REFERENCES repositories (id)"
"application_config" |o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"application_owners" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"application_owners" }o--|| "users" : "FOREIGN KEY (user_id) REFERENCES users (id)"
"artifacts" }o--|| "builds" : "FOREIGN KEY (build_id) REFERENCES builds (id)"
"builds" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"environments" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"port_publications" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"repository_auth" |o--|| "repositories" : "FOREIGN KEY (repository_id) REFERENCES repositories (id)"
"repository_owners" }o--|| "repositories" : "FOREIGN KEY (repository_id) REFERENCES repositories (id)"
"repository_owners" }o--|| "users" : "FOREIGN KEY (user_id) REFERENCES users (id)"
"runtime_images" |o--|| "builds" : "FOREIGN KEY (build_id) REFERENCES builds (id)"
"user_keys" }o--|| "users" : "FOREIGN KEY (user_id) REFERENCES users (id)"
"websites" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"

"applications" {
  char_22_ id PK
  varchar_100_ name
  varchar_22_ repository_id FK
  varchar_100_ ref_name
  char_40_ commit
  enum__runtime___static__ deploy_type
  tinyint_1_ running
  enum__missing___starting___restarting___running___idle___exited___errored___unknown__ container
  text container_message
  char_22_ current_build
  datetime_6_ created_at
  datetime_6_ updated_at
}
"application_config" {
  char_22_ application_id PK
  tinyint_1_ use_mariadb
  tinyint_1_ use_mongodb
  tinyint_1_ auto_shutdown
  enum__runtime-buildpack___runtime-cmd___runtime-dockerfile___static-buildpack___static-cmd___static-dockerfile__ build_type
  varchar_1000_ base_image
  text build_cmd
  varchar_100_ artifact_path
  tinyint_1_ spa
  varchar_100_ dockerfile_name
  varchar_100_ context
  text entrypoint
  text command
}
"application_owners" {
  char_22_ user_id PK
  char_22_ application_id PK
}
"artifacts" {
  char_22_ id PK
  varchar_1000_ name
  bigint_20_ size
  datetime_6_ created_at
  datetime_6_ deleted_at
  varchar_22_ build_id FK
}
"builds" {
  char_22_ id PK
  char_40_ commit
  char_16_ config_hash
  enum__building___succeeded___failed___canceled___queued___skipped__ status
  datetime_6_ queued_at
  datetime_6_ started_at
  datetime_6_ updated_at
  datetime_6_ finished_at
  tinyint_1_ retriable
  char_22_ application_id FK
}
"environments" {
  char_22_ application_id PK
  varchar_100_ key PK
  text value
  tinyint_1_ system
}
"port_publications" {
  char_22_ application_id FK
  int_11_ internet_port PK
  int_11_ application_port
  enum__tcp___udp__ protocol PK
}
"repositories" {
  char_22_ id PK
  varchar_256_ name
  varchar_256_ url
}
"repository_auth" {
  char_22_ repository_id PK
  enum__basic___ssh__ method
  varchar_256_ username
  varchar_256_ password
  text ssh_key
}
"repository_commits" {
  char_40_ hash PK
  varchar_256_ author_name
  varchar_256_ author_email
  datetime_6_ author_date
  varchar_256_ committer_name
  varchar_256_ committer_email
  datetime_6_ committer_date
  text message
  tinyint_1_ error
}
"repository_owners" {
  char_22_ user_id PK
  char_22_ repository_id PK
}
"runtime_images" {
  char_22_ build_id PK
  bigint_20_ size
  datetime_6_ created_at
}
"users" {
  char_22_ id PK
  varchar_255_ name
  tinyint_1_ admin
}
"user_keys" {
  char_22_ id PK
  char_22_ user_id FK
  text public_key
  varchar_255_ name
  datetime_6_ created_at
}
"websites" {
  char_22_ id PK
  varchar_100_ fqdn
  varchar_100_ path_prefix
  tinyint_1_ strip_prefix
  tinyint_1_ https
  tinyint_1_ h2c
  int_11_ http_port
  enum__off___soft___hard__ authentication
  char_22_ application_id FK
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
