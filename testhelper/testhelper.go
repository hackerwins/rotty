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

package testhelper

import (
	"fmt"
	"log"
	"runtime"
	defaultTime "time"

	"github.com/yorkie-team/yorkie/pkg/document/change"
	"github.com/yorkie-team/yorkie/pkg/document/json"
	"github.com/yorkie-team/yorkie/pkg/document/time"
	"github.com/yorkie-team/yorkie/yorkie"
	"github.com/yorkie-team/yorkie/yorkie/backend"
	"github.com/yorkie-team/yorkie/yorkie/backend/mongo"
	"github.com/yorkie-team/yorkie/yorkie/metrics"
	"github.com/yorkie-team/yorkie/yorkie/rpc"
)

var testStartedAt int64

const (
	RPCPort     = 1101
	MetricsPort = 1102

	MongoConnectionURI        = "mongodb://localhost:27017"
	MongoConnectionTimeoutSec = 5
	MongoPingTimeoutSec       = 5

	SnapshotThreshold = 10

	Collection = "test-collection"
)

func init() {
	now := defaultTime.Now()
	testStartedAt = now.Unix()
}

// TestDBName returns the name of test database with timestamp.
// timestamp is set only once on first call.
func TestDBName() string {
	return fmt.Sprintf("test-%s-%d", yorkie.DefaultMongoYorkieDatabase, testStartedAt)
}

// TextChangeContext returns the context of test change.
func TextChangeContext() *change.Context {
	return change.NewContext(
		change.InitialID,
		"",
		json.NewRoot(json.NewObject(json.NewRHTPriorityQueueMap(), time.InitialTicket)),
	)
}

// PrintMemStats prints memory stats.
func PrintMemStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Println("mem.Alloc:", ByteCountIEC(mem.Alloc))
	fmt.Println("mem.TotalAlloc:", ByteCountIEC(mem.TotalAlloc))
	fmt.Println("mem.HeapAlloc:", ByteCountIEC(mem.HeapAlloc))
	fmt.Println("mem.NumGC:", mem.NumGC)
}

func PrintBytesSize(bytes []byte) {
	byteSize := len(bytes)
	fmt.Println("sna.Bytes:", ByteCountIEC(uint64(byteSize)))
}

func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// TestYorkie is return Yorkie instance for testing.
func TestYorkie() *yorkie.Yorkie {
	y, err := yorkie.New(&yorkie.Config{
		RPC: &rpc.Config{
			Port: RPCPort,
		},
		Metrics: &metrics.Config{
			Port: MetricsPort,
		},
		Backend: &backend.Config{
			SnapshotThreshold: SnapshotThreshold,
		},
		Mongo: &mongo.Config{
			ConnectionURI:        MongoConnectionURI,
			ConnectionTimeoutSec: MongoConnectionTimeoutSec,
			PingTimeoutSec:       MongoPingTimeoutSec,
			YorkieDatabase:       TestDBName(),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return y
}
