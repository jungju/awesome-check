package main

import "testing"

func TestMustExecuteReadmemd(t *testing.T) {
	mapService := map[string]*Service{
		"test1": &Service{
			Name:            "Test1",
			StargazersCount: 103,
			LicenseName:     "apache",
		},
		"test2": &Service{
			Name:            "Test2",
			StargazersCount: 2,
			LicenseName:     "MIT",
		},
	}

	mustExecuteReadmemd(mapService)
}
