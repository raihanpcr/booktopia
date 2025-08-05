package route

import (
	"gateway-service/internal/handler"
	customMiddleware "gateway-service/internal/middleware"

	"github.com/labstack/echo/v4"
)

// SetupRoutes mendaftarkan semua route aplikasi ke instance Echo.
func SetupRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	bookHandler *handler.BookHandler,
	transactionHandler *handler.TransactionHandler,
	walletHandler *handler.WalletHandler,
	giftingHandler *handler.GiftingHandler,
) {
	api := e.Group("/api")
	{
		// === ROUTE PUBLIK ===
		// Route untuk login/register di-proxy ke auth-service
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// Route untuk buku di-proxy ke book-service
		api.GET("/books", bookHandler.GetBooks)
		api.GET("/books/:id", bookHandler.GetBookByID)

		// === ROUTE TERLINDUNGI (BUTUH LOGIN/TOKEN JWT) ===
		// Buat grup baru dan terapkan middleware otentikasi
		protected := api.Group("")
		protected.Use(customMiddleware.JwtAuthMiddleware)
		{
			// Pindahkan route transaksi ke dalam grup yang dilindungi
			protected.POST("/transactions", transactionHandler.CreateTransaction)
			protected.GET("/transactions", transactionHandler.GetTransactions)
			protected.GET("/wallet/balance", walletHandler.GetBalance)
			protected.POST("/wallet/topup", walletHandler.TopUp)
			protected.POST("/gifts", giftingHandler.SendGift)
			
			// --- ROUTE KHUSUS ADMIN ---
			// Anda bisa membuat middleware baru untuk memeriksa role 'admin'
			admin := protected.Group("/admin")
			admin.Use(customMiddleware.AdminOnlyMiddleware)
			{
				admin.POST("/books", bookHandler.CreateBook)
				admin.PUT("/books/:id", bookHandler.UpdateBook)
				admin.DELETE("/books/:id", bookHandler.DeleteBook)
			}
		}
	}
}
