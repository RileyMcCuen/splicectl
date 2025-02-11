package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/common"
)

var versionsDefaultCRCmd = &cobra.Command{
	Use:   "default-cr",
	Short: "Retrieve the versions of the default CR in the cluster.",
	Long: `EXAMPLES
	splicectl versions default-cr
`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = versionDetail.RequirementMet("versions_default-cr")

		out, err := getDefaultCRVersions()
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayVersionsDefaultCRV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayVersionsDefaultCRV2(out)
			}
		}
	},
}

func displayVersionsDefaultCRV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayVersionsDefaultCRV2(in string) {
	if strings.ToLower(outputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	crData, cerr := common.RestructureVersions(in)
	if cerr != nil {
		logrus.Fatal("Vault Version JSON conversion failed.")
	}

	if !formatOverridden {
		outputFormat = "text"
	}

	switch strings.ToLower(outputFormat) {
	case "json":
		crData.ToJSON()
	case "gron":
		crData.ToGRON()
	case "yaml":
		crData.ToYAML()
	case "text", "table":
		crData.ToTEXT(noHeaders)
	}

}

func getDefaultCRVersions() (string, error) {
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

	uri := "splicectl/v1/vault/defaultcrversions"
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	versionsCmd.AddCommand(versionsDefaultCRCmd)

}
