package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "taskfuss-cli",
	Short: "Main CLI tool",
}

// Команда database
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database operations",
}

// Команда add
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add resources to database",
}

// Команда user
var addUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Add new user to database",
	Run: func(cmd *cobra.Command, args []string) {
		// Получаем значения флагов
		name, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")

		// Здесь должна быть логика добавления в БД
		fmt.Printf(
			"Добавлен пользователь:\nName: %s\nPassword: %s\nEmail: %s\n",
			name,
			password,
			email,
		)
	},
}

func init() {
	// Иерархия команд
	rootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(addCmd)
	addCmd.AddCommand(addUserCmd)

	// Флаги для команды user
	addUserCmd.Flags().StringP("name", "n", "", "User name (required)")
	addUserCmd.Flags().StringP("password", "p", "", "User password (required)")
	addUserCmd.Flags().StringP("email", "e", "", "User email (required)")

	// Помечаем флаги как обязательные
	addUserCmd.MarkFlagRequired("name")
	addUserCmd.MarkFlagRequired("password")
	addUserCmd.MarkFlagRequired("email")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
