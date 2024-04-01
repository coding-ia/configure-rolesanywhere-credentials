package credentials

import (
	"os"
	"testing"
)

func TestCreateCredentials(t *testing.T) {
	privateKey := `
`
	certificate := `
`
	os.Setenv("INPUT_PRIVATE-KEY", privateKey)
	os.Setenv("INPUT_CERTIFICATE", certificate)
	os.Setenv("INPUT_ROLE-ARN", "arn:aws:iam::211125334931:role/TestRole")
	os.Setenv("INPUT_PROFILE-ARN", "arn:aws:rolesanywhere:us-east-1:211125334931:profile/89afc4fd-7b96-443b-9613-10bd828d4c24")
	os.Setenv("INPUT_TRUST-ANCHOR-ARN", "arn:aws:rolesanywhere:us-east-1:211125334931:trust-anchor/6dd648ad-92a7-4a44-a59c-2b1d56dfbbee")

	tempEnv, err := createTempGitHubEnviornment(".env")
	if err != nil {
		t.Fail()
	}
	tempOutput, err := createTempGitHubEnviornment(".output")
	if err != nil {
		t.Fail()
	}
	os.Setenv("GITHUB_ENV", tempEnv)
	os.Setenv("GITHUB_OUTPUT", tempOutput)
	CreateCredentials()
}

func createTempGitHubEnviornment(fileName string) (string, error) {
	tempFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return "", err
	}
	tempFile.Close()
	return tempFile.Name(), nil
}
