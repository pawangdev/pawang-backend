-- AlterTable
ALTER TABLE `users` MODIFY `phone` VARCHAR(15) NULL,
    MODIFY `gender` ENUM('male', 'female') NULL;
