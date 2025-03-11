// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"drone-infisical-secrets/plugin"

	"github.com/drone/drone-go/plugin/secret"

	infisical "github.com/infisical/go-sdk"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind                  string `envconfig:"DRONE_INFISICAL_DRONE_BIND"`
	Debug                 bool   `envconfig:"DRONE_INFISICAL_DEBUG"`
	Secret                string `envconfig:"DRONE_INFISICAL_DRONE_SECRET"`
	InfisicalUrl          string `envconfig:"DRONE_INFISICAL_URL"`
	InfisicalClientId     string `envconfig:"DRONE_INFISICAL_CLIENT_ID"`
	InfisicalClientSecret string `envconfig:"DRONE_INFISICAL_CLIENT_SECRET"`
	InfisicalProjectId    string `envconfig:"DRONE_INFISICAL_PROJECT_ID"`
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	// creates infisical client
	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl:          spec.InfisicalUrl,
		AutoTokenRefresh: true, // Wether or not to let the SDK handle the access token lifecycle. Defaults to true if not specified.
	})

	if err != nil {
		logrus.Fatalln(err)
	}

	_, err = client.Auth().UniversalAuthLogin(spec.InfisicalClientId, spec.InfisicalClientSecret)

	if err != nil {
		fmt.Printf("Authentication failed: %v", err)
		os.Exit(1)
	}

	handler := secret.Handler(
		spec.Secret,
		plugin.New(
			client,
			spec.InfisicalProjectId,
		),
		logrus.StandardLogger(),
	)
	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	http.HandleFunc("/health", healthCheckHandler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
