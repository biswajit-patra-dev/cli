package report

import (
	"github.com/debricked/cli/internal/cmd/report/license"
	"github.com/debricked/cli/internal/cmd/report/sbom"
	"github.com/debricked/cli/internal/cmd/report/vulnerability"
	licenseReport "github.com/debricked/cli/internal/report/license"
	sbomReport "github.com/debricked/cli/internal/report/sbom"
	vulnerabilityReport "github.com/debricked/cli/internal/report/vulnerability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewReportCmd(
	licenseReporter licenseReport.Reporter,
	vulnerabilityReporter vulnerabilityReport.Reporter,
	sbomReporter sbomReport.Reporter,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Generate exports for vulnerabilities, licenses, and SBOM.",
		Long: `Generate exports.
Premium is required for license and vulnerability exports. Enterprise is required for SBOM exports. Please visit https://debricked.com/pricing/ for more info.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	cmd.AddCommand(license.NewLicenseCmd(licenseReporter))
	cmd.AddCommand(vulnerability.NewVulnerabilityCmd(vulnerabilityReporter))
	cmd.AddCommand(sbom.NewSBOMCmd(sbomReporter))

	return cmd
}
