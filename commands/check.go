package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
)

type CheckCommand struct {
	UI cli.Ui
}

func (c *CheckCommand) Run(args []string) int {
	var root, token string

	currentDir, _ := os.Getwd()
	f := flag.NewFlagSet("check", flag.ExitOnError)
	f.StringVar(&token, "token", "", "Terraform Cloud token")
	f.StringVar(&root, "root-path", currentDir, "Terraform config root path (default: current directory)")
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	ws, err := InitCLI(root, token)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	currentVer, err := ws.GetCurrentVersion()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	latestVer, err := ws.GetLatestVersion()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	compatibleVer, err := ws.GetCompatibleLatestVersion()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if currentVer.String() != latestVer.String() {
		c.UI.Warn("New version is available.")
		if compatibleVer.String() != latestVer.String() {
			c.UI.Error(fmt.Sprintf("This version is not compatible with required version %v", ws.GetRequiredVersions().String()))
			c.UI.Info(fmt.Sprintf("%s -> %s (NOT COMPATIBLE)", currentVer.String(), latestVer.String()))
		} else {
			c.UI.Info(fmt.Sprintf("%s -> %s", currentVer.String(), latestVer.String()))
		}
	} else {
		c.UI.Warn("No updates available.")
	}

	return 0
}

func (c *CheckCommand) Help() string {
	return strings.TrimSpace(helpMessageCheck)
}

func (c *CheckCommand) Synopsis() string {
	return "Check if new Terraform version is available"
}

const helpMessageCheck = `
Usage: terraform-cloud-updater check [OPTION]

--token        Terraform Cloud token        (default: TFE_TOKEN env var or parse from your .terraformrc)
--root-path    Terraform config root path   (default: current directory)
`
