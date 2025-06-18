package runners

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sync"
	"time"

	"github.com/foto-leistenschneider/admin-panel/pkg/protos"
	"github.com/google/uuid"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var Register = map[string]*Runner{}
var mutex = &sync.RWMutex{}

type Runner struct {
	Jobs     map[string]*protos.Job
	LastPing time.Time
	Name     string
}

func (r Runner) GetJobs() []*protos.Job {
	mutex.RLock()
	defer mutex.RUnlock()

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

func (r *Runner) AddJob(scope string, command string) error {
	mutex.Lock()
	defer mutex.Unlock()

	jobScope, err := protos.ParseJobScope(scope)
	if err != nil {
		return fmt.Errorf("invalid scope: %s", err)
	}

	job := protos.Job{
		Id:        uuid.NewString(),
		Command:   command,
		Scope:     jobScope,
		Status:    protos.JobStatus_Pending,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}

	r.Jobs[job.Id] = &job
	return nil
}

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
			job.UpdatedAt = timestamppb.Now()
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

func FindRunners(query string) ([]*Runner, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	if len(query) == 0 {
		// empty regex would match everything
		// check can be skipped and all runners are returned
		runners := make([]*Runner, 0, len(Register))
		for _, runner := range Register {
			runners = append(runners, runner)
		}
		return runners, nil
	}

	exp, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}

	var runners []*Runner
	for name, runner := range Register {
		if exp.MatchString(name) {
			runners = append(runners, runner)
		}
	}
	return runners, nil
}
