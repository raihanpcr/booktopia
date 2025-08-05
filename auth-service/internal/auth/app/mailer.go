package app

// Mailer mendefinisikan kontrak untuk pengiriman email.
type Mailer interface {
	SendWelcomeEmail(toEmail, toName string)
}
