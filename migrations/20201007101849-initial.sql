-- +migrate Up
CREATE TABLE `repositories`
(
    `id`     VARCHAR(22) NOT NULL COMMENT 'リポジトリID',
    `remote` TEXT        NOT NULL COMMENT 'Git Remote URL',
    `refs`   TEXT        NOT NULL COMMENT '使用するGit Ref',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = 'Gitリポジトリテーブル';

CREATE TABLE `applications`
(
    `id`            VARCHAR(22)  NOT NULL COMMENT 'アプリID',
    `owner`         VARCHAR(100) NOT NULL COMMENT 'アプリ所有者',
    `name`          VARCHAR(100) NOT NULL COMMENT 'アプリ名',
    `repository_id` VARCHAR(22)  NOT NULL COMMENT 'アプリのリポジトリID',
    `created_at`    DATETIME(6)  NOT NULL COMMENT '作成日時',
    `updated_at`    DATETIME(6)  NOT NULL COMMENT '更新日時',
    `deleted_at`    DATETIME(6) COMMENT '削除日時',
    PRIMARY KEY (`id`),
    CONSTRAINT fk_applications_repository_id FOREIGN KEY (`repository_id`) REFERENCES repositories (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = 'アプリテーブル';

CREATE TABLE `sites`
(
    `id`             VARCHAR(22)                NOT NULL COMMENT 'サイトID',
    `fqdn`           VARCHAR(50)                NOT NULL COMMENT 'サイトURLのFQDN',
    `path_prefix`    VARCHAR(50)                NOT NULL COMMENT 'サイトURLのパスプレフィックス',
    `type`           ENUM ('static', 'dynamic') NOT NULL COMMENT 'サイト種類',
    `application_id` VARCHAR(22)                NOT NULL COMMENT 'アプリケーションID',
    `created_at`     DATETIME(6)                NOT NULL COMMENT '作成日時',
    `updated_at`     DATETIME(6)                NOT NULL COMMENT '更新日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`fqdn`, `path_prefix`),
    CONSTRAINT fk_sites_application_id FOREIGN KEY (`application_id`) REFERENCES applications (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = 'サイトテーブル';

CREATE TABLE `build_logs`
(
    `id`             VARCHAR(22)                                          NOT NULL COMMENT 'ビルドログID',
    `application_id` VARCHAR(22)                                          NOT NULL COMMENT 'アプリケーションID',
    `result`         ENUM ('BUILDING', 'SUCCEEDED', 'FAILED', 'CANCELED') NOT NULL COMMENT 'ビルド結果',
    `started_at`     DATETIME(6)                                          NOT NULL COMMENT 'ビルド開始日時',
    `finished_at`    DATETIME(6) COMMENT 'ビルド終了日時',
    PRIMARY KEY (`id`),
    CONSTRAINT fk_build_logs_application_id FOREIGN KEY (`application_id`) REFERENCES applications (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = 'ビルドログテーブル';

CREATE TABLE `artifacts`
(
    `id`           VARCHAR(22) NOT NULL COMMENT '生成物ID',
    `build_log_id` VARCHAR(22) NOT NULL COMMENT 'ビルドログID',
    `size`         BIGINT      NOT NULL COMMENT '生成物ファイルサイズ(tar)',
    `created_at`   DATETIME(6) NOT NULL COMMENT '作成日時',
    `deleted_at`   DATETIME(6) COMMENT '削除日時',
    PRIMARY KEY (`id`),
    CONSTRAINT fk_artifacts_buildlog_id FOREIGN KEY (`build_log_id`) REFERENCES build_logs (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = '静的ファイル生成物テーブル';

CREATE TABLE `static_site_details`
(
    `site_id`     VARCHAR(22) NOT NULL COMMENT 'サイトID',
    `artifact_id` VARCHAR(22) COMMENT '配信する静的ファイル生成物のID',
    PRIMARY KEY (`site_id`),
    CONSTRAINT fk_static_site_details_site_id FOREIGN KEY (`site_id`) REFERENCES sites (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_static_site_details_artifact_id FOREIGN KEY (`artifact_id`) REFERENCES artifacts (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = '静的サイト詳細テーブル';

CREATE TABLE `dynamic_site_details`
(
    `site_id` VARCHAR(22) NOT NULL COMMENT 'サイトID',
    PRIMARY KEY (`site_id`),
    CONSTRAINT fk_dynamic_site_details_site_id FOREIGN KEY (`site_id`) REFERENCES sites (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = '動的サイト詳細テーブル';

-- +migrate Down
DROP TABLE dynamic_site_details;
DROP TABLE static_site_details;
DROP TABLE artifacts;
DROP TABLE build_logs;
DROP TABLE sites;
DROP TABLE applications;
DROP TABLE repositories;
