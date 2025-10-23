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

	// DecisionPoint enables conditional branching based on scores or field values.
	// After this field is generated, the DecisionPoint evaluates conditions and routes to different Definitions.
	// This enables complex decision trees where the LLM's output quality or content determines the next generation step.
	DecisionPoint *DecisionPoint `json:"decisionPoint,omitempty"`

	// ScoringCriteria defines how to evaluate the quality of this field's generation.
	// Scores are calculated by an LLM evaluation and can be used by DecisionPoint for routing decisions.
	// Example: Rate technical accuracy (0-100), readability (0-100), etc.
	ScoringCriteria *ScoringCriteria `json:"scoringCriteria,omitempty"`

	// RecursiveLoop enables iterative refinement of this field's generation.
	// The field will be regenerated multiple times, with each iteration scored and potentially improved.
	// Useful for high-quality content that benefits from multiple attempts and selection of the best result.
	RecursiveLoop *RecursiveLoop `json:"recursiveLoop,omitempty"`
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

// DecisionPoint enables routing to different Definitions based on conditions.
// This allows for complex decision trees where the path forward depends on scores or field values.
// Example: "if technical_accuracy > 70 AND readability < 80 THEN improve_readability"
type DecisionPoint struct {
	// Name identifies this decision point for debugging and logging purposes.
	Name string `json:"name,omitempty"`

	// EvaluationPrompt guides the LLM in scoring the current generation.
	// This prompt should instruct the LLM on what dimensions to evaluate and how to score them.
	// Example: "Rate this content on technical_accuracy (0-100) and readability (0-100)"
	EvaluationPrompt string `json:"evaluationPrompt,omitempty"`

	// Branches define conditional paths that are evaluated in order.
	// The first branch whose conditions all evaluate to true will be executed.
	// If multiple branches could match, prioritize them using the Priority field.
	Branches []ConditionalBranch `json:"branches"`

	// Default fallback Definition if no branch conditions match.
	// If nil and no branches match, processing continues without branching.
	Default *Definition `json:"default,omitempty"`

	// Strategy determines how conditions are evaluated:
	// - RouteByScore: Use LLM evaluation scores from ScoringCriteria
	// - RouteByField: Use values from SelectFields
	// - RouteByHybrid: Combine both score-based and field-based conditions
	Strategy RoutingStrategy `json:"strategy,omitempty"`
}

// ConditionalBranch represents an if/then rule in the decision tree.
// All conditions in a branch must be true (AND logic) for the branch to execute.
type ConditionalBranch struct {
	// Name for debugging and logging. Helps identify which branch was taken.
	Name string `json:"name,omitempty"`

	// Conditions that must ALL be true for this branch to execute (AND logic).
	// Example: [{field: "accuracy", operator: "gt", value: 70}, {field: "readability", operator: "lt", value: 80}]
	// This creates: "if accuracy > 70 AND readability < 80"
	Conditions []Condition `json:"conditions"`

	// Then is the Definition to generate if all conditions match.
	// This Definition can have its own DecisionPoint, enabling nested decision trees.
	Then Definition `json:"then"`

	// Priority determines evaluation order when multiple branches could match.
	// Higher priority branches are evaluated first. Default is 0.
	Priority int `json:"priority,omitempty"`
}

// Condition represents a single evaluation rule in a conditional expression.
// Conditions can check scores (from ScoringCriteria) or field values (from SelectFields).
type Condition struct {
	// Field is the name of the score or value to check.
	// For scores: matches dimension names in ScoringCriteria (e.g., "technical_accuracy")
	// For fields: matches the final key in a SelectFields path (e.g., "is_technical")
	Field string `json:"field"`

	// Operator defines the comparison operation to perform.
	// Supported: eq, neq, gt, lt, gte, lte, in, nin, contains
	Operator ComparisonOperator `json:"operator"`

	// Value to compare the field against.
	// Type should match the field type (int for scores, bool/string for fields, etc.)
	Value interface{} `json:"value"`

	// FieldPath is the full path for nested field access when using SelectFields.
	// Example: "requirements.is_technical" to access the is_technical field from requirements object.
	// If not specified, Field is used directly.
	FieldPath string `json:"fieldPath,omitempty"`
}

// ComparisonOperator defines the type of comparison to perform in a Condition.
type ComparisonOperator string

const (
	OpEqual              ComparisonOperator = "eq"       // Equal (==)
	OpNotEqual           ComparisonOperator = "neq"      // Not equal (!=)
	OpGreaterThan        ComparisonOperator = "gt"       // Greater than (>)
	OpLessThan           ComparisonOperator = "lt"       // Less than (<)
	OpGreaterThanOrEqual ComparisonOperator = "gte"      // Greater than or equal (>=)
	OpLessThanOrEqual    ComparisonOperator = "lte"      // Less than or equal (<=)
	OpIn                 ComparisonOperator = "in"       // Value is in a list
	OpNotIn              ComparisonOperator = "nin"      // Value is not in a list
	OpContains           ComparisonOperator = "contains" // String contains substring
)

// RoutingStrategy determines how DecisionPoint evaluates conditions.
type RoutingStrategy string

const (
	// RouteByScore uses LLM evaluation scores from ScoringCriteria.
	// The LLM evaluates the generated content and produces numeric/categorical scores.
	RouteByScore RoutingStrategy = "score"

	// RouteByField uses values extracted via SelectFields.
	// Checks boolean flags, categorical values, or other generated field values.
	RouteByField RoutingStrategy = "field"

	// RouteByHybrid combines both score-based and field-based conditions.
	// Allows complex rules like "if is_technical=true AND accuracy > 70"
	RouteByHybrid RoutingStrategy = "hybrid"
)

// ScoringCriteria defines how to evaluate the quality of a generated field.
// The LLM acts as a judge, scoring the output across multiple dimensions.
// These scores can then be used by DecisionPoint for conditional branching or by RecursiveLoop for termination.
type ScoringCriteria struct {
	// Dimensions map score names to their evaluation definitions.
	// Each dimension represents one aspect to evaluate.
	// Example: {"technical_accuracy": {...}, "readability": {...}, "engagement": {...}}
	Dimensions map[string]ScoringDimension `json:"dimensions"`

	// EvaluationModel specifies which LLM model to use for scoring.
	// If empty, uses the default model. Consider using a faster/cheaper model for evaluation.
	EvaluationModel string `json:"evaluationModel,omitempty"`

	// AggregationMethod determines how to combine multiple dimension scores into a single score.
	// Used when you need one overall quality metric from multiple dimensions.
	AggregationMethod AggregationMethod `json:"aggregationMethod,omitempty"`
}

// ScoringDimension defines a single evaluation metric for quality assessment.
type ScoringDimension struct {
	// Description guides the LLM on how to evaluate this dimension.
	// Should be clear and specific about what constitutes high vs low scores.
	// Example: "Rate technical accuracy from 0-100, where 0 is completely inaccurate and 100 is perfectly accurate"
	Description string `json:"description"`

	// Scale defines the numeric range for scores (e.g., 0-100, 1-10).
	// Only applicable for numeric scores. Helps LLM understand the expected range.
	Scale *ScoreScale `json:"scale,omitempty"`

	// Type specifies the kind of score: numeric (0-100), boolean (true/false), or categorical ("low"/"medium"/"high").
	Type ScoreType `json:"type,omitempty"`

	// Weight for aggregation when combining multiple dimensions (0.0-1.0).
	// Higher weight = more influence on the aggregate score. Weights should sum to 1.0.
	Weight float64 `json:"weight,omitempty"`
}

// ScoreScale defines the numeric range for a scoring dimension.
type ScoreScale struct {
	Min int `json:"min"` // Minimum score value (e.g., 0)
	Max int `json:"max"` // Maximum score value (e.g., 100)
}

// ScoreType defines the data type of a score.
type ScoreType string

const (
	ScoreNumeric     ScoreType = "numeric"     // Numeric score within a scale (e.g., 0-100)
	ScoreBoolean     ScoreType = "boolean"     // Binary true/false evaluation
	ScoreCategorical ScoreType = "categorical" // Categorical value (e.g., "low", "medium", "high")
)

// AggregationMethod defines how to combine multiple dimension scores.
type AggregationMethod string

const (
	// AggregateWeightedAverage computes weighted average of all dimension scores.
	// Each dimension's weight determines its influence on the final score.
	AggregateWeightedAverage AggregationMethod = "weighted_average"

	// AggregateMinimum uses the lowest score across all dimensions.
	// Useful when all dimensions must meet a minimum threshold.
	AggregateMinimum AggregationMethod = "minimum"

	// AggregateMaximum uses the highest score across all dimensions.
	// Useful when any dimension meeting threshold is sufficient.
	AggregateMaximum AggregationMethod = "maximum"

	// AggregateCustom allows custom aggregation logic defined by the implementation.
	AggregateCustom AggregationMethod = "custom"
)

// RecursiveLoop enables iterative refinement of a field's generation.
// The field will be generated multiple times, with each iteration potentially improving on the last.
// Useful for high-quality content that benefits from multiple attempts and selection of the best result.
type RecursiveLoop struct {
	// MaxIterations limits the number of generation attempts.
	// Prevents infinite loops and controls costs. Processing stops after this many iterations
	// regardless of whether termination conditions are met.
	MaxIterations int `json:"maxIterations"`

	// Selection determines which iteration to keep as the final result.
	// Options: highest (best score), lowest (worst score), latest (most recent), first (first successful)
	Selection SelectionStrategy `json:"selection"`

	// TerminateWhen defines conditions that stop iteration early if met.
	// Example: Stop when quality_score >= 85, even if under MaxIterations.
	// If empty, always runs until MaxIterations is reached.
	TerminateWhen []TerminationCondition `json:"terminateWhen,omitempty"`

	// FeedbackPrompt guides improvement between iterations.
	// This prompt is shown to the LLM along with the previous attempt to guide refinement.
	// Example: "Improve the following content by addressing these weaknesses: {previous_scores}"
	FeedbackPrompt string `json:"feedbackPrompt,omitempty"`

	// IncludePreviousAttempts determines whether to show the LLM all previous iterations
	// or just the most recent one. Useful for learning from multiple attempts.
	IncludePreviousAttempts bool `json:"includePreviousAttempts,omitempty"`
}

// SelectionStrategy determines which iteration to keep from a RecursiveLoop.
type SelectionStrategy string

const (
	// SelectHighestScore keeps the iteration with the best (highest) aggregate score.
	// Requires ScoringCriteria to be defined.
	SelectHighestScore SelectionStrategy = "highest"

	// SelectLowestScore keeps the iteration with the worst (lowest) aggregate score.
	// Useful for debugging or testing edge cases.
	SelectLowestScore SelectionStrategy = "lowest"

	// SelectLatest keeps the most recent iteration.
	// Assumes each iteration improves on the last.
	SelectLatest SelectionStrategy = "latest"

	// SelectFirst keeps the first successful iteration.
	// Useful when any passing result is acceptable.
	SelectFirst SelectionStrategy = "first"

	// SelectAll returns all iterations as an array instead of selecting one.
	// Allows downstream processing to compare or combine multiple attempts.
	SelectAll SelectionStrategy = "all"
)

// TerminationCondition defines when to stop a RecursiveLoop early.
// If any termination condition is met, iteration stops even if under MaxIterations.
type TerminationCondition struct {
	// Field is the name of the score to check (from ScoringCriteria dimensions).
	Field string `json:"field"`

	// Operator defines the comparison to perform.
	// Example: "gte" for "greater than or equal to"
	Operator ComparisonOperator `json:"operator"`

	// Value is the threshold that triggers termination.
	// Example: 85 to stop when score >= 85
	Value interface{} `json:"value"`

	// StopOn determines whether to stop when condition is true or false.
	// Default (true) stops when condition matches. Set to false to stop when condition fails.
	StopOn bool `json:"stopOn,omitempty"`
}
