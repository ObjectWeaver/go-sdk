package jsonSchema

// This file contains practical examples demonstrating the Object-Oriented AI decision point system.

// ExampleSimpleScoreBasedRouting demonstrates basic conditional branching based on LLM evaluation scores.
func ExampleSimpleScoreBasedRouting() Definition {
	return Definition{
		Type: Object,
		Properties: map[string]Definition{
			"content": {
				Type:        String,
				Instruction: "Generate a technical blog post about microservices",

				// Define evaluation criteria
				ScoringCriteria: &ScoringCriteria{
					Dimensions: map[string]ScoringDimension{
						"technical_accuracy": {
							Description: "Rate the technical accuracy from 0-100, where 0 is completely inaccurate and 100 is perfectly accurate",
							Scale:       &ScoreScale{Min: 0, Max: 100},
							Type:        ScoreNumeric,
						},
						"readability": {
							Description: "Rate the readability from 0-100, where 0 is incomprehensible and 100 is crystal clear",
							Scale:       &ScoreScale{Min: 0, Max: 100},
							Type:        ScoreNumeric,
						},
					},
				},

				// Route based on scores
				DecisionPoint: &DecisionPoint{
					Name:     "ContentQualityRouter",
					Strategy: RouteByScore,
					Branches: []ConditionalBranch{
						{
							Name: "HighAccuracyLowReadability",
							// If technical_accuracy > 70 AND readability < 60
							Conditions: []Condition{
								{Field: "technical_accuracy", Operator: OpGreaterThan, Value: 70},
								{Field: "readability", Operator: OpLessThan, Value: 60},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Simplify the language and add more explanations while maintaining technical accuracy",
								SelectFields: []string{"content"}, // Pass existing content
							},
						},
						{
							Name: "LowAccuracyHighReadability",
							// If technical_accuracy < 60 AND readability > 70
							Conditions: []Condition{
								{Field: "technical_accuracy", Operator: OpLessThan, Value: 60},
								{Field: "readability", Operator: OpGreaterThan, Value: 70},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Add more technical depth, citations, and specific examples",
								SelectFields: []string{"content"},
							},
						},
						{
							Name: "HighQuality",
							// If both technical_accuracy > 80 AND readability > 80
							Conditions: []Condition{
								{Field: "technical_accuracy", Operator: OpGreaterThan, Value: 80},
								{Field: "readability", Operator: OpGreaterThan, Value: 80},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Content approved! Generate SEO meta description",
								SelectFields: []string{"content"},
							},
						},
					},
					Default: &Definition{
						Type:         String,
						Instruction:  "Perform general improvements to both accuracy and readability",
						SelectFields: []string{"content"},
					},
				},
			},
		},
	}
}

// ExampleFieldBasedRouting demonstrates routing based on generated field values rather than scores.
func ExampleFieldBasedRouting() Definition {
	return Definition{
		Type: Object,
		Properties: map[string]Definition{
			"requirements": {
				Type:        Object,
				Instruction: "Analyze the user's requirements",
				Properties: map[string]Definition{
					"is_technical": {
						Type:        Boolean,
						Instruction: "Is this technical content?",
					},
					"audience_level": {
						Type:        String,
						Instruction: "What is the audience level? (Beginner/Intermediate/Advanced)",
					},
					"needs_code_examples": {
						Type:        Boolean,
						Instruction: "Does this need code examples?",
					},
				},
			},
			"content": {
				Type:         String,
				Instruction:  "Generate appropriate content based on requirements",
				SelectFields: []string{"requirements"},

				// Route based on field values
				DecisionPoint: &DecisionPoint{
					Name:     "ContentTypeRouter",
					Strategy: RouteByField,
					Branches: []ConditionalBranch{
						{
							Name: "TechnicalWithCode",
							Conditions: []Condition{
								{Field: "is_technical", FieldPath: "requirements.is_technical", Operator: OpEqual, Value: true},
								{Field: "needs_code_examples", FieldPath: "requirements.needs_code_examples", Operator: OpEqual, Value: true},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Write technical content with detailed code examples and explanations",
								SystemPrompt: stringPtr("You are a technical writer with deep programming expertise. Include well-commented code examples."),
							},
						},
						{
							Name: "TechnicalNoCode",
							Conditions: []Condition{
								{Field: "is_technical", FieldPath: "requirements.is_technical", Operator: OpEqual, Value: true},
								{Field: "needs_code_examples", FieldPath: "requirements.needs_code_examples", Operator: OpEqual, Value: false},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Write technical content focused on concepts and architecture",
								SystemPrompt: stringPtr("You are a technical architect. Focus on concepts, patterns, and design decisions."),
							},
						},
						{
							Name: "BeginnerFriendly",
							Conditions: []Condition{
								{Field: "audience_level", FieldPath: "requirements.audience_level", Operator: OpEqual, Value: "Beginner"},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Write accessible content for beginners with simple explanations and analogies",
								SystemPrompt: stringPtr("You are an educator. Use simple language, analogies, and break down complex ideas."),
							},
						},
					},
				},
			},
		},
		ProcessingOrder: []string{"requirements", "content"},
	}
}

// ExampleHybridRouting demonstrates combining field-based and score-based conditions.
func ExampleHybridRouting() Definition {
	return Definition{
		Type: Object,
		Properties: map[string]Definition{
			"analysis": {
				Type:        Object,
				Instruction: "Analyze the requirements",
				Properties: map[string]Definition{
					"content_type": {
						Type:        String,
						Instruction: "What type of content? (blog/documentation/tutorial)",
					},
					"complexity": {
						Type:        String,
						Instruction: "What complexity level? (simple/moderate/complex)",
					},
				},
			},
			"draft": {
				Type:         String,
				Instruction:  "Generate initial draft",
				SelectFields: []string{"analysis"},

				ScoringCriteria: &ScoringCriteria{
					Dimensions: map[string]ScoringDimension{
						"quality": {
							Description: "Overall quality 0-100",
							Scale:       &ScoreScale{Min: 0, Max: 100},
							Type:        ScoreNumeric,
						},
					},
				},
			},
			"final_content": {
				Type:         String,
				Instruction:  "Finalize content based on type and quality",
				SelectFields: []string{"analysis", "draft"},

				// Hybrid routing: check both field values AND scores
				DecisionPoint: &DecisionPoint{
					Name:     "HybridRouter",
					Strategy: RouteByHybrid,
					Branches: []ConditionalBranch{
						{
							Name:     "ComplexContentLowQuality",
							Priority: 10, // High priority
							Conditions: []Condition{
								// Field condition
								{Field: "complexity", FieldPath: "analysis.complexity", Operator: OpEqual, Value: "complex"},
								// Score condition
								{Field: "quality", Operator: OpLessThan, Value: 70},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Add detailed explanations and break down complex concepts step-by-step",
								SelectFields: []string{"draft", "analysis"},
							},
						},
						{
							Name: "TutorialHighQuality",
							Conditions: []Condition{
								{Field: "content_type", FieldPath: "analysis.content_type", Operator: OpEqual, Value: "tutorial"},
								{Field: "quality", Operator: OpGreaterThanOrEqual, Value: 80},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Add interactive elements, exercises, and a summary section",
								SelectFields: []string{"draft"},
							},
						},
					},
				},
			},
		},
		ProcessingOrder: []string{"analysis", "draft", "final_content"},
	}
}

// ExampleRecursiveLoop demonstrates iterative refinement with score-based termination.
func ExampleRecursiveLoop() Definition {
	return Definition{
		Type:        String,
		Instruction: "Generate a high-quality product description",

		ScoringCriteria: &ScoringCriteria{
			Dimensions: map[string]ScoringDimension{
				"persuasiveness": {
					Description: "How persuasive is this description? (0-100)",
					Scale:       &ScoreScale{Min: 0, Max: 100},
					Type:        ScoreNumeric,
					Weight:      0.6,
				},
				"clarity": {
					Description: "How clear and easy to understand? (0-100)",
					Scale:       &ScoreScale{Min: 0, Max: 100},
					Type:        ScoreNumeric,
					Weight:      0.4,
				},
			},
			AggregationMethod: AggregateWeightedAverage,
		},

		// Iterate up to 5 times, keep the highest-scoring version
		RecursiveLoop: &RecursiveLoop{
			MaxIterations: 5,
			Selection:     SelectHighestScore,
			TerminateWhen: []TerminationCondition{
				{
					Field:    "persuasiveness",
					Operator: OpGreaterThanOrEqual,
					Value:    85, // Stop early if persuasiveness >= 85
					StopOn:   true,
				},
			},
			FeedbackPrompt:          "Improve the product description by making it more compelling and benefit-focused",
			IncludePreviousAttempts: true, // Show LLM all previous attempts to learn from them
		},
	}
}

// ExampleNestedDecisionTree demonstrates deep decision trees with multiple levels.
func ExampleNestedDecisionTree() Definition {
	return Definition{
		Type: Object,
		Properties: map[string]Definition{
			"content_analysis": {
				Type:        Object,
				Instruction: "Analyze content requirements",
				Properties: map[string]Definition{
					"topic_complexity": {
						Type:        String,
						Instruction: "Rate complexity: simple/moderate/complex",
					},
				},
			},
			"draft": {
				Type:         String,
				Instruction:  "Create initial draft",
				SelectFields: []string{"content_analysis"},

				ScoringCriteria: &ScoringCriteria{
					Dimensions: map[string]ScoringDimension{
						"depth": {Description: "Content depth 0-100", Scale: &ScoreScale{Min: 0, Max: 100}, Type: ScoreNumeric},
					},
				},

				// First level decision
				DecisionPoint: &DecisionPoint{
					Name:     "DepthEvaluation",
					Strategy: RouteByScore,
					Branches: []ConditionalBranch{
						{
							Name: "ShallowContent",
							Conditions: []Condition{
								{Field: "depth", Operator: OpLessThan, Value: 60},
							},
							Then: Definition{
								Type:         String,
								Instruction:  "Add more depth and details",
								SelectFields: []string{"draft"},

								ScoringCriteria: &ScoringCriteria{
									Dimensions: map[string]ScoringDimension{
										"improvement": {Description: "How much better? 0-100", Scale: &ScoreScale{Min: 0, Max: 100}, Type: ScoreNumeric},
									},
								},

								// Nested second-level decision
								DecisionPoint: &DecisionPoint{
									Name:     "ImprovementCheck",
									Strategy: RouteByScore,
									Branches: []ConditionalBranch{
										{
											Name: "StillNeedsWork",
											Conditions: []Condition{
												{Field: "improvement", Operator: OpLessThan, Value: 70},
											},
											Then: Definition{
												Type:         String,
												Instruction:  "Comprehensive rewrite with expert-level detail",
												SelectFields: []string{"draft"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		ProcessingOrder: []string{"content_analysis", "draft"},
	}
}

// ExampleComplexSystem demonstrates the full ComplexSystem with main thread intelligence.
func ExampleComplexSystem() ComplexSystem {
	return ComplexSystem{
		Name:                "IntelligentContentGenerator",
		Description:         "A self-monitoring content generation system with quality control",
		PrimarySystemPrompt: "You are an expert content creator focused on quality, accuracy, and engagement. Always prioritize clarity and value to the reader.",

		RootSchema: Definition{
			Type: Object,
			Properties: map[string]Definition{
				"topic_analysis": {
					Type:        String,
					Instruction: "Analyze the topic and determine key points to cover",
				},
				"content": {
					Type:         String,
					Instruction:  "Generate comprehensive content",
					SelectFields: []string{"topic_analysis"},

					ScoringCriteria: &ScoringCriteria{
						Dimensions: map[string]ScoringDimension{
							"quality":  {Description: "Overall quality 0-100", Scale: &ScoreScale{Min: 0, Max: 100}, Type: ScoreNumeric},
							"accuracy": {Description: "Factual accuracy 0-100", Scale: &ScoreScale{Min: 0, Max: 100}, Type: ScoreNumeric},
						},
					},

					DecisionPoint: &DecisionPoint{
						Name:     "QualityGate",
						Strategy: RouteByScore,
						Branches: []ConditionalBranch{
							{
								Name: "NeedsImprovement",
								Conditions: []Condition{
									{Field: "quality", Operator: OpLessThan, Value: 70},
								},
								Then: Definition{
									Type:         String,
									Instruction:  "Improve content quality",
									SelectFields: []string{"content"},
								},
							},
						},
					},
				},
			},
			ProcessingOrder: []string{"topic_analysis", "content"},
		},

		// Main thread monitors and intervenes
		MainThread: &MainThreadConfig{
			MonitoringStrategy: MonitorMilestone,
			InterventionRules: []InterventionRule{
				{
					Name: "LowQualityIntervention",
					Trigger: InterventionTrigger{
						Type:            TriggerLowScore,
						ScoreThresholds: map[string]float64{"quality": 50.0},
					},
					Action: InterventionAction{
						Type:           ActionAugment,
						FetchMemory:    true,
						OverridePrompt: stringPtr("The content quality is below standards. Please review successful examples and regenerate with higher quality."),
					},
					Priority: 10,
				},
				{
					Name: "StalledProgress",
					Trigger: InterventionTrigger{
						Type: TriggerStalled,
					},
					Action: InterventionAction{
						Type:              ActionReset,
						ResetToCheckpoint: stringPtr("topic_analysis"),
					},
					Priority: 5,
				},
			},
			HierarchicalDepth: 2, // Monitor 2 levels deep
		},

		// Memory for context retrieval
		Memory: &MemoryConfig{
			Enabled:        true,
			SearchStrategy: SearchSemantic,
			MaxContextSize: 2000,
			RetrievalTriggers: []MemoryTrigger{
				{
					ScoreThreshold: &Condition{
						Field:    "quality",
						Operator: OpLessThan,
						Value:    60,
					},
					QueryTemplate: "Find high-quality examples similar to {topic_analysis}",
				},
			},
		},

		Version: "1.0.0",
	}
}

// ExampleWithLazyConfig demonstrates automatic pattern selection and self-improvement.
func ExampleWithLazyConfig() ComplexSystem {
	return ComplexSystem{
		Name:                "SelfImprovingSystem",
		PrimarySystemPrompt: "You are an adaptive AI system that learns from experience.",

		RootSchema: Definition{
			Type: Object,
			Properties: map[string]Definition{
				"task": {
					Type:        String,
					Instruction: "Generate content for the given task",
				},
			},
		},

		LazyConfig: &LazyConfigMode{
			Enabled: true,
			PatternLibrary: &PatternLibraryConfig{
				SourceURL:      "https://patterns.objectweaver.ai/library",
				SelectionModel: "gpt-4",
				Tags:           []string{"content_generation", "quality_control"},
			},
			CreationPrompt: "Analyze the task and select or create an appropriate decision tree pattern. Optimize for quality and efficiency.",
			SelfImprovement: &SelfImprovementConfig{
				Enabled:                true,
				SaveSuccessfulPatterns: true,
				SuccessThreshold:       85.0,
				RefinementStrategy:     RefineIterative,
				SafetyMechanisms: &SafetyConfig{
					RequireApproval:   false,
					RollbackOnFailure: true,
					PerformanceBaseline: map[string]float64{
						"quality": 70.0,
						"speed":   80.0,
					},
					MaxComplexity: 5,
				},
				LearningRate: 0.1,
			},
			DeploymentMode: DeploySingleInstance,
			MaxComplexity:  4,
		},
	}
}

// Helper function for creating string pointers
func stringPtr(s string) *string {
	return &s
}
