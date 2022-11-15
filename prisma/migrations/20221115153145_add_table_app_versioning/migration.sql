-- CreateTable
CREATE TABLE `app_versioning` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `version` VARCHAR(10) NOT NULL,
    `is_force` BOOLEAN NOT NULL DEFAULT false,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
