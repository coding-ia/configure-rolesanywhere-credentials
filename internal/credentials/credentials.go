package credentials

import (
	helper "github.com/aws/rolesanywhere-credential-helper/aws_signing_helper"
	"github.com/sethvargo/go-githubactions"
	"os"
	"path/filepath"
	"strconv"
)

func CreateCredentials() {
	actions := githubactions.New()

	privateKeyContent := actions.GetInput("private-key")
	certificateContent := actions.GetInput("certificate")
	roleArn := actions.GetInput("role-arn")
	profileArn := actions.GetInput("profile-arn")
	trustAnchorArn := actions.GetInput("trust-anchor-arn")

	sessionDuration := 900
	sessionDurationOverride := actions.GetInput("session-duration")
	if sessionDurationOverride != "" {
		sessionDuration, _ = strconv.Atoi(sessionDurationOverride)
	}

	privateKeyFile, err := writeToFile(privateKeyContent, "private.key")
	if err != nil {
		actions.Errorf("Unable to process the private key.")
		os.Exit(1)
	}
	certificateFile, err := writeToFile(certificateContent, "certificate.pem")
	if err != nil {
		actions.Errorf("Unable to process the certificate.")
		os.Exit(1)
	}

	credentialsOptions := helper.CredentialsOpts{
		PrivateKeyId:      privateKeyFile,
		CertificateId:     certificateFile,
		RoleArn:           roleArn,
		ProfileArnStr:     profileArn,
		TrustAnchorArnStr: trustAnchorArn,
		SessionDuration:   sessionDuration,
	}

	signer, signingAlgorithm, err := helper.GetSigner(&credentialsOptions)
	if err != nil {
		actions.Errorf("%s", err)
		os.Exit(1)
	}
	defer signer.Close()
	credentialProcessOutput, err := helper.GenerateCredentials(&credentialsOptions, signer, signingAlgorithm)
	if err != nil {
		actions.Errorf("%s", err)
		os.Exit(1)
	}

	actions.AddMask(credentialProcessOutput.AccessKeyId)
	actions.SetEnv("AWS_ACCESS_KEY_ID", credentialProcessOutput.AccessKeyId)
	actions.SetOutput("aws-access-key-id", credentialProcessOutput.AccessKeyId)

	actions.AddMask(credentialProcessOutput.SecretAccessKey)
	actions.SetEnv("AWS_SECRET_ACCESS_KEY", credentialProcessOutput.SecretAccessKey)
	actions.SetOutput("aws-secret-access-key", credentialProcessOutput.SecretAccessKey)

	actions.AddMask(credentialProcessOutput.SessionToken)
	actions.SetEnv("AWS_SESSION_TOKEN", credentialProcessOutput.SessionToken)
	actions.SetOutput("aws-session-token", credentialProcessOutput.SessionToken)
}

func writeToFile(content, fileName string) (string, error) {
	tempFilePath := filepath.Join(os.TempDir(), fileName)
	err := os.WriteFile(tempFilePath, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return tempFilePath, nil
}
