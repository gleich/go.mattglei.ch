package pkg

import (
	"sort"
	"sync"

	"go.mattglei.ch/go.mattglei.ch/internal/github"
)

type Packages struct {
	packages []github.Repository
	mutex    sync.Mutex
}

func (p *Packages) All() []github.Repository {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	sort.Slice(
		p.packages,
		func(i, j int) bool { return p.packages[j].Updated.Before(p.packages[i].Updated) },
	)
	return p.packages
}

func (p *Packages) Get(name string) github.Repository {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, repo := range p.packages {
		if repo.Name == name {
			return repo
		}
	}
	return github.Repository{}
}

func (p *Packages) Set(repo github.Repository) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for i := range p.packages {
		if p.packages[i].Name == repo.Name {
			p.packages[i] = repo
		}
	}
}
