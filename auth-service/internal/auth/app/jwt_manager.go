package app

// JWTManager adalah antarmuka (interface) yang mendefinisikan kontrak untuk layanan manajemen JSON Web Token (JWT).
// Tujuannya adalah untuk mengabstraksi detail implementasi pembuatan token JWT,
// sehingga logika aplikasi utama (seperti di AuthApp) tidak perlu tahu bagaimana token sebenarnya dibuat.
// Ini memungkinkan penggantian implementasi JWT (misalnya, dari JWT bawaan ke library pihak ketiga)
// tanpa perlu mengubah kode inti aplikasi.
type JWTManager interface {
	// GenerateToken membuat dan mengembalikan string token JWT yang baru.
	// Token ini akan berisi informasi userID dan email pengguna sebagai klaim (claims).
	// Jika terjadi kesalahan selama pembuatan token, error akan dikembalikan.
	GenerateToken(userID, email, role string) (string, error)
}