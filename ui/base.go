package ui

import (
	"github.com/boz/kubetop/backend"
	"github.com/boz/kubetop/util"
)

type wbase struct {
	backend backend.Backend
	env     util.Env
}
