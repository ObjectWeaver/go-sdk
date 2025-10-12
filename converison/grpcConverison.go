package converison

import (
	pb "github.com/firechimp-org/go-sdk/grpc"
	"github.com/firechimp-org/go-sdk/jsonSchema"
)

// ConvertProtoToModel converts a protobuf Definition to your Go model Definition
func ConvertProtoToModel(protoDef *pb.Definition) *jsonSchema.Definition {
	if protoDef == nil {
		return nil
	}

	modelDef := &jsonSchema.Definition{
		Type:               jsonSchema.DataType(protoDef.Type),
		Instruction:        protoDef.Instruction,
		Properties:         make(map[string]jsonSchema.Definition),
		Items:              ConvertProtoToModel(protoDef.GetItems()), // Use Getters to handle nil cases
		Model:              protoDef.Model,
		ProcessingOrder:    protoDef.ProcessingOrder,
		SystemPrompt:       getStringPointer(protoDef.GetSystemPrompt()), // Safe getter for pointers
		ImprovementProcess: protoDef.ImprovementProcess,
		SelectFields:       protoDef.SelectFields,
		Voters:             protoDef.Voters,
		HashMap:            ConvertProtoToHashMap(protoDef.GetHashMap()),   // Check with Getters
		NarrowFocus:        ConvertProtoToFocus(protoDef.GetNarrowFocus()), // Handle nil safely
		Req:                ConvertProtoToRequestFormat(protoDef.GetReq()),
		Choices:            ConvertProtoToChoices(protoDef.GetChoices()),
		SpeechToText:       convertProtoSpeechToText(protoDef.GetSpeechToText()), // Safely handle nested structs
		TextToSpeech:       convertProtoTextToSpeech(protoDef.GetTextToSpeech()),
		SendImage:          convertProtoSendImage(protoDef.GetSendImage()), // Handle nil structs
		Stream:             protoDef.Stream,
	}

	// Handle Properties map
	if protoDef.Properties != nil {
		for key, protoProperty := range protoDef.Properties {
			modelDef.Properties[key] = *ConvertProtoToModel(protoProperty)
		}
	}

	return modelDef
}

// Helper function to safely get string pointers
func getStringPointer(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

// ConvertModelToProto converts your Go model Definition to a protobuf Definition
func ConvertModelToProto(modelDef *jsonSchema.Definition) *pb.Definition {
	if modelDef == nil {
		return nil
	}

	systemPrompt := ""
	if modelDef.SystemPrompt != nil {
		systemPrompt = *modelDef.SystemPrompt
	}

	protoDef := &pb.Definition{
		Type:               string(modelDef.Type),
		Instruction:        modelDef.Instruction,
		Properties:         make(map[string]*pb.Definition),
		Items:              ConvertModelToProto(modelDef.Items),
		Model:              modelDef.Model,
		ProcessingOrder:    modelDef.ProcessingOrder,
		SystemPrompt:       systemPrompt,
		ImprovementProcess: modelDef.ImprovementProcess,
		SelectFields:       modelDef.SelectFields,
		Voters:             modelDef.Voters,
		HashMap:            ConvertModelToProtoHashMap(modelDef.HashMap),
		NarrowFocus:        ConvertModelToProtoFocus(modelDef.NarrowFocus),
		Req:                ConvertModelToProtoRequestFormat(modelDef.Req),
		Choices:            ConvertModelToProtoChoices(modelDef.Choices),
		Image:              convertModelImage(modelDef.Image),
		SpeechToText:       convertModelSpeechToText(modelDef.SpeechToText),
		TextToSpeech:       convertModelTextToSpeech(modelDef.TextToSpeech),
		SendImage:          convertModelSendImage(modelDef.SendImage),
		Stream:             modelDef.Stream,
	}

	// Handle Properties map
	if modelDef.Properties != nil {
		for key, modelProperty := range modelDef.Properties {
			protoDef.Properties[key] = ConvertModelToProto(&modelProperty)
		}
	}

	return protoDef
}

// Helper functions for SpeechToText, TextToSpeech, and other nested structs

func convertProtoSpeechToText(speechToText *pb.SpeechToText) *jsonSchema.SpeechToText {
	if speechToText == nil {
		return nil
	}
	return &jsonSchema.SpeechToText{
		Model:             jsonSchema.SpeechToTextModel(speechToText.Model),
		AudioToTranscribe: speechToText.AudioToTranscribe,
		Language:          speechToText.Language,
		ToString:          speechToText.ToString,
		ToCaptions:        speechToText.ToCaptions,
		Format:            jsonSchema.AudioFormat(speechToText.Format),
	}
}

func convertProtoTextToSpeech(textToSpeech *pb.TextToSpeech) *jsonSchema.TextToSpeech {
	if textToSpeech == nil {
		return nil
	}
	return &jsonSchema.TextToSpeech{
		Model:         jsonSchema.TextToSpeechModel(textToSpeech.Model),
		Voice:         jsonSchema.Voice(textToSpeech.Voice),
		StringToAudio: textToSpeech.StringToAudio,
		Format:        jsonSchema.AudioFormat(textToSpeech.Format),
	}
}

func convertProtoSendImage(sendImage *pb.SendImage) *jsonSchema.SendImage {
	if sendImage == nil {
		return nil
	}
	return &jsonSchema.SendImage{
		ImagesData: sendImage.ImagesData,
	}
}

func convertModelSpeechToText(speechToText *jsonSchema.SpeechToText) *pb.SpeechToText {
	if speechToText == nil {
		return nil
	}
	return &pb.SpeechToText{
		Model:             string(speechToText.Model),
		AudioToTranscribe: speechToText.AudioToTranscribe,
		Language:          speechToText.Language,
		ToString:          speechToText.ToString,
		ToCaptions:        speechToText.ToCaptions,
		Format:            string(speechToText.Format),
	}
}

func convertModelTextToSpeech(textToSpeech *jsonSchema.TextToSpeech) *pb.TextToSpeech {
	if textToSpeech == nil {
		return nil
	}
	return &pb.TextToSpeech{
		Model:         string(textToSpeech.Model),
		Voice:         string(textToSpeech.Voice),
		StringToAudio: textToSpeech.StringToAudio,
		Format:        string(textToSpeech.Format),
	}
}

func convertModelSendImage(sendImage *jsonSchema.SendImage) *pb.SendImage {
	if sendImage == nil {
		return nil
	}
	return &pb.SendImage{
		ImagesData: sendImage.ImagesData,
	}
}

func convertModelImage(image *jsonSchema.Image) *pb.Image {
	if image == nil {
		return nil
	}
	return &pb.Image{
		Model: string(image.Model),
		Size:  string(image.Size),
	}
}
