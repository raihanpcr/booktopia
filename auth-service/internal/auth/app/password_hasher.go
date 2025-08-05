package app

// PasswordHasher adalah antarmuka (interface) yang mendefinisikan kontrak untuk layanan hashing kata sandi.
// Tujuannya adalah untuk mengabstraksi detail implementasi algoritma hashing (misalnya, bcrypt, scrypt).
// Dengan antarmuka ini, logika aplikasi utama (seperti di AuthApp) tidak perlu tahu
// bagaimana kata sandi di-hash atau diverifikasi, sehingga kode tetap bersih dan mudah diuji.
// Ini juga memungkinkan penggantian algoritma hashing di masa depan tanpa mengubah kode inti.
type PasswordHasher interface {
	// Hash menerima kata sandi dalam bentuk teks biasa dan mengembalikan
	// representasi kata sandi yang sudah di-hash.
	// Jika terjadi kesalahan selama proses hashing, error akan dikembalikan.
	Hash(password string) (string, error)

	// Check membandingkan kata sandi teks biasa dengan kata sandi yang sudah di-hash.
	// Ini mengembalikan true jika kata sandi cocok, dan false jika tidak cocok.
	// Implementasi harus menangani perbedaan algoritma hashing secara internal.
	Check(password, hashed string) bool
}