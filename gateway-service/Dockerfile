# =========================================================================
#  Dockerfile Multi-Stage untuk Auth Service
# =========================================================================
#
# Tahap 1: Meng-compile aplikasi Go menjadi satu file binary.
# Tahap 2: Membuat image final yang sangat kecil hanya dengan file binary tersebut.
#
# =========================================================================

# --- STAGE 1: Builder ---
# Menggunakan image Go resmi berbasis Alpine Linux yang ringan.
FROM golang:1.24 AS builder

# Menentukan direktori kerja di dalam container.
WORKDIR /app

# 1. Copy the go.work and go.work.sum files from the project root.
#    This tells Go that /app is a workspace.
#    Assuming your docker-compose.yml sets context: . (project root)
#    and dockerfile: ./<service_name>/Dockerfile
COPY go.work go.work.sum ./

# Salin go.mod dan go.sum dari setiap layanan ke dalam satu direktori
COPY ./auth-service/go.mod ./auth-service/go.sum ./auth-service/
COPY ./book-service/go.mod ./book-service/go.sum ./book-service/
COPY ./gateway-service/go.mod ./gateway-service/go.sum ./gateway-service/
COPY ./gifting-service/go.mod ./gifting-service/go.sum ./gifting-service/
COPY ./transaction-service/go.mod ./transaction-service/go.sum ./transaction-service/
COPY ./wallet-service/go.mod ./wallet-service/go.sum ./wallet-service/

# Perintah go mod download akan mengunduh dependensi dan membuat go.sum jika belum ada
RUN go mod download

# Menyalin seluruh sisa source code proyek.
COPY . .

# Masuk ke direktori layanan yang sedang dibangun sebelum menjalankan go mod tidy
WORKDIR /app/gateway-service

# Jalankan go mod tidy sebelum go build untuk memastikan go.sum sinkron
RUN go mod tidy

# Meng-compile aplikasi Go menjadi binary statis.
# Ganti './cmd/main.go' jika entrypoint Anda berada di lokasi lain.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

# --- STAGE 2: Final ---
# Menggunakan image Alpine Linux yang sangat kecil sebagai dasar.
FROM alpine:latest

# Menyalin file binary yang sudah di-compile dari stage 'builder'.
COPY --from=builder /app/main /app/main

# Memberi tahu Docker bahwa container ini akan mengekspos dua port:
# Port 8082 untuk server HTTP (termasuk Swagger UI internal)
# Port 50051 untuk server gRPC yang akan dihubungi oleh Gateway.
EXPOSE 8000

# Perintah yang akan dijalankan saat container dimulai.
CMD ["/app/main"]