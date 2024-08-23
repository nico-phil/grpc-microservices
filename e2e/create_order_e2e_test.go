package e2e

import (
	"context"
	"log"
	"testing"

	"github.com/nico-phil/microservices-proto/golang/order"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreateOrderTestSuite struct {
	suite.Suite
	compose tc.ComposeStack
}

func (c *CreateOrderTestSuite) SetupSuite(t *testing.T) {
	// composeFilePaths := []string{"resources/docker-compose.yml"}
	// identifier := strings.ToLower(uuid.New().String())

	compose, err := tc.NewDockerCompose("resources/docker-compose.yml")
	require.NoError(t, err, "NewDockerComposeAPI()")
	c.compose = compose
	// execError := compose.
	// 	WithCommand([]string{"up", "-d"}).
	// 	Invoke()
	// err := execError.Error
	// if err != nil {
	// 	log.Fatalf("Could not run compose stack: %v", err)
	// }
	// ctx, cancel := context.WithCancel(context.Background())
    // t.Cleanup(cancel)
	require.NoError(t, compose.Up(context.Background(), tc.Wait(true)), "compose.Up()")
}

func (c *CreateOrderTestSuite) Test_Should_Create_Order() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderClient(conn)
	createOrderResponse, errCreate := orderClient.Create(context.Background(), &order.CreateOrderRequest{
		UserId: 23,
		OrderItems: []*order.OrderItem{
			{
				ProductCode: "CAM123",
				Quantity:    3,
				UnitPrice:   1.23,
			},
		},
	})
	c.Nil(errCreate)

	getOrderResponse, errGet := orderClient.Get(context.Background(), &order.GetOrderRequest{OrderId: createOrderResponse.OrderId})
	c.Nil(errGet)
	c.Equal(int64(23), getOrderResponse.UserId)
	orderItem := getOrderResponse.OrderItems[0]
	c.Equal(float32(1.23), orderItem.UnitPrice)
	c.Equal(int32(3), orderItem.Quantity)
	c.Equal("CAM123", orderItem.ProductCode)
}

func (c *CreateOrderTestSuite) TearDownSuite(t *testing.T) {
	t.Cleanup(func() {
        require.NoError(t, c.compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
    })
}

func TestCreateOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}