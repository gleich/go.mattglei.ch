package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v83/github"
	"go.mattglei.ch/go.mattglei.ch/internal/secrets"
	"go.mattglei.ch/timber"
)

const webhook_name = "go.mattglei.ch"

func (r Repository) Subscribe(clients Clients) error {
	start := time.Now()
	hook := &github.Hook{
		Name:   new(webhook_name),
		Active: new(true),
		Config: &github.HookConfig{
			URL:         new("https://go.mattglei.ch/github/webhook"),
			ContentType: new("json"),
			Secret:      new(secrets.ENV.GitHubWebhookSecret),
			InsecureSSL: new("0"),
		},
		Events: []string{"push", "release", "repository"},
	}

	_, _, err := clients.REST.Repositories.CreateHook(context.Background(), r.Owner, r.Name, hook)
	if err != nil {
		if strings.Contains(err.Error(), "Hook already exists on this repository") {
			return nil
		}
		return fmt.Errorf("creating hook: %w", err)
	}

	timber.DoneSince(start, "created hook for", r.Name)
	return nil
}

func Unsubscribe(clients Clients, owner, name string) error {
	hooks, _, err := clients.REST.Repositories.ListHooks(context.Background(), owner, name, nil)
	if err != nil {
		return fmt.Errorf("fetching hooks: %w", err)
	}

	for _, hook := range hooks {
		if hook.GetName() == webhook_name {
			_, err = clients.REST.Repositories.DeleteHook(
				context.Background(),
				owner,
				name,
				*hook.ID,
			)
			if err != nil {
				return fmt.Errorf("deleting hook: %w", err)
			}
		}
	}

	return nil
}
