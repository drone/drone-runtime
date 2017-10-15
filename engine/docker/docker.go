package docker

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/drone/drone-runtime/engine"
)

type dockerEngine struct {
	client client.APIClient
}

// New returns a new Docker Engine using the given client.
func New(cli client.APIClient) engine.Engine {
	return &dockerEngine{
		client: cli,
	}
}

// NewEnv returns a new Docker Engine using the client connection
// environment variables.
func NewEnv() (engine.Engine, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return New(cli), nil
}

func (e *dockerEngine) Setup(conf *engine.Config) error {
	for _, vol := range conf.Volumes {
		_, err := e.client.VolumeCreate(context.TODO(), volume.VolumesCreateBody{
			Name:       vol.Name,
			Driver:     vol.Driver,
			DriverOpts: vol.DriverOpts,
			// Labels:     defaultLabels,
		})
		if err != nil {
			return err
		}
	}
	for _, network := range conf.Networks {
		_, err := e.client.NetworkCreate(context.TODO(), network.Name, types.NetworkCreate{
			Driver:  network.Driver,
			Options: network.DriverOpts,
			// Labels:  defaultLabels,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *dockerEngine) Create(proc *engine.Step) error {
	ctx := context.Background()

	config := toConfig(proc)
	hostConfig := toHostConfig(proc)

	// create pull options with encoded authorization credentials.
	pullopts := types.ImagePullOptions{}
	if proc.AuthConfig.Username != "" && proc.AuthConfig.Password != "" {
		pullopts.RegistryAuth = encodeAuthToBase64(proc.AuthConfig)
	}

	// automatically pull the latest version of the image if requested
	// by the process configuration.
	if proc.Pull {
		rc, perr := e.client.ImagePull(ctx, config.Image, pullopts)
		if perr == nil {
			io.Copy(ioutil.Discard, rc)
			rc.Close()
		}
		// fix for drone/drone#1917
		if perr != nil && proc.AuthConfig.Password != "" {
			return perr
		}
	}

	_, err := e.client.ContainerCreate(ctx, config, hostConfig, nil, proc.Name)
	if client.IsErrImageNotFound(err) {
		// automatically pull and try to re-create the image if the
		// failure is caused because the image does not exist.
		rc, perr := e.client.ImagePull(ctx, config.Image, pullopts)
		if perr != nil {
			return perr
		}
		io.Copy(ioutil.Discard, rc)
		rc.Close()

		_, err = e.client.ContainerCreate(ctx, config, hostConfig, nil, proc.Name)
	}
	if err != nil {
		return err
	}

	if len(proc.NetworkMode) == 0 {
		for _, net := range proc.Networks {
			err = e.client.NetworkConnect(ctx, net.Name, proc.Name, &network.EndpointSettings{
				Aliases: net.Aliases,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *dockerEngine) Start(proc *engine.Step) error {
	startOpts := types.ContainerStartOptions{}
	return e.client.ContainerStart(context.TODO(), proc.Name, startOpts)
}

func (e *dockerEngine) Kill(proc *engine.Step) error {
	return e.client.ContainerKill(context.TODO(), proc.Name, "9")
}

func (e *dockerEngine) Wait(proc *engine.Step) (*engine.State, error) {
	_, err := e.client.ContainerWait(context.TODO(), proc.Name)
	if err != nil {
		// todo
	}

	info, err := e.client.ContainerInspect(context.TODO(), proc.Name)
	if err != nil {
		return nil, err
	}
	if info.State.Running {
		// todo
	}

	return &engine.State{
		Exited:    true,
		ExitCode:  info.State.ExitCode,
		OOMKilled: info.State.OOMKilled,
	}, nil
}

func (e *dockerEngine) Tail(proc *engine.Step) (io.ReadCloser, error) {
	logsOpts := types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
		Details:    false,
		Timestamps: false,
	}

	logs, err := e.client.ContainerLogs(context.TODO(), proc.Name, logsOpts)
	if err != nil {
		return nil, err
	}
	rc, wc := io.Pipe()

	go func() {
		stdcopy.StdCopy(wc, wc, logs)
		logs.Close()
		wc.Close()
		rc.Close()
	}()
	return rc, nil
}

func (e *dockerEngine) Upload(proc *engine.Step, path string, r io.Reader) error {
	options := types.CopyToContainerOptions{}
	options.AllowOverwriteDirWithFile = false
	return e.client.CopyToContainer(context.TODO(), proc.Name, path, r, options)
}

func (e *dockerEngine) Download(proc *engine.Step, path string) (io.ReadCloser, *engine.FileInfo, error) {
	rc, stat, err := e.client.CopyFromContainer(context.TODO(), proc.Name, path)
	info := &engine.FileInfo{
		Path:  path,
		Name:  stat.Name,
		Time:  stat.Mtime.Unix(),
		Size:  stat.Size,
		IsDir: stat.Mode.IsDir(),
	}
	return rc, info, err
}

func (e *dockerEngine) Destroy(conf *engine.Config) error {
	removeOpts := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}

	for _, stage := range conf.Stages {
		for _, step := range stage.Steps {
			e.client.ContainerKill(context.TODO(), step.Name, "9")
			e.client.ContainerRemove(context.TODO(), step.Name, removeOpts)
		}
	}
	for _, volume := range conf.Volumes {
		e.client.VolumeRemove(context.TODO(), volume.Name, true)
	}
	for _, network := range conf.Networks {
		e.client.NetworkRemove(context.TODO(), network.Name)
	}
	return nil
}
