package hugo

import (
	"os/exec"

	"github.com/dfkdream/hugocms/config"
)

type Hugo struct {
	c chan chan Result
}

type Result struct {
	Result string `json:"result"`
	Code   int    `json:"code"`
	Err    error  `json:"-"`
}

func New(cfg *config.Config) *Hugo {
	ch := make(chan chan Result)
	go func() {
		for r := range ch {
			cmd := exec.Command("hugo", "-s", cfg.Dir)
			exitCode := -1
			res, err := cmd.CombinedOutput()
			if err == nil {
				exitCode = cmd.ProcessState.ExitCode()
			}
			r <- Result{Result: string(res), Code: exitCode, Err: err}
		}
	}()
	return &Hugo{c: ch}
}

func (h *Hugo) Build() Result {
	r := make(chan Result)
	h.c <- r
	return <-r
}
