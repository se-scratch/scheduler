package root

import (
	"fmt"
	"log"
	"os"
	"scheduler/internal/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int    `mapstructure:"port"`
	DbFile     string `mapstructure:"dbfile"`
}

var rootCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "A scheduler server",
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("port", cmd.Flags().Lookup("port"))
		viper.BindPFlag("dbfile", cmd.Flags().Lookup("dbfile"))
		cfg, err := ReadConfig()
		if err != nil {
			log.Fatal(err)
		}
		if err = server.RunServer(cfg.ServerPort, cfg.DbFile); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Int("port", 7540, "server port")
	rootCmd.Flags().String("dbfile", "scheduler.db", "sqlite database file name")
}

func ReadConfig() (*Config, error) {
	var config Config
	viper.SetConfigType("yaml")
	viper.SetConfigFile(".scheduler.yml")
	viper.SetEnvPrefix("todo")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config file can't read, %s", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("config file has incorrect syntax, %v", err)
	}
	return &config, nil
}
