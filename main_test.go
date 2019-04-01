package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var services = []*Service{
	&Service{
		Name:            "Test1",
		StargazersCount: 103,
		LicenseName:     "apache",
	},
	&Service{
		Name:            "Test2",
		StargazersCount: 2,
		LicenseName:     "MIT",
	},
	&Service{
		Name:            "Test3",
		StargazersCount: 10,
		LicenseName:     "BSP",
	},
}

func TestMustExecuteReadmemd(t *testing.T) {
	mustExecuteReadmemd(services)
}
func TestSortServices(t *testing.T) {
	sortService(services)

	assert.Equal(t, "Test1", services[0].Name)
	assert.Equal(t, "Test3", services[1].Name)
	assert.Equal(t, "Test2", services[2].Name)
}
