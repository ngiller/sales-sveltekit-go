# Backend Development Guidelines

Dokumen ini berisi aturan dan pola standar untuk pengembangan backend Go (Fiber + GORM).
Kamu adalah programer handal untuk membuat aplikasi backend go dengan standart go yang rapi, bersih, dan mudah dikembangkan mengikuti dokumen ini dengan sangat teliti dan akurat. Kamu juga seorang database designer profesional yang sangat rapi dan teliti dalam membuat database, tidak membuat database yang berlebihan, dan selalu normalisasi database dengan baik. Kamu juga tahu teknologi apa yang digunakan dalam project ini, kamu ahli dalam database dan ERD (Entity Relationship Diagram).

## 1. Arsitektur Repository & Transaction
Semua repositori harus menyediakan akses ke instance DB aslinya melalui metode `GetDB()` untuk mendukung transaksi lintas operasi di level Handler.

### Pola Transaksi di Handler:
Gunakan `h.repo.GetDB().Transaction` untuk membungkus operasi yang membutuhkan integritas data (terutama saat menyimpan data utama beserta relasinya).

```go
err := h.repo.GetDB().Transaction(func(tx *gorm.DB) error {
    // Jalankan operasi menggunakan tx (bukan h.repo)
    if err := tx.Create(&model).Error; err != nil {
        return err
    }
    return nil
})
```

## 2. Standardisasi List (Pagination & Sorting)
Setiap endpoint `FindAll` wajib mendukung:
- **Search**: Pencarian string pada kolom-kolom relevan (Name, Email, dsb).
- **Pagination**: Parameter `page` (default 1) dan `limit` (default 50).
- **Sorting**: Parameter `sort` (nama kolom) dan `order` (`asc` atau `desc`).

### Implementasi di Repository:
Gunakan `.Order()`, `.Offset()`, dan `.Limit()` dari GORM.
```go
func (r *MyRepository) FindAll(search string, page, limit int, sortBy, sortDir string) ([]models.MyModel, int64, error) {
    // ... logic query ...
    orderClause := sortBy + " " + sortDir
    err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&items).Error
    return items, total, err
}
```

## 3. Penanganan Joins
Saat melakukan sorting pada kolom dari tabel lain (misal: `category_name`), pastikan repositori menggunakan `Joins` dan mapping nama kolom yang benar di `sortBy` untuk menghindari ambiguitas SQL.

Contoh mapping di `UserRepository`:
```go
if sortBy == "role_name" {
    sortBy = "user_groups.name"
}
```

## 4. Format Respons
Selalu gunakan `utils.SuccessResponse` dan `utils.ErrorResponse` untuk konsistensi.
- **List**: Harus menyertakan metadata `total`, `page`, dan `limit`.
- **Error**: Gunakan pesan error yang deskriptif namun aman.

## 5. Audit Fields
Pastikan field seperti `UserCreated`, `UserUpdate`, `CreatedAt`, dan `UpdatedAt` terisi dengan benar (biasanya melalui middleware auth yang menyuntikkan `user_id` ke `c.Locals`).

## 6. Security Standards

Setiap pengembang wajib mematuhi standar keamanan berikut untuk menjaga integritas sistem:

### A. Authentication & JWT
- Gunakan `AuthMiddleware()` untuk melindungi route yang membutuhkan login.
- Middleware ini secara otomatis menyuntikkan `user_id` ke dalam `c.Locals("user_id")`.
- Pastikan token dikirim melalui Header `Authorization: Bearer <token>` atau Cookie `token`.

### B. Authorization (RBAC/Policy)
Sistem ini menggunakan Role-Based Access Control yang sangat granular berbasis tabel dan aksi.

#### 1. Komponen Data:
- **`master_table_access`**: Daftar modul/tabel yang diproteksi (menyimpan `name`, `endpoint` untuk pemetaan URL otomatis, `menu_name`, `path`, dll).
- **`user_groups`**: Definisi peran pengguna.
- **`user_group_policies`**: Tabel yang menghubungkan grup ke tabel/modul dengan aksi tertentu (`read`, `write`, `edit`, `delete`).

#### 2. Implementasi Middleware:
Gunakan middleware `RequirePolicy(db, action)` untuk mengecek izin akses.
- **Admin Bypass**: User dengan `group_id = 1` (Super Admin) memiliki akses penuh otomatis.
- **Dynamic Path Mapping**: Middleware secara dinamis mengekstrak segmen URL dari request (misal: `/api/payment-terms` menjadi `payment-terms`) dan mencocokkannya dengan kolom `endpoint` di tabel `master_table_access`. Pendekatan ini mencegah *typo* karena nama tabel tidak lagi di-hardcode.
- **Policy Check**: Middleware memverifikasi izin dengan mencari kecocokan pada tabel `user_group_policies` berdasarkan `group_id`, `action`, serta referensi ke tabel (`table_name` atau `table_id`).

#### 3. Contoh Penggunaan di Route:
```go
// URL "/customers" akan otomatis dipetakan ke endpoint "customers" di master_table_access
// Izin membaca data
api.Get("/customers", middleware.RequirePolicy(db, "read"), h.FindAll)
// Izin menambah data
api.Post("/customers", middleware.RequirePolicy(db, "create"), h.Create)
```

#### 4. Integrasi UI:
Backend menyediakan flag `can_read`, `can_write`, `can_edit`, dan `can_delete` pada objek menu untuk memudahkan frontend dalam menyembunyikan/menampilkan elemen UI sesuai izin user.

### C. Password Security
- **DILARANG** menyimpan password dalam bentuk plain text.
- Gunakan `config.HashPassword(password)` (Bcrypt) sebelum menyimpan ke database.
- Verifikasi menggunakan `config.CheckPasswordHash(password, hash)`.

### D. SQL Injection Prevention
- Gunakan GORM API (`.Where()`, `.Find()`, `.Create()`) secara benar. Hindari penggunaan string formatting (`fmt.Sprintf`) untuk membangun query SQL manual.
- Jika terpaksa menggunakan raw SQL, gunakan *parameterized queries* (misal: `db.Raw("SELECT * FROM users WHERE id = ?", id)`).

### E. CORS Configuration
- Konfigurasi CORS ada di `cmd/server/main.go`. Saat ini dibatasi ke `http://localhost:5173`.
- Jika ada perubahan domain frontend, pastikan untuk memperbarui list origin yang diizinkan.

### F. File Upload Security
- File upload (avatar, signature) disimpan di folder `uploads/`.
- Selalu generate nama file unik menggunakan UUID (`github.com/google/uuid`) untuk mencegah *file name collision* dan serangan *directory traversal*.
- Pastikan validasi tipe file (extension) dilakukan sebelum menyimpan ke disk.

### G. Input Validation
- Selalu validasi body request menggunakan `c.BodyParser` dan lakukan pengecekan field wajib sebelum memproses data ke database/transaksi.
