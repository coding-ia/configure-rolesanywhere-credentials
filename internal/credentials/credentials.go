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
	action := githubactions.New()

	privateKeyContent := action.GetInput("private-key")
	certificateContent := action.GetInput("certificate")

	credentialsOptions := createCredentialsOptions(action)

	privateKey, err := util.GetPrivateKeyFromPEMString([]byte(privateKeyContent))
	if err != nil {
		action.Errorf("%s", err)
		os.Exit(1)
	}
	certificate, err := util.GetCertificateFromPEM([]byte(certificateContent))
	if err != nil {
		action.Errorf("%s", err)
		os.Exit(1)
	}
	var certificateChain []*x509.Certificate

	signer, signingAlgorithm, err := helper.GetFileSystemSigner(privateKey, certificate, certificateChain)
	if err != nil {
		action.Errorf("%s", err)
		os.Exit(1)
	}
	defer signer.Close()
	credentialProcessOutput, err := helper.GenerateCredentials(&credentialsOptions, signer, signingAlgorithm)
	if err != nil {
		action.Errorf("%s", err)
		os.Exit(1)
	}

	action.AddMask(credentialProcessOutput.AccessKeyId)
	action.SetEnv("AWS_ACCESS_KEY_ID", credentialProcessOutput.AccessKeyId)
	action.SetOutput("aws-access-key-id", credentialProcessOutput.AccessKeyId)

	action.AddMask(credentialProcessOutput.SecretAccessKey)
	action.SetEnv("AWS_SECRET_ACCESS_KEY", credentialProcessOutput.SecretAccessKey)
	action.SetOutput("aws-secret-access-key", credentialProcessOutput.SecretAccessKey)

	action.AddMask(credentialProcessOutput.SessionToken)
	action.SetEnv("AWS_SESSION_TOKEN", credentialProcessOutput.SessionToken)
	action.SetOutput("aws-session-token", credentialProcessOutput.SessionToken)

	action.SetEnv("AWS_DEFAULT_REGION", credentialsOptions.Region)
	action.SetEnv("AWS_REGION", credentialsOptions.Region)
}

func createCredentialsOptions(action *githubactions.Action) helper.CredentialsOpts {
	sessionDuration := 900
	noVerifySSL := false
	withProxy := false

	region := action.GetInput("aws-region")
	roleArn := action.GetInput("role-arn")
	profileArn := action.GetInput("profile-arn")
	trustAnchorArn := action.GetInput("trust-anchor-arn")
	sessionDurationOverride := action.GetInput("session-duration")
	endpoint := action.GetInput("endpoint")
	ssl := action.GetInput("no-verify-ssl")
	proxy := action.GetInput("with-proxy")

	if sessionDurationOverride != "" {
		sessionDuration, _ = strconv.Atoi(sessionDurationOverride)
	}
	noVerifySSL, _ = strconv.ParseBool(ssl)
	withProxy, _ = strconv.ParseBool(proxy)

	credentialsOptions := helper.CredentialsOpts{
		RoleArn:           roleArn,
		ProfileArnStr:     profileArn,
		TrustAnchorArnStr: trustAnchorArn,
		SessionDuration:   sessionDuration,
		Region:            region,
		Endpoint:          endpoint,
		NoVerifySSL:       noVerifySSL,
		WithProxy:         withProxy,
	}

	return credentialsOptions
}
