package main

import (
	"fmt"
	"github.com/akshanshgusain/Go-Chi-DynamoDB/config"
	"github.com/akshanshgusain/Go-Chi-DynamoDB/internal/repository/adapter"
	"github.com/akshanshgusain/Go-Chi-DynamoDB/internal/repository/instance"
	"github.com/akshanshgusain/Go-Chi-DynamoDB/internal/routes"
	"github.com/akshanshgusain/Go-Chi-DynamoDB/internal/rules"
	RulesProduct "github.com/akshanshgusain/Go-Chi-DynamoDB/internal/rules/product"
	"github.com/akshanshgusain/Go-Chi-DynamoDB/utils/logger"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"net/http"
)

func main() {
	configs := config.GetConfig()

	dbConnection := instance.GetConnection()
	repository := adapter.NewAdapter(dbConnection)

	logger.INFO("Waiting service starting.... ", nil)

	errors := Migrate(dbConnection)
	if len(errors) > 0 {
		for _, err := range errors {
			logger.PANIC("Error on migrate: ", err)
		}
	}
	logger.PANIC("", checkTables(dbConnection))

	port := fmt.Sprintf(":%v", configs.Port)
	router := routes.NewRouter().SetRouters(repository)
	logger.INFO("Service running on port ", port)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.DynamoDB) []error {
	var errors []error

	callMigrateAndAppendError(&errors, connection, &RulesProduct.Rules{})

	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.DynamoDB, rule rules.Interface) {
	err := rule.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}

func checkTables(connection *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTablesInput{})
	if response != nil {
		if len(response.TableNames) == 0 {
			logger.INFO("Tables not found: ", nil)
		}
		for _, tableName := range response.TableNames {
			logger.INFO("Table found: ", *tableName)
		}
	}
	return err
}
