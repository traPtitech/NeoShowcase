# applications

## Description

アプリケーションテーブル

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE `applications` (
  `id` char(22) NOT NULL COMMENT 'アプリケーションID',
  `name` varchar(100) NOT NULL COMMENT 'アプリケーション名',
  `repository_id` varchar(22) NOT NULL COMMENT 'リポジトリID',
  `ref_name` varchar(100) NOT NULL COMMENT 'Gitブランチ・タグ名',
  `commit` char(40) NOT NULL COMMENT '解決されたコミット',
  `deploy_type` enum('runtime','static') NOT NULL COMMENT 'デプロイタイプ',
  `running` tinyint(1) NOT NULL COMMENT 'アプリを起動させるか(desired state)',
  `container` enum('missing','starting','restarting','running','idle','exited','errored','unknown') NOT NULL COMMENT 'コンテナの状態(runtime only)',
  `container_message` text NOT NULL COMMENT 'コンテナの状態の詳細な情報(runtime only)',
  `current_build` char(22) NOT NULL COMMENT 'デプロイするビルド',
  `created_at` datetime(6) NOT NULL COMMENT '作成日時',
  `updated_at` datetime(6) NOT NULL COMMENT '更新日時',
  PRIMARY KEY (`id`),
  KEY `fk_applications_repository_id` (`repository_id`),
  CONSTRAINT `fk_applications_repository_id` FOREIGN KEY (`repository_id`) REFERENCES `repositories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='アプリケーションテーブル'
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | char(22) |  | false | [application_config](application_config.md) [application_owners](application_owners.md) [builds](builds.md) [environments](environments.md) [port_publications](port_publications.md) [websites](websites.md) |  | アプリケーションID |
| name | varchar(100) |  | false |  |  | アプリケーション名 |
| repository_id | varchar(22) |  | false |  | [repositories](repositories.md) | リポジトリID |
| ref_name | varchar(100) |  | false |  |  | Gitブランチ・タグ名 |
| commit | char(40) |  | false |  |  | 解決されたコミット |
| deploy_type | enum('runtime','static') |  | false |  |  | デプロイタイプ |
| running | tinyint(1) |  | false |  |  | アプリを起動させるか(desired state) |
| container | enum('missing','starting','restarting','running','idle','exited','errored','unknown') |  | false |  |  | コンテナの状態(runtime only) |
| container_message | text |  | false |  |  | コンテナの状態の詳細な情報(runtime only) |
| current_build | char(22) |  | false |  |  | デプロイするビルド |
| created_at | datetime(6) |  | false |  |  | 作成日時 |
| updated_at | datetime(6) |  | false |  |  | 更新日時 |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| fk_applications_repository_id | FOREIGN KEY | FOREIGN KEY (repository_id) REFERENCES repositories (id) |
| PRIMARY | PRIMARY KEY | PRIMARY KEY (id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| fk_applications_repository_id | KEY fk_applications_repository_id (repository_id) USING BTREE |
| PRIMARY | PRIMARY KEY (id) USING BTREE |

## Relations

```mermaid
erDiagram

"application_config" |o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"application_owners" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"builds" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"environments" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"port_publications" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"websites" }o--|| "applications" : "FOREIGN KEY (application_id) REFERENCES applications (id)"
"applications" }o--|| "repositories" : "FOREIGN KEY (repository_id) REFERENCES repositories (id)"

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
"repositories" {
  char_22_ id PK
  varchar_256_ name
  varchar_256_ url
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
