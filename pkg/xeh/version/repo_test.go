package version

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/xeh/ehmocks"
	"testing"
	"time"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/repo/version"
	"github.com/looplab/eventhorizon/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite

	inner      *ehmocks.ReadRepoMock
	eventStore *ehmocks.EventStoreMock

	repo *Repo

	ev1 eh.Event
	ev2 eh.Event
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, &RepoTestSuite{})
}

// before each test
func (s *RepoTestSuite) SetupTest() {
	s.inner = &ehmocks.ReadRepoMock{}
	s.eventStore = &ehmocks.EventStoreMock{}
	s.repo = NewRepo(s.inner, s.eventStore)

	aggID := uuid.New()
	s.ev1 = eh.NewEvent("test", nil, time.Now(), eh.ForAggregate("test", aggID, 1))
	s.ev2 = eh.NewEvent("test-2", nil, time.Now(), eh.ForAggregate("test", aggID, 2))
}

func (s *RepoTestSuite) Test_find_calls_underling_repo() {
	// GIVEN repo has entity
	id := uuid.New()

	s.inner.On("Find", mock.Anything, id).Return(&ehmocks.EntityFake{ID: id}, nil)
	s.eventStore.On("LoadFrom", mock.Anything, id, 0).Return([]eh.Event{s.ev1}, nil)

	// WHEN we call find
	e, err := s.repo.Find(context.TODO(), id)

	// THEN the underling repo is called
	s.inner.AssertCalled(s.T(), "Find", mock.Anything, id)

	// AND the entity is returned
	s.NoError(err)
	s.Equal(&ehmocks.EntityFake{ID: id}, e, "Should return the found entity")
}

func (s *RepoTestSuite) Test_sets_context_with_version_to_find() {
	// GIVEN event store has events with version 1 and 2
	id := uuid.New()

	s.inner.On("Find", mock.Anything, id).Return(&ehmocks.EntityFake{ID: id}, nil)
	s.eventStore.On("LoadFrom", mock.Anything, id, mock.Anything).Return([]eh.Event{s.ev2, s.ev1}, nil)

	// WHEN we call find
	_, _ = s.repo.Find(context.TODO(), id)

	// THEN the event store is called with version 0
	s.eventStore.AssertCalled(s.T(), "LoadFrom", mock.Anything, id, 0)

	// AND the context was set with version 2
	s.inner.AssertCalled(s.T(), "Find", mock.MatchedBy(func(ctx context.Context) bool {
		minVersion, ok := version.MinVersionFromContext(ctx)
		s.True(ok, "Should set context with version")
		s.Equal(2, minVersion, "Should set context with version 2")
		return true
	}), id)
}

func (s *RepoTestSuite) Test_find_returns_error_if_underling_repo_returns_error() {
	// GIVEN repo returns error
	id := uuid.New()

	s.inner.On("Find", mock.Anything, id).Return(nil, &eh.RepoError{})
	s.eventStore.On("LoadFrom", mock.Anything, id, 0).Return([]eh.Event{s.ev1}, nil)

	// WHEN we call find
	_, err := s.repo.Find(context.Background(), id)

	// THEN the error is returned
	s.Error(err)
}
