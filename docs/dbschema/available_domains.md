# available_domains

## Description

利用可能ドメインテーブル

<details>
<summary><strong>Table Definition</strong></summary>

```sql
CREATE TABLE `available_domains` (
  `id` varchar(22) NOT NULL COMMENT 'テンプレートID',
  `domain` varchar(100) NOT NULL COMMENT 'ドメイン',
  `subdomain` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'サブドメインが利用可能か',
  PRIMARY KEY (`id`),
  UNIQUE KEY `domain` (`domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='利用可能ドメインテーブル'
```

</details>

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | varchar(22) |  | false |  |  | テンプレートID |
| domain | varchar(100) |  | false |  |  | ドメイン |
| subdomain | tinyint(1) | 0 | false |  |  | サブドメインが利用可能か |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| domain | UNIQUE | UNIQUE KEY domain (domain) |
| PRIMARY | PRIMARY KEY | PRIMARY KEY (id) |

## Indexes

| Name | Definition |
| ---- | ---------- |
| PRIMARY | PRIMARY KEY (id) USING BTREE |
| domain | UNIQUE KEY domain (domain) USING BTREE |

## Relations

![er](available_domains.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
