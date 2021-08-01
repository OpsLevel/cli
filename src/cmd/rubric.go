package cmd

import (
	"fmt"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"
	"github.com/spf13/cobra"
)

var createCategoryCmd = &cobra.Command{
	Use:        "category [name]",
	Short:      "Create a rubric category",
	Long:       `Create a rubric category`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"name"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		category, err := client.CreateCategory(opslevel.CategoryCreateInput{
			Name: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(category.Id)
	},
}

var getCategoryCmd = &cobra.Command{
	Use:        "category [id]",
	Short:      "Get details about a rubic category given its ID",
	Long:       `Get details about a rubic category given its ID`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"id"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		category, err := client.GetCategory(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(category)
	},
}

var listCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Lists the valid names for rubric categories",
	Long:  `Lists the valid names for rubric categories`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		list, err := client.ListCategories()
		cobra.CheckErr(err)
		w := common.NewTabWriter("NAME", "ID")
		if err == nil {
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Name, item.Id)
			}
		}
		w.Flush()
	},
}

var deleteCategoryCmd = &cobra.Command{
	Use:        "category [id]",
	Short:      "Delete a rubric category",
	Long:       `Delete a rubric category`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"id"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		err := client.DeleteCategory(args[0])
		cobra.CheckErr(err)
	},
}

var createLevelCmd = &cobra.Command{
	Use:        "level [name]",
	Short:      "Create a rubric level",
	Long:       `Create a rubric level`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"name"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		category, err := client.CreateLevel(opslevel.LevelCreateInput{
			Name: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(category.Id)
	},
}

var getLevelCmd = &cobra.Command{
	Use:        "level [id]",
	Short:      "Get details about a rubic level given its ID",
	Long:       `Get details about a rubic level given its ID`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"id"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		level, err := client.GetLevel(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(level)
	},
}

var listLevelCmd = &cobra.Command{
	Use:   "level",
	Short: "Lists the valid alias for rubric levels",
	Long:  `Lists the valid alias for rubric levels`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		list, err := client.ListLevels()
		cobra.CheckErr(err)
		w := common.NewTabWriter("Alias", "ID")
		if err == nil {
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Alias, item.Id)
			}
		}
		w.Flush()
	},
}

var deleteLevelCmd = &cobra.Command{
	Use:        "level [id]",
	Short:      "Delete a rubric level",
	Long:       `Delete a rubric level`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"id"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		err := client.DeleteLevel(args[0])
		cobra.CheckErr(err)
	},
}

func init() {
	createCmd.AddCommand(createCategoryCmd)
	getCmd.AddCommand(getCategoryCmd)
	listCmd.AddCommand(listCategoryCmd)
	deleteCmd.AddCommand(deleteCategoryCmd)

	createCmd.AddCommand(createLevelCmd)
	getCmd.AddCommand(getLevelCmd)
	listCmd.AddCommand(listLevelCmd)
	deleteCmd.AddCommand(deleteLevelCmd)
}
