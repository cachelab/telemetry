package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	svc := Service{}

	err := svc.Init()
	assert.Equal(t, err, nil)
}
