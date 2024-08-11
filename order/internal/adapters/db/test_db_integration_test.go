package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/nico-phil/grpc-microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type OrderDatabaseTestSuite struct {
	suite.Suite
	DataSourceUrl string
}

func (o *OrderDatabaseTestSuite) SetupSuite() {
	ctx := context.Background()
	port := "3306/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("root:s3cr3t@tcp(localhost:%s)/orders?charset=utf8mb4&parseTime=True&loc=Local", port.Port())
	}
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:8.0.30",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "s3cr3t",
			"MYSQL_DATABASE":      "orders",
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "mysql", dbURL),
	}
	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal("Failed to start Mysql.", err)
	}
	endpoint, _ := mysqlContainer.Endpoint(ctx, "")
	o.DataSourceUrl = fmt.Sprintf("root:s3cr3t@tcp(%s)/orders?charset=utf8mb4&parseTime=True&loc=Local", endpoint)
}

func (o *OrderDatabaseTestSuite) Test_Should_Save_Order() {
	adapter, err := NewAdapter(o.DataSourceUrl)
	o.Nil(err)
	saveErr := adapter.Save(context.Background(),&domain.Order{})
	o.Nil(saveErr)
}

func (o *OrderDatabaseTestSuite) Test_Should_Get_Order() {
	adapter, _ := NewAdapter(o.DataSourceUrl)
	order := domain.NewOder(2, []domain.OrderItem{
		{
			ProductCode: "CAM",
			Quantity:    5,
			UnitPrice:   1.32,
		},
	})

	ctx := context.Background()
	adapter.Save(ctx,  &order)
	ord, _ := adapter.Get(ctx, order.ID)
	o.Equal(int64(2), ord.CustomerID)
}

func TestOrderDatabaseTestSuite(t *testing.T) {
	suite.Run(t, &OrderDatabaseTestSuite{})
}