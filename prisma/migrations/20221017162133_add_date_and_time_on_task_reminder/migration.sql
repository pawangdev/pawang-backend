/*
  Warnings:

  - You are about to drop the column `days` on the `task_reminders` table. All the data in the column will be lost.
  - Added the required column `date` to the `task_reminders` table without a default value. This is not possible if the table is not empty.
  - Added the required column `time` to the `task_reminders` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE `task_reminders` DROP COLUMN `days`,
    ADD COLUMN `date` DATETIME NOT NULL,
    ADD COLUMN `time` TIME NOT NULL;
