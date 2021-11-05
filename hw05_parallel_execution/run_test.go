package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestEventuallyResCheck(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("tasks without errors and without timeout", func(t *testing.T) {
		tasksCount := 153
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 13
		maxErrorsCount := 0

		require.Eventually(t, func() bool {
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)
			return runTasksCount == int32(len(tasks))
		}, time.Second/100, time.Second/1000)
	})

	t.Run("tasks with errors and without timeout", func(t *testing.T) {
		tasksCount := 117
		errorCount := 11
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var runErrCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		for i, errCnt := 0, 0; errCnt < errorCount && i < tasksCount; i, errCnt = i+2, errCnt+1 {
			err := fmt.Errorf("error from task %d", i)
			tasks[i] = func() error {
				atomic.AddInt32(&runErrCount, 1)
				return err
			}
		}

		workersCount := 1
		maxErrorsCount := 5

		require.Eventually(t, func() bool {
			err := Run(tasks, workersCount, maxErrorsCount)
			require.Error(t, err)
			// fmt.Printf("maxErrorsCount = %v, runErrCount = %v, errorCount = %v\n", maxErrorsCount, runErrCount, errorCount)
			return int32(maxErrorsCount) < runErrCount && runErrCount <= int32(errorCount)
		}, time.Second/100, time.Second/1000)
	})
}
