# Applicant Management API — Sekolah Karir Back-End Assignment

REST API sederhana untuk mengelola applicant menggunakan **Go, Echo, GORM, dan SQLite**.

## Fitur

### Fitur
- Create applicant dengan validasi nama, email, dan posisi.
- List applicants.
- Update stage applicant.
- Get applicant berdasarkan ID.
- Search berdasarkan nama, email, atau posisi (`search`).
- Filter berdasarkan stage (`stage`) dan posisi (`position`).
- Pagination (`page`, `limit`, maksimum 100 item per halaman).
- Summary jumlah applicant per stage.
- Unique email untuk mencegah applicant terduplikasi.
- Response JSON konsisten, HTTP status code, request ID, logger, dan recovery middleware.
- SQLite auto-migration dan dummy seed yang idempotent.


## Arsitektur

```text
.
├── app/
│   ├── configs/       # Konfigurasi environment
│   ├── database/      # SQLite connection, migration, dummy seed
│   ├── middlewares/   # Recovery, request ID, logger, CORS
│   └── routers/       # Registrasi endpoint dan dependency wiring
├── features/
│   └── applicant/
│       ├── data/      # GORM model dan repository/query
│       ├── handler/   # HTTP handler, request DTO, response DTO
│       ├── service/   # Business logic dan validasi
│       └── entities.go # Domain entity, interface, dan filter
├── postman/          # Collection, environment, instruksi test
├── utils/responses/  # Standard response helper
└── main.go           # Entry point dan dependency injection
```

Alur request: `HTTP Request → Handler → Service → Repository → SQLite`. Service bergantung pada interface repository, bukan langsung ke GORM.

## Prasyarat dan menjalankan aplikasi

- Go 1.22 atau lebih baru.

```bash
cp .env.example .env
go mod tidy
go run .
```

Server default: `http://localhost:8080`. Database default: `./app.db`.

| Environment | Default | Keterangan |
|---|---|---|
| `APP_PORT` | `8080` | Port HTTP |
| `DB_PATH` | `./app.db` | Lokasi file SQLite |
| `SEED_DATA` | `true` | Memasukkan dummy data saat startup |

Set `SEED_DATA=false` untuk menonaktifkan seed. Seed menggunakan email unik dan `FirstOrCreate`, sehingga restart tidak menggandakan data demo.

## Endpoint API

Semua endpoint utama menggunakan prefix `/api/v1`.

| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/api/v1/applicants` | Membuat applicant (**wajib**) |
| GET | `/api/v1/applicants` | List applicant (**wajib**) |
| PATCH | `/api/v1/applicants/:id/stage` | Update stage (**wajib**) |
| GET | `/api/v1/applicants/:id` | Detail applicant |
| GET | `/api/v1/applicants/summary` | Total dan breakdown per stage |

### 1. Create applicant

`POST /api/v1/applicants`

```json
{
  "name": "Raka Pratama",
  "email": "raka@example.com",
  "position": "Backend Engineer"
}
```

Success: `201 Created`. Stage awal otomatis `APPLIED`.

Validasi: semua field wajib; nama dan posisi maksimal 120 karakter; email maksimal 160 karakter dan harus valid; email dinormalisasi menjadi lowercase; email duplikat menghasilkan `409 Conflict`.

### 2. List, search, filter, dan pagination

`GET /api/v1/applicants?page=1&limit=10&search=raka&stage=APPLIED&position=backend`

Semua query bersifat opsional. Default `page=1`, `limit=10`; limit maksimal 100. `search` mencari pada nama/email/posisi; `stage` menerima stage valid; `position` melakukan pencarian parsial tanpa membedakan kapitalisasi.

Response `200 OK` memiliki `data` array dan `meta`:

```json
{
  "success": true,
  "message": "applicants retrieved",
  "data": [],
  "meta": { "page": 1, "limit": 10, "total": 0, "total_pages": 0 }
}
```

### 3. Update stage

`PATCH /api/v1/applicants/1/stage`

```json
{ "stage": "INTERVIEW" }
```

Stage yang diterima: `APPLIED`, `SCREENING`, `INTERVIEW`, `OFFER`, `HIRED`, `REJECTED` (case-insensitive). Stage di luar daftar menghasilkan `400 Bad Request`. ID yang tidak ditemukan menghasilkan `404 Not Found`.

### 4. Get applicant detail

`GET /api/v1/applicants/1` — mengembalikan satu applicant atau `404 Not Found`.

### 5. Summary

`GET /api/v1/applicants/summary` — mengembalikan total applicant dan jumlah per stage.

## Dummy data

| Nama | Email | Posisi | Stage |
|---|---|---|---|
| Alya Putri | `alya.putri@example.com` | Backend Engineer | `APPLIED` |
| Bima Pratama | `bima.pratama@example.com` | Frontend Developer | `SCREENING` |
| Citra Maharani | `citra.maharani@example.com` | QA Engineer | `INTERVIEW` |
| Dimas Saputra | `dimas.saputra@example.com` | DevOps Engineer | `OFFER` |
| Nadia Lestari | `nadia.lestari@example.com` | Product Designer | `REJECTED` |

Semua data bersifat fiktif.

## Testing

```bash
go test ./...
go test -cover ./...
```

### Postman

Import dua file berikut ke Postman:
- `postman/Applicant_API.postman_collection.json`
- `postman/Local.postman_environment.json`

Pilih environment **Sekolah Karir - Local**, jalankan API (`go run .`), lalu klik **Run collection**. Jalankan request sesuai urutan; request Create menghasilkan email unik dan menyimpan `applicantId` serta `applicantEmail` untuk test berikutnya. Collection berisi 23 request/test case termasuk happy path dan validasi.

Jika memakai Newman:

```bash
npm install -g newman
newman run postman/Applicant_API.postman_collection.json -e postman/Local.postman_environment.json
```

### Ringkasan test case Postman

| No. | Test case | Expected status |
|---:|---|---:|
| 1 | Health check | 200 |
| 2 | Create applicant | 201 |
| 3 | Get applicant by ID | 200 |
| 4 | List default pagination | 200 |
| 5 | Search applicant | 200 |
| 6 | Filter by stage | 200 |
| 7 | Filter by position | 200 |
| 8 | Pagination limit 2 | 200 |
| 9 | Update stage | 200 |
| 10 | Verify filter after stage update | 200 |
| 11 | Summary by stage | 200 |
| 12–15 | Missing fields / invalid email | 400 |
| 16 | Duplicate email | 409 |
| 17–20 | Invalid stage/filter/pagination | 400 |
| 21 | Unknown applicant | 404 |
| 22–23 | Malformed / zero ID | 400 |

## Design decisions & batasan

- SQLite dipilih karena brief memperbolehkan database lokal ringan.
- Data minimum applicant adalah name, email, position, dan stage.
- Stage awal `APPLIED`; perubahan stage tidak menerapkan workflow transition restriction karena tidak diminta dalam brief.
- Tidak menambahkan authentication, authorization, soft delete, atau integrasi eksternal agar scope tetap kecil dan mudah dijelaskan dalam estimasi pengerjaan 90–120 menit.
- GORM AutoMigrate menyederhanakan setup lokal.

## AI/tool usage note

OpenAI ChatGPT digunakan untuk membantu interpretasi brief, menyusun struktur awal Go/Echo, dokumentasi, dan rancangan test.