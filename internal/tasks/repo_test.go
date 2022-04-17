package tasks_test

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/fylerx/fyler/internal/tasks"
	"github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	repo tasks.Repository
	task *tasks.Task
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

	s.repo = tasks.InitRepo(s.DB)

	now := time.Now()
	s.task = &tasks.Task{
		ID:        999,
		ProjectID: 777,
		TaskType:  "doc_to_pdf",
		FilePath:  "document/22/file/file.docx",
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

	now := time.Now()

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "status", "task_type", "file_path", "updated_at"}).
			AddRow(1, 500, "queued", "doc_to_pdf", "http://example.com/cat.docx", now).
			AddRow(2, 505, "progress", "video_to_mp4", "http://example.com/cat.mp4", now).
			AddRow(3, 500, "failed", "doc_to_pdf", "http://example.com/cat.xslx", now))

	res, err := s.repo.GetAll()

	expectedProjects := []*tasks.Task{
		{
			ID:        1,
			ProjectID: 500,
			Status:    "queued",
			TaskType:  "doc_to_pdf",
			FilePath:  "http://example.com/cat.docx",
			UpdatedAt: now,
		},
		{
			ID:        2,
			ProjectID: 505,
			Status:    "progress",
			TaskType:  "video_to_mp4",
			FilePath:  "http://example.com/cat.mp4",
			UpdatedAt: now,
		},
		{
			ID:        3,
			ProjectID: 500,
			Status:    "failed",
			TaskType:  "doc_to_pdf",
			FilePath:  "http://example.com/cat.xslx",
			UpdatedAt: now,
		},
	}
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(expectedProjects, res))
}

func (s *Suite) TestGetAllError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks"`)).
		WillReturnError(gorm.ErrInvalidData)

	res, err := s.repo.GetAll()

	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *Suite) TestCreate() {
	var projectID uint32 = 500
	var newTaskID uint64 = 900
	now := time.Now()

	input := &tasks.Task{
		ProjectID: projectID,
		TaskType:  "doc_to_pdf",
		FilePath:  "http://example.com/cat.docx",
	}
	expTask := &tasks.Task{
		ID:        newTaskID,
		ProjectID: projectID,
		Status:    "queued",
		TaskType:  "doc_to_pdf",
		FilePath:  "http://example.com/cat.docx",
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "tasks"
		("project_id","status","task_type","file_path","created_at","updated_at")
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`),
	).
		WithArgs(projectID, "queued", input.TaskType, input.FilePath, AnyTime{}, AnyTime{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newTaskID))
	s.mock.ExpectCommit()

	res, err := s.repo.Create(input)

	require.NoError(s.T(), err)
	require.True(s.T(), cmp.Equal(expTask, res, cmpopts.IgnoreFields(tasks.Task{}, "CreatedAt", "UpdatedAt")))
}

func (s *Suite) TestCreateError() {
	var projectID uint32 = 500
	input := &tasks.Task{
		ProjectID: projectID,
		TaskType:  "doc_to_pdf",
		FilePath:  "http://example.com/cat.docx",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "tasks"
		("project_id","status","task_type","file_path","created_at","updated_at")
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`),
	).
		WithArgs(projectID, "queued", input.TaskType, input.FilePath, AnyTime{}, AnyTime{}).
		WillReturnError(gorm.ErrInvalidData)
	s.mock.ExpectRollback()

	res, err := s.repo.Create(input)

	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}
