package runners

import (
	"errors"
	"fmt"
	"time"

	"github.com/foto-leistenschneider/admin-panel/pkg/protos"
)

var Register = map[string]Runner{}

type Runner struct {
	Jobs     []*protos.Job
	LastPing time.Time
	Name     string
}

func Ping(ping *protos.Ping) (*protos.Jobs, error) {
	if ping == nil {
		return nil, errors.New("ping data is nil")
	}
	if r, ok := Register[ping.Name]; !ok {
		Register[ping.Name] = Runner{
			Jobs:     nil,
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
			jobIndex := int(jobUpdate.JobId)
			if jobIndex >= len(r.Jobs) || jobIndex < 0 {
				return nil, fmt.Errorf("job id %d does not exist in this runner", jobIndex)
			}
			job := r.Jobs[jobIndex]
			job.Status = jobUpdate.NewStatus
			if job.Status == protos.JobStatus_JOB_STATUS_FAILED {
				job.Output = jobUpdate.Output
			} else {
				job.Output = ""
			}
		}
		var newJobs protos.Jobs
		for _, job := range r.Jobs {
			if job.Status == protos.JobStatus_JOB_STATUS_PENDING {
				newJobs.Jobs = append(newJobs.Jobs, job)
			}
		}
		return &newJobs, nil
	}
}
