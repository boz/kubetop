package util_test

import (
	"testing"
	"time"

	"github.com/boz/kubetop/ui/util"
	"github.com/stretchr/testify/assert"
)

func TestAge(t *testing.T) {
	assert.Equal(t, "5d", util.FormatAge(5*time.Hour*24))
	assert.Equal(t, "5h", util.FormatAge(5*time.Hour))
	assert.Equal(t, "5m", util.FormatAge(5*time.Minute))
	assert.Equal(t, "5s", util.FormatAge(5*time.Second))
	assert.Equal(t, "0.5s", util.FormatAge(time.Second/2))
}
