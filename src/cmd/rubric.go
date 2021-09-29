package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"
	"github.com/spf13/cobra"
)

var createCategoryCmd = &cobra.Command{
	Use:        "category NAME",
	Short:      "Create a rubric category",
	Long:       `Create a rubric category`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"NAME"},
	Run: func(cmd *cobra.Command, args []string) {
		category, err := graphqlClient.CreateCategory(opslevel.CategoryCreateInput{
			Name: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(category.Id)
	},
}

var getCategoryCmd = &cobra.Command{
	Use:        "category ID",
	Short:      "Get details about a rubic category",
	Long:       `Get details about a rubic category`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		category, err := graphqlClient.GetCategory(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(category)
	},
}

var listCategoryCmd = &cobra.Command{
	Use:     "category",
	Aliases: []string{"categories"},
	Short:   "Lists rubric categories",
	Long:    `Lists rubric categories`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := graphqlClient.ListCategories()
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Name, item.Id)
			}
			w.Flush()
		}
	},
}

var deleteCategoryCmd = &cobra.Command{
	Use:        "category ID",
	Short:      "Delete a rubric category",
	Long:       `Delete a rubric category`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := graphqlClient.DeleteCategory(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' category\n", key)
	},
}

var createLevelCmd = &cobra.Command{
	Use:        "level NAME",
	Short:      "Create a rubric level",
	Long:       `Create a rubric level`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"NAME"},
	Run: func(cmd *cobra.Command, args []string) {
		category, err := graphqlClient.CreateLevel(opslevel.LevelCreateInput{
			Name: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(category.Id)
	},
}

var getLevelCmd = &cobra.Command{
	Use:        "level ID",
	Short:      "Get details about a rubic level",
	Long:       `Get details about a rubic level`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		level, err := graphqlClient.GetLevel(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(level)
	},
}

var listLevelCmd = &cobra.Command{
	Use:     "level",
	Aliases: []string{"levels"},
	Short:   "Lists rubric levels",
	Long:    `Lists rubric levels`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := graphqlClient.ListLevels()
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("Alias", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Alias, item.Id)
			}
			w.Flush()
		}
	},
}

var deleteLevelCmd = &cobra.Command{
	Use:        "level ID",
	Short:      "Delete a rubric level",
	Long:       `Delete a rubric level`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := graphqlClient.DeleteLevel(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' level\n", key)
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
