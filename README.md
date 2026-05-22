# Zenith Pay 💳

Backend API untuk sistem Point-of-Sale (POS) berbasis Go. Dibangun dengan arsitektur bersih (clean architecture) untuk memisahkan concern antara handler, business logic, dan data access layer.

## Fitur ✨

- **Autentikasi JWT** — Login dengan token JWT (30 menit kedaluwarsa), role-based access (admin/cashier)
- **Manajemen Kategori & Produk** — CRUD kategori dan produk khusus admin, dengan generate SKU otomatis
- **Manajemen Shift** — Cashier wajib membuka shift sebelum melakukan transaksi, hanya satu shift aktif per cashier
- **Transaksi POS** — Mendukung pembayaran tunai (cash), debit, dan QRIS; validasi stok dengan row-level locking; kalkulasi kembalian otomatis
- **Laporan Penjualan** — Laporan harian dan bulanan dengan agregasi total transaksi dan revenue
- **Rate Limiting** — Perlindungan endpoint dengan pembatasan request per menit
- **Pagination** — List endpoint menggunakan pagination (page & limit query param)
- **Logging Terstruktur** — Zerolog (console development, JSON production)
- **Health Check** — Endpoint liveness (`/health`) dan readiness (`/health/ready`)

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
| **zerolog** | Structured logging (JSON production) |

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
├── dto/                  # Request/Response Data Transfer Objects (termasuk pagination DTO)
├── middlewares/          # CORS, JWT, role, rate-limiter, method validation, active shift
├── model/               # GORM models (User, Category, Product, Transaction, TransactionItem, Shift)
├── repository/          # Data access layer (interfaces + implementations)
├── usecase/             # Business logic layer (interfaces + implementations)
└── utils/
    ├── helpers/         # Bcrypt, JWT, response helper, validator
    ├── logger/          # Zerolog (console dev, JSON prod)
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
    | `CLOUDINARY_URL` | — | Cloudinary connection URL untuk upload image |
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

Semua endpoint diawali dengan prefix `/zenith-pay`. List endpoint mendukung pagination via query param `page` (default: 1) dan `limit` (default: 10).

### Health Check

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/health` | Liveness — uptime server |
| `GET` | `/health/ready` | Readiness — koneksi database |

#### Contoh Request & Response

```json
// GET /health
// Response:
{
  "code": 200,
  "status": "success",
  "message": "server is running",
  "data": {
    "uptime": "2h30m15s"
  }
}

// GET /health/ready
// Response:
{
  "code": 200,
  "status": "success",
  "message": "database is connected",
  "data": {
    "database": "connected"
  }
}
```

### Autentikasi

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `POST` | `/zenith-pay/auth/login` | Rate limit (5/menit) | Login cashier/admin |

#### Contoh Request & Response

```json
// POST /zenith-pay/auth/login
// Request:
{
  "username": "kasir1",
  "password": "secret123"
}

// Response:
{
  "code": 200,
  "status": "success",
  "message": "login berhasil",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### Users (Admin Only)

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `GET` | `/zenith-pay/admin/users` | JWT + Admin + Rate limit (20/menit) | List semua user |
| `POST` | `/zenith-pay/admin/users` | JWT + Admin + Rate limit (20/menit) | Register user baru |

#### Contoh Request & Response

```json
// POST /zenith-pay/admin/users
// Request:
{
  "username": "kasir1",
  "password": "secret123",
  "full_name": "Kasir Satu",
  "email": "kasir1@example.com",
  "role": "cashier"
}

// Response:
{
  "code": 201,
  "status": "success",
  "message": "user berhasil dibuat",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "kasir1",
    "full_name": "Kasir Satu",
    "email": "kasir1@example.com",
    "role": "cashier",
    "is_active": true
  }
}

// GET /zenith-pay/admin/users?page=1&limit=10
// Response:
{
  "code": 200,
  "status": "success",
  "message": "daftar user berhasil diambil",
  "pagination": {
    "total": 5,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  },
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "kasir1",
      "full_name": "Kasir Satu",
      "email": "kasir1@example.com",
      "role": "cashier",
      "is_active": true
    }
  ]
}
```

### Kategori

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `GET` | `/zenith-pay/categories` | JWT (100/menit) | List kategori |
| `GET` | `/zenith-pay/categories/:id` | JWT (100/menit) | Detail kategori |
| `POST` | `/zenith-pay/categories/admin` | JWT + Admin (50/menit) | Buat kategori |
| `PUT` | `/zenith-pay/categories/admin/:id` | JWT + Admin (50/menit) | Update kategori |
| `DELETE` | `/zenith-pay/categories/admin/:id` | JWT + Admin (50/menit) | Hapus kategori |

#### Contoh Request & Response

```json
// POST /zenith-pay/categories/admin
// Request:
{
  "name": "Makanan"
}

// Response:
{
  "code": 201,
  "status": "success",
  "message": "kategori berhasil dibuat",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Makanan"
  }
}

// PUT /zenith-pay/categories/admin/:id
// Request:
{
  "name": "Minuman"
}

// Response:
{
  "code": 200,
  "status": "success",
  "message": "kategori berhasil diupdate",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Minuman"
  }
}

// GET /zenith-pay/categories
// Response:
{
  "code": 200,
  "status": "success",
  "message": "daftar kategori berhasil diambil",
  "pagination": {
    "total": 5,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  },
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "name": "Makanan"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "name": "Minuman"
    }
  ]
}

// GET /zenith-pay/categories/:id
// Response:
{
  "code": 200,
  "status": "success",
  "message": "detail kategori berhasil diambil",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Makanan"
  }
}

// DELETE /zenith-pay/categories/admin/:id
// Response:
{
  "code": 200,
  "status": "success",
  "message": "kategori berhasil dihapus",
  "data": null
}
```

### Produk

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `GET` | `/zenith-pay/products` | JWT (100/menit) | List produk |
| `GET` | `/zenith-pay/products/:id` | JWT (100/menit) | Detail produk |
| `POST` | `/zenith-pay/products/admin` | JWT + Admin (50/menit) | Buat produk |
| `PUT` | `/zenith-pay/products/admin/:id` | JWT + Admin (50/menit) | Update produk |
| `DELETE` | `/zenith-pay/products/admin/:id` | JWT + Admin (50/menit) | Hapus produk |

#### Contoh Request & Response

```json
// POST /zenith-pay/products/admin
// Request:
{
  "category_id": "660e8400-e29b-41d4-a716-446655440001",
  "sku": "",
  "name": "Nasi Goreng",
  "price": 15000,
  "stock": 50
}

// Response:
{
  "code": 201,
  "status": "success",
  "message": "produk berhasil dibuat",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440003",
    "category_id": "660e8400-e29b-41d4-a716-446655440001",
    "category_name": "Makanan",
    "sku": "SKU-20260520-00001",
    "name": "Nasi Goreng",
    "price": 15000,
    "stock": 50,
    "image": "https://foodish-api.com/images/rice/rice15.jpg",
    "created_at": "2026-05-20T08:00:00Z",
    "updated_at": "2026-05-20T08:00:00Z"
  }
}

// PUT /zenith-pay/products/admin/:id
// Request:
{
  "name": "Nasi Goreng Spesial",
  "price": 18000
}

// Response:
{
  "code": 200,
  "status": "success",
  "message": "produk berhasil diupdate",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440003",
    "category_id": "660e8400-e29b-41d4-a716-446655440001",
    "category_name": "Makanan",
    "sku": "SKU-20260520-00001",
    "name": "Nasi Goreng Spesial",
    "price": 18000,
    "stock": 50,
    "image": "https://foodish-api.com/images/rice/rice15.jpg",
    "created_at": "2026-05-20T08:00:00Z",
    "updated_at": "2026-05-20T08:05:00Z"
  }
}

// GET /zenith-pay/products?page=1&limit=10
// Response:
{
  "code": 200,
  "status": "success",
  "message": "daftar produk berhasil diambil",
  "pagination": {
    "total": 10,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  },
  "data": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440003",
      "category_id": "660e8400-e29b-41d4-a716-446655440001",
      "category_name": "Makanan",
      "sku": "SKU-20260520-00001",
      "name": "Nasi Goreng",
      "price": 15000,
      "stock": 50,
      "image": "https://foodish-api.com/images/rice/rice15.jpg",
      "created_at": "2026-05-20T08:00:00Z",
      "updated_at": "2026-05-20T08:00:00Z"
    }
  ]
}

// GET /zenith-pay/products/:id
// Response:
{
  "code": 200,
  "status": "success",
  "message": "detail produk berhasil diambil",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440003",
    "category_id": "660e8400-e29b-41d4-a716-446655440001",
    "category_name": "Makanan",
    "sku": "SKU-20260520-00001",
    "name": "Nasi Goreng",
    "price": 15000,
    "stock": 50,
    "image": "https://foodish-api.com/images/rice/rice15.jpg",
    "created_at": "2026-05-20T08:00:00Z",
    "updated_at": "2026-05-20T08:00:00Z"
  }
}

// DELETE /zenith-pay/products/admin/:id
// Response:
{
  "code": 200,
  "status": "success",
  "message": "produk berhasil dihapus",
  "data": null
}
```

### Transaksi

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `POST` | `/zenith-pay/transactions` | JWT + Cashier + Active Shift (30/menit) | Buat transaksi (cashier) |
| `GET` | `/zenith-pay/transactions` | JWT + Cashier + Active Shift (30/menit) | List transaksi (cashier) |
| `GET` | `/zenith-pay/transactions/:id` | JWT + Cashier + Active Shift (30/menit) | Detail transaksi (cashier) |
| `GET` | `/zenith-pay/admin/transactions` | JWT + Admin (60/menit) | List semua transaksi (admin) |
| `GET` | `/zenith-pay/admin/transactions/:id` | JWT + Admin (60/menit) | Detail transaksi (admin) |

#### Contoh Request & Response

```json
// POST /zenith-pay/transactions
// Request:
{
  "payment_method": "cash",
  "payment_amount": 50000,
  "items": [
    { "product_id": "770e8400-e29b-41d4-a716-446655440003", "quantity": 2 },
    { "product_id": "770e8400-e29b-41d4-a716-446655440004", "quantity": 1 }
  ]
}

// Response:
{
  "code": 201,
  "status": "success",
  "message": "transaksi berhasil",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440005",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "transaction_date": "2026-05-20T10:30:00Z",
    "payment_method": "cash",
    "total_amount": 35000,
    "payment_amount": 50000,
    "change_amount": 15000,
    "items": [
      {
        "product_id": "770e8400-e29b-41d4-a716-446655440003",
        "product_name": "Nasi Goreng",
        "product_price": 15000,
        "quantity": 2,
        "sub_total": 30000
      },
      {
        "product_id": "770e8400-e29b-41d4-a716-446655440004",
        "product_name": "Es Teh",
        "product_price": 5000,
        "quantity": 1,
        "sub_total": 5000
      }
    ]
  }
}

// GET /zenith-pay/transactions?page=1&limit=10&from=2026-05-01&to=2026-05-20
// Response:
{
  "code": 200,
  "status": "success",
  "message": "daftar transaksi berhasil diambil",
  "pagination": {
    "total": 25,
    "page": 1,
    "limit": 10,
    "total_pages": 3
  },
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440005",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "transaction_date": "2026-05-20T10:30:00Z",
      "payment_method": "cash",
      "total_amount": 35000,
      "payment_amount": 50000,
      "change_amount": 15000,
      "items": [
        {
          "product_id": "770e8400-e29b-41d4-a716-446655440003",
          "product_name": "Nasi Goreng",
          "product_price": 15000,
          "quantity": 2,
          "sub_total": 30000
        }
      ]
    }
  ]
}

// GET /zenith-pay/transactions/:id
// Response:
{
  "code": 200,
  "status": "success",
  "message": "detail transaksi berhasil diambil",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440005",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "transaction_date": "2026-05-20T10:30:00Z",
    "payment_method": "cash",
    "total_amount": 35000,
    "payment_amount": 50000,
    "change_amount": 15000,
    "items": [
      {
        "product_id": "770e8400-e29b-41d4-a716-446655440003",
        "product_name": "Nasi Goreng",
        "product_price": 15000,
        "quantity": 2,
        "sub_total": 30000
      },
      {
        "product_id": "770e8400-e29b-41d4-a716-446655440004",
        "product_name": "Es Teh",
        "product_price": 5000,
        "quantity": 1,
        "sub_total": 5000
      }
    ]
  }
}
```

> **Catatan:** Endpoint `GET /zenith-pay/admin/transactions` dan `GET /zenith-pay/admin/transactions/:id` memiliki response yang sama seperti di atas, hanya saja akses admin dan tidak terbatas pada transaksi milik cashier tertentu.

### Shift (Cashier Only)

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `POST` | `/zenith-pay/shifts/open` | JWT + Cashier (10/menit) | Buka shift |
| `POST` | `/zenith-pay/shifts/close` | JWT + Cashier (10/menit) | Tutup shift |
| `GET` | `/zenith-pay/shifts/active` | JWT + Cashier (10/menit) | Cek shift aktif |

#### Contoh Request & Response

```json
// POST /zenith-pay/shifts/open
// Request:
{
  "opening_balance": 1000000
}

// Response:
{
  "code": 200,
  "status": "success",
  "message": "shift berhasil dibuka",
  "data": {
    "id": "990e8400-e29b-41d4-a716-446655440006",
    "cashier_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "open",
    "opening_balance": 1000000,
    "closing_balance": 0,
    "expected_closing_balance": 0,
    "variance": 0,
    "cash_income": 0,
    "debit_income": 0,
    "qris_income": 0,
    "opened_at": "2026-05-20T08:00:00Z",
    "closed_at": null
  }
}

// POST /zenith-pay/shifts/close
// Request:
{
  "shift_id": "990e8400-e29b-41d4-a716-446655440006",
  "closing_balance": 1178000
}

// Response:
{
  "code": 200,
  "status": "success",
  "message": "shift berhasil ditutup",
  "data": {
    "id": "990e8400-e29b-41d4-a716-446655440006",
    "cashier_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "closed",
    "opening_balance": 1000000,
    "closing_balance": 1178000,
    "expected_closing_balance": 1178000,
    "variance": 0,
    "cash_income": 178000,
    "debit_income": 0,
    "qris_income": 0,
    "opened_at": "2026-05-20T08:00:00Z",
    "closed_at": "2026-05-20T17:00:00Z"
  }
}

// GET /zenith-pay/shifts/active
// Response:
{
  "code": 200,
  "status": "success",
  "message": "shift aktif berhasil diambil",
  "data": {
    "id": "990e8400-e29b-41d4-a716-446655440006",
    "cashier_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "open",
    "opening_balance": 1000000,
    "closing_balance": 0,
    "expected_closing_balance": 0,
    "variance": 0,
    "cash_income": 0,
    "debit_income": 0,
    "qris_income": 0,
    "opened_at": "2026-05-20T08:00:00Z",
    "closed_at": null
  }
}
```

### Laporan (Admin Only)

| Method | Endpoint | Middleware | Deskripsi |
|--------|----------|------------|-----------|
| `GET` | `/zenith-pay/admin/reports/daily?date=YYYY-MM-DD` | JWT + Admin (60/menit) | Laporan harian |
| `GET` | `/zenith-pay/admin/reports/monthly?period=YYYY-MM` | JWT + Admin (60/menit) | Laporan bulanan |
| `GET` | `/zenith-pay/admin/reports/revenue?from=YYYY-MM-DD&to=YYYY-MM-DD` | JWT + Admin (60/menit) | Tren revenue |

#### Contoh Request & Response

```json
// GET /zenith-pay/admin/reports/daily?date=2026-05-20
// Response:
{
  "code": 200,
  "status": "success",
  "message": "laporan harian berhasil diambil",
  "data": {
    "date": "2026-05-20",
    "total_transactions": 25,
    "total_revenue": 875000
  }
}

// GET /zenith-pay/admin/reports/monthly?period=2026-05
// Response:
{
  "code": 200,
  "status": "success",
  "message": "laporan bulanan berhasil diambil",
  "data": {
    "month": "2026-05",
    "total_transactions": 450,
    "total_revenue": 15750000
  }
}

// GET /zenith-pay/admin/reports/revenue?from=2026-05-01&to=2026-05-20
// Response:
{
  "code": 200,
  "status": "success",
  "message": "tren revenue berhasil diambil",
  "data": [
    { "date": "2026-05-01", "total_revenue": 750000 },
    { "date": "2026-05-02", "total_revenue": 820000 },
    { "date": "2026-05-20", "total_revenue": 875000 }
  ]
}
```

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
- `shift_id` (UUID, FK → shifts)
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
- `expected_closing_balance` (int64, nullable)
- `variance` (int64, nullable)
- `cash_income` (int64, nullable)
- `debit_income` (int64, nullable)
- `qris_income` (int64, nullable)
- `opened_at` (timestamp)
- `closed_at` (timestamp, nullable)

## Middleware 🛡️

| Middleware | Deskripsi |
|---|---|
| **Recover** | Panic recovery — cegah crash server |
| **Request ID** | Tambahkan `X-Request-ID` header otomatis |
| **CORS** | Mengizinkan origin `localhost:3000` |
| **Method Validation** | Whitelist HTTP method (GET, POST, PUT, PATCH, DELETE, OPTIONS) |
| **HTTP Logger** | Log request/response (format: `[time] status - latency method path`) |
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
- Zerolog structured logging (console development, JSON production)
