package projects_test

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/fylerx/fyler/internal/projects"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type Suite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	repo    projects.Repository
	project *projects.Project
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.repo = projects.InitRepo(s.DB)

	now := time.Now()
	s.project = &projects.Project{
		ID:        999,
		Name:      "new project",
		APIKey:    "apikey",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetAll() {
	s.mock.MatchExpectationsInOrder(false)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "projects"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "api_key", "created_at", "updated_at"}).
			AddRow(1, "first", "apikey", nil, nil).
			AddRow(2, "second", "apikey", nil, nil))

	res, err := s.repo.GetAll()

	expectedProjects := []*projects.Project{
		{
			ID:     1,
			Name:   "first",
			APIKey: "apikey",
		},
		{
			ID:     2,
			Name:   "second",
			APIKey: "apikey",
		},
	}
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(expectedProjects, res))
}

func (s *Suite) TestGetAllError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "projects"`)).
		WillReturnError(gorm.ErrInvalidData)

	res, err := s.repo.GetAll()

	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *Suite) TestGetByID() {
	var ID uint32 = 10
	expProject := &projects.Project{ID: ID, Name: "new project", APIKey: "apikey"}

	project := sqlmock.NewRows([]string{"id", "name", "api_key"}).
		AddRow(expProject.ID, expProject.Name, expProject.APIKey)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "projects" WHERE "projects"."id" = $1`)).
		WithArgs(ID).
		WillReturnRows(project)

	res, err := s.repo.GetByID(ID)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(expProject, res))
}

func (s *Suite) TestGetByIDError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "projects" WHERE "projects"."id" = $1`)).
		WithArgs(s.project.ID).
		WillReturnError(gorm.ErrRecordNotFound)

	res, err := s.repo.GetByID(s.project.ID)

	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *Suite) TestGetByAPIKey() {
	expProject := &projects.Project{ID: 999, Name: "new project", APIKey: "apikey"}

	project := sqlmock.NewRows([]string{"id", "name", "api_key"}).
		AddRow(expProject.ID, expProject.Name, expProject.APIKey)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "projects" WHERE api_key = $1`)).
		WithArgs(expProject.APIKey).
		WillReturnRows(project)

	res, err := s.repo.GetByAPIKey(expProject.APIKey)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(expProject, res))
}

func (s *Suite) TestGetByAPIKeyError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "projects" WHERE api_key = $1`)).
		WithArgs(s.project.APIKey).
		WillReturnError(gorm.ErrRecordNotFound)

	res, err := s.repo.GetByAPIKey(s.project.APIKey)

	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *Suite) TestCreate() {
	var newProjectID uint32 = 123
	now := time.Now()
	input := &projects.Project{Name: "Main Project", CreatedAt: now, UpdatedAt: now}
	expProject := &projects.Project{
		Name:      input.Name,
		ID:        newProjectID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "projects"
		("name","created_at","updated_at")
		VALUES ($1,$2,$3) RETURNING "id"`),
	).
		WithArgs(input.Name, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newProjectID))
	s.mock.ExpectCommit()

	res, err := s.repo.Create(input)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(expProject, res))
}

// func (s *Suite) TestCreateError() {
// 	input := &projects.Project{Name: "Bad Project"}

// 	s.mock.ExpectBegin()
// 	s.mock.ExpectQuery(
// 		regexp.QuoteMeta(`INSERT INTO "projects" ("name") VALUES ($1) RETURNING "id"`),
// 	).
// 		WithArgs(input.Name).
// 		WillReturnError(gorm.ErrInvalidData)
// 	s.mock.ExpectRollback()

// 	res, err := s.repo.Create(input)

// 	require.Error(s.T(), err)
// 	require.Nil(s.T(), res)
// }

func (s *Suite) TestUpdate() {
	project := sqlmock.NewRows([]string{"id", "name", "api_key"}).
		AddRow(s.project.ID, s.project.Name, s.project.APIKey)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "projects" WHERE "projects"."id" = $1`)).
		WithArgs(s.project.ID).
		WillReturnRows(project)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "projects" SET "name"=$1,"updated_at"=$2 WHERE id = $3`),
	).
		WithArgs(s.project.Name, AnyTime{}, s.project.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	res, err := s.repo.Update(s.project.ID, s.project.Name)

	require.NoError(s.T(), err)
	require.True(s.T(), res)
}

func (s *Suite) TestUpdateError() {
	project := sqlmock.NewRows([]string{"id", "name", "api_key"}).
		AddRow(s.project.ID, s.project.Name, s.project.APIKey)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "projects" WHERE "projects"."id" = $1`)).
		WithArgs(s.project.ID).
		WillReturnRows(project)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "projects" SET "name"=$1,"updated_at"=$2 WHERE id = $3`),
	).
		WithArgs(s.project.Name, AnyTime{}, s.project.ID).
		WillReturnError(gorm.ErrInvalidData)
	s.mock.ExpectRollback()

	res, err := s.repo.Update(s.project.ID, s.project.Name)

	require.Error(s.T(), err)
	require.False(s.T(), res)
}

func (s *Suite) TestDelete() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "projects"`).
		WithArgs(s.project.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	res, err := s.repo.Delete(s.project.ID)

	require.NoError(s.T(), err)
	require.True(s.T(), res)
}

func (s *Suite) TestDeleteError() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "projects"`).
		WithArgs(s.project.ID).
		WillReturnError(gorm.ErrRecordNotFound)
	s.mock.ExpectRollback()

	res, err := s.repo.Delete(s.project.ID)

	require.Error(s.T(), err)
	require.False(s.T(), res)
}
