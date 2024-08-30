package converison

import (
	pb "github.com/henrylamb/object-generation-golang/grpc"
	"github.com/henrylamb/object-generation-golang/jsonSchema"
)

// ConvertProtoToHashMap converts a protobuf HashMap to the Go model HashMap
func ConvertProtoToHashMap(protoHashMap *pb.HashMap) *jsonSchema.HashMap {
	if protoHashMap == nil {
		return nil
	}

	return &jsonSchema.HashMap{
		KeyInstruction:  protoHashMap.KeyInstruction,
		FieldDefinition: ConvertProtoToModel(protoHashMap.FieldDefinition),
	}
}

// ConvertModelToProtoHashMap converts a Go model HashMap to a protobuf HashMap
func ConvertModelToProtoHashMap(modelHashMap *jsonSchema.HashMap) *pb.HashMap {
	if modelHashMap == nil {
		return nil
	}

	return &pb.HashMap{
		KeyInstruction:  modelHashMap.KeyInstruction,
		FieldDefinition: ConvertModelToProto(modelHashMap.FieldDefinition),
	}
}
