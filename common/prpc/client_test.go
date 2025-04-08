package prpc

import (
	"testing"

	"github.com/HuanXin-Chen/MyIM/common/config"

	ptrace "github.com/HuanXin-Chen/MyIM/common/prpc/trace"
	"github.com/stretchr/testify/assert"
)

func TestNewPClient(t *testing.T) {
	config.Init("../../im.yaml")
	ptrace.StartAgent()
	defer ptrace.StopAgent()

	_, err := NewPClient("im_server")
	assert.NoError(t, err)
}
