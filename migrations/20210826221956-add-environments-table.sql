-- +migrate Up
CREATE TABLE `environments`
(
    `id`        varchar(22)                              NOT NULL COMMENT '環境変数ID',
    `branch_id` varchar(22) COLLATE 'utf8mb4_general_ci' NOT NULL COMMENT 'ブランチID',
    `key`       varchar(100)                             NOT NULL COMMENT '環境変数のキー',
    `value`     text                                     NOT NULL COMMENT '環境変数の値',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_environments_branch_id` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`)
) COMMENT ='環境変数テーブル' ENGINE = 'InnoDB'
                      COLLATE 'utf8mb4_general_ci';

-- +migrate Down
DROP TABLE `environments`;
