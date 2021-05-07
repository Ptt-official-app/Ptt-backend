package usecase

import "github.com/Ptt-official-app/Ptt-backend/internal/mail"

func (usecase *usecase) UpdateMail(mail mail.Mail) error {
	usecase.mailProvider = mail
	return nil
}
