# Rangkuman: Checklist Final Project untuk Setiap Kelompok

## 1. Tema & Persiapan
- [x] Pilih dan konfirmasi Topik, tema/project ke instruktur
### Nama Kelompok: Booktopia (Alpha)
- Tema: Pendidikan Berkualitas untuk Masa Depan Gemilang: Merangkai Potensi Anak Bangsa
- Topik: Pendidikan Berkualitas

## 2. Teknologi & Stack
- [ ] Pilih stack:
  - Framework Go (Gin/Echo/gRPC)        ✅ Using Echo (REST) and gRPC
  - Database (MySQL/PostgreSQL/MongoDB) ✅ Using PostgreSQL
  - Deployment tools (Docker, GCP)      ⌛
  - Third Party API                     ✅ Using mail notification (mailtrap.io)

## 3. Desain Database
- [x] Rancang ERD sesuai kebutuhan aplikasi
  - Wajib ada: `users`, `main entity` (sesuai tema), dan `transactions`

## 4. Implementasi REST API
- [x] Register & login (JWT authentication)
- [x] CRUD entity utama (data master)
- [x] Proses transaksi
- [x] Scheduler / backup data (otomatis/manual)
- [ ] Report / laporan
- [x] Multi-role (admin & user)

## 5. Integrasi & Fitur Tambahan
- [x] Integrasi minimal 1 third party API
  - Pilih salah satu: payment gateway, product API, atau email notification

## 6. Dokumentasi & Testing
- [ ] Dokumentasi API lengkap (Swagger/OpenAPI)
- [x] Unit testing dengan mocking (setiap fitur/endpoint)

## 7. Deployment & Presentasi
- [ ] Deployment: aplikasi harus running di server GCP via Docker
- [ ] Presentasi: PowerPoint + demo aplikasi

## 8. Pengumpulan
- [ ] Kumpulkan source code, desain, dan dokumentasi testing
