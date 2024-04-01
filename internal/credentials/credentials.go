package credentials

import (
	"crypto/x509"
	helper "github.com/aws/rolesanywhere-credential-helper/aws_signing_helper"
	"github.com/sethvargo/go-githubactions"
	"main/internal/util"
	"os"
	"strconv"
)

func CreateCredentials() {
	actions := githubactions.New()

	privateKeyContent := actions.GetInput("private-key")
	certificateContent := actions.GetInput("certificate")
	region := actions.GetInput("aws-region")
	roleArn := actions.GetInput("role-arn")
	profileArn := actions.GetInput("profile-arn")
	trustAnchorArn := actions.GetInput("trust-anchor-arn")

	sessionDuration := 900
	sessionDurationOverride := actions.GetInput("session-duration")
	if sessionDurationOverride != "" {
		sessionDuration, _ = strconv.Atoi(sessionDurationOverride)
	}

	credentialsOptions := helper.CredentialsOpts{
		RoleArn:           roleArn,
		ProfileArnStr:     profileArn,
		TrustAnchorArnStr: trustAnchorArn,
		SessionDuration:   sessionDuration,
		Region:            region,
	}

	privateKey, err := util.GetPrivateKeyFromPEMString([]byte(privateKeyContent))
	if err != nil {
		actions.Errorf("%s", err)
		os.Exit(1)
	}
	certificate, err := util.GetCertificateFromPEM([]byte(certificateContent))
	if err != nil {
		actions.Errorf("%s", err)
		os.Exit(1)
	}
	var certificateChain []*x509.Certificate

	signer, signingAlgorithm, err := helper.GetFileSystemSigner(privateKey, certificate, certificateChain)
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

	actions.SetEnv("AWS_DEFAULT_REGION", region)
	actions.SetEnv("AWS_REGION", region)
}
