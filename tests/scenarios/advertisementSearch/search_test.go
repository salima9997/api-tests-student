package advertisementSearch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

const (
	searchBeginning = "Начало"
	searchMiddle    = "середина"
	searchEnd       = "конец"
	searchMixed     = "НаЧаЛо СеРеДиНа КОНЕЦ"
)

type TestSuite struct {
	suite.Suite
	loginData      models.LoginOkResponse
	createdAdID    string
	createdAdTitle string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя")
	s.loginData = auth.Login(
		s.T(),
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	s.createdAdTitle = fmt.Sprintf(
		"Начало середина конец %d",
		time.Now().UnixNano(),
	)

	photoPath := filepath.Join(
		"..",
		"advertisementCreation",
		"testdata",
		"photo.png",
	)

	created := advList.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		s.createdAdTitle,
		"Объявление для тестирования поиска",
		100,
		1,
		photoPath,
	)

	require.NotEmpty(s.T(), created.ID, "Не получен ID созданного объявления")
	s.createdAdID = created.ID
}

func (s *TestSuite) TearDownSuite() {
	if s.createdAdID != "" {
		advList.DeleteAdvertisement(
			s.T(),
			s.loginData.Token,
			s.createdAdID,
		)
	}

	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle() {
	list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.createdAdTitle)
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", searchMixed},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.createdAdTitle)
		})
	}
}
