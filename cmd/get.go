package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jonathanbeber/sys11ctl/metakube"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validResources = map[string]struct{}{
	"project":  struct{}{},
	"projects": struct{}{},
	// "clusters": struct{}{},
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Long:  `Prints a table of the most important information about the specified resources.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires resources type argument")
		}
		if _, ok := validResources[args[0]]; !ok {
			return fmt.Errorf("found not valid resource type '%s'", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		baseURL := viper.GetString("api_url")
		token := viper.GetString("token")

		if baseURL == "" {
			fmt.Fprint(os.Stderr, "Could not found api_url config on configuration file\n")
			os.Exit(1)
		}
		if token == "" {
			fmt.Fprint(os.Stderr, "Could not found API token config on configuration file\n")
			os.Exit(1)
		}

		client := metakube.NewClient(baseURL, token)
		switch args[0] {
		case "projects":
			projects, err := metakube.GetProjects(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing projects: %s!", err.Error())
				os.Exit(1)
			}
			printProjects(projects)
		case "project":
			projects, err := metakube.GetProjects(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing projects: %s!", err.Error())
				os.Exit(1)
			}
			printProjects(projects)
		default:
			fmt.Fprint(os.Stderr, "Not implemented!\n")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func printProjects(projects []metakube.Project) {
	data := [][]string{}
	for _, project := range projects {
		owners := ""
		for _, owner := range project.Owners {
			owners += owner.Name
		}
		data = append(data, []string{project.ID, project.ID, project.Status, owners})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "NAME", "STATUS", "OWNERS"})

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	table.AppendBulk(data)
	table.Render()
}
