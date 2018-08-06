package collection

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCollection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collection Suite")
}

type TestTSD struct {
	timestamp int64
}

func (t TestTSD) Timestamp() int64 {
	return t.timestamp
}
