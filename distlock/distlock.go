package distlock

import (
	"github.com/hashicorp/consul/api"
)

const lockPrefix = "lock/"

type DistLock struct {
	Lock				*api.Lock
}

func NewDistLock(address string, name string) (*DistLock, error) {
	dl := new(DistLock)

	client, err := api.NewClient(&api.Config{Address: address})
	if err != nil {
		return nil, err
	}

	agentChecks, err := client.Agent().Checks()
	if err != nil {
		return nil, err
	}
	checks := make([]string, 0)
	checks = append(checks, "serfHealth")
	for _, ac := range agentChecks {
		checks = append(checks, ac.CheckID)
	}

	opts := &api.LockOptions{
		Key:        	lockPrefix + name,
		SessionOpts: 	&api.SessionEntry{
			Checks:   checks,
			Behavior: "release",
		},
	}

	lock, err := client.LockOpts(opts)
	if err != nil {
		return nil, err
	}

	dl.Lock = lock

	return dl, nil
}

func (dl *DistLock) AquireLock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	return dl.Lock.Lock(stopCh)
}

func (dl *DistLock) ReleaseLock() error {
	return dl.Lock.Unlock()
}

func (dl *DistLock) Destroy() error {
	return dl.Lock.Destroy()
}
