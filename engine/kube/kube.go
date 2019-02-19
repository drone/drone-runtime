// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package kube

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-runtime/engine/docker/auth"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type kubeEngine struct {
	client *kubernetes.Clientset
	node   string
}

// NewFile returns a new Kubernetes engine from a
// Kubernetes configuration file (~/.kube/config).
func NewFile(url, path, node string) (engine.Engine, error) {
	config, err := clientcmd.BuildConfigFromFlags(url, path)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &kubeEngine{client: client, node: node}, nil
}

func (e *kubeEngine) Setup(ctx context.Context, spec *engine.Spec) error {
	ns := toNamespace(spec)

	// create the project namespace. all pods and
	// containers are created within the namespace, and
	// are removed when the pipeline execution completes.
	_, err := e.client.CoreV1().Namespaces().Create(ns)
	if err != nil {
		return err
	}

	// create all secrets
	for _, secret := range spec.Secrets {
		_, err := e.client.CoreV1().Secrets(ns.Name).Create(
			toSecret(spec, secret),
		)
		if err != nil {
			return err
		}
	}

	// create all registry credentials as secrets.
	if spec.Docker != nil && len(spec.Docker.Auths) > 0 {
		out, err := auth.Marshal(spec.Docker.Auths)
		if err != nil {
			return err
		}
		_, err = e.client.CoreV1().Secrets(ns.Name).Create(
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "docker-auth-config",
				},
				Type: "kubernetes.io/dockerconfigjson",
				StringData: map[string]string{
					".dockerconfigjson": string(out),
				},
			},
		)
		if err != nil {
			return err
		}
	}

	// create all files as config maps.
	for _, file := range spec.Files {
		_, err := e.client.CoreV1().ConfigMaps(ns.Name).Create(
			&v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: file.Metadata.UID,
				},
				Data: map[string]string{
					file.Metadata.UID: string(file.Data),
				},
			},
		)
		if err != nil {
			return err
		}
	}

	// pv := toPersistentVolume(e.node, spec.Metadata.Namespace, spec.Metadata.Namespace, filepath.Join("/tmp", spec.Metadata.Namespace))
	// _, err = e.client.CoreV1().PersistentVolumes().Create(pv)
	// if err != nil {
	// 	return err
	// }

	// pvc := toPersistentVolumeClaim(spec.Metadata.Namespace, spec.Metadata.Namespace)
	// _, err = e.client.CoreV1().PersistentVolumeClaims(spec.Metadata.Namespace).Create(pvc)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (e *kubeEngine) Create(_ context.Context, _ *engine.Spec, _ *engine.Step) error {
	// no-op
	return nil
}

func (e *kubeEngine) Start(ctx context.Context, spec *engine.Spec, step *engine.Step) error {
	pod := toPod(spec, step)
	if len(step.Docker.Ports) != 0 {
		service := toService(spec, step)
		_, err := e.client.CoreV1().Services(spec.Metadata.Namespace).Create(service)
		if err != nil {
			return err
		}
	}

	if e.node != "" {
		pod.Spec.Affinity = &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
					NodeSelectorTerms: []v1.NodeSelectorTerm{{
						MatchExpressions: []v1.NodeSelectorRequirement{{
							Key:      "kubernetes.io/hostname",
							Operator: v1.NodeSelectorOpIn,
							Values:   []string{e.node},
						}},
					}},
				},
			},
		}
	}

	_, err := e.client.CoreV1().Pods(spec.Metadata.Namespace).Create(pod)
	return err
}

func (e *kubeEngine) Wait(ctx context.Context, spec *engine.Spec, step *engine.Step) (*engine.State, error) {
	stopper := make(chan struct{})
	updater := func(old interface{}, new interface{}) {
		pod := new.(*v1.Pod)
		// ignore events that do not come from the
		// current pod namespace.
		if pod.ObjectMeta.Namespace != step.Metadata.Namespace {
			return
		}
		if pod.Name == step.Metadata.UID {
			switch pod.Status.Phase {
			case v1.PodSucceeded, v1.PodFailed, v1.PodUnknown:
				// TODO need to understand if this could be
				// invoked multiple times.
				select {
				case stopper <- struct{}{}:
				default:
				}
			}
		}
	}

	factory := informers.NewSharedInformerFactory(e.client, time.Second)
	informer := factory.Core().V1().Pods().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: updater,
		},
	)
	factory.Start(wait.NeverStop)

	// TODO Cancel on ctx.Done
	<-stopper

	pod, err := e.client.CoreV1().Pods(spec.Metadata.Namespace).Get(step.Metadata.UID, metav1.GetOptions{
		IncludeUninitialized: true,
	})
	if err != nil {
		return nil, err
	}

	state := &engine.State{
		ExitCode:  int(pod.Status.ContainerStatuses[0].State.Terminated.ExitCode),
		Exited:    true,
		OOMKilled: false,
	}
	return state, nil
}

func (e *kubeEngine) Tail(ctx context.Context, spec *engine.Spec, step *engine.Step) (io.ReadCloser, error) {
	ns := spec.Metadata.Namespace
	podName := step.Metadata.UID

	up := make(chan bool)

	var podUpdated = func(old interface{}, new interface{}) {
		pod := new.(*v1.Pod)
		if pod.Name == podName {
			switch pod.Status.Phase {
			case v1.PodRunning, v1.PodSucceeded, v1.PodFailed:
				up <- true
			}
		}
	}

	si := informers.NewSharedInformerFactory(e.client, 5*time.Minute)
	si.Core().V1().Pods().Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: podUpdated,
		},
	)
	si.Start(wait.NeverStop)

	select {
	case <-up:
	case <-ctx.Done():
	}

	opts := &v1.PodLogOptions{
		Follow: true,
	}

	return e.client.CoreV1().RESTClient().Get().
		Namespace(ns).
		Name(podName).
		Resource("pods").
		SubResource("log").
		VersionedParams(opts, scheme.ParameterCodec).
		Stream()
}

func (e *kubeEngine) Destroy(ctx context.Context, spec *engine.Spec) error {
	// err := e.client.CoreV1().PersistentVolumes().Delete(spec.Metadata.Namespace, nil)
	// if err != nil {
	// 	// TODO show error message
	// }

	// err = e.client.CoreV1().PersistentVolumeClaims(spec.Metadata.Namespace).Delete("workspace", nil)
	// if err != nil {
	// 	// TODO show error message
	// }

	// this is a complete hack. we are creating a host machine
	// directory which should be handled by a persistent volume.
	// I am planning to switch to a persistent volume, but am
	// leaving this in place as a temporary workaround in the short
	// term.
	os.RemoveAll(
		filepath.Join(
			"/tmp",
			"drone",
			spec.Metadata.Namespace,
		),
	)

	// deleting the namespace should destroy all secrets,
	// volumes, configuration files and more.
	return e.client.CoreV1().Namespaces().Delete(
		spec.Metadata.Namespace,
		&metav1.DeleteOptions{},
	)
}
