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

package types

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yorkie-team/yorkie/api"
	"github.com/yorkie-team/yorkie/api/converter"
	"github.com/yorkie-team/yorkie/pkg/document/change"
	"github.com/yorkie-team/yorkie/pkg/document/operation"
	"github.com/yorkie-team/yorkie/pkg/document/time"
	"github.com/yorkie-team/yorkie/pkg/log"
)

// ChangeInfo is a structure representing information of a change.
type ChangeInfo struct {
	DocID      primitive.ObjectID `bson:"doc_id"`
	ServerSeq  uint64             `bson:"server_seq"`
	ClientSeq  uint32             `bson:"client_seq"`
	Lamport    uint64             `bson:"lamport"`
	Actor      primitive.ObjectID `bson:"actor"`
	Message    string             `bson:"message"`
	Operations [][]byte           `bson:"operations"`
}

// EncodeOperations encodes the given operations into bytes array.
func EncodeOperations(operations []operation.Operation) ([][]byte, error) {
	var encodedOps [][]byte

	for _, pbOp := range converter.ToOperations(operations) {
		encodedOp, err := pbOp.Marshal()
		if err != nil {
			log.Logger.Error(err)
			return nil, errors.New("fail to encode operation")
		}
		encodedOps = append(encodedOps, encodedOp)
	}

	return encodedOps, nil
}

// EncodeActorID encodes the given ActorID into object ID.
func EncodeActorID(id *time.ActorID) primitive.ObjectID {
	objectID := primitive.ObjectID{}
	copy(objectID[:], id[:])
	return objectID
}

// ToChange creates Change model from this ChangeInfo.
func (i *ChangeInfo) ToChange() (*change.Change, error) {
	actorID := time.ActorID{}
	copy(actorID[:], i.Actor[:])
	changeID := change.NewID(i.ClientSeq, i.Lamport, &actorID)

	var pbOps []*api.Operation
	for _, bytesOp := range i.Operations {
		pbOp := api.Operation{}
		if err := pbOp.Unmarshal(bytesOp); err != nil {
			log.Logger.Error(err)
			return nil, err
		}
		pbOps = append(pbOps, &pbOp)
	}

	c := change.New(changeID, i.Message, converter.FromOperations(pbOps))
	c.SetServerSeq(i.ServerSeq)

	return c, nil
}
