package runners

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/foto-leistenschneider/admin-panel/pkg/protos"
)

var Register = map[string]*Runner{}

type Runner struct {
	Jobs     map[string]*protos.Job
	LastPing time.Time
	Name     string
}

func (r Runner) GetJobs() []*protos.Job {
	jobs := make([]*protos.Job, len(r.Jobs))
	i := 0
	for _, job := range r.Jobs {
		jobs[i] = job
		i++
	}
	slices.SortFunc(jobs, func(a, b *protos.Job) int {
		tsA := a.GetCreatedAt().AsTime()
		tsB := b.GetCreatedAt().AsTime()
		return tsB.Compare(tsA)
	})
	return jobs
}

var mutex = &sync.Mutex{}

func Ping(ping *protos.Ping) (*protos.Jobs, error) {
	if ping == nil {
		return nil, errors.New("ping data is nil")
	}
	mutex.Lock()
	defer mutex.Unlock()

	if r, ok := Register[ping.Name]; !ok {
		Register[ping.Name] = &Runner{
			Jobs:     make(map[string]*protos.Job),
			LastPing: time.Now(),
			Name:     ping.Name,
		}
		if len(ping.JobUpdates) > 0 {
			return nil, errors.New("runner does not exist but job updates were sent, this should not happen as it can't have any jobs")
		}
		return nil, nil
	} else {
		r.LastPing = time.Now()
		for _, jobUpdate := range ping.JobUpdates {
			job, ok := r.Jobs[jobUpdate.JobId]
			if !ok {
				return nil, fmt.Errorf("job with id %s not found in runner %s", jobUpdate.JobId, ping.Name)
			}
			job.Status = jobUpdate.NewStatus
			job.Output = jobUpdate.Output
		}
		Register[ping.Name] = r

		var newJobs protos.Jobs
		for _, job := range r.Jobs {
			if job.Status == protos.JobStatus_Pending {
				newJobs.Jobs = append(newJobs.Jobs, job)
			}
		}
		return &newJobs, nil
	}
}
