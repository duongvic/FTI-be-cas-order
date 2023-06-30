package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm/clause"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/spf13/viper"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"casorder/cmd"
	"casorder/db"
	"casorder/utils/mgrpc"

	"casorder/api/v1/common/services"
	"casorder/db/migrations"
	"casorder/db/models"
)

type Command struct {
}

func (c *Command) MigrateDB() {
	DB := db.GetDB()

	fmt.Printf("Migrating database %v", DB.Migrator().CurrentDatabase())

	migrations.Upgrate(DB)
}

func (c *Command) SeedDB() {
	DB := db.GetDB()

	fmt.Printf("Seeding database %v", DB.Migrator().CurrentDatabase())
	migrations.Seed(DB)
}

func (c *Command) TestDB() {
	DB := db.GetDB()
	var user models.User
	var users []*models.User
	var serv services.CommonApi
	var pagination models.Pagination

	result, err := serv.List(user, &users, &pagination, DB)
	if err != nil {
		fmt.Printf("Error getting users: %v", err)
		return
	}

	j, _ := json.Marshal(users)

	fmt.Printf("%v %v", result, string(j))
}

func (c *Command) DowngrateDB() {
	fmt.Println("Dowgrating database")
}

func (c *Command) StartGrpc() {
	host := viper.GetString("grpc.host")
	port := viper.GetString("grpc.port")
	grpcServer := mgrpc.New(host, port)
	grpcServer.Start()
}

func (c *Command) TestGrpc() {
	fmt.Println("Test Grpc Server")
	grpcAddr := fmt.Sprintf("%v:%v", viper.GetString("grpc.host"), viper.GetString("grpc.port"))
	fmt.Printf("GRPC Server Address: %v\n", grpcAddr)
	mgrpc.TestClient(grpcAddr)
}

func (c *Command) TestVerifyUser() {
	fmt.Println("Trying to verify token...")
	mgrpc.VerifyToken("129bb35564b5ad76c47cd9f832133502")
}

func (c *Command) TestApproval() {
	fmt.Println("Testing Order Approval...\n")
	var testOrder *models.Order
	DB := db.GetDB()
	if err := DB.Model(&testOrder).Preload(clause.Associations).First(&testOrder, "id = ?", 1).Error; err != nil { fmt.Println("Test failed")}
	approved, err := mgrpc.GetApproval(testOrder, DB)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n", approved)
}

func main() {
	cmd.Initialize()

	app := kingpin.New("CasOrder", "My Order application.")
	app.Version("1.0.0")

	dbMigrate := app.Command("db-migrate", "Migrate database")
	dbDowngrate := app.Command("db-downgrate", "Downgrate database")
	dbSeed := app.Command("db-seed", "Seed database")
	startGrpc := app.Command("grpc-start", "Start GRPC server")
	testGrpc := app.Command("test-grpc", "Test grpc server")
	testDB := app.Command("test-db", "Test DB")
	testVerifyUser := app.Command("test-verify", "Test Verify User")
	testOrderApproval := app.Command("test-approve", "Test Order Approval")

	checkConfigCmd := app.Command("config", "Check if the config files are valid or not.")
	configFiles := checkConfigCmd.Arg(
		"config-files",
		"The config files to check.",
	).Required().ExistingFiles()

	command := Command{}

	parsedCmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch parsedCmd {
	case checkConfigCmd.FullCommand():
		if *configFiles == nil {
			os.Exit(0)
		}
	case dbMigrate.FullCommand():
		command.MigrateDB()
	case dbDowngrate.FullCommand():
		command.DowngrateDB()
	case dbSeed.FullCommand():
		command.SeedDB()
	case startGrpc.FullCommand():
		command.StartGrpc()
	case testGrpc.FullCommand():
		command.TestGrpc()
	case testDB.FullCommand():
		command.TestDB()
	case testVerifyUser.FullCommand():
		command.TestVerifyUser()
	case testOrderApproval.FullCommand():
		command.TestApproval()
	}
}
