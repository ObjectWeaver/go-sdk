package jsonSchema


// ComplexSystem orchestrates hierarchical AI decision-making workflows.
// This represents the top-level container for Object-Oriented AI patterns,
// combining primary prompts, decision trees, memory, and meta-intelligence.
type ComplexSystem struct {
	// Name identifies this system for logging and debugging.
	Name string `json:"name"`

	// Description explains the purpose and behavior of this system.
	Description string `json:"description,omitempty"`

	// PrimarySystemPrompt is inherited by all child LLM interactions in the decision tree.
	// Defines overarching goals, constraints, and behavioral guidelines.
	// Child definitions can override this with their own SystemPrompt.
	PrimarySystemPrompt string `json:"primarySystemPrompt"`

	// RootSchema defines the top-level object structure and decision tree entry point.
	// This is where the generation process begins.
	RootSchema Definition `json:"rootSchema"`

	// MainThread provides meta-intelligence that monitors and intervenes across decision points.
	// Enables hierarchical supervision similar to a "Gaia-like" architecture.
	MainThread *MainThreadConfig `json:"mainThread,omitempty"`

	// Memory configuration for context retrieval and knowledge management.
	// Allows the system to fetch relevant historical information during generation.
	Memory *MemoryConfig `json:"memory,omitempty"`

	// LazyConfig enables automatic decision tree selection and self-improvement.
	// When enabled, the system can auto-generate appropriate decision structures.
	LazyConfig *LazyConfigMode `json:"lazyConfig,omitempty"`

	// Version for tracking system schema versions.
	Version string `json:"version,omitempty"`
}

// MainThreadConfig defines supervisory meta-intelligence that oversees the generation process.
// The main thread monitors decision points, scores, and can intervene to improve outcomes.
type MainThreadConfig struct {
	// MonitoringStrategy determines how frequently the main thread observes the generation process.
	// Options: continuous (real-time), milestone (at decision points), final (only at completion)
	MonitoringStrategy MonitoringStrategy `json:"monitoringStrategy"`

	// InterventionRules define when and how the main thread should intervene.
	// Rules are evaluated in priority order; first matching rule triggers intervention.
	InterventionRules []InterventionRule `json:"interventionRules"`

	// HierarchicalDepth defines how many levels of nested decision trees to supervise.
	// 0 = no supervision, 1 = immediate children only, -1 = all levels
	HierarchicalDepth int `json:"hierarchicalDepth,omitempty"`

	// ParentThread for recursive supervision in multi-level hierarchies.
	// Enables "main threads monitoring main threads" for complex systems.
	ParentThread *MainThreadConfig `json:"parentThread,omitempty"`

	// AggregationPrompt for synthesizing outputs from multiple sub-modules.
	// Used when the main thread needs to combine results from different branches.
	AggregationPrompt string `json:"aggregationPrompt,omitempty"`

	// Model to use for main thread intelligence. Can differ from generation models.
	Model string `json:"model,omitempty"`
}

// MonitoringStrategy defines how the main thread observes the generation process.
type MonitoringStrategy string

const (
	// MonitorContinuous provides real-time monitoring at every generation step.
	// Highest oversight but most expensive.
	MonitorContinuous MonitoringStrategy = "continuous"

	// MonitorMilestone observes only at decision points and key checkpoints.
	// Balanced approach for most use cases.
	MonitorMilestone MonitoringStrategy = "milestone"

	// MonitorFinal observes only the final output.
	// Minimal oversight, useful for validation rather than guidance.
	MonitorFinal MonitoringStrategy = "final"
)

// InterventionRule defines conditions under which the main thread intervenes.
type InterventionRule struct {
	// Name identifies this rule for logging and debugging.
	Name string `json:"name,omitempty"`

	// Trigger defines the condition that activates this intervention.
	Trigger InterventionTrigger `json:"trigger"`

	// Action defines what the main thread should do when triggered.
	Action InterventionAction `json:"action"`

	// Priority determines evaluation order. Higher priority rules are checked first.
	Priority int `json:"priority,omitempty"`
}

// InterventionTrigger defines conditions that activate main thread intervention.
type InterventionTrigger struct {
	// Type specifies the kind of trigger condition.
	Type TriggerType `json:"type"`

	// ScoreThresholds map score names to threshold values.
	// Used with low_score or high_score trigger types.
	// Example: {"quality": 50.0} triggers when quality < 50
	ScoreThresholds map[string]float64 `json:"scoreThresholds,omitempty"`

	// PathDepth triggers at a specific hierarchy level in nested decision trees.
	// Useful for catching deep recursion or complexity issues.
	PathDepth *int `json:"pathDepth,omitempty"`

	// CustomCondition for complex logic expressed as a string.
	// Implementation-specific; could be evaluated by an LLM or custom code.
	CustomCondition string `json:"customCondition,omitempty"`
}

// TriggerType defines categories of intervention triggers.
type TriggerType string

const (
	TriggerLowScore      TriggerType = "low_score"      // Score falls below threshold
	TriggerHighScore     TriggerType = "high_score"     // Score exceeds threshold
	TriggerStalled       TriggerType = "stalled"        // No progress after N iterations
	TriggerDivergence    TriggerType = "divergence"     // Results inconsistent across attempts
	TriggerDepthExceeded TriggerType = "depth_exceeded" // Decision tree too deep
	TriggerCustom        TriggerType = "custom"         // Custom condition
)

// InterventionAction defines what the main thread does when intervening.
type InterventionAction struct {
	// Type specifies the kind of intervention action.
	Type ActionType `json:"type"`

	// OverridePrompt replaces or augments the current generation prompt.
	// Used to redirect generation with additional context or constraints.
	OverridePrompt *string `json:"overridePrompt,omitempty"`

	// FetchMemory triggers context retrieval before continuing.
	// Useful when low scores indicate missing information.
	FetchMemory bool `json:"fetchMemory,omitempty"`

	// ResetToCheckpoint rolls back to a previous state in the generation tree.
	// Checkpoint name must match a previously saved checkpoint.
	ResetToCheckpoint *string `json:"resetToCheckpoint,omitempty"`

	// CustomHandler name for implementation-specific intervention logic.
	CustomHandler string `json:"customHandler,omitempty"`

	// ModifyDefinition allows the main thread to alter the current Definition structure.
	// Enables dynamic decision tree modification based on observed behavior.
	ModifyDefinition *Definition `json:"modifyDefinition,omitempty"`
}

// ActionType defines categories of intervention actions.
type ActionType string

const (
	ActionOverride  ActionType = "override"  // Replace current generation prompt/definition
	ActionAugment   ActionType = "augment"   // Add context and retry current step
	ActionReset     ActionType = "reset"     // Roll back to checkpoint
	ActionTerminate ActionType = "terminate" // Stop processing immediately
	ActionEscalate  ActionType = "escalate"  // Send to parent thread for handling
	ActionCustom    ActionType = "custom"    // Custom action handler
)

// MemoryConfig defines context management and retrieval for the system.
type MemoryConfig struct {
	// Enabled determines whether memory system is active.
	Enabled bool `json:"enabled"`

	// SearchStrategy defines how to retrieve relevant context.
	SearchStrategy SearchStrategy `json:"searchStrategy"`

	// MaxContextSize limits the amount of retrieved context (in tokens or characters).
	// Prevents context overflow while ensuring relevant information is included.
	MaxContextSize int `json:"maxContextSize,omitempty"`

	// RetrievalTriggers define when to fetch memory.
	RetrievalTriggers []MemoryTrigger `json:"retrievalTriggers"`

	// KnowledgeGraph enables graph-based memory for relationship tracking.
	KnowledgeGraph *KnowledgeGraphConfig `json:"knowledgeGraph,omitempty"`

	// StorageBackend specifies where memory is stored (e.g., "vector_db", "sql", "file").
	StorageBackend string `json:"storageBackend,omitempty"`
}

// SearchStrategy defines how to retrieve relevant memory.
type SearchStrategy string

const (
	SearchSemantic SearchStrategy = "semantic" // Vector similarity search
	SearchKeyword  SearchStrategy = "keyword"  // Text matching/keyword search
	SearchGraph    SearchStrategy = "graph"    // Knowledge graph traversal
	SearchHybrid   SearchStrategy = "hybrid"   // Combination of multiple strategies
)

// MemoryTrigger defines when to retrieve context from memory.
type MemoryTrigger struct {
	// DecisionPoint identifies where to trigger retrieval (by name).
	// If empty, applies to all decision points when other conditions match.
	DecisionPoint string `json:"decisionPoint,omitempty"`

	// ScoreThreshold triggers retrieval when score condition is met.
	// Example: Fetch memory when quality < 70
	ScoreThreshold *Condition `json:"scoreThreshold,omitempty"`

	// AlwaysFetch retrieves memory at every generation step.
	// Expensive but ensures maximum context availability.
	AlwaysFetch bool `json:"alwaysFetch,omitempty"`

	// QueryTemplate for constructing the search query.
	// Can include placeholders like {field_name} that get replaced with actual values.
	// Example: "Find examples of {topic} with quality > 80"
	QueryTemplate string `json:"queryTemplate,omitempty"`
}

// KnowledgeGraphConfig enables graph-based memory representation.
type KnowledgeGraphConfig struct {
	// Enabled determines whether knowledge graph is active.
	Enabled bool `json:"enabled"`

	// NodeTypes define categories of entities in the graph.
	// Example: ["concept", "example", "pattern", "constraint"]
	NodeTypes []string `json:"nodeTypes,omitempty"`

	// RelationshipTypes define categories of connections between nodes.
	// Example: ["depends_on", "similar_to", "improves_upon", "contradicts"]
	RelationshipTypes []string `json:"relationshipTypes,omitempty"`

	// TraversalDepth limits how far to explore the graph during retrieval.
	// Higher depth finds more distantly related information but is slower.
	TraversalDepth int `json:"traversalDepth,omitempty"`
}

// LazyConfigMode enables automatic decision tree selection and self-improvement.
type LazyConfigMode struct {
	// Enabled determines whether lazy configuration is active.
	Enabled bool `json:"enabled"`

	// PatternLibrary references available pre-built decision patterns.
	PatternLibrary *PatternLibraryConfig `json:"patternLibrary"`

	// CreationPrompt guides auto-generation of decision structures.
	// Used by the system to generate appropriate DecisionPoints and branches.
	CreationPrompt string `json:"creationPrompt"`

	// SelfImprovement enables learning from successful task completions.
	SelfImprovement *SelfImprovementConfig `json:"selfImprovement,omitempty"`

	// DeploymentMode determines state management approach.
	// single_instance: Stateful learning in one container
	// distributed: Delegate to multiple ObjectWeaver instances
	DeploymentMode DeploymentMode `json:"deploymentMode,omitempty"`

	// MaxComplexity limits auto-generated decision tree depth.
	// Prevents runaway complexity in lazy configuration.
	MaxComplexity int `json:"maxComplexity,omitempty"`
}

// PatternLibraryConfig manages reusable decision tree patterns.
type PatternLibraryConfig struct {
	// SourceURL points to the pattern repository (local file, API endpoint, etc.).
	SourceURL string `json:"sourceUrl"`

	// SelectionModel specifies which LLM to use for pattern selection.
	// Can be different from generation models; optimize for reasoning.
	SelectionModel string `json:"selectionModel,omitempty"`

	// Tags for filtering patterns by use case.
	// Example: ["data_analysis", "content_generation", "code_review"]
	Tags []string `json:"tags,omitempty"`

	// CustomPatterns are user-defined patterns specific to this system.
	CustomPatterns []DecisionPattern `json:"customPatterns,omitempty"`
}

// DecisionPattern represents a reusable decision tree template.
type DecisionPattern struct {
	// ID uniquely identifies this pattern.
	ID string `json:"id"`

	// Name is a human-readable identifier.
	Name string `json:"name"`

	// Description explains when and how to use this pattern.
	Description string `json:"description"`

	// Tags for categorization and filtering.
	Tags []string `json:"tags"`

	// Schema is the decision tree structure (ComplexSystem or Definition).
	// Can be adapted when selected for a specific use case.
	Schema interface{} `json:"schema"` // ComplexSystem or Definition

	// UseCases documents successful applications of this pattern.
	UseCases []string `json:"useCases,omitempty"`

	// PerformanceMetrics tracks historical success rates.
	// Example: {"average_quality": 8.5, "success_rate": 0.92}
	PerformanceMetrics map[string]float64 `json:"performanceMetrics,omitempty"`

	// CreatedAt timestamp for versioning and tracking.
	CreatedAt string `json:"createdAt,omitempty"`

	// Version for pattern evolution tracking.
	Version string `json:"version,omitempty"`
}

// SelfImprovementConfig enables system learning and adaptation.
type SelfImprovementConfig struct {
	// Enabled determines whether the system learns from experience.
	Enabled bool `json:"enabled"`

	// SaveSuccessfulPatterns stores working configurations to the pattern library.
	SaveSuccessfulPatterns bool `json:"saveSuccessfulPatterns"`

	// SuccessThreshold defines what qualifies as a successful pattern.
	// Aggregate score must meet or exceed this value to be saved.
	SuccessThreshold float64 `json:"successThreshold"`

	// RefinementStrategy determines how patterns are improved over time.
	RefinementStrategy RefinementStrategy `json:"refinementStrategy,omitempty"`

	// SafetyMechanisms prevent system degradation from bad learning.
	SafetyMechanisms *SafetyConfig `json:"safetyMechanisms,omitempty"`

	// LearningRate controls how quickly patterns are updated (0.0-1.0).
	// Lower values = more conservative, higher = faster adaptation.
	LearningRate float64 `json:"learningRate,omitempty"`
}

// RefinementStrategy defines how patterns are improved.
type RefinementStrategy string

const (
	// RefineIterative makes gradual improvements based on feedback.
	RefineIterative RefinementStrategy = "iterative"

	// RefineComparative uses A/B testing to compare pattern variants.
	RefineComparative RefinementStrategy = "comparative"

	// RefineEvolutionary applies genetic algorithm-like evolution.
	RefineEvolutionary RefinementStrategy = "evolutionary"
)

// SafetyConfig prevents system degradation during self-improvement.
type SafetyConfig struct {
	// RequireApproval mandates human review before adding new patterns.
	RequireApproval bool `json:"requireApproval"`

	// RollbackOnFailure automatically reverts to previous pattern on poor performance.
	RollbackOnFailure bool `json:"rollbackOnFailure"`

	// PerformanceBaseline defines minimum acceptable performance metrics.
	// System won't adopt changes that fall below these thresholds.
	PerformanceBaseline map[string]float64 `json:"performanceBaseline,omitempty"`

	// MaxComplexity limits auto-generated decision tree depth.
	// Prevents runaway complexity that becomes unmaintainable.
	MaxComplexity int `json:"maxComplexity,omitempty"`

	// ValidationSet specifies test cases for validating new patterns.
	ValidationSet []string `json:"validationSet,omitempty"`
}

// DeploymentMode defines how the system manages state and scaling.
type DeploymentMode string

const (
	// DeploySingleInstance runs in a single container with persistent state.
	// Enables true learning and memory accumulation but limits scalability.
	DeploySingleInstance DeploymentMode = "single_instance"

	// DeployDistributed delegates object generation to multiple instances.
	// Better scalability but requires external state management.
	DeployDistributed DeploymentMode = "distributed"
)
