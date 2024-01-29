package tests

import (
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DEFAULT_POSTGRES_PASSWORD = "password"
	DEFAULT_POSTGRES_USERNAME = "rootuser"
	DEFAULT_POSTGRES_DB       = "test"
	DEFAULT_POSTGRES_SCHEMA   = "block_feed"
	POSTGRES_PORT             = "5432/tcp"
	REDIS_PORT                = "6379/tcp"
)

type (
	HostConnectionInfo struct {
		Url  string
		Host string
		Port nat.Port
	}

	ContainerWithConnectionInfo struct {
		Container testcontainers.Container
		Conn      *HostConnectionInfo
	}
)

func NewPostgresContainer(ctx context.Context, t *testing.T, version string) (*ContainerWithConnectionInfo, error) {
	// Gets the directory that this file exists in
	dir, err := GetCurrentDir()
	if err != nil {
		return nil, err
	}

	// Creates the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			ExposedPorts: []string{POSTGRES_PORT},
			WaitingFor:   wait.ForExposedPort(),
			Cmd:          []string{"postgres", "-c", "log_statement=all"},
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    path.Join(*dir, ".."),
				Dockerfile: path.Join("tests", "assets", "postgres", "Dockerfile"),
				BuildArgs: map[string]*string{
					"POSTGRES_VERSION": &version,
				},
			},
			Env: map[string]string{
				"POSTGRES_PASSWORD": DEFAULT_POSTGRES_PASSWORD,
				"POSTGRES_USER":     DEFAULT_POSTGRES_USERNAME,
				"POSTGRES_DB":       DEFAULT_POSTGRES_DB,
			},
		},
		Started: true,
	})

	// Schedules the container for termination once the test case is completed
	if err != nil {
		return nil, err
	} else {
		ScheduleContainerTermination(t, container)
	}

	// Gets the connection info of the container
	conn, err := GetConnectionInfo(ctx, container, nat.Port(POSTGRES_PORT))
	if err != nil {
		return nil, err
	}

	// Replace the default URL with the correct postgres connection string
	conn.Url = fmt.Sprintf(
		"postgres://%s:%s@%s:%s",
		DEFAULT_POSTGRES_USERNAME,
		DEFAULT_POSTGRES_PASSWORD,
		"host.docker.internal",
		conn.Port.Port(),
	)

	// Returns the container info
	return &ContainerWithConnectionInfo{
		Container: container,
		Conn:      conn,
	}, nil
}

func NewRedisContainer(ctx context.Context, t *testing.T, version string, cmd []string) (*ContainerWithConnectionInfo, error) {
	// Creates the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        fmt.Sprintf("redis:%s", version),
			ExposedPorts: []string{REDIS_PORT},
			WaitingFor:   wait.ForExposedPort(),
			Cmd:          cmd,
		},
		Started: true,
	})

	// Schedules the container for termination once the test case is completed
	if err != nil {
		return nil, err
	} else {
		ScheduleContainerTermination(t, container)
	}

	// Gets the connection info of the container
	conn, err := GetConnectionInfo(ctx, container, nat.Port(REDIS_PORT))
	if err != nil {
		return nil, err
	}

	// Returns the container info
	return &ContainerWithConnectionInfo{
		Container: container,
		Conn:      conn,
	}, nil
}

func RedisDefaultCmd() []string {
	return []string{
		"redis-server",
		"--port", "6379",
		"--loglevel", "debug",
	}
}

func RedisCacheCmd() []string {
	return append(RedisDefaultCmd(), []string{
		"--maxmemory", "100mb",
		"--maxmemory-policy", "allkeys-lfu",
		"--activedefrag", "yes",
		"--save", "\"\"",
	}...)
}

func RedisQueueCmd() []string {
	return append(RedisDefaultCmd(), []string{
		"--maxmemory", "0", // no memory limit
		"--maxmemory-policy", "noeviction",
		"--appendonly", "yes",
	}...)
}

func GetConnectionInfo(ctx context.Context, container testcontainers.Container, containerPort nat.Port) (*HostConnectionInfo, error) {
	// Fetches the container host
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	// Fetches the port on the host machine that the container is listening on
	port, err := container.MappedPort(ctx, containerPort)
	if err != nil {
		return nil, err
	}

	// Returns the connection info to use for testing
	return &HostConnectionInfo{
		Url:  fmt.Sprintf("%s:%s", host, port.Port()),
		Host: host,
		Port: port,
	}, nil
}

func ScheduleContainerTermination(t *testing.T, container testcontainers.Container) {
	// Clean up the container when the test case finishes
	t.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Log(err)
		}
	})
}

func Dexec(ctx context.Context, container testcontainers.Container, cmd []string) (*string, error) {
	// Prints the unified command
	fmt.Println(strings.Join(cmd, " "))

	// Executes the unified command
	_, reader, err := container.Exec(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Fetches log output
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}

	// Returns output as a string
	output := buf.String()
	return &output, nil
}

func GetCurrentDir() (*string, error) {
	// Get the file path and line number of the calling function
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get file path")
	}

	// Convert the relative path to an absolute path
	absDirPath, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return nil, err
	}

	// Returns the absolute path
	return &absDirPath, nil
}