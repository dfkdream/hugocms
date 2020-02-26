package main

import (
	"os/exec"
)

type hugo struct {
	c chan chan hugoResult
}

type hugoResult struct {
	Result string `json:"result"`
	Code   int    `json:"code"`
	err    error
}

func newHugo(cfg *config) *hugo {
	ch := make(chan chan hugoResult)
	go func() {
		for r := range ch {
			cmd := exec.Command("hugo", "-s", cfg.Dir)
			exitCode := -1
			res, err := cmd.CombinedOutput()
			if err == nil {
				exitCode = cmd.ProcessState.ExitCode()
			}
			r <- hugoResult{Result: string(res), Code: exitCode, err: err}
		}
	}()
	return &hugo{c: ch}
}

func (h *hugo) build() hugoResult {
	r := make(chan hugoResult)
	h.c <- r
	return <-r
}
