-- ================================================================
-- SQL ALTER TABLE untuk menyelaraskan struktur magnum_sales
-- dengan magnum_sales_svelte_go tanpa kehilangan data
-- ================================================================
-- Database sumber: magnum_sales_svelte_go (1).sql
-- Database target: magnum_sales (1).sql
-- ================================================================

-- ================================================================
-- BAGIAN 1: MODIFY/ADD COLUMN pada tabel yang sudah ada
-- ================================================================

-- ALTER TABLE `counter_id`
--   MODIFY COLUMN `counter_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
--   MODIFY COLUMN `ym` varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

-- ALTER TABLE `currency`
--   MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

-- ALTER TABLE `customer`
--   MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
--   MODIFY COLUMN `address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
--   MODIFY COLUMN `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `npwp` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `contact` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
--   MODIFY COLUMN `allow_to_vendor` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '1';

-- ALTER TABLE `customer_category`
--   MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

-- ALTER TABLE `customer_contact`
--   MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
--   MODIFY COLUMN `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `position` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL;

-- ALTER TABLE `master_departements`
--   MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
--   MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
--   MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL;

-- ALTER TABLE `master_properties`
--   MODIFY COLUMN `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
--   MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
--   MODIFY COLUMN `address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
--   MODIFY COLUMN `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `contact` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `enable` tinyint(1) DEFAULT '1';

-- NOTE: master_table_access TIDAK diubah (dibiarkan utuh untuk aplikasi lama).
-- Sebagai gantinya, dibuat tabel baru menu_access (copy struct + data) untuk aplikasi baru.

CREATE TABLE IF NOT EXISTS `menu_access` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `parent_id` bigint UNSIGNED DEFAULT NULL,
  `menu_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `endpoint` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `sort_order` bigint DEFAULT '0',
  `is_active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `menu_access` (`id`, `name`, `parent_id`, `menu_name`, `path`, `endpoint`, `icon`, `sort_order`, `is_active`) VALUES
(1, 'Departments', 22, 'Departments', '/departments', 'departements', 'folder', 6, 1),
(2, 'Users', 22, 'Users', '/users', 'users', 'layout', 20, 1),
(3, 'Roles', 22, 'Roles', '/roles', 'roles', 'layout', 10, 1),
(7, 'Customers', 23, 'Customers', '/customers', 'customers', 'layout', 104, 1),
(8, 'Quotations', NULL, 'Quotations', '/quotations', 'quotations', 'layout', 10, 1),
(10, 'Reports', NULL, 'Reports', '', '', 'scroll-text', 40, 1),
(18, 'Payment Terms', 22, 'Payment Terms', '/payment-terms', 'project-terms', 'creadi-card', 130, 1),
(20, 'Customer Category', 23, 'Customer Category', '/customer-categories', 'customer-categories', 'users', 103, 1),
(22, 'Seetings', NULL, 'Seetings', NULL, '', 'settings', 500, 1),
(23, 'Master Data', NULL, 'Master Data', NULL, '', 'folder', 100, 1),
(24, 'Project Levels', 22, 'Project Levels', '/project-levels', '', 'layout', 120, 1),
(25, 'Project Priorities', 22, 'Project Priorities', '/project-priorities', 'project-priorities', 'layout', 140, 1),
(49, 'Quotation Status', 22, 'Quotation Status', '/quotation-statuses', 'quotation-statuses', 'layout', 150, 1),
(50, 'Quotation Progress', 22, 'Quotation Progress', '/quotation-progress', 'quotation-progress', 'layout', 145, 1),
(51, 'Units', 23, 'Units', '/units', 'units', 'layout', 150, 1),
(58, 'Live Stocks', NULL, 'Live Stocks', '/live-stocks', 'live-stocks', 'circle-pile', 30, 1),
(67, 'usergroupspolicies', NULL, 'usergroupspolicies', NULL, 'policies', NULL, 0, 1),
(69, 'usergroups', NULL, 'usergroups', NULL, 'roles', NULL, 0, 1),
(70, 'project priority', NULL, 'project priority', NULL, 'project-priorities', NULL, 0, 1),
(71, 'departements', NULL, 'departements', NULL, 'departements', NULL, 0, 1),
(72, 'payment term', NULL, 'payment term', NULL, 'payment-terms', NULL, 0, 1),
(73, 'project level', NULL, 'project level', NULL, 'project-levels', NULL, 0, 1),
(75, 'Quotations Report', 10, 'Quotations Report', '/reports/quotations-report', 'quotations', 'file-text', 1, 1),
(76, 'Top Customers', 10, 'Top Customers', '/reports/top-customers', 'quotations', 'award', 2, 1),
(77, 'Followups Report', 10, 'Followups Report', '/reports/followups-report', 'quotation-followups', 'clipboard-list', 3, 1),
(78, 'Sales Summary By Customer', 10, 'Sales Summary By Customer', '/reports/sales-summary-by-customer', 'quotations', 'bar-chart-3', 4, 1),
(79, 'Sales Detail By Customer', 10, 'Sales Detail By Customer', '/reports/sales-detail-by-customer', 'quotations', 'list', 5, 1),
(80, 'Sales Detail By Sales Person', 10, 'Sales Detail By Sales Person', '/reports/sales-detail-by-sales-person/', '/reports/sales-detail-by-sales-person/', 'layout', 6, 1),
(81, 'Sales Item By Customer', 10, 'Sales Item By Customer', '/reports/sales-item-by-customer', '/reports/sales-item-by-customer', 'layout', 7, 1),
(82, 'Sales Item By Sales Person', 10, 'Sales Item By Sales Person', '/reports/sales-item-by-sales-person', '/reports/sales-item-by-sales-person', 'layout', 8, 1),
(83, 'Sales Summary By Sales Person', 10, 'Sales Summary By Sales Person', '/reports/sales-summary-by-sales-person', 'quotations', 'bar-chart-3', 5, 1),
(84, 'Sales Charts By Customer', 10, 'Sales Charts By Customer', '/reports/sales-charts-by-customer', 'quotations', 'bar-chart-3', 6, 1),
(85, 'Sales Charts By Sales Person', 10, 'Sales Charts By Sales Person', '/reports/sales-charts-by-sales-person', 'quotations', 'bar-chart-3', 7, 1),
(86, 'Kanban Boards', NULL, 'Kanban Boards', '/kanban', 'kanban', 'kanban', 35, 1);

ALTER TABLE `payment_term`
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `day` bigint DEFAULT NULL,
  MODIFY COLUMN `property_id` bigint UNSIGNED DEFAULT NULL,
  MODIFY COLUMN `user_created` bigint UNSIGNED DEFAULT NULL,
  MODIFY COLUMN `user_update` bigint UNSIGNED DEFAULT NULL,
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL;

ALTER TABLE `project_level`
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

ALTER TABLE `project_priority`
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

ALTER TABLE `quotation`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `quotation_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `subject` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  MODIFY COLUMN `status_review` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  MODIFY COLUMN `folder` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `po_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `po_file` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `po_assign_to` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `sales_id` int DEFAULT NULL;

ALTER TABLE `quotation_detail`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `part_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `descriptions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

ALTER TABLE `quotation_files`
  MODIFY COLUMN `id` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `link` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

-- ALTER TABLE `quotation_followup`
--   MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
--   MODIFY COLUMN `po_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
--   MODIFY COLUMN `po_file` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL;

ALTER TABLE `quotation_master`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `subject` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  MODIFY COLUMN `po_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `po_file` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `po_assign_to` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `sales_id` int DEFAULT NULL;

ALTER TABLE `quotation_progress`
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

ALTER TABLE `quotation_status`
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

ALTER TABLE `quotation_subdetail`
  MODIFY COLUMN `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `part_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `descriptions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

ALTER TABLE `setting`
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

ALTER TABLE `units`
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

ALTER TABLE `user_group_policies`
  ADD COLUMN `table_id` int NOT NULL,
  MODIFY COLUMN `table_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `action` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL;

ALTER TABLE `user_groups`
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL;

ALTER TABLE `user_policies`
  MODIFY COLUMN `table_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `action` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL;

ALTER TABLE `users`
  ADD COLUMN `role_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  ADD COLUMN `dept_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  MODIFY COLUMN `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `phone_no` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  MODIFY COLUMN `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `sign` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `enable` tinyint(1) DEFAULT '1',
  MODIFY COLUMN `inisial` varchar(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  MODIFY COLUMN `created_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `updated_at` datetime(3) DEFAULT NULL,
  MODIFY COLUMN `user_created` bigint NOT NULL,
  MODIFY COLUMN `user_update` bigint NOT NULL;


-- ================================================================
-- BAGIAN 2: CREATE TABLE untuk tabel baru (hanya di svelte_go)
-- ================================================================

CREATE TABLE IF NOT EXISTS `kanban_attachments` (
  `id` int NOT NULL,
  `card_id` int NOT NULL,
  `file_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `file_path` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `file_size` bigint DEFAULT NULL,
  `mime_type` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `uploaded_by` int DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_boards` (
  `id` int NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `description` text COLLATE utf8mb4_general_ci,
  `background` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `is_archived` tinyint(1) DEFAULT '0',
  `property_id` int DEFAULT NULL,
  `user_created` int DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_card_labels` (
  `kanban_card_id` int UNSIGNED NOT NULL,
  `kanban_label_id` int UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_card_members` (
  `kanban_card_id` int UNSIGNED NOT NULL,
  `user_id` int UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_cards` (
  `id` int NOT NULL,
  `list_id` int NOT NULL,
  `board_id` int NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `description` text COLLATE utf8mb4_general_ci,
  `position` int DEFAULT '0',
  `due_date` datetime DEFAULT NULL,
  `start_date` datetime DEFAULT NULL,
  `cover_image` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `is_archived` tinyint(1) DEFAULT '0',
  `user_created` int DEFAULT NULL,
  `user_updated` int DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_checklist_items` (
  `id` int NOT NULL,
  `checklist_id` int NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `is_checked` tinyint(1) DEFAULT '0',
  `position` int DEFAULT '0',
  `assignee_id` int DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_checklists` (
  `id` int NOT NULL,
  `card_id` int NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `position` int DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_comments` (
  `id` int NOT NULL,
  `card_id` int NOT NULL,
  `user_id` int DEFAULT NULL,
  `content` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_labels` (
  `id` int NOT NULL,
  `board_id` int NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_general_ci NOT NULL,
  `color` varchar(50) COLLATE utf8mb4_general_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `kanban_lists` (
  `id` int NOT NULL,
  `board_id` int NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `position` int DEFAULT '0',
  `color` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- ================================================================
-- BAGIAN 3: INSERT DATA untuk tabel kanban* (dari svelte_go)
-- ================================================================

INSERT INTO `kanban_attachments` (`id`, `card_id`, `file_name`, `file_path`, `file_size`, `mime_type`, `uploaded_by`, `created_at`) VALUES
(1, 1, 'CONFIDENTIAL MEDICAL REPORT Pedan Vitalii.docx', 'uploads/kanban/8075019b-4b5a-41a8-b901-4d228ca15c46.docx', 210349, 'application/vnd.openxmlformats-officedocument.wordprocessingml.document', 109, '2026-06-02 22:33:51');

INSERT INTO `kanban_boards` (`id`, `name`, `description`, `background`, `is_archived`, `property_id`, `user_created`, `created_at`, `updated_at`) VALUES
(1, 'sales board', NULL, '#7c3aed', 0, NULL, 109, '2026-06-02 20:57:55', '2026-06-02 20:57:55');

INSERT INTO `kanban_cards` (`id`, `list_id`, `board_id`, `title`, `description`, `position`, `due_date`, `start_date`, `cover_image`, `is_archived`, `user_created`, `user_updated`, `created_at`, `updated_at`) VALUES
(1, 2, 1, 'buat quotation', 'sdsdd sdadasd2345', 1, '2026-06-27 20:58:00', '2026-06-05 20:58:00', NULL, 0, 109, 109, '2026-06-02 20:58:54', '2026-06-02 22:45:23'),
(2, 1, 1, 'Test Card', 'Testing card creation', 0, NULL, NULL, NULL, 0, 109, 109, '2026-06-02 21:19:57', '2026-06-02 22:45:26'),
(3, 4, 1, 'sdsdsdsdsds', 'dsdsdsd', 0, NULL, NULL, NULL, 0, 109, 109, '2026-06-02 21:36:35', '2026-06-02 21:36:35');

INSERT INTO `kanban_card_labels` (`kanban_card_id`, `kanban_label_id`) VALUES
(1, 3),
(1, 4),
(2, 1),
(2, 2);

INSERT INTO `kanban_card_members` (`kanban_card_id`, `user_id`) VALUES
(1, 122),
(2, 109),
(2, 111);

INSERT INTO `kanban_checklists` (`id`, `card_id`, `name`, `position`) VALUES
(2, 1, '123', 1),
(3, 2, 'fdfdfdf', 0);

INSERT INTO `kanban_checklist_items` (`id`, `checklist_id`, `name`, `is_checked`, `position`, `assignee_id`) VALUES
(1, 3, 'bvbvbv', 1, 0, NULL),
(2, 3, 'vbvbvb', 0, 1, NULL);

INSERT INTO `kanban_comments` (`id`, `card_id`, `user_id`, `content`, `created_at`, `updated_at`) VALUES
(1, 2, 109, 'test', '2026-06-02 22:19:08', '2026-06-02 22:19:08'),
(2, 1, 109, 'test', '2026-06-02 22:43:23', '2026-06-02 22:43:23');

INSERT INTO `kanban_labels` (`id`, `board_id`, `name`, `color`) VALUES
(1, 1, 'Bug', '#ef4444'),
(2, 1, 'Feature', '#3b82f6'),
(3, 1, 'Improvement', '#10b981'),
(4, 1, 'Question', '#f59e0b'),
(5, 1, 'Urgent', '#ec4899'),
(6, 1, 'Test', '#6b7280');

INSERT INTO `kanban_lists` (`id`, `board_id`, `name`, `position`, `color`) VALUES
(1, 1, 'To Do', 0, NULL),
(2, 1, 'In Progress', 1, NULL),
(3, 1, 'Done', 2, NULL),
(4, 1, 'test', 3, NULL);

-- ================================================================
-- BAGIAN 3b: CREATE TABLE + DATA untuk group_policies
-- group_policies adalah kopi dari user_group_policies dengan table_id
-- yang sudah dikoreksi (JOIN ke menu_access.id) untuk backend Go.
-- user_group_policies tetap utuh untuk aplikasi lama.
-- ================================================================

CREATE TABLE IF NOT EXISTS `group_policies` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `group_id` bigint NOT NULL,
  `table_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `action` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `property_id` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `table_id` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `group_policies` (`id`, `group_id`, `table_name`, `action`, `property_id`, `created_at`, `updated_at`, `table_id`) VALUES
(1, 1, 'departements', 'read', 1, '2021-08-09 16:01:42.000', '2024-10-26 11:31:05.000', 1),
(2, 1, 'departements', 'create', 1, '2021-08-09 16:01:47.000', '2024-10-26 11:31:05.000', 1),
(3, 1, 'departements', 'update', 1, '2021-08-09 16:01:54.000', '2024-10-26 11:31:05.000', 1),
(4, 1, 'departements', 'delete', 1, '2021-08-09 16:01:59.000', '2024-10-26 11:31:05.000', 1),
(5, 1, 'users', 'read', 1, '2021-08-20 17:22:04.000', '2024-10-26 11:31:05.000', 2),
(6, 1, 'users', 'create', 1, '2021-08-20 17:22:10.000', '2024-10-26 11:31:05.000', 2),
(7, 1, 'users', 'update', 1, '2021-08-20 17:22:18.000', '2024-10-26 11:31:05.000', 2),
(8, 1, 'users', 'delete', 1, '2021-08-20 17:22:25.000', '2024-10-26 11:31:05.000', 2),
(9, 1, 'usergroups', 'read', 1, '2021-09-07 11:38:41.000', '2024-10-26 11:31:05.000', 3),
(10, 1, 'usergroups', 'create', 1, '2021-09-07 11:38:49.000', '2024-10-26 11:31:05.000', 3),
(11, 1, 'usergroups', 'delete', 1, '2021-09-07 11:39:08.000', '2024-10-26 11:31:05.000', 3),
(12, 1, 'property', 'read', 1, '2021-09-07 11:39:37.000', '2024-10-26 11:31:05.000', 68),
(13, 1, 'property', 'create', 1, '2021-09-07 11:39:43.000', '2024-10-26 11:31:05.000', 68),
(14, 1, 'property', 'update', 1, '2021-09-07 11:39:49.000', '2024-10-26 11:31:05.000', 68),
(15, 1, 'property', 'delete', 1, '2021-09-07 11:39:55.000', '2024-10-26 11:31:05.000', 68),
(16, 1, 'usergroups', 'update', 1, '2021-09-07 11:41:07.000', '2024-10-26 11:31:05.000', 3),
(17, 1, 'userpolicies', 'read', 1, '2021-09-07 11:41:30.000', '2024-10-26 11:31:05.000', 0),
(18, 1, 'userpolicies', 'create', 1, '2021-09-07 11:41:38.000', '2024-10-26 11:31:05.000', 0),
(19, 1, 'userpolicies', 'update', 1, '2021-09-07 11:41:45.000', '2024-10-26 11:31:05.000', 0),
(20, 1, 'userpolicies', 'delete', 1, '2021-09-07 11:41:59.000', '2024-10-26 11:31:05.000', 0),
(21, 1, 'usergroupspolicies', 'read', 1, '2021-09-07 11:42:10.000', '2024-10-26 11:31:05.000', 67),
(22, 1, 'usergroupspolicies', 'create', 1, '2021-09-07 11:42:16.000', '2024-10-26 11:31:05.000', 67),
(23, 1, 'usergroupspolicies', 'update', 1, '2021-09-07 11:42:21.000', '2024-10-26 11:31:05.000', 67),
(24, 1, 'usergroupspolicies', 'delete', 1, '2021-09-07 11:42:27.000', '2024-10-26 11:31:05.000', 67),
(25, 2, 'departements', 'read', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 1),
(26, 2, 'departements', 'create', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 1),
(27, 2, 'departements', 'update', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 1),
(28, 2, 'departements', 'delete', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 1),
(29, 2, 'users', 'read', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 2),
(30, 2, 'users', 'create', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 2),
(31, 2, 'users', 'update', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 2),
(32, 2, 'users', 'delete', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 2),
(33, 2, 'usergroups', 'read', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 3),
(34, 2, 'usergroups', 'create', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 3),
(35, 2, 'usergroups', 'delete', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 3),
(36, 2, 'usergroups', 'update', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 3),
(37, 2, 'userpolicies', 'read', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 0),
(38, 2, 'userpolicies', 'create', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 0),
(39, 2, 'userpolicies', 'update', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 0),
(40, 2, 'userpolicies', 'delete', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 0),
(41, 2, 'usergroupspolicies', 'read', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 67),
(42, 2, 'usergroupspolicies', 'create', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 67),
(43, 2, 'usergroupspolicies', 'update', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 67),
(44, 2, 'usergroupspolicies', 'delete', 1, '2021-12-08 11:30:20.000', '2024-10-26 11:31:05.000', 67),
(45, 2, 'quotations', 'read', 1, '2021-12-08 11:30:52.000', '2024-10-26 11:31:05.000', 8),
(46, 2, 'quotations', 'create', 1, '2021-12-08 11:30:59.000', '2024-10-26 11:31:05.000', 8),
(47, 2, 'quotations', 'update', 1, '2021-12-08 11:31:08.000', '2024-10-26 11:31:05.000', 8),
(48, 2, 'quotations', 'delete', 1, '2021-12-08 11:31:19.000', '2024-10-26 11:31:05.000', 8),
(49, 2, 'quotation follow up', 'read', 1, '2021-12-08 11:32:06.000', '2024-10-26 11:31:05.000', 0),
(50, 2, 'quotation follow up', 'create', 1, '2021-12-08 11:32:11.000', '2024-10-26 11:31:05.000', 0),
(51, 2, 'quotation follow up', 'update', 1, '2021-12-08 11:32:17.000', '2024-10-26 11:31:05.000', 0),
(52, 2, 'quotation follow up', 'delete', 1, '2021-12-08 11:32:23.000', '2024-10-26 11:31:05.000', 0),
(53, 2, 'customers', 'read', 1, '2021-12-08 11:32:49.000', '2024-10-26 11:31:05.000', 7),
(54, 2, 'customers', 'create', 1, '2021-12-08 11:32:59.000', '2024-10-26 11:31:05.000', 7),
(55, 2, 'customers', 'update', 1, '2021-12-08 11:33:05.000', '2024-10-26 11:31:05.000', 7),
(56, 2, 'customers', 'delete', 1, '2021-12-08 11:33:10.000', '2024-10-26 11:31:05.000', 7),
(57, 1, 'vendors', 'read', 1, '2022-01-26 10:37:04.000', '2024-10-26 11:31:05.000', 0),
(58, 1, 'vendors', 'create', 1, '2022-01-26 10:37:12.000', '2024-10-26 11:31:05.000', 0),
(59, 1, 'vendors', 'update', 1, '2022-01-26 10:37:26.000', '2024-10-26 11:31:05.000', 0),
(60, 1, 'vendors', 'delete', 1, '2022-01-26 10:37:32.000', '2024-10-26 11:31:05.000', 0),
(61, 1, 'productcategory', 'read', 1, '2022-01-26 15:09:34.000', '2024-10-26 11:31:05.000', 0),
(62, 1, 'productcategory', 'create', 1, '2022-01-26 15:09:40.000', '2024-10-26 11:31:05.000', 0),
(63, 1, 'productcategory', 'update', 1, '2022-01-26 15:09:46.000', '2024-10-26 11:31:05.000', 0),
(64, 1, 'productcategory', 'delete', 1, '2022-01-26 15:09:52.000', '2024-10-26 11:31:05.000', 0),
(65, 1, 'productunits', 'read', 1, '2022-01-27 14:24:02.000', '2024-10-26 11:31:05.000', 0),
(66, 1, 'productunits', 'create', 1, '2022-01-27 14:24:08.000', '2024-10-26 11:31:05.000', 0),
(67, 1, 'productunits', 'update', 1, '2022-01-27 14:24:14.000', '2024-10-26 11:31:05.000', 0),
(68, 1, 'productunits', 'delete', 1, '2022-01-27 14:24:20.000', '2024-10-26 11:31:05.000', 0),
(69, 1, 'storages', 'read', 1, '2022-01-27 16:23:51.000', '2024-10-26 11:31:05.000', 0),
(70, 1, 'storages', 'create', 1, '2022-01-27 16:23:57.000', '2024-10-26 11:31:05.000', 0),
(71, 1, 'storages', 'update', 1, '2022-01-27 16:24:04.000', '2024-10-26 11:31:05.000', 0),
(72, 1, 'storages', 'delete', 1, '2022-01-27 16:24:10.000', '2024-10-26 11:31:05.000', 0),
(73, 1, 'products', 'read', 1, '2022-01-31 15:50:13.000', '2024-10-26 11:31:05.000', 0),
(74, 1, 'products', 'create', 1, '2022-01-31 15:50:33.000', '2024-10-26 11:31:05.000', 0),
(75, 1, 'products', 'update', 1, '2022-01-31 15:50:38.000', '2024-10-26 11:31:05.000', 0),
(76, 1, 'products', 'delete', 1, '2022-01-31 15:50:43.000', '2024-10-26 11:31:05.000', 0),
(77, 1, 'shipping', 'read', 1, '2022-05-09 14:13:33.000', '2024-10-26 11:31:05.000', 0),
(78, 1, 'shipping', 'create', 1, '2022-05-09 14:13:39.000', '2024-10-26 11:31:05.000', 0),
(79, 1, 'shipping', 'update', 1, '2022-05-09 14:13:45.000', '2024-10-26 11:31:05.000', 0),
(80, 1, 'shipping', 'delete', 1, '2022-05-09 14:13:50.000', '2024-10-26 11:31:05.000', 0),
(81, 1, 'purchase', 'read', 1, '2022-05-09 17:17:22.000', '2024-10-26 11:31:05.000', 0),
(82, 1, 'purchase', 'create', 1, '2022-05-09 17:17:29.000', '2024-10-26 11:31:05.000', 0),
(83, 1, 'purchase', 'update', 1, '2022-05-09 17:17:36.000', '2024-10-26 11:31:05.000', 0),
(84, 1, 'purchase', 'delete', 1, '2022-05-09 17:17:42.000', '2024-10-26 11:31:05.000', 0),
(85, 1, 'payment term', 'read', 1, '2022-05-12 12:12:30.000', '2024-10-26 11:31:05.000', 18),
(86, 1, 'payment term', 'create', 1, '2022-05-12 12:12:36.000', '2024-10-26 11:31:05.000', 18),
(87, 1, 'payment term', 'update', 1, '2022-05-12 12:12:42.000', '2024-10-26 11:31:05.000', 18),
(88, 1, 'payment term', 'delete', 1, '2022-05-12 12:12:47.000', '2024-10-26 11:31:05.000', 18),
(89, 72, 'departements', 'read', 1, '2023-05-04 10:21:51.000', '2024-10-26 11:31:05.000', 1),
(90, 72, 'users', 'read', 1, '2023-05-04 10:22:06.000', '2024-10-26 11:31:05.000', 2),
(91, 72, 'usergroups', 'read', 1, '2023-05-04 10:22:14.000', '2024-10-26 11:31:05.000', 3),
(92, 72, 'property', 'read', 1, '2023-05-04 10:22:23.000', '2024-10-26 11:31:05.000', 68),
(93, 72, 'usergroupspolicies', 'read', 1, '2023-05-04 10:22:34.000', '2024-10-26 11:31:05.000', 67),
(94, 72, 'customers', 'read', 1, '2023-05-04 10:22:45.000', '2024-10-26 11:31:05.000', 7),
(95, 72, 'customers', 'create', 1, '2023-05-04 10:22:57.000', '2024-10-26 11:31:05.000', 7),
(96, 72, 'quotations', 'read', 1, '2023-05-04 10:23:09.000', '2024-10-26 11:31:05.000', 8),
(97, 72, 'quotations', 'create', 1, '2023-05-04 10:23:21.000', '2024-10-26 11:31:05.000', 8),
(98, 72, 'quotation follow up', 'read', 1, '2023-05-04 10:23:40.000', '2024-10-26 11:31:05.000', 0),
(99, 72, 'quotation follow up', 'create', 1, '2023-05-04 10:23:50.000', '2024-10-26 11:31:05.000', 0),
(100, 72, 'quotation report', 'read', 1, '2023-05-04 10:23:59.000', '2024-10-26 11:31:05.000', 10),
(101, 1, 'Followup Report', 'read', 1, '2024-11-12 11:36:13.000', '2024-11-12 11:36:13.000', 77),
(102, 1, 'Followup Report', 'create', 1, '2024-11-12 11:36:23.000', '2024-11-12 11:36:23.000', 77),
(103, 1, 'Followup Report', 'update', 1, '2024-11-12 11:36:30.000', '2024-11-12 11:36:30.000', 77),
(104, 1, 'Followup Report', 'delete', 1, '2024-11-12 11:36:36.000', '2024-11-12 11:36:36.000', 77),
(105, 73, 'customers', 'read', 1, '2025-03-03 13:44:14.000', '2025-03-03 13:44:14.000', 7),
(106, 73, 'customers', 'create', 1, '2025-03-03 13:44:22.000', '2025-03-03 13:44:22.000', 7),
(107, 73, 'customers', 'update', 1, '2025-03-03 13:44:33.000', '2025-03-03 13:44:33.000', 7),
(128, 1, 'Reports', 'read', 1, '2026-06-03 18:06:36.000', '2026-06-03 18:06:36.000', 10),
(129, 1, 'Seetings', 'read', 1, '2026-06-03 18:06:36.000', '2026-06-03 18:06:36.000', 22),
(130, 1, 'Master Data', 'read', 1, '2026-06-03 18:06:36.000', '2026-06-03 18:06:36.000', 23);


-- ================================================================
-- BAGIAN 4: INDEX/KEY untuk tabel baru
-- ================================================================

ALTER TABLE `kanban_attachments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_card_id` (`card_id`);

ALTER TABLE `kanban_boards`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `kanban_card_labels`
  ADD PRIMARY KEY (`kanban_card_id`,`kanban_label_id`);

ALTER TABLE `kanban_card_members`
  ADD PRIMARY KEY (`kanban_card_id`,`user_id`);

ALTER TABLE `kanban_cards`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_list_id` (`list_id`),
  ADD KEY `idx_board_id` (`board_id`);

ALTER TABLE `kanban_checklist_items`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_checklist_id` (`checklist_id`);

ALTER TABLE `kanban_checklists`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_card_id` (`card_id`);

ALTER TABLE `kanban_comments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_card_id` (`card_id`);

ALTER TABLE `kanban_labels`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_board_id` (`board_id`);

ALTER TABLE `kanban_lists`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_board_id` (`board_id`);


-- ================================================================
-- BAGIAN 5: INDEX/KEY untuk tabel existing (optional)
-- ================================================================
-- Berikut adalah perbedaan PRIMARY KEY yang perlu disesuaikan:
--
-- Table: quotation
--   Target (svelte_go): ['ALTER TABLE `quotation`\n  ADD PRIMARY KEY (`id`);']
--   Current (magnum):  ['ALTER TABLE `quotation`\n  ADD PRIMARY KEY (`property_id`,`id`),\n  ADD KEY `id` (`id`),\n  ADD KEY `fromdate_todate` (`quotation_date`);']

-- Table: quotation_detail
--   Target (svelte_go): ['ALTER TABLE `quotation_detail`\n  ADD PRIMARY KEY (`id`,`rev_id`,`line`);']
--   Current (magnum):  ['ALTER TABLE `quotation_detail`\n  ADD PRIMARY KEY (`rev_id`,`id`,`line`) USING BTREE;']

-- Table: quotation_subdetail
--   Target (svelte_go): ['ALTER TABLE `quotation_subdetail`\n  ADD PRIMARY KEY (`id`,`rev_id`,`line`,`subline`);']
--   Current (magnum):  ['ALTER TABLE `quotation_subdetail`\n  ADD PRIMARY KEY (`rev_id`,`id`,`line`,`subline`) USING BTREE;']

-- Table: quotation_followup
--   Target (svelte_go): []
--   Current (magnum):  ['ALTER TABLE `quotation_followup`\n  ADD PRIMARY KEY (`property_id`,`id`,`line_id`) USING BTREE,\n  ADD KEY `id` (`id`);']


-- ================================================================
-- BAGIAN 6: Konversi COLLATION ke utf8mb4_general_ci
-- ================================================================

ALTER TABLE `counter_id` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `currency` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `customer` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `customer_category` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `customer_contact` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `master_departements` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `master_properties` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `group_policies` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `menu_access` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
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
ALTER TABLE `user_group_policies` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `user_groups` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `user_policies` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
ALTER TABLE `users` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- ================================================================
-- BAGIAN 7: TABEL YANG ADA DI MAGNUM TAPI TIDAK DI SVELTE_GO
-- ================================================================
-- Tabel berikut ada di magnum_sales tapi TIDAK di svelte_go:
-- customer_asli, followup_notification, migrations, password_resets,
-- po_detail, po_master, product_bundle_items, product_category,
-- product_unit_items, products, purchase, purchase_item,
-- quotation_pivot, quotation_target, reminder, reminder_participant,
-- reminder_type, setting_bak1, shipping, store, vendor
-- Tabel-tabel ini TIDAK dihapus untuk menjaga data.

-- ================================================================
-- BAGIAN 8: PROCEDURE (ada 14 di magnum_sales, 0 di svelte_go)
-- ================================================================
-- Stored Procedure TIDAK dihapus untuk menjaga fungsionalitas.

-- ================================================================
-- NOTES / CATATAN PENTING:
-- ================================================================
-- 1. PERHATIAN: master_properties.code berubah varchar(255) DEFAULT NULL -> varchar(50) NOT NULL
--    Pastikan tidak ada nilai NULL di kolom code sebelum menjalankan.
-- 2. PERHATIAN: quotation.id, quotation_detail.id, quotation_master.id,
--    quotation_subdetail.id, quotation_followup.id berubah varchar(15) -> varchar(20)
--    Perubahan ini aman (widening), tidak akan memotong data.
-- 3. PERHATIAN: payment_term.created_at berubah dari datetime NOT NULL -> datetime(3) NULL
--    Ini mengubah NOT NULL menjadi nullable, data existing tetap aman.
-- 4. Untuk PK changes (BAGIAN 5), jangan dijalankan otomatis.
--    Drop/add PK merebuild tabel dan berisiko pada tabel besar.
-- 5. master_table_access TIDAK diubah — dibuat tabel baru menu_access sebagai copy
--    untuk digunakan oleh aplikasi baru (backend Go). Aplikasi lama tetap pakai master_table_access.
-- 6. user_group_policies TIDAK diubah — dibuat tabel baru group_policies sebagai copy
--    dengan table_id yang sudah dikoreksi (JOIN ke menu_access.id).
-- 7. Urutan eksekusi: BAGIAN 1 -> BAGIAN 2 -> BAGIAN 3 -> BAGIAN 3b -> BAGIAN 4 -> BAGIAN 6
-- ================================================================
