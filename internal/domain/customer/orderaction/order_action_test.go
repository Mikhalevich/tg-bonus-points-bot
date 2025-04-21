package orderaction_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/customer/orderaction"
	"github.com/Mikhalevich/tg-bonus-points-bot/internal/domain/port"
)

type OrderActionSuite struct {
	*suite.Suite
	ctrl *gomock.Controller

	sender     *port.MockMessageSender
	repository *port.MockCustomerOrderActionRepository

	orderAction *orderaction.OrderAction
}

func TestProcessorSuit(t *testing.T) {
	t.Parallel()
	suite.Run(t, &OrderActionSuite{
		Suite: new(suite.Suite),
	})
}

func (s *OrderActionSuite) SetupSuite() {
	s.ctrl = gomock.NewController(s.T())

	s.sender = port.NewMockMessageSender(s.ctrl)
	s.repository = port.NewMockCustomerOrderActionRepository(s.ctrl)

	s.orderAction = orderaction.New(s.sender, s.repository, nil, nil)
}

func (s *OrderActionSuite) TearDownSuite() {
}

func (s *OrderActionSuite) TearDownTest() {
}

func (s *OrderActionSuite) TearDownSubTest() {
}
