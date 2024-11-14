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
		"AddStakeRequest":                   {Type: AddStakeEvent},
		"RemoveStakeRequest":                {Type: RemoveStakeEvent},
		"CancelRemoveStakeRequest":          {Type: CancelRemoveStakeEvent},
		"DelegateStakeRequest":              {Type: DelegateStakeEvent},
		"RemoveDelegateStakeRequest":        {Type: RemoveDelegateStakeEvent},
		"CancelRemoveDelegateStakeRequest":  {Type: CancelRemoveDelegateStakeEvent},
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
						{Type: "emissions.v5.AddStakeRequest"},
						{Type: "emissions.v5.RemoveStakeRequest"},
						{Type: "emissions.v5.CancelRemoveStakeRequest"},
						{Type: "emissions.v5.DelegateStakeRequest"},
						{Type: "emissions.v5.RemoveDelegateStakeRequest"},
						{Type: "emissions.v5.CancelRemoveDelegateStakeRequest"},
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
				{Type: "emissions.v5.AddStakeRequest"},
				{Type: "emissions.v5.RemoveStakeRequest"},
				{Type: "emissions.v5.CancelRemoveStakeRequest"},
				{Type: "emissions.v5.DelegateStakeRequest"},
				{Type: "emissions.v5.RemoveDelegateStakeRequest"},
				{Type: "emissions.v5.CancelRemoveDelegateStakeRequest"},
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
			name: "Staking events match",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "emissions.v5.AddStakeRequest"},
						{Type: "emissions.v5.RemoveStakeRequest"},
						{Type: "emissions.v5.CancelRemoveStakeRequest"},
						{Type: "emissions.v5.DelegateStakeRequest"},
						{Type: "emissions.v5.RemoveDelegateStakeRequest"},
						{Type: "emissions.v5.CancelRemoveDelegateStakeRequest"},
					},
				},
			},
			expected: []Event{
				{Type: "emissions.v5.AddStakeRequest"},
				{Type: "emissions.v5.RemoveStakeRequest"},
				{Type: "emissions.v5.CancelRemoveStakeRequest"},
				{Type: "emissions.v5.DelegateStakeRequest"},
				{Type: "emissions.v5.RemoveDelegateStakeRequest"},
				{Type: "emissions.v5.CancelRemoveDelegateStakeRequest"},
			},
		},
		{
			name: "Mixed staking and other events",
			events: &BlockResult{
				Result: struct {
					Height              string    `json:"height"`
					FinalizeBlockEvents []Event   `json:"finalize_block_events"`
					TxsBlockEvents      []TxEvent `json:"txs_results"`
				}{
					FinalizeBlockEvents: []Event{
						{Type: "emissions.v5.AddStakeRequest"},
						{Type: "emissions.v1.EventScoresSet"},
						{Type: "emissions.v5.DelegateStakeRequest"},
						{Type: "emissions.v1.EventRewardsSettled"},
					},
				},
			},
			expected: []Event{
				{Type: "emissions.v5.AddStakeRequest"},
				{Type: "emissions.v1.EventScoresSet"},
				{Type: "emissions.v5.DelegateStakeRequest"},
				{Type: "emissions.v1.EventRewardsSettled"},
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
