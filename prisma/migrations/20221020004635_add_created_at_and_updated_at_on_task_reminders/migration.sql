/*
  Warnings:

  - You are about to alter the column `date` on the `task_reminders` table. The data in that column could be lost. The data in that column will be cast from `DateTime(0)` to `DateTime`.
  - Added the required column `updated_at` to the `task_reminders` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE `task_reminders` ADD COLUMN `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    ADD COLUMN `updated_at` DATETIME(3) NOT NULL,
    MODIFY `date` DATETIME NOT NULL,
    MODIFY `time` TIME NOT NULL;
