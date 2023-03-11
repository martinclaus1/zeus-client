package restrictedPage

import (
	"github.com/golang/mock/gomock"
	"github.com/playwright-community/playwright-go"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"zeus-client/mocks"
)

func Test_That_NavigateToMyZeusView_Works(t *testing.T) {
	// given
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPage := mocks.NewMockPage(mockCtrl)
	page := playwright.Page(mockPage)
	logger, _ := test.NewNullLogger()
	restrictedPage := Instance(&page, logger)

	gomock.InOrder(
		mockPage.EXPECT().WaitForSelector(restrictedPage.refreshIcon, gomock.Any()).Return(nil, nil).Times(1),
		mockPage.EXPECT().GetAttribute(restrictedPage.refreshIcon, "style").Return("display: none;", nil).Times(1),
		mockPage.EXPECT().Click(restrictedPage.myZeusButton).Return(nil).Times(1),
	)

	// when
	result := restrictedPage.NavigateToMyZeusView()

	// then
	assert.NotNil(t, result)
	assert.IsType(t, &MyZeusView{}, result)
}

func Test_That_GetStatus_Works(t *testing.T) {
	// given
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPage := mocks.NewMockPage(mockCtrl)
	page := playwright.Page(mockPage)
	logger, _ := test.NewNullLogger()
	restrictedPage := Instance(&page, logger)
	myZeusView := instance(restrictedPage)
	testCases := []struct{ zeusStatus, expectedStatus string }{
		{"Abwesend", "abwesend"}, {"Anwesend", "anwesend"},
	}

	for _, testCase := range testCases {
		gomock.InOrder(
			mockPage.EXPECT().Click(myZeusView.refreshButton).Return(nil).Times(1),
			mockPage.EXPECT().WaitForSelector(restrictedPage.refreshIcon, gomock.Any()).Return(nil, nil).Times(1),
			mockPage.EXPECT().GetAttribute(restrictedPage.refreshIcon, "style").Return("display: none;", nil).Times(1),
			mockPage.EXPECT().GetAttribute(myZeusView.presenceStatus, "title").Return(testCase.zeusStatus, nil).Times(1),
		)

		result := myZeusView.GetStatus()

		assert.Equal(t, testCase.expectedStatus, result)
	}
}

func Test_That_ToggleStatus_Works(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPage := mocks.NewMockPage(mockCtrl)
	page := playwright.Page(mockPage)
	logger, _ := test.NewNullLogger()
	restrictedPage := Instance(&page, logger)
	myZeusView := instance(restrictedPage)
	testCases := []struct{ zeusStatusBefore, statusButton, zeusStatusAfter string }{
		{"Abwesend", myZeusView.presentButton, "Anwesend"},
		{"Anwesend", myZeusView.absentButton, "Abwesend"},
	}

	for index, testCase := range testCases {
		t.Logf("Test case %d %s", index, testCase)
		gomock.InOrder(
			mockPage.EXPECT().Click(myZeusView.terminalButton).Return(nil).Times(1),
			// wait for terminal click to finish
			mockPage.EXPECT().WaitForSelector(restrictedPage.refreshIcon, gomock.Any()).Return(nil, nil).Times(1),
			mockPage.EXPECT().GetAttribute(restrictedPage.refreshIcon, "style").Return("display: none;", nil).Times(1),
			// refresh UI and get status
			mockPage.EXPECT().Click(myZeusView.refreshButton).Return(nil).Times(1),
			mockPage.EXPECT().WaitForSelector(restrictedPage.refreshIcon, gomock.Any()).Return(nil, nil).Times(1),
			mockPage.EXPECT().GetAttribute(restrictedPage.refreshIcon, "style").Return("display: none;", nil).Times(1),
			mockPage.EXPECT().GetAttribute(myZeusView.presenceStatus, "title").Return(testCase.zeusStatusBefore, nil).Times(1),
			// toggle status
			mockPage.EXPECT().Click(testCase.statusButton).Return(nil).Times(1),
			// get new status
			mockPage.EXPECT().Click(myZeusView.refreshButton).Return(nil).Times(1),
			mockPage.EXPECT().WaitForSelector(restrictedPage.refreshIcon, gomock.Any()).Return(nil, nil).Times(1),
			mockPage.EXPECT().GetAttribute(restrictedPage.refreshIcon, "style").Return("display: none;", nil).Times(1),
			mockPage.EXPECT().GetAttribute(myZeusView.presenceStatus, "title").Return(testCase.zeusStatusAfter, nil).Times(1),
		)

		myZeusView.ToggleStatus()
	}
}
