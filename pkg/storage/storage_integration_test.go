package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var datastore *Datastore

func TestMain(m *testing.M) {
	datastore = setupRedisGraph()
	code := m.Run()
	os.Exit(code)
}

func TestInsertNewApp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	err := datastore.Upsert(App{Name: "Los Appos", Version: "1", Libs: []Library{{Name: "mylib", Version: "1.2.3"}}})
	if err != nil {
		t.Errorf("error while querying redis: %v", err)
	}
}

func TestLibrariesAreNotDuplicated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	_ = datastore.Upsert(App{Name: "firstapp", Version: "1", Libs: []Library{{Name: "mylib", Version: "same"}}})
	_ = datastore.Upsert(App{Name: "secondapp", Version: "2", Libs: []Library{{Name: "mylib", Version: "same"}}})
	_ = datastore.Upsert(App{Name: "thirdapp", Version: "3", Libs: []Library{{Name: "mylib", Version: "different"}}})

	apps, _ := datastore.AppsUsing(Library{Name: "mylib", Version: "same"})
	if len(apps) != 2 {
		t.Errorf("found to many nodes (%v), duplication is going on", len(apps))
	}
}

func TestAppsAreNotDuplicated(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow test")
	}

	_ = datastore.Upsert(App{Name: "sameapp", Version: "same", Libs: []Library{{Name: "samelib", Version: "same"}}})
	_ = datastore.Upsert(App{Name: "sameapp", Version: "same", Libs: []Library{{Name: "samelib", Version: "same"}}})
	_ = datastore.Upsert(App{Name: "sameapp", Version: "different", Libs: []Library{{Name: "samelib", Version: "same"}}})

	apps, _ := datastore.AppsUsing(Library{Name: "samelib", Version: "same"})
	if len(apps) != 2 {
		t.Errorf("found to many nodes (%v), duplication is going on", len(apps))
	}
}

func setupRedisGraph() *Datastore {
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
		log.Fatalf("error while setting up redis: %v", err)
	}

	ip, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("error while setting up redis: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatalf("error while setting up redis: %v", err)
	}

	hostAndPort := fmt.Sprintf("%s:%s", ip, mappedPort.Port())
	_ = os.Setenv("REDIS_ADDRESS", hostAndPort)
	_ = os.Setenv("REDIS_PASSWORD", "")

	ds, err := NewDatastore()
	if err != nil {
		log.Fatalf("error while setting up redis: %v", err)
	}

	return ds
}
