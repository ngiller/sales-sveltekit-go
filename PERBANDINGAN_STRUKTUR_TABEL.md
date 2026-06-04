# Perbandingan Struktur Tabel Database

## Sumber: `magnum_sales_svelte_go (1).sql` vs `magnum_sales (1).sql`

> Hanya tabel yang **ADA** di `magnum_sales_svelte_go (1).sql` yang dibandingkan.

---

## Ringkasan

| Item | magnum_sales_svelte_go | magnum_sales |
|------|------------------------|-------------|
| **Collation** | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Jumlah Tabel** | 34 tabel | 52 tabel |
| **View sebagai Tabel** | Ya (CREATE TABLE) | Ya (CREATE VIEW) |

---

## Daftar Tabel

| No | Tabel | Status di magnum_sales |
|----|-------|------------------------|
| 1 | counter_id | ✅ Ada |
| 2 | currency | ✅ Ada |
| 3 | customer | ✅ Ada |
| 4 | customer_category | ✅ Ada |
| 5 | customer_contact | ✅ Ada |
| 6 | followup_notification | ✅ Ada (sebagai VIEW) |
| 7 | get_pivot_cancel_quotation | ✅ Ada (sebagai VIEW) |
| 8 | get_pivot_decline_quotation | ✅ Ada (sebagai VIEW) |
| 9 | get_pivot_followup_quotation | ✅ Ada (sebagai VIEW) |
| 10 | get_pivot_po_quotation | ✅ Ada (sebagai VIEW) |
| 11 | get_quotation_year_to_date | ✅ Ada (sebagai VIEW) |
| 12 | get_sales_project_retail | ✅ Ada (sebagai VIEW) |
| 13 | master_departements | ✅ Ada |
| 14 | master_properties | ✅ Ada |
| 15 | master_table_access | ✅ Ada |
| 16 | **master_user_groups** | ❌ **Tidak Ada** |
| 17 | payment_term | ✅ Ada |
| 18 | project_level | ✅ Ada |
| 19 | project_priority | ✅ Ada |
| 20 | quotation | ✅ Ada |
| 21 | quotation_detail | ✅ Ada |
| 22 | quotation_files | ✅ Ada |
| 23 | quotation_followup | ✅ Ada |
| 24 | quotation_master | ✅ Ada |
| 25 | quotation_progress | ✅ Ada |
| 26 | quotation_status | ✅ Ada |
| 27 | quotation_subdetail | ✅ Ada |
| 28 | setting | ✅ Ada |
| 29 | units | ✅ Ada |
| 30 | users | ✅ Ada |
| 31 | user_groups | ✅ Ada |
| 32 | user_group_policies | ✅ Ada |
| 33 | user_policies | ✅ Ada |
| 34 | vendor | ✅ Ada |

---

## Perbedaan Detail Per Tabel

### 1. `counter_id`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Index** | Sama ✅ | `PRIMARY KEY (property_id, counter_name, ym, type)` |
| **Struktur kolom** | Sama ✅ | 5 kolom identik |

### 2. `currency`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Index** | Sama ✅ | `PRIMARY KEY (id), UNIQUE KEY (name)` |
| **Struktur kolom** | Sama ✅ | 7 kolom identik |

### 3. `customer`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `allow_to_vendor` | `varchar(255)` — collation `utf8mb4_general_ci` | `varchar(255)` — collation `utf8mb4_0900_ai_ci` |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)`, **`UNIQUE KEY name`**, **`KEY category_id`**, **`KEY sales_id`** |
| **Struktur kolom** | Sama ✅ | 21 kolom identik |

### 4. `customer_category`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Index** | Sama ✅ | `PRIMARY KEY (id), UNIQUE KEY (name)` |
| **Struktur kolom** | Sama ✅ | 7 kolom identik |

### 5. `customer_contact`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)`, **`UNIQUE KEY cust_contact_name (customer_id, id)`** |
| **Struktur kolom** | Sama ✅ | 6 kolom identik |

### 6. `followup_notification`

| Aspek | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| **Tipe** | **Table** (CREATE TABLE) | **View** (CREATE VIEW) |
| Collation | `utf8mb4_general_ci` | — (view) |
| **Struktur kolom** | Sama ✅ | 9 kolom identik |

### 7. `get_pivot_cancel_quotation`

| Aspek | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| **Tipe** | **Table** (CREATE TABLE) | **View** (CREATE VIEW) |
| Collation | `utf8mb4_general_ci` | — (view) |
| **Urutan kolom** | `April, August, December, February, January, July, June, March, May, November, October, September` | `January, February, March, April, May, June, July, August, September, October, November, December` |
| **Struktur kolom** | Semua `double DEFAULT NULL` | Semua `double` (tanpa DEFAULT) |

### 8. `get_pivot_decline_quotation`

| Aspek | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| **Tipe** | **Table** (CREATE TABLE) | **View** (CREATE VIEW) |
| Collation | `utf8mb4_general_ci` | — (view) |
| **Urutan kolom** | `April, August, December, February, January, July, June, March, May, November, October, September` | `January, February, March, April, May, June, July, August, September, October, November, December` |
| **Struktur kolom** | Semua `double DEFAULT NULL` | Semua `double` (tanpa DEFAULT) |

### 9. `get_pivot_followup_quotation`

| Aspek | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| **Tipe** | **Table** (CREATE TABLE) | **View** (CREATE VIEW) |
| Collation | `utf8mb4_general_ci` | — (view) |
| **Urutan kolom** | `April, August, December, Febuary, January, July, June, March, May, November, October, September` | `January, Febuary, March, April, May, June, July, August, September, October, November, December` |
| **Struktur kolom** | Semua `double DEFAULT NULL` | Semua `double` (tanpa DEFAULT) |
| **Catatan** | Keduanya memiliki typo yang sama: `Febuary` (seharusnya `February`) |

### 10. `get_pivot_po_quotation`

| Aspek | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| **Tipe** | **Table** (CREATE TABLE) | **View** (CREATE VIEW) |
| Collation | `utf8mb4_general_ci` | — (view) |
| **Urutan kolom** | `April, August, December, Febuary, January, July, June, March, May, November, October, September` | `January, Febuary, March, April, May, June, July, August, September, October, November, December` |
| **Struktur kolom** | Semua `double DEFAULT NULL` | Semua `double` (tanpa DEFAULT) |
| **Catatan** | Keduanya memiliki typo yang sama: `Febuary` |

### 11. `get_quotation_year_to_date`

| Aspek | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| **Tipe** | **Table** (CREATE TABLE) | **View** (CREATE VIEW) |
| Collation | `utf8mb4_general_ci` | — (view) |
| **Struktur kolom** | Sama ✅ | `name varchar(255), total double` |

### 12. `get_sales_project_retail`

| Aspek | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| **Tipe** | **Table** (CREATE TABLE) | **View** (CREATE VIEW) |
| Collation | `utf8mb4_general_ci` | — (view) |
| **Struktur kolom** | Sama ✅ | `quotation_type int, total double` |

### 13. `master_departements`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `created_at` | `datetime(3) DEFAULT NULL` | `timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP` |
| `updated_at` | `datetime(3) DEFAULT NULL` | `timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP` |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)`, **`UNIQUE KEY name_property (name, property_id)`** |

### 14. `master_properties`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `code` | `varchar(50) **NOT NULL**` | `varchar(255) **DEFAULT NULL**` |
| `enable` | `tinyint(1) **DEFAULT '1'**` | `tinyint(1) **NOT NULL**` (tanpa default) |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)`, **`UNIQUE KEY name`**, **`UNIQUE KEY code`** |

### 15. `master_table_access`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `name` | ✅ Ada | ✅ Ada |
| `parent_id` | `bigint UNSIGNED DEFAULT NULL` | ❌ **Tidak Ada** |
| `menu_name` | `varchar(255) NOT NULL` | ❌ **Tidak Ada** |
| `path` | `varchar(255) DEFAULT NULL` | ❌ **Tidak Ada** |
| `endpoint` | `varchar(255) NOT NULL` | ❌ **Tidak Ada** |
| `icon` | `varchar(255) DEFAULT NULL` | ❌ **Tidak Ada** |
| `sort_order` | `bigint DEFAULT '0'` | ❌ **Tidak Ada** |
| `is_active` | `tinyint(1) DEFAULT '1'` | ❌ **Tidak Ada** |
| **Total kolom** | **9 kolom** | **2 kolom** |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)` |

> ⚠️ **Perbedaan signifikan**: `master_table_access` di `svelte_go` memiliki 7 kolom tambahan untuk mendukung menu navigasi (`parent_id`, `menu_name`, `path`, `endpoint`, `icon`, `sort_order`, `is_active`).

### 16. `master_user_groups` ⚠️

| Aspek | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| **Status** | ✅ Ada | ❌ **Tidak Ada** |
| Kolom | `id`, `code`, `name` | — |

> ⚠️ **Tabel ini hanya ada di `magnum_sales_svelte_go`**, tidak ditemukan di `magnum_sales`.

### 17. `payment_term`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `day` | `**bigint** DEFAULT NULL` | `**int** DEFAULT NULL` |
| `property_id` | `**bigint UNSIGNED** DEFAULT NULL` | `**int** DEFAULT NULL` |
| `user_created` | `**bigint UNSIGNED** DEFAULT NULL` | `**int** DEFAULT NULL` |
| `user_update` | `**bigint UNSIGNED** DEFAULT NULL` | `**int** DEFAULT NULL` |
| `created_at` | `datetime(3) **DEFAULT NULL**` | `datetime **NOT NULL DEFAULT CURRENT_TIMESTAMP**` |
| `updated_at` | `datetime(3) **DEFAULT NULL**` | `datetime **NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP**` |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)`, **`UNIQUE KEY (property_id, name)`** |

### 18. `project_level`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_unicode_ci` |
| **Struktur kolom** | Sama ✅ | `id int, name varchar(255)` |

### 19. `project_priority`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_unicode_ci` |
| **Struktur kolom** | Sama ✅ | `id int, name varchar(255)` |

### 20. `quotation`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `id` | **`varchar(20)`** NOT NULL | **`varchar(15)`** NOT NULL |
| `margin` | ❌ **Tidak Ada** | `double DEFAULT NULL` ✅ |
| `sales_id` | `**int** DEFAULT NULL` (posisi setelah `user_created`) | `**bigint UNSIGNED** DEFAULT NULL` (posisi di akhir) |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (property_id, id)`, **`KEY (id)`**, **`KEY fromdate_todate (quotation_date)`** |

### 21. `quotation_detail`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `id` | **`varchar(20)`** NOT NULL | **`varchar(15)`** NOT NULL |
| **Index** | `PRIMARY KEY (id, rev_id, line)` | `PRIMARY KEY (rev_id, id, line)` (urutan berbeda) |

### 22. `quotation_files`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Struktur kolom** | Sama ✅ | `id varchar(15), file_name varchar(255), link varchar(255)` |
| **Index** | Tidak ada index tercantum | **`UNIQUE KEY (id, file_name)`** |

### 23. `quotation_followup`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `id` | `varchar(20) **DEFAULT NULL**` | `varchar(15) **NOT NULL**` |
| **Index** | Tidak ada index tercantum | `PRIMARY KEY (property_id, id, line_id)`, **`KEY (id)`** |

### 24. `quotation_master`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `id` | **`varchar(20)`** NOT NULL | **`varchar(15)`** NOT NULL |
| `margin` | ❌ **Tidak Ada** | `double DEFAULT NULL` ✅ |
| `sales_id` | `**int** DEFAULT NULL` | `**bigint UNSIGNED** DEFAULT NULL` |
| **Index** | Sama ✅ | `PRIMARY KEY (id, rev_id)` |

### 25. `quotation_progress`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Struktur kolom** | Sama ✅ | `id int, name varchar(255), progress double` |

### 26. `quotation_status`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Struktur kolom** | Sama ✅ | `id int, name varchar(255)` |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)`, **`KEY (name)`** |

### 27. `quotation_subdetail`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `id` | **`varchar(20)`** NOT NULL | **`varchar(15)`** NOT NULL |
| **Index** | `PRIMARY KEY (id, rev_id, line, subline)` | `PRIMARY KEY (rev_id, id, line, subline)` (urutan berbeda) |

### 28. `setting`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Struktur kolom** | Sama ✅ | `id int, property_id int, code int, name varchar(255), value text` |
| **Index** | Tidak ada index tercantum | `PRIMARY KEY (id)`, **`UNIQUE KEY setting_code (property_id, code)`** |

### 29. `units`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Struktur kolom** | Sama ✅ | 7 kolom identik |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)`, **`UNIQUE KEY (property_id, name)`** |

### 30. `users`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_unicode_ci` |
| `enable` | `tinyint(1) **DEFAULT '1'**` | `tinyint(1) **NOT NULL DEFAULT '1'**` |
| `created_at` | `**datetime(3) DEFAULT NULL**` | `**timestamp NULL DEFAULT CURRENT_TIMESTAMP**` |
| `updated_at` | `**datetime(3) DEFAULT NULL**` | `**timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP**` |
| `user_created` | `**bigint** NOT NULL` | `**int** NOT NULL` |
| `user_update` | `**bigint** NOT NULL` | `**int** NOT NULL` |
| `role_name` | `longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci` | ❌ **Tidak Ada** |
| `dept_name` | `longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci` | ❌ **Tidak Ada** |
| **Index** | `PRIMARY KEY (id)`, `UNIQUE KEY email`, `UNIQUE KEY idx_users_email` | `PRIMARY KEY (id)`, `UNIQUE KEY users_email_unique` |

> ⚠️ **Perbedaan**: `svelte_go` memiliki 2 kolom tambahan (`role_name`, `dept_name`) yang tidak ada di `magnum_sales`.

### 31. `user_groups`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `created_at` | `**datetime(3) DEFAULT NULL**` | `**timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP**` |
| `updated_at` | `**datetime(3) DEFAULT NULL**` | `**timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP**` |
| **Index** | Sama ✅ | `PRIMARY KEY (id)` |

> Catatan: Struktur kolom sama (id, name, property_id, user_created, user_update, created_at, updated_at), hanya tipe data timestamp yang berbeda.

### 32. `user_group_policies`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| `table_id` | `int **NOT NULL**` | ❌ **Tidak Ada** |
| `created_at` | `**datetime(3) DEFAULT NULL**` | `**timestamp NULL DEFAULT CURRENT_TIMESTAMP**` |
| `updated_at` | `**datetime(3) DEFAULT NULL**` | `**timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP**` |
| **Index** | `PRIMARY KEY (id)`, `KEY fk_table_access (table_id)` | `PRIMARY KEY (id)`, **`UNIQUE KEY tblname (property_id, group_id, table_name, action)`**, **`KEY group_id (group_id)`** |
| **Foreign Key** | `fk_table_access` → `master_table_access(id)` ON DELETE CASCADE ON UPDATE CASCADE | ❌ Tidak ada |

### 33. `user_policies`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Struktur kolom** | Sama ✅ | 7 kolom identik |
| **Index** | Sama ✅ | `PRIMARY KEY (id)` |

### 34. `vendor`

| Kolom | magnum_sales_svelte_go | magnum_sales |
|-------|------------------------|-------------|
| Collation | `utf8mb4_general_ci` | `utf8mb4_0900_ai_ci` |
| **Struktur kolom** | Sama ✅ | 13 kolom identik |
| **Index** | `PRIMARY KEY (id)` | `PRIMARY KEY (id)`, **`UNIQUE KEY name`** |

---

## Ringkasan Perbedaan Umum

### 1. Collation / Charset
- **`magnum_sales_svelte_go`**: `utf8mb4_general_ci` (lebih lama, kurang akurat dalam sorting)
- **`magnum_sales`**: `utf8mb4_0900_ai_ci` (lebih baru, lebih akurat)

### 2. Tipe Data Timestamp
- **`magnum_sales_svelte_go`**: cenderung menggunakan `datetime(3) DEFAULT NULL` untuk created_at/updated_at
- **`magnum_sales`**: cenderung menggunakan `timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP`

### 3. Perbedaan Tipe Data Integer
- **`magnum_sales_svelte_go`**: sering menggunakan `bigint UNSIGNED`, `bigint`
- **`magnum_sales`**: sering menggunakan `int` untuk foreign keys

### 4. Kolom `margin`
- **`magnum_sales`** memiliki kolom `margin` di tabel `quotation` dan `quotation_master`
- **`magnum_sales_svelte_go`** tidak memiliki kolom `margin` di kedua tabel tersebut

### 5. View vs Table
- **`magnum_sales_svelte_go`**: Membuat `CREATE TABLE` untuk tabel pivot (followup_notification, get_pivot_*, dll.)
- **`magnum_sales`**: Membuat `CREATE VIEW` untuk tabel yang sama (dengan struktur view yang sebenarnya)

### 6. Perbedaan Index
- **`magnum_sales`** secara umum memiliki lebih banyak index (UNIQUE KEY, FOREIGN KEY) dibanding `magnum_sales_svelte_go`

### 7. Tabel Khusus `magnum_sales_svelte_go`
| Tabel | Keterangan |
|-------|------------|
| `master_user_groups` | Tidak ada di `magnum_sales` |

### 8. Tabel yang Hanya Ada di `magnum_sales` (tidak dibandingkan)
| Tabel |
|-------|
| `customer_asli` |
| `migrations` |
| `password_resets` |
| `po_detail` |
| `po_master` |
| `products` |
| `product_bundle_items` |
| `product_category` |
| `product_unit_items` |
| `purchase` |
| `purchase_item` |
| `quotation_pivot` |
| `quotation_target` |
| `reminder` |
| `reminder_participant` |
| `reminder_type` |
| `setting_bak1` |
| `shipping` |
| `store` |

---

*Dokumen ini dibuat pada: 29 Mei 2026*
