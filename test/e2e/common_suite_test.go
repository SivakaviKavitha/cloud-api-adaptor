// (C) Copyright Confidential Containers Contributors
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	envconf "sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

const WAIT_POD_RUNNING_TIMEOUT = time.Second * 300

// doTestCreateSimplePod tests a simple peer-pod can be created.
func doTestCreateSimplePod(t *testing.T, assert CloudAssert) {
	// TODO: generate me.
	namespace := "default"
	name := "simple-peer-pod"
	pod := newPod(namespace, name, "nginx", "kata")

	simplePodFeature := features.New("Simple Peer Pod").
		WithSetup("Create pod", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			if err = client.Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}
			if err = wait.For(conditions.New(client.Resources()).PodRunning(pod), wait.WithTimeout(WAIT_POD_RUNNING_TIMEOUT)); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).
		Assess("PodVM is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			assert.HasPodVM(t, name)

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			if err = client.Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()
	testEnv.Test(t, simplePodFeature)
}
func doTestCreatePublicPod(t *testing.T, assert CloudAssert) {
	// TODO: generate me.
	namespace := "default"
	deploymentname := "nginx"
	podname := "public-pod"
	servicename := "nginx-service"
	ipaddress := "8.8.8.8"
	deployment := newDeployment(namespace, deploymentname, "kata")
	service := newService(namespace, deploymentname, servicename, ipaddress)
	PublipodFeature := features.New("Public Pod").
		WithSetup("Create pod", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			if err = client.Resources(namespace).Create(ctx, deployment); err != nil {
				t.Fatal(err)
			}
			if err = client.Resources(namespace).Create(ctx, service); err != nil {
				t.Fatal(err)
			}
			if err = wait.For(conditions.New(client.Resources()).PodRunning(deployment), wait.WithTimeout(WAIT_POD_RUNNING_TIMEOUT)); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).
		Assess("PodVM is created", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var stdout, stderr bytes.Buffer
			if err := cfg.Client().Resources().ExecInPod(ctx, namespace, podname, deploymentname, []string{"curl", ipaddress}, &stdout, &stderr); err != nil {
				t.Errorf("Error: %s", stderr.String())
				t.Fatal(err)
			}
			logrus.Info(stdout.String())
			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			if err = client.Resources().Delete(ctx, deployment); err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()
	testEnv.Test(t, PublipodFeature)
}
