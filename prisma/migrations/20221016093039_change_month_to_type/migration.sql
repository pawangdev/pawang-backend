/*
  Warnings:

  - You are about to drop the column `month` on the `task_reminders` table. All the data in the column will be lost.
  - Added the required column `type` to the `task_reminders` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE `task_reminders` DROP COLUMN `month`,
    ADD COLUMN `type` ENUM('daily', 'weekly', 'monthly', 'yearly') NOT NULL,
    MODIFY `is_active` BOOLEAN NOT NULL DEFAULT true;
