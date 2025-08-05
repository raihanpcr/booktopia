package hasher

import "golang.org/x/crypto/bcrypt"

// Bcrypt adalah implementasi dari antarmuka PasswordHasher (yang didefinisikan di internal/auth/app).
// Struct ini menggunakan algoritma hashing bcrypt untuk mengamankan kata sandi.
// Dengan mengimplementasikan antarmuka, kita bisa dengan mudah mengganti algoritma hashing
// di masa depan tanpa memengaruhi logika aplikasi inti.
type Bcrypt struct{}

// NewBcrypt adalah konstruktor untuk Bcrypt.
// Fungsi ini mengembalikan instance Bcrypt yang siap digunakan.
func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}

// Hash menerima kata sandi teks biasa dan mengembalikannya dalam bentuk hash bcrypt.
// bcrypt.DefaultCost adalah nilai biaya default yang direkomendasikan,
// yang menentukan kompleksitas komputasi hashing.
func (b *Bcrypt) Hash(password string) (string, error) {
	// Mengubah kata sandi string menjadi slice of byte
	// Menggunakan bcrypt.DefaultCost untuk kekuatan hashing default
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// Mengembalikan hash sebagai string dan error jika ada
	return string(bytes), err
}

// Check membandingkan kata sandi teks biasa dengan kata sandi yang sudah di-hash.
// Fungsi ini mengembalikan true jika kata sandi cocok dengan hash, dan false jika tidak.
// Ini adalah cara yang aman untuk memverifikasi kata sandi tanpa perlu mendekripsi hash.
func (b *Bcrypt) Check(password, hashed string) bool {
	// Membandingkan hash dengan kata sandi yang diberikan.
	// Jika kata sandi cocok, err akan nil. Jika tidak, err akan non-nil.
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	// Mengembalikan true jika tidak ada error (kata sandi cocok), false jika ada error
	return err == nil
}