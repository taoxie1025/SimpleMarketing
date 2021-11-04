package mail

import (
	"email_action/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	log     = logging.NewZapLogger()
	charSet = "UTF-8"
)

type SesAdapter struct {
	sesSvc *ses.SES
}

func NewSesAdapter() *SesAdapter {
	log.Infof("NewSesAdapter():")
	awsConfig := &aws.Config{
		Region: aws.String("us-west-2"),
	}
	sess := session.Must(session.NewSession(awsConfig))
	svc := ses.New(sess)
	return &SesAdapter{
		sesSvc: svc,
	}
}

func (s *SesAdapter) SendEmail(from, to, subject, htmlBody, textBody string) error {
	log.Infof("SendEmail(): from %s, to %s", from, to)
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}
	_, err := s.sesSvc.SendEmail(input)
	if err != nil {
		log.Errorf("SendEmail(): error = %v", err)
		return err
	}
	return nil
}

func (s *SesAdapter) VerifyEmail(email string) error {
	log.Infof("VerifyEmail(): %s", email)
	// TODO: create custom template for verification
	verifyReq := ses.SendCustomVerificationEmailInput{
		EmailAddress: &email,
		TemplateName: aws.String("EmailVerificationTemplate"),
	}
	_, err := s.sesSvc.SendCustomVerificationEmail(&verifyReq)
	if err != nil {
		log.Errorf("VerifyEmail(): error = %v", err)
		return err
	}
	return nil
}

func (s *SesAdapter) IsEmailVerified(email string) (bool, error) {
	log.Infof("IsEmailVerified(): email = %s", email)
	req := ses.GetIdentityVerificationAttributesInput{
		Identities: []*string{&email},
	}
	output, err := s.sesSvc.GetIdentityVerificationAttributes(&req)
	if err != nil {
		log.Errorf("IsEmailVerified(): error = %v", err)
		return false, err
	}

	if output == nil || output.VerificationAttributes[email] == nil ||
		*output.VerificationAttributes[email].VerificationStatus != ses.VerificationStatusSuccess {
		return false, nil
	}
	return true, nil
}
