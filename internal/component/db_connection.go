package component

import (
	"database/sql"
	"e-wallet/internal/config"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/lib/pq"
)

func GetDatabaseConnectiom(cnf *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s "+
			"port=%s "+
			"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"sslmode=disable",
		cnf.Database.Host,
		cnf.Database.Port,
		cnf.Database.User,
		cnf.Database.Password,
		cnf.Database.Name)

	conncetion, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("unbale to connect database %s", err.Error())
	}

	err = conncetion.Ping()
	if err != nil {
		log.Fatalf("unbale to connect database %s", err.Error())
	}
	return conncetion
}
