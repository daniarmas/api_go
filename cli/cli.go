package cli

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/daniarmas/api_go/datasource"
	"github.com/daniarmas/api_go/models"
	"github.com/daniarmas/api_go/repository"
	"github.com/daniarmas/api_go/utils"
	"gorm.io/gorm"
)

const (
	unknownCommand     = "unknown command\nRun './main help' for usage."
	commandHelpUnknown = `unknown help topic. Run './main help'.`
	migrateHelp        = `usage: ./main migrate`
	createAppHelp      = `usage: ./main createapp`
	help               = `This command-line tool is for administrative tasks

Usage:

	./main <command> [arguments]

The commands are:

	createapp         create app for consume the api

Use "./main help <command>" for more information about a command.
`
	createApp = "createapp"
)

func HandleCli(args []string, db *gorm.DB, config *utils.Config, repositoryDao repository.DAO) {
	if len(args) != 1 {
		createAppCmd := flag.NewFlagSet("createapp", flag.ExitOnError)
		flag.NewFlagSet("--config", flag.ExitOnError)
		helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
		getHelpCommand := helpCmd.String("command", "", "Command")
		switch args[1] {
		case "help":
			if len(args) == 2 {
				fmt.Println(help)
				os.Exit(0)
			} else {
				handleHelp(helpCmd, getHelpCommand)
			}
		case "createapp":
			handleCreateApp(createAppCmd, repositoryDao, db)
		case "--config":
		default:
			fmt.Println(unknownCommand)
			os.Exit(0)
		}
	}
}

func handleCreateApp(getCreateAdmin *flag.FlagSet, repositoryDao repository.DAO, db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		expTime := time.Now().UTC().AddDate(0, 0, 1)
		createAppRes, createAppErr := repositoryDao.NewApplicationRepository().CreateApplication(tx, &models.Application{Name: "Initial", Version: "1.0.0", ExpirationTime: expTime})
		if createAppErr != nil {
			log.Error(createAppErr)
			os.Exit(0)
		}
		jwtAcessToken := datasource.JsonWebTokenMetadata{TokenId: createAppRes.ID}
		jwtAcessTokenErr := repository.Datasource.NewJwtTokenDatasource().CreateJwtAuthorizationToken(&jwtAcessToken)
		if jwtAcessTokenErr != nil {
			log.Error(jwtAcessTokenErr)
			os.Exit(0)
		}
		log.Info("AccessToken: ", *jwtAcessToken.Token)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func handleHelp(getHelp *flag.FlagSet, command *string) {
	getHelp.Parse(os.Args[2:])
	switch os.Args[2] {
	case "createapp":
		fmt.Println(createAppHelp)
		os.Exit(0)
	default:
		fmt.Println(commandHelpUnknown)
		os.Exit(0)
	}
}
