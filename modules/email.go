package modules

import (
	"context"
	"log"

	"github.com/trycourier/courier-go/v2"
)

func SendEmailMessage(email string, verificationCode string) {

	client := courier.CreateClient("pk_prod_J10WMS1MPF4HRHPZ2X4NP8KQ2W8D", nil)

	requestID, err := client.SendMessage(
		context.Background(),
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to": map[string]string{
					"email": email,
				},
				"content": map[string]string{
					"title": "Website",
					"body":  "Verification code: {{code}}",
				},
				"data": map[string]string{
					"code": verificationCode,
				},
			},
		})

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(requestID)
}
