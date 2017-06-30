package grpcservice_test

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/quaponatech/golang-extensions/grpcservice"
)

/* Test Data */

var serverName = "Test Server"

var mainPort int
var portCounter int
var service *grpcservice.GRPCService
var client *grpcservice.GRPCClient

/* Test Funtions */

func TestMain(m *testing.M) {
	flag.Parse()

	// Do setup
	log.Println("GRPC Service: TestMain SetUp")

	rand.Seed(time.Now().Unix())
	mainPort = rand.Intn(65500-10000) + 10000

	// Run tests
	testreturn := m.Run()

	// Do teardown
	log.Println("GRPC Service: TestMain TearDown")

	os.Exit(testreturn)
}
