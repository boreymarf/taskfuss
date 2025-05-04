// TODO: REWORK THIS LATER OK?
// Use separate project with cobra-cli as a reference

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "taskfuss-cli",
	Short: "Main CLI tool",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(); err != nil {
			logger.Log.Fatal().Err(err).Msg("Error loading .env file")
		}

		// If not verbose, only log errors
		if verbose {
			os.Setenv("LOG_LEVEL", "debug")
		} else {
			os.Setenv("LOG_LEVEL", "error")
		}

		// Changing level of global logger
		if level := os.Getenv("LOG_LEVEL"); level != "" {
			if logLevel, err := zerolog.ParseLevel(level); err == nil {
				logger.Log = logger.Log.Level(logLevel)
			}
		}
	},
}

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database operations",
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add resources to database",
}

var addUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Add new user to database",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")

		logger.Log.Info().Msg("Starting server...")

		database, err := db.InitDB()
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to connect to the database")
		} else {
			logger.Log.Info().Msg("Connected to the database successfully!")
		}
		defer database.Close()

		userRepository, err := db.InitUserRepository(database)
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to create userRepository")
		}

		user := db.User{}
		user.Name = name
		user.PasswordHash = password
		user.Email = email

		err = userRepository.CreateUser(&user)

		if err != nil {
			logger.Log.Fatal().
				Err(err).
				Str("name", name).
				Str("email", email).
				Str("password", user.PasswordHash).
				Msg("Failed to create user")
		} else {
			logger.Log.Info().
				Str("name", name).
				Str("email", email).
				Str("password", user.PasswordHash).
				Str("created_at", user.CreatedAt.String()).
				Msg("Successfully added a new user")
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources from database",
}

var getUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get user from database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]

		database, err := db.InitDB()
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to connect to the database")
		} else {
			logger.Log.Info().Msg("Connected to the database successfully!")
		}
		defer database.Close()

		userRepository, err := db.InitUserRepository(database)
		if err != nil {
			logger.Log.Fatal().Err(err).Msg("Failed to create userRepository")
		}

		if input == "all" {
			users, err := userRepository.GetAllUsers()
			if err != nil {
				logger.Log.Fatal().Err(err).Send()
			}

			jsonData, err := json.Marshal(users)
			if err != nil {
				logger.Log.Fatal().
					Err(err).
					Msg("Serialization error")
			}

			fmt.Println(string(jsonData))
			os.Exit(0)
		}

		id, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			logger.Log.Fatal().Msg("Incorrect input, must be 'all' or numeric ID")
		}

		user := db.User{}
		err = userRepository.GetUserByID(id, &user)
		if err != nil {
			logger.Log.Fatal().
				Err(err).
				Int64("id", id).
				Msg("Failed to get User")
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			logger.Log.Fatal().
				Err(err).
				Msg("Serialization error")
		}

		fmt.Println(string(jsonData))
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)

	// addCmd
	databaseCmd.AddCommand(addCmd)
	addCmd.AddCommand(addUserCmd)

	addUserCmd.Flags().StringP("name", "n", "", "User name (required)")
	addUserCmd.Flags().StringP("password", "p", "", "User password (required)")
	addUserCmd.Flags().StringP("email", "e", "", "User email (required)")

	addUserCmd.MarkFlagRequired("name")
	addUserCmd.MarkFlagRequired("password")
	addUserCmd.MarkFlagRequired("email")

	// getCmd
	databaseCmd.AddCommand(getCmd)
	getCmd.AddCommand(getUserCmd)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
