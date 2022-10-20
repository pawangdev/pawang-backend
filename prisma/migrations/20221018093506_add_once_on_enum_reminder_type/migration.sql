/*
  Warnings:

  - You are about to alter the column `date` on the `task_reminders` table. The data in that column could be lost. The data in that column will be cast from `DateTime(0)` to `DateTime`.

*/
-- AlterTable
ALTER TABLE `task_reminders` MODIFY `type` ENUM('once', 'daily', 'weekly', 'monthly', 'yearly') NOT NULL,
    MODIFY `date` DATETIME NOT NULL,
    MODIFY `time` TIME NOT NULL;
