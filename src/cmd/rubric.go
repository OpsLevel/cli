package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2025"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var exampleCategoryCmd = &cobra.Command{
	Use:     "category",
	Aliases: []string{"cat"},
	Short:   "Example rubric category",
	Long:    `Example rubric category`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample(opslevel.CategoryCreateInput{
			Name: "example_name",
		}))
	},
}

var createCategoryCmd = &cobra.Command{
	Use:        "category NAME",
	Aliases:    []string{"cat"},
	Short:      "Create a rubric category",
	Long:       `Create a rubric category`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"NAME"},
	Run: func(cmd *cobra.Command, args []string) {
		category, err := getClientGQL().CreateCategory(opslevel.CategoryCreateInput{
			Name: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(category.Id)
	},
}

var getCategoryCmd = &cobra.Command{
	Use:        "category ID",
	Aliases:    []string{"cat"},
	Short:      "Get details about a rubric category",
	Long:       `Get details about a rubric category`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		category, err := getClientGQL().GetCategory(opslevel.ID(key))
		cobra.CheckErr(err)
		common.PrettyPrint(category)
	},
}

var listCategoryCmd = &cobra.Command{
	Use:     "category",
	Aliases: []string{"categories", "cat", "cats"},
	Short:   "Lists rubric categories",
	Long:    `Lists rubric categories`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListCategories(nil)
		cobra.CheckErr(err)
		list := resp.Nodes
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Alias(), item.Id)
			}
			w.Flush()
		}
	},
}

var deleteCategoryCmd = &cobra.Command{
	Use:        "category ID",
	Aliases:    []string{"cat"},
	Short:      "Delete a rubric category",
	Long:       `Delete a rubric category`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteCategory(opslevel.ID(key))
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' category\n", key)
	},
}

var exampleLevelCmd = &cobra.Command{
	Use:   "level",
	Short: "Example rubric level",
	Long:  `Example rubric level`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample(opslevel.LevelCreateInput{
			Name: "example_name",
		}))
	},
}

var createLevelCmd = &cobra.Command{
	Use:        "level NAME",
	Short:      "Create a rubric level",
	Long:       `Create a rubric level`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"NAME"},
	Run: func(cmd *cobra.Command, args []string) {
		category, err := getClientGQL().CreateLevel(opslevel.LevelCreateInput{
			Name: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(category.Id)
	},
}

var getLevelCmd = &cobra.Command{
	Use:        "level ID",
	Short:      "Get details about a rubric level",
	Long:       `Get details about a rubric level`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		level, err := getClientGQL().GetLevel(opslevel.ID(key))
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
		resp, err := getClientGQL().ListLevels(nil)
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(resp.Nodes, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "ALIAS", "ID")
			for _, item := range resp.Nodes {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Alias, item.Id)
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
		err := getClientGQL().DeleteLevel(opslevel.ID(key))
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' level\n", key)
	},
}

func init() {
	exampleCmd.AddCommand(exampleCategoryCmd)
	createCmd.AddCommand(createCategoryCmd)
	getCmd.AddCommand(getCategoryCmd)
	listCmd.AddCommand(listCategoryCmd)
	deleteCmd.AddCommand(deleteCategoryCmd)

	exampleCmd.AddCommand(exampleLevelCmd)
	createCmd.AddCommand(createLevelCmd)
	getCmd.AddCommand(getLevelCmd)
	listCmd.AddCommand(listLevelCmd)
	deleteCmd.AddCommand(deleteLevelCmd)
}
