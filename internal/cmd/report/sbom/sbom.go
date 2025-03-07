package sbom

import (
	"fmt"

	"github.com/debricked/cli/internal/report"
	"github.com/debricked/cli/internal/report/sbom"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commitId string
var repositoryId string
var branch string
var format string
var vulnerabilities bool
var licenses bool
var output string

const CommitFlag = "commit"
const RepositorylFlag = "repository"
const TokenFlag = "token"
const BranchFlag = "branch"
const VulnerabilitiesFlag = "vulnerabilities"
const LicensesFlag = "licenses"
const OutputFlag = "output"
const FormatFlag = "format"

func NewSBOMCmd(reporter report.IReporter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sbom",
		Short: "Generate SBOM export",
		Long: `Generate SBOM export for chosen commit and repository. 
For an example of the SBOM format see https://github.com/debricked/blog-snippets/blob/main/example-sbom-report/SBOM_2022-12-14.json.

This is an enterprise feature. Please visit https://debricked.com/pricing/ for more info.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(reporter),
	}

	cmd.Flags().StringVarP(&commitId, CommitFlag, "c", "", "The commit that you want an SBOM export for")
	_ = cmd.MarkFlagRequired(CommitFlag)
	viper.MustBindEnv(CommitFlag)

	cmd.Flags().StringVarP(&repositoryId, RepositorylFlag, "r", "", "The repository that you want an SBOM export for")
	_ = cmd.MarkFlagRequired(RepositorylFlag)
	viper.MustBindEnv(RepositorylFlag)

	cmd.Flags().StringVarP(&branch, BranchFlag, "b", "", "The branch that you want an SBOM export for")
	viper.MustBindEnv(BranchFlag)

	cmd.Flags().StringVarP(&format, FormatFlag, "f", "", `The format that you want the SBOM export in.

Supported options are: 'CycloneDX', 'SPDX'`,
	)
	viper.MustBindEnv(FormatFlag)

	cmd.Flags().BoolVar(&vulnerabilities, VulnerabilitiesFlag, true, "Toggle SBOM vulnerability data inclusion")
	viper.MustBindEnv(VulnerabilitiesFlag)

	cmd.Flags().BoolVar(&licenses, LicensesFlag, true, "Toggle SBOM license data inclusion")
	viper.MustBindEnv(LicensesFlag)

	cmd.Flags().StringVarP(&output, OutputFlag, "o", "", `Set output path for downloaded SBOM json file.

If no output path is set the file is created in the format <repository_id>-<commit_id>.sbom.json`,
	)
	viper.MustBindEnv(OutputFlag)

	return cmd
}

func RunE(r report.IReporter) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		orderArgs := sbom.OrderArgs{
			RepositoryID:    viper.GetString(RepositorylFlag),
			CommitID:        viper.GetString(CommitFlag),
			Branch:          viper.GetString(BranchFlag),
			Vulnerabilities: viper.GetBool(VulnerabilitiesFlag),
			Licenses:        viper.GetBool(LicensesFlag),
			Output:          viper.GetString(OutputFlag),
			Format:          viper.GetString(FormatFlag),
		}

		if err := r.Order(orderArgs); err != nil {
			return fmt.Errorf("%s %s", color.RedString("⨯"), err.Error())
		}

		return nil
	}
}
