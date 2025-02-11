package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"

	"github.com/spf13/cobra"
)

var getAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Get a list of Cloud Manager Accounts",
	Long: `EXAMPLES
	splicectl get accounts

	    * if no accounts are listed, you will need to logon to the Ops Center
`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = versionDetail.RequirementMet("get_accounts")

		out, err := getAccounts()
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.1.7"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetAccountsV1(out)
			}
		}
	},
}

func displayGetAccountsV1(in string) {
	if strings.ToLower(outputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}

	var accounts objects.AccountList

	marshErr := json.Unmarshal([]byte(in), &accounts)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !formatOverridden {
		outputFormat = "text"
	}

	switch strings.ToLower(outputFormat) {

	case "json":
		accounts.ToJSON()
	case "gron":
		accounts.ToGRON()
	case "yaml":
		accounts.ToYAML()
	case "text", "table":
		accounts.ToTEXT(noHeaders)
	}

}

func getAccounts() (string, error) {
	restClient := resty.New()
	// Check if we've set a caBundle (via --ca-cert parameter)
	if len(caBundle) > 0 {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(caBundle))
		if !ok {
			logrus.Info("Failed to parse CABundle")
		}
		restClient.SetTLSClientConfig(&tls.Config{RootCAs: roots})
	}

	uri := "splicectl/v1/cm/accounts"
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Account List Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	getCmd.AddCommand(getAccountsCmd)
}
