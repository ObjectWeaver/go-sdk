package jsonSchema

// Definition is a struct for describing a JSON Schema.
// It is fairly limited, and you may have better luck using a third-party library.
type Definition struct {
	// Type specifies the data type of the schema.
	Type DataType `json:"type,omitempty"`
	// Instruction is the instruction for what to generate.
	Instruction string `json:"instruction,omitempty"`
	// Properties describes the properties of an object, if the schema type is Object.
	Properties map[string]Definition `json:"properties"`

	// Items specifies which data type an array contains, if the schema type is Array.
	Items *Definition `json:"items,omitempty"`
	// Model this needs to match the exact name of the model that you will be sending to AI as a Service provider. Such as OpenAI or Google Gemini
	Model string `json:"model,omitempty"`
	// ProcessingOrder this is the order of strings ie the fields of the parent property keys that need to be processed first before this field is processed
	ProcessingOrder []string `json:"processingOrder,omitempty"`

	// SystemPrompt allows the developer to spefificy their own system prompt so the processing. It operates current at the properties level.
	SystemPrompt *string `json:"systemPrompt,omitempty"`
	// Req allows a developer to send out a request at a given point to ensure that additional information can be extracted from external databases

	// ImprovementProcess --> so that the user can specify when a super high quality completion is needed and it can be improved upon
	ImprovementProcess bool `json:"improvementProcess,omitempty"`

	//Map is used here as so that a map of values can be created and then returned -- useful in the instruction creation process -- not sure how useful it is otherwise
	HashMap *HashMap

	//the other data types that need to be filled for the object to be generated within GoR
	TextToSpeech *TextToSpeech `json:"textToSpeech,omitempty"`
	SpeechToText *SpeechToText `json:"speechToText,omitempty"`
	Image        *Image        `json:"image,omitempty"`

	//Utility fields:
	Req *RequestFormat `json:"req,omitempty"`
	// NarrowFocus
	NarrowFocus *Focus `json:"narrowFocus,omitempty"`

	// SelectFields has the aim of being able to select multiple pieice of information and when they are all present then continue with processing. Such that the selection of information can work like so:
	//The system works as an absolute path that has to be selected. So starting from the top most object then down to the selected field(s)
	//"car.color" --> this would fetch the information from the car field and then the color field.
	//"cars.color" --> Would return the entire list of colours that have been generated so far
	SelectFields []string `json:"selectFields,omitempty"`

	// Choices For determining which of the property fields should be generated
	Choices *Choices `json:"choices,omitempty"`

	// Voters this is used for determining whether you want to have voters determine the qulaity of completions. Increases costs but improves quality. If avialible to your tier then turned on automatically.
	Voters bool `json:"voters,omitempty"`

	//Image URL --> if the LLM supports reading an image due to it being multi-model then the image URL will be passed in here
	SendImage *SendImage `json:"sendImage,omitempty"`

	//Stream - used for instructing when the information should be streamed. Please visit the documentation for more information for which types are supported.
	Stream bool `json:"stream,omitempty"`

	//Temp - used for passing in a temperature value for the prompt request
	Temp float64 `json:"temp,omitempty"`

	//OverridePrompt - used for overriding the prompt that is passed in
	OverridePrompt *string `json:"overridePrompt,omitempty"`

	//Priority - used for setting the priority of the request
	Priority int32 `json:"priority,omitempty"`
}

type Choices struct {
	Number  int      `json:"number,omitempty"`  //this denotes the number of choices that should be selected
	Options []string `json:"options,omitempty"` //this is the list of fields that will be chosen from
	/*
		How this works is that it needs to be in a of definitions which match with the properties field. From the properties fields the choice of those keys will be selected
		the information of what the overall object, the properties being selected along with the instruction and their type and the types that they contain if the object goes down further.
		the prompt will also be pass in so that the agent can make the best decesion possible

		Once the choices have been selected the choices that haven't been selected will be deleted from the remaining keys avialible in both the ordered and unordered keys.
	*/
}

// HashMap this can output a map of values and so whilst it may take up a single field it could output many fields
type HashMap struct {
	KeyInstruction  string      `json:"keyInstruction,omitempty"`
	FieldDefinition *Definition `json:"fieldDefinition,omitempty"`
}

// Focus the idea for this is so that when a narrow focus request needs to be sent out to an LLM without needing all the additional information. From prior generation.
type Focus struct {
	Prompt string `json:"prompt"`
	//the fields value denotes the properties that will be extracted from the properties fields. These will only operate at a single level of generation.
	//the order in which the fields that are listed will be the order for which the currently generated information will be presented below the prompt value.
	Fields []string `json:"fields"`

	//KeepOriginal -- for keeping the original prompt in cases for lists where it would otherwise be removed from the context
	KeepOriginal bool `json:"keepOriginal,omitempty"`
}

// RequestFormat defines the structure of the request
type RequestFormat struct {
	URL           string                 `json:"url"`
	Method        HTTPMethod             `json:"method"`
	Headers       map[string]string      `json:"headers,omitempty"`
	Body          map[string]interface{} `json:"body,omitempty"`
	Authorization string                 `json:"authorization,omitempty"`
	RequireFields []string               `json:"requirFields,omitempty"`
}

const (
	UrgentPriority   = int32(2)  //the highest priority and will take precedence over all other requests
	StandardPriority = int32(1)  //standard priority for most requests
	LowPriority      = int32(0)  //low priority requests
	EventualPriority = int32(-1) //for batch LLM usage
)
