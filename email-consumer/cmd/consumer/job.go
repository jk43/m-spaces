package main

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func (app *Application) Job(msg *sarama.ConsumerMessage) {
	//json to map
	//val := kafka.SESTemplate{}
	val := make(map[string]any)
	err := json.Unmarshal(msg.Value, &val)
	if err != nil {
		app.Logger.Error().Err(err).Send()
	}

	data, _ := val["value"].(map[string]any)
	template := aws.String(data["template"].(string))
	templateData, _ := json.Marshal(data["templateData"])
	source := aws.String(data["source"].(string))
	var toAddresses []*string
	tas, _ := data["toAddresses"].([]any)
	for _, i := range tas {
		toAddresses = append(toAddresses, aws.String(i.(string)))
	}
	conf := aws.Config{
		Region: aws.String("us-east-1"),
	}
	sess, _ := session.NewSession(&conf)
	svc := ses.New(sess)

	tei := ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: toAddresses,
		},
		Template:     template,
		TemplateData: aws.String(string(templateData)),
		Source:       source,
	}

	result, err := svc.SendTemplatedEmail(&tei)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return
	}

	// fmt.Println("Email Sent to address: " + Recipient)
	fmt.Println(result)
}
