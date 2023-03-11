package loginPage

import (
	"github.com/golang/mock/gomock"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"zeus-client/mocks"
	"zeus-client/restrictedPage"
)

func Test_That_Login_Executes_The_Commands_In_The_Correct_Order(t *testing.T) {
	// given
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPage := mocks.NewMockPage(mockCtrl)
	page := playwright.Page(mockPage)
	logger, _ := test.NewNullLogger()
	login := Instance(&page, logger)

	gomock.InOrder(
		mockPage.EXPECT().Goto(login.url).Return(nil, nil).Times(1),
		mockPage.EXPECT().Type(login.usernameInput, "username").Return(nil).Times(1),
		mockPage.EXPECT().Type(login.passwordInput, "password").Return(nil).Times(1),
		mockPage.EXPECT().Click(login.loginButton).Return(nil).Times(1),
	)

	// when
	result := login.Login("username", "password")

	// then
	assert.NotNil(t, result)
	assert.IsType(t, &restrictedPage.RestrictedPage{}, result)
}
