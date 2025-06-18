package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/foto-leistenschneider/admin-panel/internal/db"
	"github.com/foto-leistenschneider/admin-panel/internal/runners"
	"github.com/robfig/cron/v3"
)

var (
	c        *cron.Cron
	registry = make(map[int64]cron.EntryID)
	mut      = sync.Mutex{}
)

func init() {
	mut.Lock()
	c = cron.New(
		cron.WithLocation(time.UTC),
		cron.WithParser(cron.NewParser(
			cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)),
	)

	c.Start()
	mut.Unlock()

	go func() {
		ts, err := db.Q.GetTasks(context.Background())
		if err != nil {
			log.Fatal("get tasks", "err", err)
		}
		for _, t := range ts {
			if err := Add(t); err != nil {
				log.Error("add task", "err", err)
			}
		}
	}()
}

func Close() {
	mut.Lock()
	defer mut.Unlock()
	Clear()
	c.Stop()
}

func Add(task db.Task) error {
	mut.Lock()
	defer mut.Unlock()

	if _, ok := registry[task.ID]; ok {
		c.Remove(registry[task.ID])
	}
	entryId, err := c.AddFunc(task.Schedule, func() {
		runners, err := runners.FindRunners(task.Selector)
		log.Info("creating jobs", "task", task.Description, "selector", task.Selector, "scope", task.Scope, "runners", len(runners))
		if err != nil {
			log.Error("find runners", "err", err)
			return
		}
		for _, runner := range runners {
			if err := runner.AddJob(task.Scope, task.Command); err != nil {
				log.Error("add job", "runner", runner, "err", err)
			}
		}
	})
	if err != nil {
		return err
	}
	registry[task.ID] = entryId
	return nil
}

func Remove(id int64) error {
	mut.Lock()
	defer mut.Unlock()

	if entryId, ok := registry[id]; ok {
		c.Remove(entryId)
		delete(registry, id)
		return nil
	}
	return fmt.Errorf("task with id %d not found", id)
}

func Clear() {
	mut.Lock()
	defer mut.Unlock()

	for id := range registry {
		c.Remove(registry[id])
	}
	registry = make(map[int64]cron.EntryID)
}
