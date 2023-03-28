-- +migrate Up
CREATE TABLE `available_domains`
(
    `domain` VARCHAR(100) NOT NULL COMMENT 'ドメイン',
    PRIMARY KEY (`domain`)
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

CREATE TABLE `application_state`
(
    `state` VARCHAR(10) NOT NULL COMMENT 'デプロイの状態',
    PRIMARY KEY (`state`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'デプロイの状態';

INSERT INTO `application_state` (`state`)
VALUES ('IDLE'),
       ('DEPLOYING'),
       ('RUNNING'),
       ('ERRORED');

CREATE TABLE `applications`
(
    `id`             VARCHAR(22)                NOT NULL COMMENT 'アプリケーションID',
    `name`           VARCHAR(100)               NOT NULL COMMENT 'アプリケーション名',
    `repository_id`  VARCHAR(22)                NOT NULL COMMENT 'リポジトリID',
    `branch_name`    VARCHAR(100)               NOT NULL COMMENT 'Gitブランチ・タグ名',
    `build_type`     ENUM ('runtime', 'static') NOT NULL COMMENT 'ビルドタイプ',
    `state`          VARCHAR(10)                NOT NULL COMMENT 'デプロイの状態',
    `current_commit` CHAR(40)                   NOT NULL COMMENT 'デプロイされたコミット',
    `want_commit`    CHAR(40)                   NOT NULL COMMENT 'デプロイを待つコミット',
    `created_at`     DATETIME(6)                NOT NULL COMMENT '作成日時',
    `updated_at`     DATETIME(6)                NOT NULL COMMENT '更新日時',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_applications_repository_id`
        FOREIGN KEY (`repository_id`) REFERENCES `repositories` (`id`),
    CONSTRAINT `fk_applications_state`
        FOREIGN KEY (`state`) REFERENCES `application_state` (`state`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT 'アプリケーションテーブル';

CREATE TABLE `application_config`
(
    `application_id`  VARCHAR(22)                  NOT NULL COMMENT 'アプリケーションID',
    `use_mariadb`     TINYINT(1)                   NOT NULL COMMENT 'MariaDBを使用するか',
    `use_mongodb`     TINYINT(1)                   NOT NULL COMMENT 'MongoDBを使用するか',
    `base_image`      VARCHAR(1000)                NOT NULL COMMENT 'ベースイメージの名前',
    `dockerfile_name` VARCHAR(100)                 NOT NULL COMMENT 'Dockerfile名',
    `artifact_path`   VARCHAR(100)                 NOT NULL COMMENT '静的成果物のパス',
    `build_cmd`       TEXT                         NOT NULL COMMENT 'ビルドコマンド',
    `entrypoint_cmd`  TEXT                         NOT NULL COMMENT 'コンテナのエントリポイント',
    `authentication`  ENUM ('off', 'soft', 'hard') NOT NULL COMMENT 'traP部員認証タイプ',
    PRIMARY KEY (`application_id`),
    CONSTRAINT `fk_application_config_application_id`
        FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE InnoDB
  DEFAULT CHARACTER SET = `utf8mb4`
    COMMENT '';

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
    `commit`         CHAR(40)    NOT NULL COMMENT 'コミットハッシュ',
    `status`         VARCHAR(10) NOT NULL COMMENT 'ビルドの状態',
    `started_at`     DATETIME(6) NULL COMMENT 'ビルド開始日時',
    `updated_at`     DATETIME(6) NULL COMMENT 'ビルド更新日時',
    `finished_at`    DATETIME(6) NULL COMMENT 'ビルド終了日時',
    `retriable`      TINYINT(1)  NOT NULL COMMENT '再ビルド可能フラグ',
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
    `fqdn`           VARCHAR(100)   NOT NULL COMMENT 'サイトURLのFQDN',
    `path_prefix`    VARCHAR(100)   NOT NULL COMMENT 'サイトPathのPrefix',
    `strip_prefix`   TINYINT(1)     NOT NULL COMMENT 'PathのPrefixを落とすかどうか',
    `https`          TINYINT(1)     NOT NULL COMMENT 'httpsの接続かどうか',
    `http_port`      INT DEFAULT 80 NOT NULL COMMENT 'コンテナhttpポート番号',
    `created_at`     DATETIME(6)    NOT NULL COMMENT '作成日時',
    `updated_at`     DATETIME(6)    NOT NULL COMMENT '更新日時',
    `application_id` VARCHAR(22)    NOT NULL COMMENT 'アプリケーションID',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`fqdn`, `path_prefix`),
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

# TODO: delete once user setup / auth is implemented
INSERT INTO `users`
VALUES ('tmp-user', 'toki');

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
DROP TABLE `application_config`;
DROP TABLE `applications`;
DROP TABLE `application_state`;
DROP TABLE `repositories`;
DROP TABLE `available_domains`;
