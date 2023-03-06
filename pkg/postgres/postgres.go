package dbpostgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

// структура Client
type ConnectPostgres struct {
	Client   *pgx.Conn
	User     string
	Password string
	Url      string
	Port     string
	DBName   string
}

// Метод для коннекта с Postgres
func (c *ConnectPostgres) Connect() error {

	//Инициализируем .env
	viper.SetConfigFile("../../internal/cfg/.env")
	viper.ReadInConfig()

	c.User = viper.GetString("DB_USERNAME")
	c.Password = viper.GetString("DB_PASSWORD")
	c.Url = viper.GetString("DB_HOST")
	c.Port = viper.GetString("DB_PORT")
	c.DBName = viper.GetString("DB_NAME")

	urlDataDase := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Url, c.Port, c.DBName)

	client, err := pgx.Connect(context.Background(), urlDataDase)
	if err != nil {
		fmt.Println("Ошибка подключения к Postgres", err)
	}

	c.Client = client

	return err

}

// Метод для закрытия соденинения с Postgres
func (c *ConnectPostgres) Close() {

	err := c.Client.Close(context.Background())
	if err != nil {
		fmt.Println("Ошибка закрытия соединения Postgres", err)
	}

}
