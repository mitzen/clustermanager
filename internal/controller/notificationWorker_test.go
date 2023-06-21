package controller

import (
	"testing"

	"cdx.foc/clusterwatch/mocks"
	"github.com/golang/mock/gomock"
)

func Send_Test(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockMessenger := mocks.NewMockMessenger(ctl)
	gomock.InOrder(
		mockMessenger.EXPECT().SendMessage("test").Return(1),
	)

	nw := NewNotificationWorker(mockMessenger)
	nw.SendMessage("test message")
}
