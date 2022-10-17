/*
  Warnings:

  - A unique constraint covering the columns `[email]` on the table `user_reset_passwords` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX `user_reset_passwords_email_key` ON `user_reset_passwords`(`email`);
