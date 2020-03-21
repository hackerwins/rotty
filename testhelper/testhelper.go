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
	"time"

	"github.com/yorkie-team/yorkie/yorkie"
)

const (
	TestPort               = 1101
	TestMongoConnectionURI = "mongodb://localhost:27017"
)

// TestDBName returns the name of test database with timestamp.
func TestDBName() string {
	now := time.Now()
	return fmt.Sprintf("test-%s-%d", yorkie.DefaultYorkieDatabase, now.Unix())
}
