# Dokumentasi Book-Service

## Description
Tujuan: Menyimpan informasi lengkap mengenai setiap buku yang tersedia di platform.

Stok Numerik vs. Status: Anda menyebutkan "Tidak pakai stok numerik, hanya status available/unavailable". Ini berarti sistem Anda mengelola ketersediaan buku secara biner: buku itu ada atau tidak ada. Tidak ada jumlah spesifik dari buku yang sama. Ini cocok jika setiap buku dianggap unik atau ini adalah platform e-book di mana ketersediaan tidak terbatas secara fisik.

### Kolom Penting:
- book_id: Identifikasi unik untuk setiap buku.,
- status: Menunjukkan ketersediaan buku (available atau unavailable),
- is_donation_only: Indikator penting apakah buku tersebut gratis/donasi atau dijual.

### Relasi Penting
- books memiliki relasi one-to-many dengan transaction_details (satu buku bisa muncul di banyak detail transaksi).,
- books memiliki relasi one-to-many dengan ebook_gift_log (satu buku bisa diberikan sebagai hadiah berkali-kali).

## Database
```
CREATE TABLE books (
    book_id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL,
    author VARCHAR(100),
    publisher VARCHAR(100),
    year_published YEAR,
    category VARCHAR(100),
    price DECIMAL(10, 2) NOT NULL,
    status ENUM('available', 'unavailable') DEFAULT 'available',
    is_donation_only BOOLEAN DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

# Struktur Folder

## CMD
- /main.go

## Internal

### Internal / dto
- / book_request.go
- / book_response.go
- / error_message.go
- / mapper.go

### interal / handler
- / book_handler

## Model
- / book.go

## Repository
- / book_repository.go

## Service
- / book_service.go

| Service                 | Cocok Pakai                | Alasan                                                                                                                      |
| ----------------------- | -------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| **auth-service**        | **PostgreSQL**             | - Data struktural (email, password, user)<br>- Butuh constraint seperti unique, foreign key<br>- Query kompleks             |
| **book-service**        | **MongoDB** ✅ | - Kalau data buku fleksibel (bisa punya genre, tag dinamis): MongoDB<br> |
| **transaction-service** | **PostgreSQL** ✅           | - Transaksi harus ACID<br>- Butuh join antar user, wallet, book<br>- Sangat cocok untuk relational DB                       |
| **wallet-service**      | **PostgreSQL** ✅           | - Saldo perlu akurasi tinggi<br>- Transaksi harus konsisten<br>- Decimal types lebih aman di SQL                            |
