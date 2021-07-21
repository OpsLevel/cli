package cmd

import (
	"fmt"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCategoriesCmd = &cobra.Command{
	Use:   "categories",
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

var createCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Create a rubric category",
	Long:  `Create a rubric category`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		category, err := client.CreateCategory(opslevel.CategoryCreateInput{
			Name: viper.GetString("name"),
		})
		cobra.CheckErr(err)
		fmt.Println(category.Id)
	},
}

var deleteCategoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Delete a rubric category",
	Long:  `Delete a rubric category`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		err := client.DeleteCategory(viper.GetString("id"))
		cobra.CheckErr(err)
	},
}

var getLevelsCmd = &cobra.Command{
	Use:   "levels",
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

var createLevelCmd = &cobra.Command{
	Use:   "level",
	Short: "Create a rubric level",
	Long:  `Create a rubric level`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		category, err := client.CreateLevel(opslevel.LevelCreateInput{
			Name: viper.GetString("name"),
		})
		cobra.CheckErr(err)
		fmt.Println(category.Id)
	},
}

var deleteLevelCmd = &cobra.Command{
	Use:   "level",
	Short: "Delete a rubric level",
	Long:  `Delete a rubric level`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		err := client.DeleteLevel(viper.GetString("id"))
		cobra.CheckErr(err)
	},
}

func init() {
	getCmd.AddCommand(getCategoriesCmd)
	createCmd.AddCommand(createCategoryCmd)
	deleteCmd.AddCommand(deleteCategoryCmd)

	createCategoryCmd.Flags().StringP("name", "n", "", "the name for the category")
	viper.BindPFlags(createCategoryCmd.Flags())

	deleteCategoryCmd.Flags().String("id", "", "the id for the category")
	viper.BindPFlags(deleteCategoryCmd.Flags())

	getCmd.AddCommand(getLevelsCmd)
	createCmd.AddCommand(createLevelCmd)
	deleteCmd.AddCommand(deleteLevelCmd)

	createLevelCmd.Flags().StringP("name", "n", "", "the name for the category")
	viper.BindPFlags(createLevelCmd.Flags())

	deleteLevelCmd.Flags().String("id", "", "the id for the category")
	viper.BindPFlags(deleteLevelCmd.Flags())
}
