package mail

import (
	"github.com/auronvila/simple-bank/util"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	err := godotenv.Load("../app.env")
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	sender := NewGmailSender(config.SmtpSenderName, config.SmtpUsername, config.SmtpPass)

	subject := "A test email"
	contet := `
<h1>Hello World</h1>
<p> This is a test message from AV </p>
`
	to := []string{"auronvila.dev@gmail.com"}
	attachFiles := []string{"../Makefile"}

	err = sender.SendEmail(subject, contet, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
