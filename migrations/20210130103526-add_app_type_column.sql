-- +migrate Up
ALTER TABLE `applications`
    ADD COLUMN `build_type` ENUM ('image', 'static') COMMENT 'ビルドタイプ';

-- +migrate Down
ALTER TABLE `applications`
    DROP COLUMN `build_type`;
