CREATE TABLE `users`
(
    `id`    CHAR(22)     NOT NULL COMMENT 'ユーザーID',
    `name`  VARCHAR(255) NOT NULL COMMENT 'ユーザー名',
    `admin` TINYINT(1)   NOT NULL COMMENT 'Admin Flag',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='ユーザーテーブル';

CREATE TABLE `user_keys`
(
    `id`         CHAR(22) NOT NULL COMMENT 'ID',
    `user_id`    CHAR(22) NOT NULL COMMENT 'ユーザーID',
    `public_key` TEXT     NOT NULL COMMENT 'SSH Public Key',
    `name`       VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'キー名',
    `created_at` DATETIME(6) NOT NULL DEFAULT '1970-01-01 00:00:00.000000' COMMENT '作成日時',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_user_keys_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT = 'ユーザーSSHキーテーブル';

CREATE TABLE `repositories`
(
    `id`   CHAR(22)     NOT NULL COMMENT 'リポジトリID',
    `name` VARCHAR(256) NOT NULL COMMENT 'リポジトリ名',
    `url`  VARCHAR(256) NOT NULL COMMENT 'Git Remote URL',
    PRIMARY KEY (`id`),
    UNIQUE KEY `url` (`url`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='Gitリポジトリテーブル';

CREATE TABLE `repository_auth`
(
    `repository_id` CHAR(22)              NOT NULL COMMENT 'リポジトリID',
    `method`        ENUM ('basic', 'ssh') NOT NULL COMMENT '認証方法',
    `username`      VARCHAR(256)          NOT NULL COMMENT '(basic)ユーザー名',
    `password`      VARCHAR(256)          NOT NULL COMMENT '(basic)パスワード',
    `ssh_key`       TEXT                  NOT NULL COMMENT '(ssh)PEM encoded private key',
    PRIMARY KEY (`repository_id`),
    CONSTRAINT `fk_repository_auth_repository_id` FOREIGN KEY (`repository_id`) REFERENCES `repositories` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='Gitリポジトリ認証情報テーブル';

CREATE TABLE `repository_owners`
(
    `user_id`       CHAR(22) NOT NULL COMMENT 'ユーザーID',
    `repository_id` CHAR(22) NOT NULL COMMENT 'リポジトリID',
    PRIMARY KEY (`user_id`, `repository_id`),
    KEY `fk_repository_owners_repository_id` (`repository_id`),
    CONSTRAINT `fk_repository_owners_repository_id` FOREIGN KEY (`repository_id`) REFERENCES `repositories` (`id`),
    CONSTRAINT `fk_repository_owners_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='リポジトリ所有者テーブル';

CREATE TABLE `repository_commits`
(
    `hash` CHAR(40) PRIMARY KEY COMMENT 'Commit SHA-1 Hash',
    `author_name` VARCHAR(256) NOT NULL COMMENT 'Author Name',
    `author_email` VARCHAR(256) NOT NULL COMMENT 'Author Email',
    `author_date` DATETIME(6) NOT NULL COMMENT 'Author Date',
    `committer_name` VARCHAR(256) NOT NULL COMMENT 'Committer Name',
    `committer_email` VARCHAR(256) NOT NULL COMMENT 'Committer Email',
    `committer_date` DATETIME(6) NOT NULL COMMENT 'Commit Date',
    `message` TEXT NOT NULL COMMENT 'Commit Message',
    `error` TINYINT(1) NOT NULL COMMENT 'メタ情報取得に失敗したか'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT = 'コミットメタ情報テーブル';

CREATE TABLE `applications`
(
    `id`                CHAR(22)                   NOT NULL COMMENT 'アプリケーションID',
    `name`              VARCHAR(100)               NOT NULL COMMENT 'アプリケーション名',
    `repository_id`     VARCHAR(22)                NOT NULL COMMENT 'リポジトリID',
    `ref_name`          VARCHAR(100)               NOT NULL COMMENT 'Gitブランチ・タグ名',
    `commit`            CHAR(40)                   NOT NULL COMMENT '解決されたコミット',
    `deploy_type`       ENUM ('runtime', 'static') NOT NULL COMMENT 'デプロイタイプ',
    `running`           TINYINT(1)                 NOT NULL COMMENT 'アプリを起動させるか(desired state)',
    `container`         ENUM (
        'missing',
        'starting',
        'restarting',
        'running',
        'exited',
        'errored',
        'unknown'
        )                                          NOT NULL COMMENT 'コンテナの状態(runtime only)',
    `container_message` TEXT                       NOT NULL COMMENT 'コンテナの状態の詳細な情報(runtime only)',
    `current_build`     CHAR(22)                   NOT NULL COMMENT 'デプロイするビルド',
    `created_at`        DATETIME(6)                NOT NULL COMMENT '作成日時',
    `updated_at`        DATETIME(6)                NOT NULL COMMENT '更新日時',
    PRIMARY KEY (`id`),
    KEY `fk_applications_repository_id` (`repository_id`),
    CONSTRAINT `fk_applications_repository_id` FOREIGN KEY (`repository_id`) REFERENCES `repositories` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='アプリケーションテーブル';

CREATE TABLE `application_config`
(
    `application_id`  CHAR(22)      NOT NULL COMMENT 'アプリケーションID',
    `use_mariadb`     TINYINT(1)    NOT NULL COMMENT 'MariaDBを使用するか',
    `use_mongodb`     TINYINT(1)    NOT NULL COMMENT 'MongoDBを使用するか',
    `build_type`      ENUM (
        'runtime-buildpack',
        'runtime-cmd',
        'runtime-dockerfile',
        'static-buildpack',
        'static-cmd',
        'static-dockerfile'
        )                           NOT NULL COMMENT 'ビルドタイプ',
    `base_image`      VARCHAR(1000) NOT NULL COMMENT 'ベースイメージの名前',
    `build_cmd`       TEXT          NOT NULL COMMENT 'ビルドコマンド',
    `artifact_path`   VARCHAR(100)  NOT NULL COMMENT '静的成果物のパス',
    `spa`             TINYINT(1)    NOT NULL COMMENT '静的成果物をSPAとして配信するか',
    `dockerfile_name` VARCHAR(100)  NOT NULL COMMENT 'Dockerfile名',
    `context`         VARCHAR(100)  NOT NULL COMMENT 'ビルド時のcontext',
    `entrypoint`      TEXT          NOT NULL COMMENT 'Entrypoint(args)',
    `command`         TEXT          NOT NULL COMMENT 'Command(args)',
    PRIMARY KEY (`application_id`),
    CONSTRAINT `fk_application_config_application_id` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT = 'アプリケーション詳細設定テーブル';

CREATE TABLE `websites`
(
    `id`             CHAR(22)                   NOT NULL COMMENT 'サイトID',
    `fqdn`           VARCHAR(100)               NOT NULL COMMENT 'サイトURLのFQDN',
    `path_prefix`    VARCHAR(100)               NOT NULL COMMENT 'サイトPathのPrefix',
    `strip_prefix`   TINYINT(1)                 NOT NULL COMMENT 'PathのPrefixを落とすかどうか',
    `https`          TINYINT(1)                 NOT NULL COMMENT 'httpsの接続かどうか',
    `h2c`            TINYINT(1)                 NOT NULL COMMENT '(advanced)プロキシとアプリの通信にh2c protocolを使うかどうか',
    `http_port`      INT(11)                    NOT NULL DEFAULT 80 COMMENT '(runtime only)コンテナhttpポート番号',
    `authentication` ENUM ('off','soft','hard') NOT NULL COMMENT 'traP部員認証タイプ',
    `application_id` CHAR(22)                   NOT NULL COMMENT 'アプリケーションID',
    PRIMARY KEY (`id`),
    KEY `fk_websites_application_id` (`application_id`),
    CONSTRAINT `fk_websites_application_id` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='Webサイトテーブル';

CREATE TABLE `port_publications`
(
    `application_id`   CHAR(22) NOT NULL COMMENT 'アプリケーションID',
    `internet_port`    INT(11)  NOT NULL COMMENT '公開側ポート',
    `application_port` INT(11)  NOT NULL COMMENT 'アプリケーション側ポート',
    `protocol`         ENUM ('tcp', 'udp') COMMENT 'プロトコル',
    PRIMARY KEY (`internet_port`, `protocol`),
    CONSTRAINT `fk_port_publications_application_id` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT = '公開ポートテーブル';

CREATE TABLE `application_owners`
(
    `user_id`        CHAR(22) NOT NULL COMMENT 'ユーザーID',
    `application_id` CHAR(22) NOT NULL COMMENT 'アプリケーションID',
    PRIMARY KEY (`user_id`, `application_id`),
    KEY `fk_application_owners_application_id` (`application_id`),
    CONSTRAINT `fk_application_owners_application_id` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`),
    CONSTRAINT `fk_application_owners_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='アプリケーション所有者テーブル';

CREATE TABLE `builds`
(
    `id`             CHAR(22)    NOT NULL COMMENT 'ビルドID',
    `commit`         CHAR(40)    NOT NULL COMMENT 'コミットハッシュ',
    `config_hash`    CHAR(16)    NOT NULL COMMENT 'ビルド設定のハッシュ',
    `status`         ENUM (
        'building',
        'succeeded',
        'failed',
        'canceled',
        'queued',
        'skipped'
        )                        NOT NULL COMMENT 'ビルドの状態',
    `queued_at`      DATETIME(6) NOT NULL COMMENT 'ビルド追加日時',
    `started_at`     DATETIME(6) NULL COMMENT 'ビルド開始日時',
    `updated_at`     DATETIME(6) NULL COMMENT 'ビルド更新日時',
    `finished_at`    DATETIME(6) NULL COMMENT 'ビルド終了日時',
    `retriable`      TINYINT(1)  NOT NULL COMMENT '再ビルド可能フラグ',
    `application_id` CHAR(22)    NOT NULL COMMENT 'アプリケーションID',
    PRIMARY KEY (`id`),
    KEY `fk_builds_status` (`status`),
    KEY `fk_builds_application_id` (`application_id`),
    CONSTRAINT `fk_builds_application_id` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='ビルドテーブル';

CREATE TABLE `artifacts`
(
    `id`         CHAR(22)      NOT NULL COMMENT '成果物ID',
    `name`       VARCHAR(1000) NOT NULL COMMENT '成果物名',
    `size`       BIGINT(20)    NOT NULL COMMENT '成果物ファイルサイズ',
    `created_at` DATETIME(6)   NOT NULL COMMENT '作成日時',
    `deleted_at` DATETIME(6) DEFAULT NULL COMMENT '削除日時',
    `build_id`   VARCHAR(22)   NOT NULL COMMENT 'ビルドID',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_artifacts_build_id` FOREIGN KEY (`build_id`) REFERENCES `builds` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='静的ファイル生成物テーブル';

CREATE TABLE `runtime_images`
(
  `id`          CHAR(22)     NOT NULL COMMENT 'イメージID',
  `size`        BIGINT       NOT NULL COMMENT 'イメージサイズ',
  `created_at`  DATETIME(6)  NOT NULL COMMENT '作成日時',
  `build_id`    CHAR(22)     NOT NULL COMMENT 'ビルドID',
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_runtime_images_build_id` FOREIGN KEY (`build_id`) REFERENCES `builds` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='ランタイムイメージテーブル';

CREATE TABLE `environments`
(
    `application_id` CHAR(22)     NOT NULL COMMENT 'アプリケーションID',
    `key`            VARCHAR(100) NOT NULL COMMENT '環境変数のキー',
    `value`          TEXT         NOT NULL COMMENT '環境変数の値',
    `system`         TINYINT(1)   NOT NULL COMMENT 'システムによって設定された環境変数かどうか',
    PRIMARY KEY (`application_id`, `key`),
    CONSTRAINT `fk_environments_application_id` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='環境変数テーブル';
