-- +migrate Up
CREATE TABLE `available_domains`
(
    `id`        VARCHAR(22)          NOT NULL COMMENT 'ドメインID',
    `domain`    VARCHAR(100)         NOT NULL COMMENT 'ドメイン',
    `subdomain` TINYINT(1) DEFAULT 0 NOT NULL COMMENT 'サブドメインが利用可能か',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`domain`)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT '利用可能ドメインテーブル';

CREATE TABLE `repositories`
(
    `id`   VARCHAR(22)  NOT NULL COMMENT 'リポジトリID',
    `name` VARCHAR(256) NOT NULL COMMENT 'リポジトリ名',
    `url`  VARCHAR(256) NOT NULL COMMENT 'Git Remote URL',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`url`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'Gitリポジトリテーブル';

CREATE TABLE `applications`
(
    `id`            VARCHAR(22)              NOT NULL COMMENT 'アプリケーションID',
    `repository_id` VARCHAR(22)              NOT NULL COMMENT 'リポジトリID',
    `branch_name`   VARCHAR(100)             NOT NULL COMMENT 'Gitブランチ・タグ名',
    `build_type`    ENUM ('image', 'static') NOT NULL COMMENT 'ビルドタイプ',
    `created_at`    DATETIME(6)              NOT NULL COMMENT '作成日時',
    `updated_at`    DATETIME(6)              NOT NULL COMMENT '更新日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`repository_id`, `branch_name`),
    CONSTRAINT `fk_applications_repository_id`
        FOREIGN KEY (`repository_id`) REFERENCES `repositories` (`id`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'アプリケーションテーブル';

CREATE TABLE `build_status`
(
    `status` VARCHAR(10) NOT NULL COMMENT 'ビルドの状態',
    PRIMARY KEY (`status`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'ビルドの状態';

INSERT INTO `build_status` (`status`)
VALUES ('BUILDING'),
       ('SUCCEEDED'),
       ('FAILED'),
       ('CANCELED'),
       ('QUEUED'),
       ('SKIPPED');

CREATE TABLE `builds`
(
    `id`             VARCHAR(22) NOT NULL COMMENT 'ビルドID',
    `status`         VARCHAR(10) NOT NULL COMMENT 'ビルドの状態',
    `started_at`     DATETIME(6) NOT NULL COMMENT 'ビルド開始日時',
    `finished_at`    DATETIME(6) NULL COMMENT 'ビルド終了日時',
    `application_id` VARCHAR(22) NOT NULL COMMENT 'アプリケーションID',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_builds_status`
        FOREIGN KEY (`status`) REFERENCES `build_status` (`status`),
    CONSTRAINT `fk_builds_application_id`
        FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'ビルドテーブル';

CREATE TABLE `artifacts`
(
    `id`         VARCHAR(22) NOT NULL COMMENT '生成物ID',
    `size`       BIGINT      NOT NULL COMMENT '生成物ファイルサイズ(tar)',
    `created_at` DATETIME(6) NOT NULL COMMENT '作成日時',
    `deleted_at` DATETIME(6) NULL COMMENT '削除日時',
    `build_id`   VARCHAR(22) NOT NULL COMMENT 'ビルドID',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`build_id`),
    CONSTRAINT `fk_artifacts_build_id`
        FOREIGN KEY (`build_id`) REFERENCES `builds` (`id`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT '静的ファイル生成物テーブル';

CREATE TABLE `environments`
(
    `id`             VARCHAR(22)  NOT NULL COMMENT '環境変数ID',
    `application_id` VARCHAR(22)  NOT NULL COMMENT 'アプリケーションID',
    `key`            VARCHAR(100) NOT NULL COMMENT '環境変数のキー',
    `value`          TEXT         NOT NULL COMMENT '環境変数の値',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`application_id`, `key`),
    CONSTRAINT `fk_environments_application_id`
        FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT '環境変数テーブル';

CREATE TABLE `websites`
(
    `id`             VARCHAR(22)    NOT NULL COMMENT 'サイトID',
    `fqdn`           VARCHAR(50)    NOT NULL COMMENT 'サイトURLのFQDN',
    `http_port`      INT DEFAULT 80 NOT NULL COMMENT 'HTTPポート番号',
    `created_at`     DATETIME(6)    NOT NULL COMMENT '作成日時',
    `updated_at`     DATETIME(6)    NOT NULL COMMENT '更新日時',
    `application_id` VARCHAR(22)    NOT NULL COMMENT 'アプリケーションID',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`fqdn`),
    UNIQUE KEY (`application_id`),
    CONSTRAINT `fk_websites_application_id`
        FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'Webサイトテーブル';

CREATE TABLE `users`
(
    `id`   CHAR(36)     NOT NULL COMMENT 'ユーザーID',
    `name` VARCHAR(255) NOT NULL COMMENT 'ユーザー名',
    PRIMARY KEY (`id`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'ユーザーテーブル';

CREATE TABLE `owners`
(
    `user_id`        CHAR(36) NOT NULL COMMENT 'ユーザーID',
    `application_id` CHAR(36) NOT NULL COMMENT 'アプリケーションID',
    PRIMARY KEY (`user_id`, `application_id`),
    CONSTRAINT `fk_owners_user_id`
        FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_owners_application_id`
        FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'アプリケーション所有者テーブル';

-- +migrate Down
DROP TABLE `owners`;
DROP TABLE `users`;
DROP TABLE `websites`;
DROP TABLE `environments`;
DROP TABLE `artifacts`;
DROP TABLE `builds`;
DROP TABLE `build_status`;
DROP TABLE `applications`;
DROP TABLE `repositories`;
DROP TABLE `providers`;
DROP TABLE `available_domains`;