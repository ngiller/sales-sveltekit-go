-- =============================================================================
-- Script: alter_to_match_svelte_go.sql
-- Tujuan: Mengubah struktur database magnum_sales agar sesuai dengan 
--         magnum_sales_svelte_go (membandingkan tabel yang ada di svelte_go)
-- 
-- Sumber referensi: magnum_sales_svelte_go (1).sql (target)
-- Database target:  magnum_sales (1).sql (yang akan diubah)
--
-- Peringatan: Script ini akan mengubah struktur tabel yang sudah ada.
--             Backup data terlebih dahulu sebelum menjalankan!
-- =============================================================================

USE `magnum_sales`;

-- =============================================================================
-- BAGIAN 0: Nonaktifkan pengecekan foreign key sementara
-- =============================================================================
SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS;
SET FOREIGN_KEY_CHECKS = 0;


-- =============================================================================
-- BAGIAN 3: Ubah PRIMARY KEY tabel tertentu
-- =============================================================================

-- quotation_subdetail: dari PRIMARY KEY (rev_id, id, line, subline) menjadi PRIMARY KEY (id, rev_id, line, subline)
ALTER TABLE `quotation_subdetail` DROP PRIMARY KEY;
ALTER TABLE `quotation_subdetail` ADD PRIMARY KEY (`id`, `rev_id`, `line`, `subline`);

-- =============================================================================
-- BAGIAN 4: Modifikasi kolom tabel
-- =============================================================================

-- -----------------------------------------------------
-- master_departements
-- -----------------------------------------------------
-- Ubah created_at, updated_at dari timestamp menjadi datetime(3) DEFAULT NULL
ALTER TABLE `master_departements` 
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL;

-- -----------------------------------------------------
-- master_properties
-- -----------------------------------------------------
-- Ubah code dari varchar(255) DEFAULT NULL menjadi varchar(50) NOT NULL
-- Ubah enable dari tinyint(1) NOT NULL menjadi tinyint(1) DEFAULT '1'
ALTER TABLE `master_properties`
  MODIFY COLUMN `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `enable` tinyint(1) DEFAULT '1';

-- -----------------------------------------------------
-- master_table_access: Tambah kolom dari svelte_go
-- -----------------------------------------------------
ALTER TABLE `master_table_access`
  ADD COLUMN `parent_id` bigint UNSIGNED DEFAULT NULL AFTER `name`,
  ADD COLUMN `menu_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL AFTER `parent_id`,
  ADD COLUMN `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL AFTER `menu_name`,
  ADD COLUMN `endpoint` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL AFTER `path`,
  ADD COLUMN `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL AFTER `endpoint`,
  ADD COLUMN `sort_order` bigint DEFAULT '0' AFTER `icon`,
  ADD COLUMN `is_active` tinyint(1) DEFAULT '1' AFTER `sort_order`;

-- -----------------------------------------------------
-- payment_term
-- -----------------------------------------------------
-- Ubah tipe data kolom
ALTER TABLE `payment_term`
  MODIFY COLUMN `day` bigint DEFAULT NULL,
  MODIFY COLUMN `property_id` bigint UNSIGNED DEFAULT NULL,
  MODIFY COLUMN `user_created` bigint UNSIGNED DEFAULT NULL,
  MODIFY COLUMN `user_update` bigint UNSIGNED DEFAULT NULL,
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL;

-- -----------------------------------------------------
-- quotation
-- -----------------------------------------------------
-- Ubah id dari varchar(15) menjadi varchar(20)
ALTER TABLE `quotation`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL FIRST;

-- Hapus kolom margin (hanya ada di magnum_sales)
ALTER TABLE `quotation` DROP COLUMN `margin`;

-- Ubah sales_id dari bigint UNSIGNED menjadi int, pindahkan posisi
ALTER TABLE `quotation`
  DROP COLUMN `sales_id`,
  ADD COLUMN `sales_id` int DEFAULT NULL AFTER `folder`;

-- -----------------------------------------------------
-- quotation_detail
-- -----------------------------------------------------
-- Ubah id dari varchar(15) menjadi varchar(20)
ALTER TABLE `quotation_detail`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL FIRST;

-- -----------------------------------------------------
-- quotation_files
-- -----------------------------------------------------
-- Ubah id dari varchar(15) menjadi varchar(20) [svelte_go pakai varchar(15) jg, tapi dengan charset]
ALTER TABLE `quotation_files`
  MODIFY COLUMN `id` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL FIRST,
  MODIFY COLUMN `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `link` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

-- -----------------------------------------------------
-- quotation_followup
-- -----------------------------------------------------
-- Ubah id dari varchar(15) NOT NULL menjadi varchar(20) DEFAULT NULL
ALTER TABLE `quotation_followup`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL FIRST;

-- -----------------------------------------------------
-- quotation_master
-- -----------------------------------------------------
-- Ubah id dari varchar(15) menjadi varchar(20)
ALTER TABLE `quotation_master`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL FIRST;

-- Hapus kolom margin
ALTER TABLE `quotation_master` DROP COLUMN `margin`;

-- Ubah sales_id dari bigint UNSIGNED menjadi int
ALTER TABLE `quotation_master`
  DROP COLUMN `sales_id`,
  ADD COLUMN `sales_id` int DEFAULT NULL AFTER `notes`;

-- -----------------------------------------------------
-- quotation_subdetail
-- -----------------------------------------------------
-- Ubah id dari varchar(15) menjadi varchar(20)
ALTER TABLE `quotation_subdetail`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL FIRST;

-- -----------------------------------------------------
-- setting: sudah sama strukturnya, hanya perlu change charset
-- -----------------------------------------------------

-- -----------------------------------------------------
-- users: Tambah kolom, ubah tipe data
-- -----------------------------------------------------
-- Ubah created_at, updated_at ke datetime(3)
-- Ubah enable ke DEFAULT '1'
-- Ubah user_created, user_update ke bigint
ALTER TABLE `users`
  MODIFY COLUMN `enable` tinyint(1) DEFAULT '1',
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `user_created` bigint NOT NULL,
  MODIFY COLUMN `user_update` bigint NOT NULL;

-- Tambah kolom role_name dan dept_name
ALTER TABLE `users`
  ADD COLUMN `role_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci AFTER `user_update`,
  ADD COLUMN `dept_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci AFTER `role_name`;

-- -----------------------------------------------------
-- user_groups: ubah created_at, updated_at
-- -----------------------------------------------------
ALTER TABLE `user_groups`
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL;

-- -----------------------------------------------------
-- user_group_policies: Tambah kolom table_id, ubah timestamp
-- -----------------------------------------------------
ALTER TABLE `user_group_policies`
  ADD COLUMN `table_id` int NOT NULL AFTER `property_id`,
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL;

-- =============================================================================
-- BAGIAN 5: Buat tabel baru yang tidak ada di magnum_sales
-- =============================================================================

-- -----------------------------------------------------
-- Tabel: master_user_groups
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `master_user_groups` (
  `id` bigint UNSIGNED NOT NULL,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- =============================================================================
-- BAGIAN 6: Buat TABLE dari VIEW yang sudah di-drop
-- (Di svelte_go berupa tabel, bukan view)
-- =============================================================================

-- -----------------------------------------------------
-- Tabel: followup_notification
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `followup_notification` (
  `customer_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `id` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `last_followup` date DEFAULT NULL,
  `quotation_date` date DEFAULT NULL,
  `quotation_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `subject` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `total` varchar(63) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


ALTER TABLE `users` DROP INDEX `users_email_unique`;
ALTER TABLE `users` ADD UNIQUE KEY `email` (`email`);
ALTER TABLE `users` ADD UNIQUE KEY `idx_users_email` (`email`);

-- user_groups: PRIMARY KEY (id) [sudah]

-- user_group_policies: PRIMARY KEY (id), KEY fk_table_access (table_id)
ALTER TABLE `user_group_policies` ADD KEY `fk_table_access` (`table_id`);

-- =============================================================================
-- BAGIAN 8: Tambah FOREIGN KEY (yang ada di svelte_go)
-- =============================================================================

-- user_group_policies: fk_table_access references master_table_access(id)
ALTER TABLE `user_group_policies`
  ADD CONSTRAINT `fk_table_access` FOREIGN KEY (`table_id`) REFERENCES `master_table_access` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- =============================================================================
-- BAGIAN 9: Ubah COLLATION seluruh tabel ke utf8mb4_general_ci
-- =============================================================================

ALTER TABLE `counter_id` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `currency` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `customer` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `customer_category` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `customer_contact` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `master_departements` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `master_properties` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `master_table_access` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `master_user_groups` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `payment_term` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `project_level` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `project_priority` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `quotation` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `quotation_detail` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `quotation_files` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `quotation_followup` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `quotation_master` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `quotation_progress` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `quotation_status` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `quotation_subdetail` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `setting` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `units` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `users` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `user_groups` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `user_group_policies` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `user_policies` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `vendor` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `followup_notification` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `get_pivot_cancel_quotation` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `get_pivot_decline_quotation` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `get_pivot_followup_quotation` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `get_pivot_po_quotation` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `get_quotation_year_to_date` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `get_sales_project_retail` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- =============================================================================
-- BAGIAN 10: Hapus AUTO_INCREMENT dari tabel yang tidak memilikinya di svelte_go
-- =============================================================================

-- master_properties: svelte_go tidak punya AUTO_INCREMENT
ALTER TABLE `master_properties` MODIFY COLUMN `id` bigint NOT NULL;

-- setting: svelte_go tidak punya AUTO_INCREMENT
ALTER TABLE `setting` MODIFY COLUMN `id` int NOT NULL;

-- user_policies: svelte_go tidak punya AUTO_INCREMENT
ALTER TABLE `user_policies` MODIFY COLUMN `id` bigint NOT NULL;

-- vendor: svelte_go tidak punya AUTO_INCREMENT
ALTER TABLE `vendor` MODIFY COLUMN `id` bigint NOT NULL;

-- quotation_files: svelte_go tidak punya AUTO_INCREMENT
-- (tidak terdaftar di AUTO_INCREMENT kedua file, tidak perlu perubahan)

-- =============================================================================
-- BAGIAN 11: Kembalikan pengaturan foreign key
-- =============================================================================
SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS;

-- =============================================================================
-- SELESAI
-- =============================================================================
