package dbredis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// Интерфейс для работы с несколькими Redis
type IRedis interface {
	GetItem(key string) string
	GetItems(key string, indexStart int64, indexStop int64) (datar []string, err error)
	DelItem(key string)
	Connect() error
	Close() error
}

// Redis Client
type RedisConnect struct {
	Client   *redis.Client
	Url      string
	Port     string
	Password string
}

// Метод для получения одного элемента из Redis[key](метод удаляет елемент и возвращает его)
func (r *RedisConnect) GetItem(key string) string {

	result := r.Client.LPop(key).Val()

	return result
}

// Метод для получение нескольких элементов из Redis[key]
func (r *RedisConnect) GetItems(RedisCountBatch int64) (result []string, err error) {

	//Инициализируем .env
	viper.SetConfigFile("../../internal/cfg/.env")
	viper.ReadInConfig()

	key := viper.GetString("REDISKEY")
	result, err = r.Client.LRange(key, 0, RedisCountBatch).Result()
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}
	return

}

// Метод для удаления нескольких элементов из Redis[key]
func (r *RedisConnect) DelItems(indexStart int64) (status string, err error) {

	//Инициализируем .env
	viper.SetConfigFile("../../internal/cfg/.env")
	viper.ReadInConfig()

	key := viper.GetString("REDISKEY")
	status, err = r.Client.LTrim(key, indexStart, -1).Result()
	if err != nil {
		fmt.Println(err)
	}

	return
}

// Метод для закрытия соединения Redis
func (r *RedisConnect) Close() (err error) {

	r.Client.Close()

	return

}

// Метод для соединения с Redis
func (r *RedisConnect) Connect() error {

	r.Url = viper.GetString("REDIS_DB_HOST")
	r.Port = viper.GetString("REDIS_PORT")
	r.Password = viper.GetString("REDIS_PASSWORD")

	connect := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Url, r.Port),
		Password: r.Password,
		DB:       0,
	})
	_, err := connect.Ping().Result()
	if err != nil {

		return fmt.Errorf(err.Error())
	}

	r.Client = connect

	return err
}
