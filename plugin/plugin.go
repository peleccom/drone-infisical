// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"
	infisical "github.com/infisical/go-sdk"
	"github.com/sirupsen/logrus"
)

func New(client infisical.InfisicalClientInterface, projectId string) secret.Plugin {
	return &plugin{
		client:    client,
		projectId: projectId,
	}
}

type plugin struct {
	client    infisical.InfisicalClientInterface
	projectId string
}

func (p *plugin) Find(ctx context.Context, req *secret.Request) (*drone.Secret, error) {
	path := req.Path
	name := req.Name

	apiKeySecret, err := p.client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		SecretKey:   name,
		Environment: "dev",
		ProjectID:   p.projectId,
		SecretPath:  path,
	})

	if err != nil {
		logrus.Debugf("Error: %v", err)
		return nil, nil
	}

	return &drone.Secret{
		Name:        name,
		Data:        apiKeySecret.SecretValue,
		PullRequest: false, // never expose to pulls requests
	}, nil
}
