package core

import "testing"
import "github.com/stretchr/testify/assert"

func TestDataPathAllocation(t *testing.T) {
	path := []string{"/disk1", "/disk2", "/disk3"}

	allocation := dataPathAllocation(path, 2, 3)
	assert.Equal(t, "/disk3", allocation[0])

	path = []string{"/disk1", "/disk2", "/disk3", "/disk4", "/disk5", "/disk6", "/disk7", "/disk8"}

	allocation = dataPathAllocation(path, 1, 2)
	except := []string{"/disk5", "/disk6", "/disk7", "/disk8"}
	assert.Equal(t, except, allocation)
}
