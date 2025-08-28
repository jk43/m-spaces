package models

// import (
// 	"log"
// 	"os"
// 	"testing"

// 	"github.com/moly-space/molylibs"
// 	"github.com/ory/dockertest/v3"
// 	"github.com/ory/dockertest/v3/docker"
// )

// func setEnv() {
// 	os.Setenv("MONGO_USERNAME", "mongo")
// 	os.Setenv("MONGO_PASSWORD", "mongo")
// 	os.Setenv("MONGO_URI", "mongodb://mongo@0.0.0.0:27018")
// 	os.Setenv("MONGO_DATABASE", "dev-user")
// 	database = "dev-user"
// }

// var resource *dockertest.Resource

// // var pool *dockertest.Pool
// var repo *DBRepo

// func TestMain(m *testing.M) {
// 	setEnv()
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("could not connect to docker; is it running? %s", err)
// 	}

// 	opts := dockertest.RunOptions{
// 		Repository: "mongo",
// 		Tag:        "6.0.3",
// 		Env: []string{
// 			"MONGO_INITDB_ROOT_USERNAME=" + os.Getenv("MONGO_USERNAME"),
// 			"MONGO_INITDB_ROOT_PASSWORD=" + os.Getenv("MONGO_PASSWORD"),
// 			"MONGO_INITDB_DATABASE=" + os.Getenv("MONGO_DATABASE"),
// 		},
// 		ExposedPorts: []string{"27017"},
// 		PortBindings: map[docker.Port][]docker.PortBinding{
// 			"27017": {
// 				{HostIP: "0.0.0.0", HostPort: "27018"},
// 			},
// 		},
// 	}

// 	// get a resource (docker image)
// 	resource, err = pool.RunWithOptions(&opts)
// 	if err != nil {
// 		_ = pool.Purge(resource)
// 		log.Fatalf("could not start resource: %s", err)
// 	}

// 	//start the image and wait until it's ready
// 	if err := pool.Retry(func() error {
// 		var err error
// 		mongoClient, err := molylibs.Mongo()
// 		if err != nil {
// 			log.Println("Error:", err)
// 			return err
// 		}
// 		repo = &DBRepo{
// 			Mongo: mongoClient,
// 		}
// 		return err
// 	}); err != nil {
// 		_ = pool.Purge(resource)
// 		log.Fatalf("could not connect to database: %s", err)
// 	}

// 	code := m.Run()
// 	//clean up
// 	if err := pool.Purge(resource); err != nil {
// 		log.Fatalf("could not purge resource: %s", err)
// 	}

// 	os.Exit(code)
// }

// func Test_InsertUser(t *testing.T) {
// 	user := User{
// 		FirstName: "firstname",
// 		LastName:  "lastname",
// 		Email:     "email",
// 	}

// 	_, err := repo.InsertUser(&user)
// 	if err != nil {
// 		t.Errorf("Got error inserting user %s", err)
// 	}
// }

// func Test_GetUserByEmail(t *testing.T) {
// 	user := repo.GetUserByEmail("email")

// 	if user.Email == "" {
// 		t.Errorf("User not found")
// 	}
// 	if user.FirstName != "firstname" {
// 		t.Errorf("Incorrect first name")
// 	}
// 	if user.LastName != "lastname" {
// 		t.Errorf("Incorrect last name")
// 	}
// }
