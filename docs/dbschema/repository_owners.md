# repository_owners

## Description

リポジトリ所有者テーブル

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE `repository_owners` (
  `user_id` char(22) NOT NULL COMMENT 'ユーザーID',
  `repository_id` char(22) NOT NULL COMMENT 'リポジトリID',
  PRIMARY KEY (`user_id`,`repository_id`),
  KEY `fk_repository_owners_repository_id` (`repository_id`),
  CONSTRAINT `fk_repository_owners_repository_id` FOREIGN KEY (`repository_id`) REFERENCES `repositories` (`id`),
  CONSTRAINT `fk_repository_owners_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='リポジトリ所有者テーブル'
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| user_id | char(22) |  | false |  | [users](users.md) | ユーザーID |
| repository_id | char(22) |  | false |  | [repositories](repositories.md) | リポジトリID |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| fk_repository_owners_repository_id | FOREIGN KEY | FOREIGN KEY (repository_id) REFERENCES repositories (id) |
| fk_repository_owners_user_id | FOREIGN KEY | FOREIGN KEY (user_id) REFERENCES users (id) |
| PRIMARY | PRIMARY KEY | PRIMARY KEY (user_id, repository_id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| fk_repository_owners_repository_id | KEY fk_repository_owners_repository_id (repository_id) USING BTREE |
| PRIMARY | PRIMARY KEY (user_id, repository_id) USING BTREE |

## Relations

```mermaid
erDiagram

"repository_owners" }o--|| "users" : "FOREIGN KEY (user_id) REFERENCES users (id)"
"repository_owners" }o--|| "repositories" : "FOREIGN KEY (repository_id) REFERENCES repositories (id)"

"repository_owners" {
  char_22_ user_id PK
  char_22_ repository_id PK
}
"users" {
  char_22_ id PK
  varchar_255_ name
  tinyint_1_ admin
}
"repositories" {
  char_22_ id PK
  varchar_256_ name
  varchar_256_ url
}
```

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
