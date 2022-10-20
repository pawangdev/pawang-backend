/*
  Warnings:

  - You are about to drop the column `time` on the `task_reminders` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE `task_reminders` DROP COLUMN `time`,
    MODIFY `date` DATETIME(3) NOT NULL;
