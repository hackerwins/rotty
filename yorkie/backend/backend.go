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

package backend

import (
	"github.com/yorkie-team/yorkie/pkg/document/time"
	"github.com/yorkie-team/yorkie/pkg/sync"
	"github.com/yorkie-team/yorkie/yorkie/backend/mongo"
	"github.com/yorkie-team/yorkie/yorkie/pubsub"
)

type Config struct {
	// SnapshotThreshold is the threshold that determines if changes should be
	// sent with snapshot when the number of changes is greater than this value.
	SnapshotThreshold uint64 `json:"SnapshotThreshold"`
}

// Backend manages Yorkie's remote states such as data store, distributed lock
// and etc.
type Backend struct {
	Config   *Config
	Mongo    *mongo.Client
	mutexMap *sync.MutexMap
	pubSub   *pubsub.PubSub
}

// New creates a new instance of Backend.
func New(conf *Config, mongoConf *mongo.Config) (*Backend, error) {
	client, err := mongo.NewClient(mongoConf)
	if err != nil {
		return nil, err
	}

	return &Backend{
		Config:   conf,
		Mongo:    client,
		mutexMap: sync.NewMutexMap(),
		pubSub:   pubsub.NewPubSub(),
	}, nil
}

// Close closes all resources of this instance.
func (b *Backend) Close() error {
	if err := b.Mongo.Close(); err != nil {
		return err
	}

	return nil
}

func (b *Backend) Lock(k string) error {
	return b.mutexMap.Lock(k)
}

func (b *Backend) Unlock(k string) error {
	return b.mutexMap.Unlock(k)
}

func (b *Backend) Subscribe(actor *time.ActorID, topics []string) (*pubsub.Subscription, error) {
	return b.pubSub.Subscribe(actor, topics)
}

func (b *Backend) Unsubscribe(topics []string, subscription *pubsub.Subscription) {
	b.pubSub.Unsubscribe(topics, subscription)
}

func (b *Backend) Publish(actor *time.ActorID, topic string, event pubsub.Event) {
	b.pubSub.Publish(actor, topic, event)
}
