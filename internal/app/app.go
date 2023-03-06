package app

import (
	"DefaultEx1/internal/models"
	dbpostgres "DefaultEx1/pkg/postgres"
	dbredis "DefaultEx1/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func Run() {

	fmt.Println("Приложение запущено")

	//Инициализируем .env
	viper.SetConfigFile("../../internal/cfg/.env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	//По сколько считывать/удалять элементов из Redis
	RedisCountBatch := viper.GetInt64("REDISCOUNTBATCH")

	//Инициализация клиента Redis
	var ClientRedis dbredis.RedisConnect

	ClientRedis.Connect()
	defer ClientRedis.Close()

	//Инициализация клиента Postgres
	var ClientPostgres dbpostgres.ConnectPostgres

	ClientPostgres.Connect()
	defer ClientPostgres.Close()

	//Процесс записи данных в Postgres (через транзакции)
	for {
		//Берем данные из редис
		value, err := ClientRedis.GetItems(RedisCountBatch - 1)
		if err != nil {
			fmt.Println("Ошибка загрузки", err)
		}

		if len(value) == 0 {
			time.Sleep(8 * time.Second)
			continue
		}

		//Открывам транзакцию
		tx, err := ClientPostgres.Client.Begin(context.Background())
		if err != nil {
			fmt.Println("Ошибка создания транзакции", err)
		}

		//Идем циклом по взятому batch
		for _, val := range value {

			//Стректура для данных из Redis
			var Data models.ClientInfo

			json.Unmarshal([]byte(val), &Data)

			//Подготавливем sql запрос
			query := fmt.Sprintf("INSERT INTO clientinfo (name, phone, city, address, region, email) VALUES ('%s', '%s', '%s', '%s', '%s', '%s')",
				Data.Name, Data.Phone, Data.City, Data.Address, Data.Region, Data.Email)

			_, err = tx.Exec(context.Background(), query)
			if err != nil {
				tx.Rollback(context.Background())
				fmt.Println("Ошибка подготовки транзакции (tx.Exec)", err)
			}

		}

		//Удаляем данные из Redis(уже подготовлены для записи в Postgres)
		status, err := ClientRedis.DelItems(RedisCountBatch - 1)
		if err != nil {
			tx.Rollback(context.Background())
			fmt.Println("Ошибка удаления записей из Redis", err)
		} else {
			fmt.Println("Удаление из Redis прошло успешно, status:", status)
		}

		//Записываем данные в Postgres
		err = tx.Commit(context.Background())
		if err != nil {
			tx.Rollback(context.Background())
			fmt.Println("Ошибка записи данных в Postgres", err)
		}
		fmt.Println("Ошибка равна", err, "клиент внесен в базу данных")

	}

}
