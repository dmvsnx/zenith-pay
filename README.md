# Zenith Pay 💳

Backend API untuk sistem Point-of-Sale (POS) berbasis Go. Dibangun dengan arsitektur bersih (clean architecture) untuk memisahkan concern antara handler, business logic, dan data access layer.

## Fitur ✨

- **Autentikasi JWT** — Login dengan token JWT (30 menit kedaluwarsa), role-based access (admin/cashier)
- **Manajemen Kategori & Produk** — CRUD kategori dan produk khusus admin, dengan generate SKU otomatis
- **Manajemen Shift** — Cashier wajib membuka shift sebelum melakukan transaksi, hanya satu shift aktif per cashier
- **Transaksi POS** — Mendukung pembayaran tunai (cash), debit, dan QRIS; validasi stok dengan row-level locking; kalkulasi kembalian otomatis
- **Laporan Penjualan** — Laporan harian dan bulanan dengan agregasi total transaksi dan revenue
- **Rate Limiting** — Perlindungan endpoint dengan pembatasan request per menit

## Tech Stack 🛠️

| Teknologi | Keterangan |
|---|---|
| **Go** 1.25.3 | Bahasa pemrograman |
| **Fiber** v2 | HTTP framework (Express-like, cepat) |
| **GORM** v2 | ORM untuk PostgreSQL |
| **PostgreSQL** | Database relasional |
| **JWT** (HS256) | Autentikasi token |
| **godotenv** | Manajemen konfigurasi environment |
| **go-playground/validator** | Validasi input |

## Arsitektur 🏗️

```
cmd/
└── main.go              # Entry point aplikasi

config/
└── config.go            # Load konfigurasi dari environment

internal/
├── app/
│   └── app.go           # Bootstrap server, middleware global, graceful shutdown
├── database/
│   ├── postgres.go       # Koneksi & migrasi database
│   └── seed/
│       └── user_seed.go  # Seed admin user (development)
├── delivery/
│   ├── handlers/         # HTTP handlers (user, category, product, transaction, shift, report)
│   ├── routes/           # Route registrations & grup middleware
├── dto/                  # Request/Response Data Transfer Objects
├── middlewares/          # CORS, JWT, role, rate-limiter, method validation, active shift
├── model/               # GORM models (User, Category, Product, Transaction, TransactionItem, Shift)
├── repository/          # Data access layer (interfaces + implementations)
├── usecase/             # Business logic layer (interfaces + implementations)
└── utils/
    ├── helpers/         # Bcrypt, JWT, response helper, validator
    └── sku.go           # Generator SKU produk
```

## Persyaratan 📋

- Go 1.25+
- PostgreSQL (local atau container)

## Instalasi & Menjalankan 🚀

1. Clone repositori:
   ```bash
   git clone https://github.com/dmvsnx/zenith-pay.git
   cd zenith-pay
   ```

2. Salin `.env.sample` menjadi `.env` dan sesuaikan konfigurasi:
   ```bash
   cp .env.sample .env
   ```

   Variable environment yang dibutuhkan:

   | Variable | Default | Keterangan |
   |---|---|---|
   | `APP_NAME` | `zenith-pay` | Nama aplikasi |
   | `APP_ENV` | `development` | Environment (`development` / `production`) |
   | `APP_PORT` | `3000` | Port server |
   | `DB_HOST` | `localhost` | Host database |
   | `DB_PORT` | `5432` | Port database |
   | `DB_USER` | `postgres` | User database |
   | `DB_PASSWORD` | `root` | Password database |
   | `DB_NAME` | `zenith-pay_db` | Nama database |
   | `JWT_SECRET` | — | Secret key untuk JWT |
   | `ADMIN_USERNAME` | — | Username admin (development seed) |
   | `ADMIN_PASSWORD` | — | Password admin (development seed) |
   | `ADMIN_EMAIL` | — | Email admin (development seed) |
   | `ADMIN_FULL_NAME` | — | Nama lengkap admin (development seed) |

3. Jalankan aplikasi:
   ```bash
   go run cmd/main.go
   ```

   Server akan berjalan di `http://localhost:3000`. Pada mode development, tabel akan otomatis dimigrasi dan admin user akan di-seed jika environment variable `ADMIN_USERNAME` dan `ADMIN_PASSWORD` diisi.

## API Endpoints 📡

Semua endpoint diawali dengan prefix `/zenith-pay`.

### Autentikasi

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `POST` | `/zenith-pay/auth/login` | Rate limit (5/menit) | Login cashier/admin |

### Users (Admin Only)

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `POST` | `/zenith-pay/admin/users` | JWT + Admin + Rate limit (20/menit) | Register user baru |

### Kategori

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `GET` | `/zenith-pay/categories` | JWT (100/menit) | List kategori |
| `GET` | `/zenith-pay/categories/:id` | JWT (100/menit) | Detail kategori |
| `POST` | `/zenith-pay/categories/admin` | JWT + Admin (50/menit) | Buat kategori |
| `PUT` | `/zenith-pay/categories/admin/:id` | JWT + Admin (50/menit) | Update kategori |
| `DELETE` | `/zenith-pay/categories/admin/:id` | JWT + Admin (50/menit) | Hapus kategori |

### Produk

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `GET` | `/zenith-pay/products` | JWT (100/menit) | List produk |
| `GET` | `/zenith-pay/products/:id` | JWT (100/menit) | Detail produk |
| `POST` | `/zenith-pay/products/admin` | JWT + Admin (50/menit) | Buat produk |
| `PUT` | `/zenith-pay/products/admin/:id` | JWT + Admin (50/menit) | Update produk |
| `DELETE` | `/zenith-pay/products/admin/:id` | JWT + Admin (50/menit) | Hapus produk |

### Transaksi (Cashier Only)

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `POST` | `/zenith-pay/transactions` | JWT + Cashier + Active Shift (30/menit) | Buat transaksi |
| `GET` | `/zenith-pay/transactions` | JWT + Cashier + Active Shift (30/menit) | List transaksi |
| `GET` | `/zenith-pay/transactions/:id` | JWT + Cashier + Active Shift (30/menit) | Detail transaksi |

### Shift (Cashier Only)

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `POST` | `/zenith-pay/shifts/open` | JWT + Cashier (10/menit) | Buka shift |
| `POST` | `/zenith-pay/shifts/close` | JWT + Cashier (10/menit) | Tutup shift |
| `GET` | `/zenith-pay/shifts/active` | JWT + Cashier (10/menit) | Cek shift aktif |

### Laporan (Admin Only)

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `GET` | `/zenith-pay/admin/reports/daily?date=YYYY-MM-DD` | JWT + Admin (60/menit) | Laporan harian |
| `GET` | `/zenith-pay/admin/reports/monthly?period=YYYY-MM` | JWT + Admin (60/menit) | Laporan bulanan |

## Model Database 📦

### User
- `id` (UUID, PK) — `gen_random_uuid()`
- `username` (string, unique)
- `password` (string, bcrypt hash)
- `full_name` (string)
- `email` (string, unique)
- `role` (enum: `admin` / `cashier`)
- `is_active` (boolean)
- `created_at`, `updated_at`

### Category
- `id` (UUID, PK)
- `name` (string, unique)
- `created_at`, `updated_at`

### Product
- `id` (UUID, PK)
- `category_id` (UUID, FK → categories)
- `sku` (string, unique) — format: `SKU-YYYYMMDD-XXXXX`
- `name` (string)
- `price` (int64, satuan terkecil/sen)
- `stock` (int)
- `created_at`, `updated_at`

### Transaction
- `id` (UUID, PK)
- `user_id` (UUID, FK → users)
- `transaction_date` (timestamp)
- `payment_method` (enum: `cash` / `debit` / `qris`)
- `total_amount` (int64)
- `payment_amount` (int64)
- `change_amount` (int64)
- `created_at`, `updated_at`

### Transaction Item
- `id` (UUID, PK)
- `transaction_id` (UUID, FK → transactions)
- `product_id` (UUID, FK → products)
- `product_name` (string)
- `product_price` (int64)
- `quantity` (int)
- `subtotal` (int64)
- `created_at`, `updated_at`

### Shift
- `id` (UUID, PK)
- `cashier_id` (string)
- `status` (enum: `open` / `closed`)
- `opening_balance` (int64)
- `closing_balance` (int64, nullable)
- `opened_at` (timestamp)
- `closed_at` (timestamp, nullable)

## Middleware 🛡️

| Middleware | Deskripsi |
|---|---|
| **CORS** | Mengizinkan origin `localhost:3000` |
| **Method Validation** | Whitelist HTTP method (GET, POST, PUT, PATCH, DELETE, OPTIONS) |
| **JWT Auth** | Ekstrak & validasi Bearer token, set `userID`, `username`, `role`, `claims` di locals |
| **Role-Based Access** | Batasi akses berdasarkan role (`admin` / `cashier`) |
| **Rate Limiter** | Batasi request per endpoint (key berdasarkan `userID` atau IP + User-Agent) |
| **Active Shift** | Verifikasi cashier memiliki shift aktif sebelum transaksi |

## Keamanan 🔒

- Password di-hash menggunakan bcrypt
- Token JWT HS256 dengan expiry 30 menit
- Rate limiting per endpoint
- Row-level locking (`FOR UPDATE`) pada transaksi untuk mencegah race condition stok
- Validasi input menggunakan `go-playground/validator`
- Secret key JWT melalui environment variable (tidak di-hardcode)
