-- +migrate Up
ALTER TABLE `applications`
    MODIFY COLUMN build_type ENUM ('image', 'static') NOT NULL COMMENT 'ビルドタイプ';

-- +migrate Down
ALTER TABLE `applications`
    MODIFY COLUMN build_type ENUM ('image', 'static') COMMENT 'ビルドタイプ';
