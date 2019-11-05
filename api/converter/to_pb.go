package converter

import (
	"github.com/gogo/protobuf/types"

	"github.com/hackerwins/rottie/api"
	"github.com/hackerwins/rottie/pkg/document/change"
	"github.com/hackerwins/rottie/pkg/document/checkpoint"
	"github.com/hackerwins/rottie/pkg/document/json"
	"github.com/hackerwins/rottie/pkg/document/json/datatype"
	"github.com/hackerwins/rottie/pkg/document/key"
	"github.com/hackerwins/rottie/pkg/document/operation"
	"github.com/hackerwins/rottie/pkg/document/time"
)

func ToChangePack(pack *change.Pack) *api.ChangePack {
	return &api.ChangePack{
		DocumentKey: toDocumentKey(pack.DocumentKey),
		Checkpoint:  toCheckpoint(pack.Checkpoint),
		Changes:     toChanges(pack.Changes),
	}
}

func toDocumentKey(key *key.Key) *api.DocumentKey {
	return &api.DocumentKey{
		Collection: key.Collection,
		Document:   key.Document,
	}
}

func toCheckpoint(cp *checkpoint.Checkpoint) *api.Checkpoint {
	return &api.Checkpoint{
		ServerSeq: cp.ServerSeq,
		ClientSeq: cp.ClientSeq,
	}
}

func toChanges(changes []*change.Change) []*api.Change {
	var pbChanges []*api.Change
	for _, c := range changes {
		pbChanges = append(pbChanges, &api.Change{
			Id:         toChangeID(c.ID()),
			Message:    c.Message(),
			Operations: ToOperations(c.Operations()),
		})
	}

	return pbChanges
}

func toChangeID(id *change.ID) *api.ChangeID {
	return &api.ChangeID{
		ClientSeq: id.ClientSeq(),
		Lamport:   id.Lamport(),
		ActorId:   id.Actor().String(),
	}
}

func ToOperations(operations []operation.Operation) []*api.Operation {
	var pbOperations []*api.Operation

	for _, o := range operations {
		pbOperation := &api.Operation{}
		switch op := o.(type) {
		case *operation.Set:
			pbOperation.Body = &api.Operation_Set_{
				Set: &api.Operation_Set{
					Key:             op.Key(),
					Value:           toJSONElement(op.Value()),
					ParentCreatedAt: toTimeTicket(op.ParentCreatedAt()),
					ExecutedAt:      toTimeTicket(op.ExecutedAt()),
				},
			}
		case *operation.Add:
			pbOperation.Body = &api.Operation_Add_{
				Add: &api.Operation_Add{
					Value:           toJSONElement(op.Value()),
					PrevCreatedAt:   toTimeTicket(op.PrevCreatedAt()),
					ParentCreatedAt: toTimeTicket(op.ParentCreatedAt()),
					ExecutedAt:      toTimeTicket(op.ExecutedAt()),
				},
			}
		}
		pbOperations = append(pbOperations, pbOperation)
	}

	return pbOperations
}

func toJSONElement(element datatype.Element) *api.JSONElement {
	switch elem := element.(type) {
	case *json.Object:
		return &api.JSONElement{
			Type:      api.ValueType_JSON_OBJECT,
			CreatedAt: toTimeTicket(element.CreatedAt()),
		}
	case *json.Array:
		return &api.JSONElement{
			Type:      api.ValueType_JSON_ARRAY,
			CreatedAt: toTimeTicket(element.CreatedAt()),
		}
	case *datatype.Primitive:
		return &api.JSONElement{
			Type:      api.ValueType_STRING,
			CreatedAt: toTimeTicket(element.CreatedAt()),
			Value:     toValue(elem.Value()),
		}
	}
	panic("fail to encode JSONElement to protobuf")
}

func toValue(bytes []byte) *types.Any {
	return &types.Any{
		Value: bytes,
	}
}

func toTimeTicket(ticket *time.Ticket) *api.TimeTicket {
	return &api.TimeTicket{
		Lamport:   ticket.Lamport(),
		Delimiter: ticket.Delimiter(),
		ActorId:   ticket.ActorID().String(),
	}
}
