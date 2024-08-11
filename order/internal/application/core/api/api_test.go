package api

import (
	"context"
	"errors"
	"testing"

	"github.com/nico-phil/grpc-microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockedPayment struct {
	mock.Mock
}

func(p *mockedPayment) Charge(ctx context.Context, order *domain.Order) error{
	args := p.Called(ctx, order)
	return args.Error(0)
}

type mockedDb struct {
	mock.Mock
}

func(d *mockedDb) Save(ctx context.Context, order *domain.Order) error {
	args := d.Called(ctx, order)
	return args.Error(0)
}

func(d *mockedDb) Get(ctx context.Context, id int64) (domain.Order, error) {
	args := d.Called(ctx, id)
	return args.Get(0).(domain.Order), args.Error(1)
}

func TestPlaceOder(t *testing.T) {
	payment := new(mockedPayment)
	db := new(mockedDb)

	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{ ProductCode: "camera",  UnitPrice: 12.3, Quantity: 2 },
		},
	})

	assert.Nil(t, err)
}

func Test_Should_Return_Error_When_Db_Persistent_Fail(t *testing.T){
	payment := new(mockedPayment)
	db := new(mockedDb)

	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(errors.New("connection error"))

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{ 
				ProductCode: "camera",  
				UnitPrice: 12.3, 
				Quantity: 2,
			},
		},
	})

	assert.EqualError(t, err, "connection error")
}

func Test_Should_Return_error_When_Payment_service_Fail(t *testing.T){
	payment := new(mockedPayment)
	db := new(mockedDb)

	payment.On("Charge", mock.Anything, mock.Anything).Return(errors.New("insufficient balance"))
	db.On("Save", mock.Anything, mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{ 
				ProductCode: "camera",  
				UnitPrice: 12.3, 
				Quantity: 2,
			},
		},
	})

	st, _ := status.FromError(err)
	assert.Equal(t, "order creation failed", st.Message())
	assert.Equal(t, st.Details()[0].(*errdetails.BadRequest).FieldViolations[0].Description, "insufficient balance")
	assert.Equal(t, codes.InvalidArgument, st.Code())

}