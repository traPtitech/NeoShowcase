-- +migrate Up
DROP TABLE `static_site_details`;
DROP TABLE `dynamic_site_details`;
DROP TABLE `sites`;

CREATE TABLE `websites`
(
    `id`             VARCHAR(22) NOT NULL COMMENT 'サイトID',
    `fqdn`           VARCHAR(50) NOT NULL UNIQUE COMMENT 'サイトURLのFQDN',
    `application_id` VARCHAR(22) NOT NULL UNIQUE COMMENT 'アプリケーションID',
    `build_id`       VARCHAR(22) COMMENT '稼働中のサイトのビルドID',
    `http_port`      INT         NOT NULL DEFAULT 80 COMMENT 'HTTPポート番号',
    `created_at`     DATETIME(6) NOT NULL COMMENT '作成日時',
    `updated_at`     DATETIME(6) NOT NULL COMMENT '更新日時',
    PRIMARY KEY (`id`),
    CONSTRAINT fk_websites_application_id FOREIGN KEY (`application_id`) REFERENCES applications (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT fk_websites_build_id FOREIGN KEY (`build_id`) REFERENCES build_logs (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = 'Webサイトテーブル';

ALTER TABLE `artifacts`
    ADD UNIQUE `uk_build_log_id` (`build_log_id`);

-- +migrate Down
ALTER TABLE `artifacts`
    DROP FOREIGN KEY `fk_artifacts_buildlog_id`;
ALTER TABLE `artifacts`
    DROP INDEX `uk_build_log_id`;
ALTER TABLE `artifacts`
    ADD CONSTRAINT fk_artifacts_buildlog_id FOREIGN KEY (`build_log_id`) REFERENCES build_logs (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;

DROP TABLE `websites`;

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

CREATE TABLE `dynamic_site_details`
(
    `site_id` VARCHAR(22) NOT NULL COMMENT 'サイトID',
    PRIMARY KEY (`site_id`),
    CONSTRAINT fk_dynamic_site_details_site_id FOREIGN KEY (`site_id`) REFERENCES sites (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = '動的サイト詳細テーブル';

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
