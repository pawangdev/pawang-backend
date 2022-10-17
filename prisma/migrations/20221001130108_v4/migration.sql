-- AddForeignKey
ALTER TABLE `transactions` ADD CONSTRAINT `transactions_subcategory_id_fkey` FOREIGN KEY (`subcategory_id`) REFERENCES `sub_categories`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;
