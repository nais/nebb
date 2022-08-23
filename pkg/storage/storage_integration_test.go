package storage

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInsertApp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	err := setupRedisGraph()
	if err != nil {
		t.Fatal(err)
	}

	datastore, _ := NewDatastore()
	err = datastore.Upsert(App{
		Name:    "Los Appos",
		Version: "1",
		Libs: []Library{
			{
				Name:    "mylib",
				Version: "1.2.3",
			},
		},
	})
	if err != nil {
		t.Errorf("error while querying redis: %v", err)
	}
}

func setupRedisGraph() error {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redislabs/redisgraph",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return err
	}

	hostAndPort := fmt.Sprintf("%s:%s", ip, mappedPort.Port())
	_ = os.Setenv("REDIS_ADDRESS", hostAndPort)
	_ = os.Setenv("REDIS_PASSWORD", "")

	return nil
}
