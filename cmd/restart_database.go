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
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var restartDatabaseCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"database"},
	Short:   "Restart a specific database in the Splice DB Cluster.",
	Long: `EXAMPLES
	splicectl list workspace
	splicectl restart workspace --database-name splicedb

	Note: --database-name and -d are the preferred way to supply the database name.
	However, --database and --workspace can also be used as well. In the event that
	more than one of them is supplied database-name and d are preferred over all
	and workspace is preferred over database. The most preferred option that is
	supplied will be used and a message will be displayed letting you know which
	option was chosen if more than one were supplied.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		var sv semver.Version

		_, sv = versionDetail.RequirementMet("restart_database")

		databaseName := common.DatabaseName(cmd)
		forceRestart, _ := cmd.Flags().GetBool("force")
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}
		out, err := restartDatabase(databaseName, forceRestart)
		if err != nil {
			logrus.WithError(err).Error("Error restarting database")
		}

		if semverV1, err := semver.ParseRange(">=0.1.6"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayRestartDatabaseV1(out)
			}
		}
	},
}

func displayRestartDatabaseV1(in string) {
	if strings.ToLower(outputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var asData objects.ActionStatus
	marshErr := json.Unmarshal([]byte(in), &asData)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !formatOverridden {
		outputFormat = "text"
	}

	switch strings.ToLower(outputFormat) {
	case "json":
		asData.ToJSON()
	case "gron":
		asData.ToGRON()
	case "yaml":
		asData.ToYAML()
	case "text", "table":
		asData.ToTEXT(noHeaders)
	}
}

func restartDatabase(dbname string, force bool) (string, error) {
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

	uri := fmt.Sprintf("splicectl/v1/splicedb/splicedatabaserestart?database-name=%s&force=%t", dbname, force)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error restarting the database")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	restartCmd.AddCommand(restartDatabaseCmd)

	// add database name and aliases
	restartDatabaseCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	restartDatabaseCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	restartDatabaseCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	restartDatabaseCmd.Flags().BoolP("force", "f", false, "Force the restart")

}
