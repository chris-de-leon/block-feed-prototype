package testutils

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DEFAULT_MYSQL_ROOT_PASSWORD = "password"
	DEFAULT_MYSQL_ROOT_USERNAME = "root"
	DEFAULT_MYSQL_DB            = "test"
	MYSQL_PORT                  = "3306/tcp"
	REDIS_PORT                  = "6379/tcp"
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

func NewMySqlContainer(ctx context.Context, t *testing.T, version string) (*ContainerWithConnectionInfo, error) {
	// Gets the directory that this file exists in
	dir, err := GetCurrentDir()
	if err != nil {
		return nil, err
	}

	// Creates the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			ExposedPorts: []string{MYSQL_PORT},
			WaitingFor:   wait.ForExposedPort(),
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    path.Join(*dir, "..", "..", "..", "db"),
				Dockerfile: path.Join("Dockerfile"),
				BuildArgs: map[string]*string{
					"MYSQL_VERSION": &version,
				},
			},
			Env: map[string]string{
				"MYSQL_ROOT_PASSWORD": DEFAULT_MYSQL_ROOT_PASSWORD,
				"MYSQL_DATABASE":      DEFAULT_MYSQL_DB,
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
	conn, err := GetConnectionInfo(ctx, container, nat.Port(MYSQL_PORT))
	if err != nil {
		return nil, err
	}

	// The default URL allows superuser access
	conn.Url = MySqlUrl(*conn, DEFAULT_MYSQL_ROOT_USERNAME, DEFAULT_MYSQL_ROOT_PASSWORD)

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

func MySqlUrl(conn HostConnectionInfo, uname string, pword string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		uname,
		pword,
		"host.docker.internal",
		conn.Port.Port(),
		DEFAULT_MYSQL_DB,
	)
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
	// Cleans up the container when the test case finishes
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
