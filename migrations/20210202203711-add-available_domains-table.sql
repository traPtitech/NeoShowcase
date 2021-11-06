-- +migrate Up
CREATE TABLE `available_domains`
(
    `id`        VARCHAR(22)  NOT NULL COMMENT 'テンプレートID',
    `domain`    VARCHAR(100) NOT NULL COMMENT 'ドメイン',
    `subdomain` BOOLEAN      NOT NULL DEFAULT false COMMENT 'サブドメインが利用可能か',
    PRIMARY KEY (`id`),
    UNIQUE (`domain`)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = '利用可能ドメインテーブル';

INSERT INTO `available_domains`
VALUES ('WKKbSIz9WiKN6EDhfHU2uT', 'local.tokyotech.org', true);

-- +migrate Down
DROP TABLE `available_domains`;
