package architecture

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTube_ProcessDelayedQueueWhenLessJobsReadyThanLimit(t *testing.T) {
	// === given
	testTube := getTestTube(t)
	delayed := testTube.Delayed.(*MockPriorityQueue)
	ready := testTube.Ready.(*MockPriorityQueue)

	// return a job two times when Peeking
	maxPeeks := 0
	delayed.On("Peek").Return(func() PriorityQueueItem {
		maxPeeks += 1
		if maxPeeks > 2 {
			return nil
		}
		return getDelayedJob()
	})

	// return a job two times when Dequeueing
	maxDequeues := 0
	delayed.On("Dequeue").Return(func() PriorityQueueItem {
		maxDequeues += 1
		if maxDequeues > 2 {
			return nil
		}
		return getDelayedJob()
	})

	// === when
	// jobs in the delayed queue(2) < limit(5)
	testTube.ProcessDelayedQueue(5)

	// === then
	ready.AssertNumberOfCalls(t, "Enqueue", 2)
}

func TestTube_ProcessDelayedQueueWhenMoreJobsReadyThanLimit(t *testing.T) {
	// === given
	testTube := getTestTube(t)
	delayed := testTube.Delayed.(*MockPriorityQueue)
	ready := testTube.Ready.(*MockPriorityQueue)

	// return a job infinite times when Peeking
	delayed.On("Peek").Return(func() PriorityQueueItem {
		return getDelayedJob()
	})

	// return a job infinite times when Dequeueing
	delayed.On("Dequeue").Return(func() PriorityQueueItem {
		return getDelayedJob()
	})

	// === when
	// jobs in the delayed queue(infinit) > limit(5)
	testTube.ProcessDelayedQueue(5)

	// === then
	ready.AssertNumberOfCalls(t, "Enqueue", 5)
}

func TestTube_ProcessReservedQueueWhenLessJobsReadyThanLimit(t *testing.T) {
	// === given
	testTube := getTestTube(t)
	reserved := testTube.Reserved.(*MockPriorityQueue)
	ready := testTube.Ready.(*MockPriorityQueue)
	// return a job two times when Peeking
	maxPeeks := 0
	reserved.On("Peek").Return(func() PriorityQueueItem {
		maxPeeks += 1
		if maxPeeks > 2 {
			return nil
		}
		return getReservedJob()
	})

	// return a job two times when Dequeueing
	maxDequeues := 0
	reserved.On("Dequeue").Return(func() PriorityQueueItem {
		maxDequeues += 1
		if maxDequeues > 2 {
			return nil
		}
		return getReservedJob()
	})

	// === when
	// jobs in the reserved queue(infinit) > limit(5)
	testTube.ProcessReservedQueue(5)

	// === then
	ready.AssertNumberOfCalls(t, "Enqueue", 2)
}

func TestTube_ProcessReservedQueueWhenMoreJobsReadyThanLimit(t *testing.T) {
	// === given
	testTube := getTestTube(t)
	reserved := testTube.Reserved.(*MockPriorityQueue)
	ready := testTube.Ready.(*MockPriorityQueue)

	// return a job infinite times when Peeking
	reserved.On("Peek").Return(func() PriorityQueueItem {
		return getReservedJob()
	})

	// return a job infinite times when Dequeueing
	reserved.On("Dequeue").Return(func() PriorityQueueItem {
		return getReservedJob()
	})

	// === when
	// jobs in the reserved queue(infinit) > limit(5)
	testTube.ProcessReservedQueue(5)

	// === then
	ready.AssertNumberOfCalls(t, "Enqueue", 5)
}

func getTestTube(t *testing.T) *Tube {
	ready := MockPriorityQueue{}
	reserved := MockPriorityQueue{}
	delayed := MockPriorityQueue{}
	buried := MockPriorityQueue{}

	// ready queue must accept the job
	ready.On("Enqueue", mock.AnythingOfTypeArgument("*architecture.Job")).Run(func(args mock.Arguments) {
		// jobs put to ready queue must have the correct state
		job := args.Get(0).(*Job)
		assert.Equal(t, READY, job.State())
	})

	return &Tube{
		Name:     "test_tube",
		Ready:    &ready,
		Reserved: &reserved,
		Delayed:  &delayed,
		Buried:   &buried,
	}
}

func getDelayedJob() *Job {
	testJob := NewJob("dummy_job", 0, 0, 1, 1, "dummy data")
	testJob.SetState(DELAYED)
	return testJob
}

func getReservedJob() *Job {
	testJob := NewJob("dummy_job", 0, 0, 1, 1, "dummy data")
	testJob.SetState(RESERVED)
	return testJob
}
