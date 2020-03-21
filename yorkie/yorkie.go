/*
 * Copyright 2020 The Yorkie Authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package yorkie

import (
	"sync"

	"github.com/yorkie-team/yorkie/yorkie/backend"
	"github.com/yorkie-team/yorkie/yorkie/rpc"
)

type Yorkie struct {
	lock sync.Mutex

	backend   *backend.Backend
	rpcServer *rpc.Server

	shutdown   bool
	shutdownCh chan struct{}
	config     *Config
}

func New(conf *Config) (*Yorkie, error) {
	be, err := backend.New(conf.Mongo)
	if err != nil {
		return nil, err
	}

	rpcServer, err := rpc.NewRPCServer(conf.RPCPort, be)
	if err != nil {
		return nil, err
	}

	return &Yorkie{
		backend:    be,
		rpcServer:  rpcServer,
		shutdownCh: make(chan struct{}),
		config:     conf,
	}, nil
}

func (r *Yorkie) Start() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.rpcServer.Start()
}

func (r *Yorkie) Shutdown(graceful bool) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.shutdown {
		return nil
	}

	if err := r.backend.Close(); err != nil {
		return err
	}

	r.rpcServer.Shutdown(graceful)

	close(r.shutdownCh)
	r.shutdown = true
	return nil
}

func (r *Yorkie) ShutdownCh() <-chan struct{} {
	return r.shutdownCh
}

func (r *Yorkie) RPCAddr() string {
	return r.config.RPCAddr()
}
