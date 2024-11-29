package main

import (
	"testing"
)

func TestFilterEvents(t *testing.T) {
	whitelist := map[string]EventProcessing{
		"EventScoresSet":                    {Type: ScoreEvent},
		"EventRewardsSettled":               {Type: RewardEvent},
		"EventNetworkLossSet":               {Type: NetworkLossEvent},
		"EventForecastTaskScoreSet":         {Type: ForecastTaskScoreEvent},
		"EventWorkerLastCommitSet":          {Type: ActorLastCommitEvent},
		"EventReputerLastCommitSet":         {Type: ActorLastCommitEvent},
		"EventTopicRewardsSet":              {Type: TopicRewardEvent},
		"EventEMAScoresSet":                 {Type: EMAScoreEvent},
		"EventTokenomicsSet":                {Type: TokenomicsEvent},
		"EventEcosystemTokenMintSet":        {Type: EcosystemTokenMintEvent},
		"EventRewardCurrentBlockEmission":   {Type: RewardCurrentBlockEmissionEvent},
		"EventListeningCoefficientsSet":     {Type: ListeningCoefficientsEvent},
		"EventInfererNetworkRegretSet":      {Type: InfererNetworkRegretEvent},
		"EventForecasterNetworkRegretSet":   {Type: ForecasterNetworkRegretEvent},
		"EventNaiveInfererNetworkRegretSet": {Type: NaiveInfererNetworkRegretEvent},
		"EventTopicInitialRegretSet":        {Type: TopicInitialRegretEvent},
		"rewards":                           {Type: ValidatorRewardsEvent},
		"commission":                        {Type: ValidatorCommissionEvent},
		"withdraw_rewards":                  {Type: ValidatorWithdrawRewardsEvent},
	}

	tests := []struct {
		name     string
		events   *BlockResult
		expected []Event
	}{
		{
			name: "All events match",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "emissions.v1.EventScoresSet"},
						{Type: "emissions.v1.EventRewardsSettled"},
						{Type: "emissions.v1.EventNetworkLossSet"},
						{Type: "emissions.v4.EventForecastTaskScoreSet"},
						{Type: "emissions.v4.EventWorkerLastCommitSet"},
						{Type: "emissions.v4.EventReputerLastCommitSet"},
						{Type: "emissions.v4.EventTopicRewardsSet"},
						{Type: "emissions.v4.EventEcosystemTokenMintSet"},
						{Type: "emissions.v4.EventRewardCurrentBlockEmission"},
						{Type: "emissions.v5.EventListeningCoefficientsSet"},
						{Type: "emissions.v5.EventInfererNetworkRegretSet"},
						{Type: "emissions.v5.EventForecasterNetworkRegretSet"},
						{Type: "emissions.v5.EventNaiveInfererNetworkRegretSet"},
						{Type: "emissions.v5.EventTopicInitialRegretSet"},
					},
					TxsBlockEvents: []TxEvent{
						{Events: []Event{
							{Type: "emissions.v2.EventScoresSet"},
							{Type: "emissions.v4.EventEMAScoresSet"},
						}},
					},
				},
			},
			expected: []Event{
				{Type: "emissions.v1.EventScoresSet"},
				{Type: "emissions.v1.EventRewardsSettled"},
				{Type: "emissions.v1.EventNetworkLossSet"},
				{Type: "emissions.v4.EventForecastTaskScoreSet"},
				{Type: "emissions.v4.EventWorkerLastCommitSet"},
				{Type: "emissions.v4.EventReputerLastCommitSet"},
				{Type: "emissions.v4.EventTopicRewardsSet"},
				{Type: "emissions.v4.EventEcosystemTokenMintSet"},
				{Type: "emissions.v4.EventRewardCurrentBlockEmission"},
				{Type: "emissions.v5.EventListeningCoefficientsSet"},
				{Type: "emissions.v5.EventInfererNetworkRegretSet"},
				{Type: "emissions.v5.EventForecasterNetworkRegretSet"},
				{Type: "emissions.v5.EventNaiveInfererNetworkRegretSet"},
				{Type: "emissions.v5.EventTopicInitialRegretSet"},
				{Type: "emissions.v2.EventScoresSet"},
				{Type: "emissions.v4.EventEMAScoresSet"},
			},
		},
		{
			name: "Some events match",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "emissions.v1.EventScoresSet"},
						{Type: "emissions.v1.UnknownEvent"},
					},
					TxsBlockEvents: []TxEvent{
						{Events: []Event{
							{Type: "emissions.v3.EventRewardsSettled"},
							{Type: "emissions.v1.AnotherUnknownEvent"},
						}},
					},
				},
			},
			expected: []Event{
				{Type: "emissions.v1.EventScoresSet"},
				{Type: "emissions.v3.EventRewardsSettled"},
			},
		},
		{
			name: "No events match",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "emissions.v1.UnknownEvent"},
					},
					TxsBlockEvents: []TxEvent{
						{Events: []Event{
							{Type: "emissions.v1.AnotherUnknownEvent"},
						}},
					},
				},
			},
			expected: []Event{},
		},
		{
			name: "Two-digit version event",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "emissions.v12.EventScoresSet"},
					},
					TxsBlockEvents: []TxEvent{
						{Events: []Event{
							{Type: "emissions.v12.EventRewardsSettled"},
						}},
					},
				},
			},
			expected: []Event{
				{Type: "emissions.v12.EventScoresSet"},
				{Type: "emissions.v12.EventRewardsSettled"},
			},
		},
		{
			name: "Event with no version",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "EventScoresSet"}, // No version
					},
					TxsBlockEvents: []TxEvent{
						{Events: []Event{
							{Type: "EventRewardsSettled"}, // No version
						}},
					},
				},
			},
			expected: []Event{}, // Should not match any
		},
		{
			name: "Validator commission events",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "rewards", Attributes: []Attribute{
							{Key: "amount", Value: "100uallo"},
							{Key: "validator", Value: "allovaloper1..."},
						}},
						{Type: "commission", Attributes: []Attribute{
							{Key: "amount", Value: "10uallo"},
							{Key: "validator", Value: "allovaloper1..."},
						}},
						{Type: "withdraw_rewards", Attributes: []Attribute{
							{Key: "amount", Value: "90uallo"},
							{Key: "validator", Value: "allovaloper1..."},
						}},
					},
				},
			},
			expected: []Event{
				{Type: "rewards"},
				{Type: "commission"},
				{Type: "withdraw_rewards"},
			},
		},
		{
			name: "Mixed validator and other events",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "rewards"},
						{Type: "emissions.v5.EventScoresSet"},
						{Type: "commission"},
						{Type: "emissions.v5.EventRewardsSettled"},
					},
				},
			},
			expected: []Event{
				{Type: "rewards"},
				{Type: "emissions.v5.EventScoresSet"},
				{Type: "commission"},
				{Type: "emissions.v5.EventRewardsSettled"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterEvents(tt.events, whitelist)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d events, got %d", len(tt.expected), len(result))
			}
			for i, event := range result {
				if event.Type != tt.expected[i].Type {
					t.Errorf("expected event type %s, got %s", tt.expected[i].Type, event.Type)
				}
			}
		})
	}
}

func TestGetBaseEventType(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		expected  string
	}{
		{
			name:      "Valid event type with version",
			eventType: "emissions.v1.EventScoresSet",
			expected:  "EventScoresSet",
		},
		{
			name:      "Valid event type with two-digit version",
			eventType: "emissions.v12.EventRewardsSettled",
			expected:  "EventRewardsSettled",
		},
		{
			name:      "Valid event type without version",
			eventType: "EventNetworkLossSet",
			expected:  string(InvalidType),
		},
		{
			name:      "Invalid event type",
			eventType: "InvalidEventType",
			expected:  string(InvalidType), // Expecting InvalidType
		},
		{
			name:      "Valid AddStake event type with version",
			eventType: "emissions.v5.AddStakeRequest",
			expected:  "AddStakeRequest",
		},
		{
			name:      "Valid DelegateStake event type with version",
			eventType: "emissions.v5.DelegateStakeRequest",
			expected:  "DelegateStakeRequest",
		},
		{
			name:      "Valid RemoveStake event type with two-digit version",
			eventType: "emissions.v5.RemoveStakeRequest",
			expected:  "RemoveStakeRequest",
		},
		{
			name:      "Valid CancelRemoveStake event type with version",
			eventType: "emissions.v5.CancelRemoveStakeRequest",
			expected:  "CancelRemoveStakeRequest",
		},
		{
			name:      "Valid RemoveDelegateStake event type with version",
			eventType: "emissions.v5.RemoveDelegateStakeRequest",
			expected:  "RemoveDelegateStakeRequest",
		},
		{
			name:      "Valid CancelRemoveDelegateStake event type with version",
			eventType: "emissions.v5.CancelRemoveDelegateStakeRequest",
			expected:  "CancelRemoveDelegateStakeRequest",
		},
		{
			name:      "Valid rewards event type",
			eventType: "rewards",
			expected:  "rewards",
		},
		{
			name:      "Valid commission event type",
			eventType: "commission",
			expected:  "commission",
		},
		{
			name:      "Valid withdraw_rewards event type",
			eventType: "withdraw_rewards",
			expected:  "withdraw_rewards",
		},
		{
			name:      "Mixed validator and emissions event",
			eventType: "emissions.v5.rewards",
			expected:  "rewards",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBaseEventType(tt.eventType)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFilterProductionEvents(t *testing.T) {
	// Setup the actual production whitelist
	whitelist := map[string]EventProcessing{
		"EventScoresSet":                    {Type: ScoreEvent},
		"EventRewardsSettled":               {Type: RewardEvent},
		"EventNetworkLossSet":               {Type: NetworkLossEvent},
		"EventForecastTaskScoreSet":         {Type: ForecastTaskScoreEvent},
		"EventWorkerLastCommitSet":          {Type: ActorLastCommitEvent},
		"EventReputerLastCommitSet":         {Type: ActorLastCommitEvent},
		"EventTopicRewardsSet":              {Type: TopicRewardEvent},
		"EventEMAScoresSet":                 {Type: EMAScoreEvent},
		"EventTokenomicsSet":                {Type: TokenomicsEvent},
		"EventEcosystemTokenMintSet":        {Type: EcosystemTokenMintEvent},
		"EventRewardCurrentBlockEmission":   {Type: RewardCurrentBlockEmissionEvent},
		"EventListeningCoefficientsSet":     {Type: ListeningCoefficientsEvent},
		"EventInfererNetworkRegretSet":      {Type: InfererNetworkRegretEvent},
		"EventForecasterNetworkRegretSet":   {Type: ForecasterNetworkRegretEvent},
		"EventNaiveInfererNetworkRegretSet": {Type: NaiveInfererNetworkRegretEvent},
		"EventTopicInitialRegretSet":        {Type: TopicInitialRegretEvent},
	}

	// Create test data with actual production events
	events := &BlockResult{
		Result: struct {
			Height              string    `json:"height"`
			FinalizeBlockEvents []Event   `json:"finalize_block_events"`
			TxsBlockEvents      []TxEvent `json:"txs_results"`
		}{
			FinalizeBlockEvents: []Event{
				{
					Type: "emissions.v5.EventInfererNetworkRegretSet",
					Attributes: []Attribute{
						{
							Key:   "addresses",
							Value: "[\"allo10aq5ue6xzz2a8dwxnfkhjxy2uypkctl8sqz2x7\"]",
						},
						{
							Key:   "block_height",
							Value: "\"1689777\"",
						},
					},
				},
				{
					Type: "emissions.v5.EventNaiveInfererNetworkRegretSet",
					Attributes: []Attribute{
						{
							Key:   "addresses",
							Value: "[\"allo10aq5ue6xzz2a8dwxnfkhjxy2uypkctl8sqz2x7\"]",
						},
						{
							Key:   "block_height",
							Value: "\"1689777\"",
						},
					},
				},
				{
					Type: "emissions.v5.EventTopicInitialRegretSet",
					Attributes: []Attribute{
						{
							Key:   "block_height",
							Value: "\"1689777\"",
						},
						{
							Key:   "topic_id",
							Value: "\"1\"",
						},
					},
				},
			},
		},
	}

	// Run the filter
	filtered := FilterEvents(events, whitelist)

	// Verify results
	if len(filtered) != 3 {
		t.Errorf("Expected 3 filtered events, got %d", len(filtered))
	}

	// Check each event type was properly filtered
	expectedTypes := []string{
		"emissions.v5.EventInfererNetworkRegretSet",
		"emissions.v5.EventNaiveInfererNetworkRegretSet",
		"emissions.v5.EventTopicInitialRegretSet",
	}

	for i, expectedType := range expectedTypes {
		if i >= len(filtered) {
			t.Errorf("Missing expected event: %s", expectedType)
			continue
		}
		if filtered[i].Type != expectedType {
			t.Errorf("Expected event type %s, got %s", expectedType, filtered[i].Type)
		}
	}
}
