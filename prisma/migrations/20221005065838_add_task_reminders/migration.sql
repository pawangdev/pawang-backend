/*
  Warnings:

  - You are about to drop the `user_notifications` table. If the table is not empty, all the data it contains will be lost.
  - A unique constraint covering the columns `[google_id]` on the table `users` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[onesignal_id]` on the table `users` will be added. If there are existing duplicate values, this will fail.

*/
-- DropForeignKey
ALTER TABLE `user_notifications` DROP FOREIGN KEY `fk_users_user_notification`;

-- AlterTable
ALTER TABLE `users` ADD COLUMN `google_id` VARCHAR(255) NULL,
    ADD COLUMN `onesignal_id` VARCHAR(255) NULL;

-- DropTable
DROP TABLE `user_notifications`;

-- CreateTable
CREATE TABLE `task_reminders` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `month` VARCHAR(255) NOT NULL,
    `days` VARCHAR(255) NOT NULL,
    `is_active` BOOLEAN NOT NULL,
    `user_id` INTEGER NOT NULL,

    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateIndex
CREATE UNIQUE INDEX `users_google_id_key` ON `users`(`google_id`);

-- CreateIndex
CREATE UNIQUE INDEX `users_onesignal_id_key` ON `users`(`onesignal_id`);

-- AddForeignKey
ALTER TABLE `task_reminders` ADD CONSTRAINT `task_reminders_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;
