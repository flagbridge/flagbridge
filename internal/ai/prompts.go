package ai

import "fmt"

const systemPromptBase = `You are FlagBridge AI, an assistant specialized in feature flag management and product intelligence.
You help teams manage feature flags, analyze rollout strategies, and make data-driven decisions about feature releases.

You have access to the project's current state including flags, targeting rules, product cards, team members, and recent activity.

Guidelines:
- Be concise and actionable. Product managers and engineers both use this tool.
- When suggesting flag changes, always explain the impact and risk.
- Reference specific flags by their key when relevant.
- When creating flags, suggest kebab-case keys (e.g., "dark-mode", "new-checkout-flow").
- For rollout strategies, consider the team size and flag complexity.
- Always mention if a suggested action requires admin privileges.

%s`

// SystemPrompt builds the full system prompt with project context injected.
func SystemPrompt(projectContext string) string {
	contextSection := fmt.Sprintf("Here is the current state of the project:\n\n%s", projectContext)
	return fmt.Sprintf(systemPromptBase, contextSection)
}

// AnalyzeFlagPrompt creates a user message for analyzing a specific flag.
func AnalyzeFlagPrompt(flagKey string) string {
	return fmt.Sprintf(`Analyze the feature flag "%s". Consider:
1. Current configuration and targeting rules
2. Product card context (hypothesis, success metrics, go/no-go criteria)
3. Rollout risk assessment
4. Suggestions for improvement

Provide a structured analysis with clear recommendations.`, flagKey)
}

// SuggestRolloutPrompt creates a prompt for rollout strategy suggestions.
func SuggestRolloutPrompt(flagKey string) string {
	return fmt.Sprintf(`Suggest a rollout strategy for the feature flag "%s". Include:
1. Recommended rollout percentage progression (e.g., 5%% → 25%% → 50%% ��� 100%%)
2. What targeting rules to set up
3. Key metrics to monitor at each stage
4. Go/no-go criteria for advancing to the next stage
5. Estimated timeline`, flagKey)
}

// SuggestProductCardPrompt creates a prompt for generating product card content.
func SuggestProductCardPrompt(flagKey, flagName, flagDescription string) string {
	return fmt.Sprintf(`Generate product card content for the feature flag "%s" (%s).
Description: %s

Generate:
1. A clear hypothesis statement (format: "If we [change], then [expected outcome] because [reasoning]")
2. Success metrics (2-3 measurable KPIs)
3. Go/no-go criteria (what signals indicate we should proceed or rollback)

Keep it concise — each section should be 1-3 sentences.`, flagKey, flagName, flagDescription)
}

// CreateFlagPrompt handles natural language flag creation requests.
func CreateFlagPrompt(userRequest string) string {
	return fmt.Sprintf(`The user wants to create a feature flag based on this request: "%s"

Based on the request, suggest:
1. flag_key (kebab-case, descriptive, max 50 chars)
2. name (human-readable)
3. description (one sentence)
4. type (boolean, string, number, or json)
5. default_value (appropriate for the type)
6. tags (1-3 relevant tags)

Respond with a JSON object containing these fields, followed by a brief explanation of your choices.
Format the JSON in a code block.`, userRequest)
}
