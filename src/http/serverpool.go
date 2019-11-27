package http

import "github.com/shaojintian/load_balancer/src/core"

type ServerPool struct {

	backends	[]*core.Backend
}