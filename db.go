package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/allora-network/allora-indexer/types"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type DBConsensusParams struct {
	MaxBytes         string
	MaxGas           string
	MaxAgeDuration   string
	MaxAgeNumBlocks  string
	EvidenceMaxBytes string
	PubKeyTypes      string // This can be a JSON-encoded array or a comma-separated list
}

type DBBlockInfo struct {
	BlockHash                  string
	BlockTotalParts            int
	BlockPartSetHeaderHash     string
	BlockVersion               string
	ChainID                    string
	Height                     uint64
	BlockTime                  time.Time
	LastBlockHash              string
	LastBlockTotalParts        int
	LastBlockPartSetHeaderHash string
	LastCommitHash             string
	DataHash                   string
	ValidatorsHash             string
	NextValidatorsHash         string
	ConsensusHash              string
	AppHash                    string
	LastResultsHash            string
	EvidenceHash               string
	ProposerAddress            string
}

const (
	TB_EVENTS                        = "events"
	TB_MESSAGES                      = "messages"
	TB_BLOCK_INFO                    = "block_info"
	TB_CONSENSUS_PARAMS              = "consensus_params"
	TB_TOPICS                        = "topics"
	TB_ADDRESSES                     = "addresses"
	TB_WORKER_REGISTRATIONS          = "worker_registrations"
	TB_TRANSFERS                     = "transfers"
	TB_INFERENCES                    = "inferences"
	TB_FORECASTS                     = "forecasts"
	TB_FORECAST_VALUES               = "forecast_values"
	TB_REPUTER_PAYLOAD               = "reputer_payload"
	TB_REPUTER_BUNDLES               = "reputer_bundles"
	TB_BUNDLE_VALUES                 = "bundle_values"
	TB_REWARDS                       = "rewards"
	TB_SCORES                        = "scores"
	TB_NETWORKLOSSES                 = "networklosses"
	TB_NETWORKLOSS_BUNDLE_VALUES     = "networkloss_bundle_values"
	TB_EMASCORES                     = "ema_scores"
	TB_ACTOR_LAST_COMMIT             = "last_commit_values"
	TB_TOKENOMICS                    = "tokenomics"
	TB_TOPIC_REWARD                  = "topic_rewards"
	TB_TOPIC_FORECASTING_SCORES      = "topic_forecasting_scores"
	TB_ECOSYSTEM_TOKEN_MINT          = "ecosystem_token_mint"
	TB_REWARD_CURRENT_BLOCK_EMISSION = "reward_current_block_emission"
	TB_LISTENING_COEFFICIENTS        = "listening_coefficients"
	TB_INFERER_NETWORK_REGRET        = "inferer_network_regret"
	TB_FORECASTER_NETWORK_REGRET     = "forecaster_network_regret"
	TB_NAIVE_INFERER_NETWORK_REGRET  = "naive_inferer_network_regret"
	TB_TOPIC_INITIAL_REGRET          = "topic_initial_regret"
	TB_REPUTER_STAKES                = "reputer_stakes"
	TB_VALIDATOR_REWARDS             = "validator_rewards"
	TB_VALIDATOR_COMMISSION          = "validator_commission"
	TB_VALIDATOR_WITHDRAW_COMMISSION = "validator_withdraw_commission"
	TB_TOPIC_INITIAL_EMA_SCORES      = "topic_initial_ema_scores"
)

var dbPool *pgxpool.Pool //*pgx.Conn

func verifyUri(originalPath string) string {
	var res = ""
	// Split the URL into components at the '@' symbol
	parts := strings.Split(originalPath, "@")
	if len(parts) != 2 {
		fmt.Println("Invalid URL format")
		return res
	}

	// Extract the userinfo (username:password) part by splitting the protocol part
	protocolAndUserInfo := strings.Split(parts[0], "//")
	if len(protocolAndUserInfo) != 2 {
		fmt.Println("Invalid userinfo format")
		return res
	}

	protocol := protocolAndUserInfo[0]
	userInfo := protocolAndUserInfo[1]

	// Now, split the userInfo into username and password using the first ':' as a delimiter
	credParts := strings.SplitN(userInfo, ":", 2)
	if len(credParts) != 2 {
		fmt.Println("Invalid credentials format")
		return res
	}

	username := credParts[0]
	password := credParts[1]

	// Encode the password
	encodedPassword := url.QueryEscape(password)

	// Reconstruct the new userInfo part
	newUserInfo := fmt.Sprintf("%s:%s", username, encodedPassword)

	// Reconstruct the final URL
	res = fmt.Sprintf("%s//%s@%s", protocol, newUserInfo, parts[1])
	return res
}

func initDB(dataSourceName string) error {
	var err error
	// dbPool, err = pgx.Connect(context.Background(), dataSourceName)

	dbConfig, err := pgxpool.ParseConfig(verifyUri(dataSourceName))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create a config, error: ")
		return err
	}
	dbPool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	return nil
}

func closeDB() {
	if dbPool != nil {
		dbPool.Close()
	}
}

func setupDB() error {
	// Execute all the SQL statements from the files in the sql-migrations folder
	sqlFiles := []string{"001_initial-schema-dump.sql", "002_initial_ema_scores.sql"}
	for _, file := range sqlFiles {
		sql, err := os.ReadFile(fmt.Sprintf("sql-migrations/%s", file))
		if err != nil {
			return err
		}
		err = executeSQL(string(sql))
		if err != nil {
			return err
		}
	}
	return nil
}

func executeSQL(sqlStatement string) error {
	if _, err := dbPool.Exec(context.Background(), sqlStatement); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute SQL statement: %v\n", err)
		return err
	}
	return nil
}

func insertBlockInfo(blockInfo DBBlockInfo) error {
	_, err := dbPool.Exec(context.Background(), `
		INSERT INTO `+TB_BLOCK_INFO+` (
			block_hash,
			block_total_parts,
			block_part_set_header_hash,
			block_version,
			chain_id,
			height,
			block_time,
			last_block_hash,
			last_block_total_parts,
			last_block_part_set_header_hash,
			last_commit_hash,
			data_hash,
			validators_hash,
			next_validators_hash,
			consensus_hash,
			app_hash,
			last_results_hash,
			evidence_hash,
			proposer_address
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
		blockInfo.BlockHash, blockInfo.BlockTotalParts, blockInfo.BlockPartSetHeaderHash,
		blockInfo.BlockVersion, blockInfo.ChainID, blockInfo.Height, blockInfo.BlockTime,
		blockInfo.LastBlockHash, blockInfo.LastBlockTotalParts, blockInfo.LastBlockPartSetHeaderHash,
		blockInfo.LastCommitHash, blockInfo.DataHash, blockInfo.ValidatorsHash,
		blockInfo.NextValidatorsHash, blockInfo.ConsensusHash, blockInfo.AppHash,
		blockInfo.LastResultsHash, blockInfo.EvidenceHash, blockInfo.ProposerAddress,
	)
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if isUniqueViolation(err) {
			log.Info().Msgf("Block height %d already exists in the database. Skipping insert.", blockInfo.Height)
			return nil // or return an error if you prefer
		}
		// Handle other types of errors
		return err
	}

	return nil
}

func insertMessage(height uint64, mtype string, sender string, data string, result string, txHash string) (uint64, error) {
	// Write Topic to the database
	var id uint64
	var dataHash = hash(data)
	log.Info().Msgf("Inserting message, height: %d,hash: %d", height, dataHash)
	err := dbPool.QueryRow(context.Background(), `
		INSERT INTO `+TB_MESSAGES+` (
			height,
			type,
			sender,
			data,
			hash,
			result,
			tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		height,
		mtype,
		sender,
		data,
		dataHash,
		result,
		txHash,
	).Scan(&id)
	if err != nil {
		log.Error().Msgf("Failed inserting message, height:%d, hash: %d", height, dataHash)
		return 0, err
	}

	return id, nil
}

func isUniqueViolation(err error) bool {
	// This function depends on your database driver
	// For example, with PostgreSQL using pq driver:
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // 23505 is the code for unique violation in PostgreSQL
	}
	return false
}

// Events
type EventRecord struct {
	Height uint64
	Type   string
	Sender string
	Data   json.RawMessage
}

func isEventType(eventType, prefix, suffix string) bool {
	return strings.HasPrefix(eventType, prefix) && strings.HasSuffix(eventType, suffix)
}

// isScoreEvent checks if the event is a score event based on its type.
func isScoreEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventScoresSet")
}

// isRewardEvent checks if the event is a reward event based on its type.
func isRewardEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventRewardsSettled")
}

// isNetworkLossEvent checks if the event is a network loss event based on its type.
func isNetworkLossEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventNetworkLossSet")
}

// isForecastTaskScoreEvent checks if the event is a forecast task score event based on its type.
func isForecastTaskScoreEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventForecastTaskScoreSet")
}

// isWorkerLastCommitEvent checks if the event is a worker/reputer last commit event based on its type.
func isWorkerLastCommitEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventWorkerLastCommitSet")
}

func isReputerLastCommitEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventReputerLastCommitSet")
}

func isTopicRewardEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventTopicRewardsSet")
}

func isEMAScoreEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventEMAScoresSet")
}

func isTokenomicsEvent(event EventRecord) bool {
	return isEventType(event.Type, "mint.v", "EventTokenomicsSet")
}

func isEcosystemTokenMintEvent(event EventRecord) bool {
	return isEventType(event.Type, "mint.v", "EventEcosystemTokenMintSet")
}

func isRewardCurrentBlockEmissionEvent(event EventRecord) bool {
	return isEventType(event.Type, "mint.v", "EventRewardCurrentBlockEmission")
}

func isListeningCoefficientsEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventListeningCoefficientsSet")
}

func isInfererNetworkRegretEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventInfererNetworkRegretSet")
}

func isForecasterNetworkRegretEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventForecasterNetworkRegretSet")
}

func isNaiveInfererNetworkRegretEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventNaiveInfererNetworkRegretSet")
}

func isTopicInitialRegretEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventTopicInitialRegretSet")
}

func isValidatorRewardsEvent(event EventRecord) bool {
	// cosmos SDK events don't have the same prefix as allora events
	return event.Type == "rewards"
}

func isValidatorCommissionEvent(event EventRecord) bool {
	// cosmos SDK events don't have the same prefix as allora events
	return event.Type == "commission"
}

func isValidatorWithdrawRewardsEvent(event EventRecord) bool {
	// cosmos SDK events don't have the same prefix as allora events
	return event.Type == "withdraw_rewards"
}

func isTopicInitialEmaScoreEvent(event EventRecord) bool {
	return isEventType(event.Type, "emissions.v", "EventTopicInitialEmaScoreSet")
}

func insertEvents(events []EventRecord) error {
	var scoreEvents []EventRecord
	var rewardEvents []EventRecord
	var networkLossEvents []EventRecord
	var forecastTaskScoreEvents []EventRecord
	var actorLastCommitEvents []EventRecord
	var topicRewardEvents []EventRecord
	var emaScoreEvents []EventRecord
	var tokenomicsEvents []EventRecord
	var ecosystemTokenMintEvents []EventRecord
	var rewardCurrentBlockEmissionEvents []EventRecord
	var listeningCoefficientsEvents []EventRecord
	var infererNetworkRegretEvents []EventRecord
	var forecasterNetworkRegretEvents []EventRecord
	var naiveInfererNetworkRegretEvents []EventRecord
	var topicInitialRegretEvents []EventRecord
	var validatorRewardsEvents []EventRecord
	var validatorCommissionEvents []EventRecord
	var validatorWithdrawRewardsEvents []EventRecord
	var topicInitialEmaScoreEvents []EventRecord

	// For inserting events in batch:
	var insertStatements []string
	var values []interface{}
	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL

	for _, event := range events {
		// Determine the type of event and accumulate accordingly
		if isScoreEvent(event) { // Function to check if it's a score event
			scoreEvents = append(scoreEvents, event)
		} else if isRewardEvent(event) { // Function to check if it's a reward event
			rewardEvents = append(rewardEvents, event)
		} else if isNetworkLossEvent(event) { // Function to check if it's a network loss event
			networkLossEvents = append(networkLossEvents, event)
		} else if isForecastTaskScoreEvent(event) {
			forecastTaskScoreEvents = append(forecastTaskScoreEvents, event) // Function to check if it's a forecast task score event
		} else if isWorkerLastCommitEvent(event) || isReputerLastCommitEvent(event) {
			actorLastCommitEvents = append(actorLastCommitEvents, event) // Function to check if it's an actor last commit event
		} else if isTopicRewardEvent(event) {
			topicRewardEvents = append(topicRewardEvents, event) // Function to check if it's a topic reward event
		} else if isEMAScoreEvent(event) {
			emaScoreEvents = append(emaScoreEvents, event) // Function to check if it's an ema score event
		} else if isTokenomicsEvent(event) {
			tokenomicsEvents = append(tokenomicsEvents, event) // Function to check if it's a tokenomics event
		} else if isEcosystemTokenMintEvent(event) {
			ecosystemTokenMintEvents = append(ecosystemTokenMintEvents, event) // Function to check if it's an ecosystem token mint event
		} else if isRewardCurrentBlockEmissionEvent(event) {
			rewardCurrentBlockEmissionEvents = append(rewardCurrentBlockEmissionEvents, event) // Function to check if it's a reward current block emission event
		} else if isListeningCoefficientsEvent(event) {
			listeningCoefficientsEvents = append(listeningCoefficientsEvents, event) // Function to check if it's a listening coefficients event
		} else if isInfererNetworkRegretEvent(event) {
			infererNetworkRegretEvents = append(infererNetworkRegretEvents, event) // Function to check if it's an inferer network regret event
		} else if isForecasterNetworkRegretEvent(event) {
			forecasterNetworkRegretEvents = append(forecasterNetworkRegretEvents, event) // Function to check if it's a forecaster network regret event
		} else if isNaiveInfererNetworkRegretEvent(event) {
			naiveInfererNetworkRegretEvents = append(naiveInfererNetworkRegretEvents, event) // Function to check if it's a naive inferer network regret event
		} else if isTopicInitialRegretEvent(event) {
			topicInitialRegretEvents = append(topicInitialRegretEvents, event) // Function to check if it's a topic initial regret event
		} else if isValidatorRewardsEvent(event) {
			validatorRewardsEvents = append(validatorRewardsEvents, event)
		} else if isValidatorCommissionEvent(event) {
			validatorCommissionEvents = append(validatorCommissionEvents, event)
		} else if isValidatorWithdrawRewardsEvent(event) {
			validatorWithdrawRewardsEvents = append(validatorWithdrawRewardsEvents, event)
		} else if isTopicInitialEmaScoreEvent(event) {
			topicInitialEmaScoreEvents = append(topicInitialEmaScoreEvents, event)
		} else {
			log.Info().Msg("Unrecognized event, ignoring")
			continue
		}

		// Prepare data for batch insert into TB_EVENTS
		dataHash := hash(string(event.Data))
		newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2, placeholderCounter+3, placeholderCounter+4)
		insertStatements = append(insertStatements, newStmt)
		values = append(values, event.Height, event.Type, event.Sender, string(event.Data), dataHash)
		placeholderCounter += 5 // Increase counter for next row
	}

	// Batch insert into TB_EVENTS
	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
				INSERT INTO %s (height, type, sender, data, hash) 
				VALUES %s
				ON CONFLICT (height, hash, type) DO NOTHING`, TB_EVENTS, strings.Join(insertStatements, ","))

		// Log the SQL statement and values for debugging
		log.Debug().Str("SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for events")

		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("event insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No events data to insert")
	}

	// Insert scores if any
	if len(scoreEvents) > 0 {
		err := insertScore(scoreEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert scores")
		}
	}

	// Insert rewards if any
	if len(rewardEvents) > 0 {
		err := insertReward(rewardEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert rewards")
		}
	}

	// Insert network losses if any
	if len(networkLossEvents) > 0 {
		err := insertNetworkLoss(networkLossEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert network losses")
		}
	}

	// Insert forecast task score if any
	if len(forecastTaskScoreEvents) > 0 {
		err := insertForecastTaskScore(forecastTaskScoreEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert forecast task score")
		}
	}

	// Insert actor last commit if any
	if len(actorLastCommitEvents) > 0 {
		err := insertActorLastCommit(actorLastCommitEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert last commit")
		}
	}

	// Insert topic reward if any
	if len(topicRewardEvents) > 0 {
		err := insertTopicReward(topicRewardEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert topic reward")
		}
	}

	// Insert ema score if any
	if len(emaScoreEvents) > 0 {
		err := insertEMAScore(emaScoreEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert ema score")
		}
	}

	// Insert tokenomics if any
	if len(tokenomicsEvents) > 0 {
		err := insertTokenomics(tokenomicsEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert tokenomics")
		}
	}

	// Insert ecosystem token mint if any
	if len(ecosystemTokenMintEvents) > 0 {
		err := insertEcosystemTokenMint(ecosystemTokenMintEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert ecosystem token mint")
		}
	}

	// Insert reward current block emission if any
	if len(rewardCurrentBlockEmissionEvents) > 0 {
		err := insertRewardCurrentBlockEmission(rewardCurrentBlockEmissionEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert reward current block emission")
		}
	}

	// Insert listening coefficients if any
	if len(listeningCoefficientsEvents) > 0 {
		err := insertListeningCoefficients(listeningCoefficientsEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert listening coefficients")
		}
	}

	// Insert inferer network regret if any
	if len(infererNetworkRegretEvents) > 0 {
		err := insertInfererNetworkRegret(infererNetworkRegretEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert inferer network regret")
		}
	}

	// Insert forecaster network regret if any
	if len(forecasterNetworkRegretEvents) > 0 {
		err := insertForecasterNetworkRegret(forecasterNetworkRegretEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert forecaster network regret")
		}
	}

	// Insert naive inferer network regret if any
	if len(naiveInfererNetworkRegretEvents) > 0 {
		err := insertNaiveInfererNetworkRegret(naiveInfererNetworkRegretEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert naive inferer network regret")
		}
	}

	// Insert topic initial regret if any
	if len(topicInitialRegretEvents) > 0 {
		err := insertTopicInitialRegret(topicInitialRegretEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert topic initial regret")
		}
	}

	if len(validatorRewardsEvents) > 0 {
		err := insertValidatorRewards(validatorRewardsEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert validator rewards events")
		}
	}

	if len(validatorCommissionEvents) > 0 {
		err := insertValidatorCommission(validatorCommissionEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert validator commission events")
		}
	}

	if len(validatorWithdrawRewardsEvents) > 0 {
		err := insertValidatorWithdrawRewards(validatorWithdrawRewardsEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert validator withdraw rewards events")
		}
	}

	if len(topicInitialEmaScoreEvents) > 0 {
		err := insertTopicInitialEmaScore(topicInitialEmaScoreEvents)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert topic initial ema scores")
		}
	}

	return nil
}

func insertScore(events []EventRecord) error {
	log.Info().Msg("Inserting scores in batch")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL

	for _, event := range events {
		log.Trace().Interface("Event score", event).Msg("Processing event score")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var topicID int
		var actorType string
		var addresses []string
		var scores []big.Float
		var blockHeight int

		for _, attr := range attributes {
			switch attr.Key {
			case "topic_id":
				cleanedValue := strings.Trim(attr.Value, "\"")
				topicID, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to convert topic_id to int: %w", err)
				}
			case "actor_type":
				actorType = strings.Trim(attr.Value, "\"")
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			case "addresses":
				err = json.Unmarshal([]byte(attr.Value), &addresses)
				if err != nil {
					return fmt.Errorf("failed to unmarshal addresses: %w", err)
				}
			case "scores":
				var rawScores []string
				err = json.Unmarshal([]byte(attr.Value), &rawScores)
				if err != nil {
					return fmt.Errorf("failed to unmarshal scores: %w", err)
				}

				for _, rawScore := range rawScores {
					rawScoreClean := strings.Trim(rawScore, "\"")
					if isInvalidNumericValue(rawScoreClean) {
						log.Error().Str("rawScore", rawScore).Msg("Failed to convert score to big.Float")
						return fmt.Errorf("Invalid Score: %s", rawScoreClean)
					} else {
						score := new(big.Float)
						score, ok := score.SetString(rawScoreClean)
						if !ok {
							log.Error().Str("rawScore", rawScore).Msg("Failed to convert score to big.Float")
							return fmt.Errorf("Invalid Score: %s", rawScoreClean)
						}
						scores = append(scores, *score)
					}
				}
			}
		}

		if len(addresses) != len(scores) {
			return fmt.Errorf("mismatch in length of addresses and scores")
		}

		for i := range addresses {
			// Generate the placeholders for this row
			newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2, placeholderCounter+3, placeholderCounter+4, placeholderCounter+5)
			insertStatements = append(insertStatements, newStmt)
			scoreValue := scores[i].Text('f', -1)
			values = append(values, event.Height, blockHeight, topicID, actorType, addresses[i], scoreValue)
			placeholderCounter += 6 // Increase counter for next row
		}
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, height, topic_id, type, address, value) 
			VALUES %s`, TB_SCORES, strings.Join(insertStatements, ","))
		log.Trace().Str("Event - Score SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for scores")
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("score insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No scores data to insert")
	}

	return nil
}

func insertReward(events []EventRecord) error {
	log.Info().Msg("Inserting rewards in batch")
	var insertStatements []string
	var values []interface{}
	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL

	for _, event := range events {
		log.Trace().Interface("Event reward", event).Msg("Processing event reward")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var topicID int
		var rewardType string
		var addresses []string
		var rewards []big.Float
		var blockHeight int

		for _, attr := range attributes {
			switch attr.Key {
			case "topic_id":
				cleanedValue := strings.Trim(attr.Value, "\"")
				topicID, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to convert topic_id to int: %w", err)
				}
			case "actor_type":
				rewardType = strings.Trim(attr.Value, "\"")
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			case "addresses":
				err = json.Unmarshal([]byte(attr.Value), &addresses)
				if err != nil {
					return fmt.Errorf("failed to unmarshal addresses: %w", err)
				}
			case "rewards":
				err = json.Unmarshal([]byte(attr.Value), &rewards)
				if err != nil {
					return fmt.Errorf("failed to unmarshal rewards: %w", err)
				}
			}
		}

		if len(addresses) != len(rewards) {
			return fmt.Errorf("mismatch in length of addresses and rewards")
		}

		for i := range addresses {
			// Generate the placeholders for this row
			newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2, placeholderCounter+3, placeholderCounter+4, placeholderCounter+5)
			insertStatements = append(insertStatements, newStmt)
			values = append(values, event.Height, blockHeight, topicID, rewardType, addresses[i], rewards[i].Text('f', -1))
			placeholderCounter += 6 // Increase counter for next row
		}
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, height, topic_id, type, address, value) 
			VALUES %s`, TB_REWARDS, strings.Join(insertStatements, ","))

		log.Trace().Str("Event - Reward SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for scores")

		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("rewards insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No rewards data to insert")
	}

	return nil
}

func insertNetworkLoss(events []EventRecord) error {
	for _, event := range events {
		log.Debug().Interface("Event network loss", event).Msg("inserting event network loss ")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return err
		}

		var topicID int
		var block_height int
		var valueBundle types.MsgValueBundle

		for _, attr := range attributes {
			cleanedValue := strings.Trim(attr.Value, "\"")
			switch attr.Key {
			case "topic_id":
				topicID, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return err
				}
			case "block_height":
				block_height, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return err
				}
			case "value_bundle":
				err = json.Unmarshal([]byte(cleanedValue), &valueBundle)
				if err != nil {
					return err
				}
			}
		}

		var bundleId uint64
		err = dbPool.QueryRow(context.Background(), `
				INSERT INTO `+TB_NETWORKLOSSES+` (height_tx, height, topic_id, naive_value, combined_value) VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT (height_tx, height, topic_id) DO NOTHING returning id`,
			event.Height, block_height, topicID, valueBundle.NaiveValue, valueBundle.CombinedValue).Scan(&bundleId)

		if err != nil {
			return fmt.Errorf("network loss event insert failed: %v", err)
		}

		log.Debug().Msgf("Inserting NetworkLoss bundle: %d, %v", bundleId, valueBundle)
		err = insertValueBundle(bundleId, valueBundle, TB_NETWORKLOSS_BUNDLE_VALUES)
		if err != nil {
			return fmt.Errorf("network loss bundle insert failed: %v", err)
		}
	}
	return nil
}

func insertForecastTaskScore(events []EventRecord) error {
	log.Info().Msg("Updating topic forecasting task score")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL

	for _, event := range events {
		log.Trace().Interface("Event topic forecast task score", event).Msg("Processing event topic forecast task score")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var topicID int
		var score string
		for _, attr := range attributes {
			switch attr.Key {
			case "topic_id":
				cleanedValue := strings.Trim(attr.Value, "\"")
				topicID, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to get topic id: %w", err)
				}
			case "score":
				score = strings.Trim(attr.Value, "\"")
			}
		}
		newStmt := fmt.Sprintf("($%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2)
		insertStatements = append(insertStatements, newStmt)
		values = append(values, event.Height, topicID, score)
		placeholderCounter += 3 // Increase counter for next row
	}
	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx,topic_id, score) 
			VALUES %s`, TB_TOPIC_FORECASTING_SCORES, strings.Join(insertStatements, ","))
		log.Trace().Str("Event - Topic forecasting task score SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for topic forecast task score")
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("topic forecasting task score failed: %v", err)
		}
	} else {
		log.Info().Msg("No topic forecasting task score to insert")
	}
	return nil
}

func insertActorLastCommit(events []EventRecord) error {
	log.Info().Msg("Inserting actor last commit in batch")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL
	for _, event := range events {
		log.Trace().Interface("Event last commit", event).Msg("Processing event last commit")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var topicID int
		var height int
		var nonce int
		var isWorker = true
		if isReputerLastCommitEvent(event) {
			isWorker = false
		}
		for _, attr := range attributes {
			switch attr.Key {
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				height, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to get block height: %w", err)
				}
			case "nonce":
				var cleanedValue map[string]string
				err = json.Unmarshal([]byte(attr.Value), &cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to unmarshal nonce: %w", err)
				}
				nonce, err = strconv.Atoi(cleanedValue["block_height"])
				if err != nil {
					return fmt.Errorf("failed to convert getting nonce: %w", err)
				}
			case "topic_id":
				cleanedValue := strings.Trim(attr.Value, "\"")
				topicID, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to get topic id: %w", err)
				}
			}
		}

		newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2, placeholderCounter+3)
		insertStatements = append(insertStatements, newStmt)
		values = append(values, height, nonce, topicID, isWorker)
		placeholderCounter += 4 // Increase counter for next row
	}
	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, height, topic_id, is_worker) 
			VALUES %s ON CONFLICT (topic_id, is_worker) 
			DO UPDATE SET height_tx=EXCLUDED.height_tx, height=EXCLUDED.height`, TB_ACTOR_LAST_COMMIT, strings.Join(insertStatements, ","))
		log.Info().Str("Event - Last commit SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for last commit")
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("failed to insert last commit event: %w", err)
		}
	} else {
		log.Info().Msg("No last commit event to insert")
	}
	return nil
}

func insertTopicReward(events []EventRecord) error {
	log.Info().Msg("Updating topic reward")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL
	for _, event := range events {
		log.Trace().Interface("Event topic reward", event).Msg("Processing event topic reward")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var topicIDs []string
		var rewards []big.Float
		for _, attr := range attributes {
			switch attr.Key {
			case "topic_ids":
				err = json.Unmarshal([]byte(attr.Value), &topicIDs)
				if err != nil {
					return fmt.Errorf("failed to unmarshal topics: %w", err)
				}
			case "rewards":
				err = json.Unmarshal([]byte(attr.Value), &rewards)
				if err != nil {
					return fmt.Errorf("failed to unmarshal rewards: %w", err)
				}
			}
		}
		if len(topicIDs) != len(rewards) {
			return fmt.Errorf("mismatch in length of topic ids and rewards")
		}

		for i := range topicIDs {
			// Generate the placeholders for this row
			newStmt := fmt.Sprintf("($%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2)
			insertStatements = append(insertStatements, newStmt)
			rewardValue := rewards[i].Text('f', -1)
			values = append(values, event.Height, topicIDs[i], rewardValue)
			placeholderCounter += 3 // Increase counter for next row
		}
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx,topic_id, reward) 
			VALUES %s`, TB_TOPIC_REWARD, strings.Join(insertStatements, ","))
		log.Trace().Str("Event - Topic reward SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for topic reward")
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("topic reward insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No topic reward data to insert")
	}
	return nil
}

func insertEMAScore(events []EventRecord) error {
	log.Info().Msg("Inserting ema scores in batch")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL

	for _, event := range events {
		log.Trace().Interface("Event score", event).Msg("Processing event score")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var topicID int
		var actorType string
		var addresses []string
		var scores []big.Float
		var activations []bool
		var blockHeight int

		for _, attr := range attributes {
			switch attr.Key {
			case "topic_id":
				cleanedValue := strings.Trim(attr.Value, "\"")
				topicID, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to convert topic_id to int: %w", err)
				}
			case "actor_type":
				actorType = strings.Trim(attr.Value, "\"")
			case "nonce":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			case "addresses":
				err = json.Unmarshal([]byte(attr.Value), &addresses)
				if err != nil {
					return fmt.Errorf("failed to unmarshal addresses: %w", err)
				}
			case "scores":
				var rawScores []string
				err = json.Unmarshal([]byte(attr.Value), &rawScores)
				if err != nil {
					return fmt.Errorf("failed to unmarshal scores: %w", err)
				}

				for _, rawScore := range rawScores {
					rawScoreClean := strings.Trim(rawScore, "\"")
					if isInvalidNumericValue(rawScoreClean) {
						log.Error().Str("rawScore", rawScore).Msg("Failed to convert score to big.Float")
						return fmt.Errorf("Invalid Score: %s", rawScoreClean)
					} else {
						score := new(big.Float)
						score, ok := score.SetString(rawScoreClean)
						if !ok {
							log.Error().Str("rawScore", rawScore).Msg("Failed to convert score to big.Float")
							return fmt.Errorf("Invalid Score: %s", rawScoreClean)
						}
						scores = append(scores, *score)
					}
				}
			case "is_active":
				err = json.Unmarshal([]byte(attr.Value), &activations)
				if err != nil {
					return fmt.Errorf("failed to unmarshal activation: %w", err)
				}
			}
		}

		if len(addresses) != len(scores) || len(scores) != len(activations) {
			return fmt.Errorf("mismatch in length of addresses, scores, activations")
		}

		for i := range addresses {
			// Generate the placeholders for this row
			newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1,
				placeholderCounter+2, placeholderCounter+3, placeholderCounter+4, placeholderCounter+5, placeholderCounter+6)
			insertStatements = append(insertStatements, newStmt)
			scoreValue := scores[i].Text('f', -1)
			values = append(values, event.Height, blockHeight, topicID, actorType, addresses[i], scoreValue, activations[i])
			placeholderCounter += 7 // Increase counter for next row
		}
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, height, topic_id, type, address, score, is_active) 
			VALUES %s ON CONFLICT (topic_id, "type", address, height)
			DO UPDATE SET score=EXCLUDED.score, is_active=EXCLUDED.is_active`, TB_EMASCORES,
			strings.Join(insertStatements, ","))
		log.Trace().Str("Event - Score SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for scores")
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("score insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No scores data to insert")
	}

	return nil
}

func insertTokenomics(events []EventRecord) error {
	log.Info().Msg("Inserting tokenomics")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL

	for _, event := range events {
		log.Trace().Interface("Event tokenomics", event).Msg("Processing tokenomic event")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var stakedTokenAmount = new(big.Float)
		var circulatingAmount = new(big.Float)
		var emissionsAmount = new(big.Float)
		for _, attr := range attributes {
			switch attr.Key {
			case "circulating_supply":
				cleanedValue := strings.Trim(attr.Value, "\"")
				_, ok := circulatingAmount.SetString(cleanedValue)
				if !ok {
					return fmt.Errorf("failed to get circulating supply: %w", err)
				}
			case "emissions_amount":
				cleanedValue := strings.Trim(attr.Value, "\"")
				_, ok := emissionsAmount.SetString(cleanedValue)
				if !ok {
					return fmt.Errorf("failed to get emissions total amount supply: %w", err)
				}
			case "staked_token_amount":
				cleanedValue := strings.Trim(attr.Value, "\"")
				_, ok := stakedTokenAmount.SetString(cleanedValue)
				if !ok {
					return fmt.Errorf("failed to get staked token amount supply: %w", err)
				}
			}
		}

		newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2, placeholderCounter+3)
		insertStatements = append(insertStatements, newStmt)
		values = append(values, event.Height, stakedTokenAmount, circulatingAmount, emissionsAmount)
		placeholderCounter += 4 // Increase counter for next row
	}
	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, staked_amount, circulating_supply, emissions_amount) 
			VALUES %s`, TB_TOKENOMICS, strings.Join(insertStatements, ","))
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("failed to insert tokenomics event")
		}
	} else {
		log.Info().Msg("No tokenomics event to insert")
	}
	return nil
}

func insertEcosystemTokenMint(events []EventRecord) error {
	log.Info().Msg("Inserting ecosystem token mint")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL
	for _, event := range events {
		log.Trace().Interface("Event ecosystem token mint", event).Msg("Processing event ecosystem token mint")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var tokenAmount = new(big.Float)
		var blockHeight uint64
		for _, attr := range attributes {
			switch attr.Key {
			case "token_amount":
				cleanedValue := strings.Trim(attr.Value, "\"")
				_, ok := tokenAmount.SetString(cleanedValue)
				if !ok {
					return fmt.Errorf("failed to get token amount: %w", err)
				}
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.ParseUint(cleanedValue, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			}
		}
		newStmt := fmt.Sprintf("($%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2)
		insertStatements = append(insertStatements, newStmt)
		values = append(values, event.Height, blockHeight, tokenAmount)
		placeholderCounter += 3 // Increase counter for next row
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, block_height, token_amount) 
			VALUES %s`, TB_ECOSYSTEM_TOKEN_MINT, strings.Join(insertStatements, ","))
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("failed to insert ecosystem token mint event")
		}
	} else {
		log.Info().Msg("No ecosystem token mint event to insert")
	}
	return nil
}

func insertRewardCurrentBlockEmission(events []EventRecord) error {
	log.Info().Msg("Inserting reward current block emission")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL
	for _, event := range events {
		log.Trace().Interface("Event reward current block emission", event).Msg("Processing event reward current block emission")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var tokenAmount = new(big.Float)
		var blockHeight uint64
		for _, attr := range attributes {
			switch attr.Key {
			case "token_amount":
				cleanedValue := strings.Trim(attr.Value, "\"")
				_, ok := tokenAmount.SetString(cleanedValue)
				if !ok {
					return fmt.Errorf("failed to get token amount: %w", err)
				}
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.ParseUint(cleanedValue, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			}
		}
		newStmt := fmt.Sprintf("($%d, $%d, $%d)", placeholderCounter, placeholderCounter+1, placeholderCounter+2)
		insertStatements = append(insertStatements, newStmt)
		values = append(values, event.Height, blockHeight, tokenAmount)
		placeholderCounter += 3 // Increase counter for next row
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, block_height, token_amount) 
			VALUES %s`, TB_REWARD_CURRENT_BLOCK_EMISSION, strings.Join(insertStatements, ","))
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("failed to insert reward current block emission event")
		}
	} else {
		log.Info().Msg("No reward current block emission event to insert")
	}
	return nil
}

func insertListeningCoefficients(events []EventRecord) error {
	log.Info().Msg("Inserting listening coefficients")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1
	for _, event := range events {
		log.Trace().Interface("Event listening coefficients", event).Msg("Processing event listening coefficients")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var actorType string
		var topicID int64
		var blockHeight uint64
		var addresses []string
		var coefficients []string

		for _, attr := range attributes {
			switch attr.Key {
			case "actor_type":
				actorType = strings.Trim(attr.Value, "\"")
			case "topic_id":
				topicID, err = strconv.ParseInt(strings.Trim(attr.Value, "\""), 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert topic_id to int: %w", err)
				}
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.ParseUint(cleanedValue, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			case "addresses":
				err = json.Unmarshal([]byte(attr.Value), &addresses)
				if err != nil {
					return fmt.Errorf("failed to unmarshal addresses: %w", err)
				}
			case "coefficients":
				err = json.Unmarshal([]byte(attr.Value), &coefficients)
				if err != nil {
					return fmt.Errorf("failed to unmarshal coefficients: %w", err)
				}
			}
		}

		newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			placeholderCounter,
			placeholderCounter+1,
			placeholderCounter+2,
			placeholderCounter+3,
			placeholderCounter+4)
		insertStatements = append(insertStatements, newStmt)
		values = append(values,
			actorType,    // actor_type
			topicID,      // topic_id
			blockHeight,  // block_height
			addresses,    // addresses
			coefficients, // coefficients
		)
		placeholderCounter += 5 // Adjusted to match the number of values
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (actor_type, topic_id, block_height, addresses, coefficients) 
			VALUES %s`, TB_LISTENING_COEFFICIENTS, strings.Join(insertStatements, ","))

		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("failed to insert listening coefficients event: %v", err)
		}
	} else {
		log.Info().Msg("No listening coefficients event to insert")
	}

	return nil
}

func insertNetworkRegret(events []EventRecord, tableName string) error {
	log.Info().Msg("Inserting network regret")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1
	for _, event := range events {
		log.Trace().Interface("Event network regret", event).Msg("Processing event network regret")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var blockHeight uint64
		var topicID int64
		var addresses []string
		var regrets []string

		for _, attr := range attributes {
			switch attr.Key {
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.ParseUint(cleanedValue, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			case "topic_id":
				cleanedValue := strings.Trim(attr.Value, "\"")
				topicID, err = strconv.ParseInt(cleanedValue, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert topic_id to int: %w", err)
				}
			case "addresses":
				err = json.Unmarshal([]byte(attr.Value), &addresses)
				if err != nil {
					return fmt.Errorf("failed to unmarshal addresses: %w", err)
				}
			case "regrets":
				err = json.Unmarshal([]byte(attr.Value), &regrets)
				if err != nil {
					return fmt.Errorf("failed to unmarshal regrets: %w", err)
				}
			}
		}

		newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			placeholderCounter,
			placeholderCounter+1,
			placeholderCounter+2,
			placeholderCounter+3,
			placeholderCounter+4)
		insertStatements = append(insertStatements, newStmt)
		values = append(values,
			event.Height, // height_tx
			blockHeight,  // block_height
			topicID,      // topic_id
			addresses,    // text[]
			regrets,      // numeric(72,18)[]
		)
		placeholderCounter += 5
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, block_height, topic_id, addresses, regrets) 
			VALUES %s`, tableName, strings.Join(insertStatements, ","))

		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("failed to insert network regret event: %v", err)
		}
	} else {
		log.Info().Msg("No network regret event to insert")
	}

	return nil
}

func insertInfererNetworkRegret(events []EventRecord) error {
	return insertNetworkRegret(events, TB_INFERER_NETWORK_REGRET)
}

func insertForecasterNetworkRegret(events []EventRecord) error {
	return insertNetworkRegret(events, TB_FORECASTER_NETWORK_REGRET)
}

func insertNaiveInfererNetworkRegret(events []EventRecord) error {
	return insertNetworkRegret(events, TB_NAIVE_INFERER_NETWORK_REGRET)
}

func insertTopicInitialRegret(events []EventRecord) error {
	log.Info().Msg("Inserting topic initial regret")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1
	for _, event := range events {
		log.Trace().Interface("Event topic initial regret", event).Msg("Processing event topic initial regret")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var blockHeight uint64
		var topicID uint64
		var regret *big.Float

		for _, attr := range attributes {
			switch attr.Key {
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.ParseUint(cleanedValue, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			case "topic_id":
				cleanedValue := strings.Trim(attr.Value, "\"")
				topicID, err = strconv.ParseUint(cleanedValue, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert topic_id to int: %w", err)
				}
			case "regret":
				cleanedValue := strings.Trim(attr.Value, "\"")
				regret = new(big.Float)
				_, ok := regret.SetString(cleanedValue)
				if !ok {
					return fmt.Errorf("failed to parse regret: %s", cleanedValue)
				}
			}
		}

		newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d)",
			placeholderCounter,
			placeholderCounter+1,
			placeholderCounter+2,
			placeholderCounter+3)
		insertStatements = append(insertStatements, newStmt)
		values = append(values,
			topicID,              // topic_id BIGINT
			event.Height,         // height_tx BIGINT
			blockHeight,          // block_height BIGINT
			regret.Text('f', 18), // regret NUMERIC(72,18)
		)
		placeholderCounter += 4
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (topic_id, height_tx, block_height, regret) 
			VALUES %s`, TB_TOPIC_INITIAL_REGRET, strings.Join(insertStatements, ","))
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("failed to insert topic initial regret event: %v", err)
		}
	} else {
		log.Info().Msg("No topic initial regret event to insert")
	}
	return nil
}

func insertValidatorRewards(events []EventRecord) error {
	log.Info().Msg("Inserting validator rewards events")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1
	for _, event := range events {
		log.Trace().Interface("Event validator rewards", event).Msg("Processing event validator rewards")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var validator string
		var amount *big.Float

		for _, attr := range attributes {
			switch attr.Key {
			case "validator":
				validator = strings.Trim(attr.Value, "\"")
			case "amount":
				var err error
				amount, err = parseAmount(attr.Value)
				if err != nil {
					log.Error().Str("raw_amount", attr.Value).Err(err).Msg("Failed to parse amount")
					return fmt.Errorf("failed to parse amount: %w", err)
				}
			}
		}

		// Skip if we don't have both validator and amount
		if validator == "" || amount == nil {
			log.Warn().Interface("Event validator rewards", event).Msg("Skipping event due to missing validator or amount")
			continue
		}

		newStmt := fmt.Sprintf("($%d, $%d, $%d)",
			placeholderCounter,
			placeholderCounter+1,
			placeholderCounter+2)
		insertStatements = append(insertStatements, newStmt)
		values = append(values,
			event.Height,         // height_tx BIGINT
			validator,            // validator VARCHAR(255)
			amount.Text('f', 18), // amount NUMERIC(72,18)
		)
		placeholderCounter += 3
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (height_tx, validator, amount) 
			VALUES %s`, TB_VALIDATOR_REWARDS, strings.Join(insertStatements, ","))

		log.Debug().Str("SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for validator rewards")

		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("validator rewards insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No validator rewards events to insert")
	}

	return nil
}

func insertValidatorCommission(events []EventRecord) error {
	log.Info().Msg("Inserting validator commission events")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1
	for _, event := range events {
		log.Trace().Interface("Event validator commission", event).Msg("Processing event validator commission")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var validator string
		var amount *big.Float

		for _, attr := range attributes {
			switch attr.Key {
			case "validator":
				validator = strings.Trim(attr.Value, "\"")
			case "amount":
				var err error
				amount, err = parseAmount(attr.Value)
				if err != nil {
					log.Error().Str("raw_amount", attr.Value).Err(err).Msg("Failed to parse amount")
					return fmt.Errorf("failed to parse amount: %w", err)
				}
			}
		}

		// Skip if we don't have both validator and amount
		if validator == "" || amount == nil {
			log.Warn().Interface("Event validator commission", event).Msg("Skipping event due to missing validator or amount")
			continue
		}

		newStmt := fmt.Sprintf("($%d, $%d, $%d)",
			placeholderCounter,
			placeholderCounter+1,
			placeholderCounter+2)
		insertStatements = append(insertStatements, newStmt)
		values = append(values,
			event.Height,         // height_tx BIGINT
			validator,            // validator VARCHAR(255)
			amount.Text('f', 18), // amount NUMERIC(72,18)
		)
		placeholderCounter += 3
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
            INSERT INTO %s (height_tx, validator, amount) 
            VALUES %s`, TB_VALIDATOR_COMMISSION, strings.Join(insertStatements, ","))

		log.Debug().Str("SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for validator commission")

		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("validator commission insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No validator commission events to insert")
	}

	return nil
}

func insertValidatorWithdrawRewards(events []EventRecord) error {
	log.Info().Msg("Inserting validator withdraw rewards events")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1
	for _, event := range events {
		log.Trace().Interface("Event validator withdraw rewards", event).Msg("Processing event validator withdraw rewards")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var validator string
		var amount *big.Float

		for _, attr := range attributes {
			switch attr.Key {
			case "validator":
				validator = strings.Trim(attr.Value, "\"")
			case "amount":
				var err error
				amount, err = parseAmount(attr.Value)
				if err != nil {
					log.Error().Str("raw_amount", attr.Value).Err(err).Msg("Failed to parse amount")
					return fmt.Errorf("failed to parse amount: %w", err)
				}
			}
		}

		// Skip if we don't have both validator and amount
		if validator == "" || amount == nil {
			log.Warn().Interface("Event validator withdraw rewards", event).Msg("Skipping event due to missing validator or amount")
			continue
		}

		newStmt := fmt.Sprintf("($%d, $%d, $%d)",
			placeholderCounter,
			placeholderCounter+1,
			placeholderCounter+2)
		insertStatements = append(insertStatements, newStmt)
		values = append(values,
			event.Height,         // height_tx BIGINT
			validator,            // validator VARCHAR(255)
			amount.Text('f', 18), // amount NUMERIC(72,18)
		)
		placeholderCounter += 3
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
            INSERT INTO %s (height_tx, validator, amount) 
            VALUES %s`, TB_VALIDATOR_WITHDRAW_COMMISSION, strings.Join(insertStatements, ","))

		log.Debug().Str("SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for validator withdraw rewards")

		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("validator withdraw rewards insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No validator withdraw rewards events to insert")
	}

	return nil
}

func isDataEmpty(table string) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	err := dbPool.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to check if table %s is empty", table)
		return false, err
	}
	return count == 0, nil
}

func tableExists(tableName string) (bool, error) {
	var exists bool
	err := dbPool.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = $1
		)`, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func insertValueBundle(
	bundleId uint64,
	valueBundle types.MsgValueBundle,
	tableName string,
) error {

	//Insert InfererValues
	for _, val := range valueBundle.InfererValues {
		_, err := dbPool.Exec(context.Background(), `
				INSERT INTO `+tableName+` (
					bundle_id,
					reputer_value_type,
					worker,
					value
				) VALUES ($1, $2, $3, $4)`,
			bundleId, "InfererValues", val.Worker, val.Value,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert InfererValues bundle_values")
			return err
		}
	}
	//Insert ForecasterValues
	for _, val := range valueBundle.ForecasterValues {
		_, err := dbPool.Exec(context.Background(), `
				INSERT INTO `+tableName+` (
					bundle_id,
					reputer_value_type,
					worker,
					value
				) VALUES ($1, $2, $3, $4)`,
			bundleId, "ForecasterValues", val.Worker, val.Value,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert ForecasterValues bundle_values")
			return err
		}
	}
	// Insert OneOutInfererValues
	for _, val := range valueBundle.OneOutInfererValues {
		_, err := dbPool.Exec(context.Background(), `
				INSERT INTO `+tableName+` (
					bundle_id,
					reputer_value_type,
					worker,
					value
				) VALUES ($1, $2, $3, $4)`,
			bundleId, "OneOutInfererValues", val.Worker, val.Value,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert OneOutInfererValues bundle_values")
			return err
		}
	}
	// Insert OneInForecasterValues
	for _, val := range valueBundle.OneInForecasterValues {
		_, err := dbPool.Exec(context.Background(), `
				INSERT INTO `+tableName+` (
					bundle_id,
					reputer_value_type,
					worker,
					value
				) VALUES ($1, $2, $3, $4)`,
			bundleId, "OneInForecasterValues", val.Worker, val.Value,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert OneInForecasterValues bundle_values")
			return err
		}
	}
	// Insert OneOutForecasterValues
	for _, val := range valueBundle.OneOutForecasterValues {
		_, err := dbPool.Exec(context.Background(), `
				INSERT INTO `+tableName+` (
					bundle_id,
					reputer_value_type,
					worker,
					value
				) VALUES ($1, $2, $3, $4)`,
			bundleId, "OneOutForecasterValues", val.Worker, val.Value,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert OneOutForecasterValues bundle_values")
			return err
		}
	}
	// Insert OneOutInfererForecasterValues
	for _, val := range valueBundle.OneOutInfererForecasterValues {
		oneOutInfererStrValues := ""
		if len(val.OneOutInfererValues) != 0 {
			mjson, err := json.Marshal(val.OneOutInfererValues)
			if err != nil {
				log.Error().Err(err).Msg("Failed to insert OneOutInfererForecasterValues bundle_values")
				return err
			}
			oneOutInfererStrValues = string(mjson)
		}
		_, err := dbPool.Exec(context.Background(), `
				INSERT INTO `+tableName+` (
					bundle_id,
					reputer_value_type,
					worker,
					value
				) VALUES ($1, $2, $3, $4)`,
			bundleId, "OneOutInfererForecasterValues", val.Forecaster, oneOutInfererStrValues,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to insert OneOutInfererForecasterValues bundle_values")
			return err
		}
	}
	return nil
}

func insertTopicInitialEmaScore(events []EventRecord) error {
	log.Info().Msg("Inserting topic initial ema scores in batch")
	var insertStatements []string
	var values []interface{}

	placeholderCounter := 1 // Placeholder index starts at 1 in PostgreSQL

	for _, event := range events {
		log.Trace().Interface("Event score", event).Msg("Processing event score")
		var attributes []Attribute
		err := json.Unmarshal(event.Data, &attributes)
		if err != nil {
			return fmt.Errorf("failed to unmarshal event data: %w", err)
		}

		var topicID int
		var actorType string
		var score big.Float
		var blockHeight int

		for _, attr := range attributes {
			switch attr.Key {
			case "topic_id":
				cleanedValue := strings.Trim(attr.Value, "\"")
				topicID, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to convert topic_id to int: %w", err)
				}
			case "actor_type":
				actorType = strings.Trim(attr.Value, "\"")
			case "score":
				var rawScore string
				err = json.Unmarshal([]byte(attr.Value), &rawScore)
				if err != nil {
					return fmt.Errorf("failed to unmarshal scores: %w", err)
				}

				rawScoreClean := strings.Trim(rawScore, "\"")
				if isInvalidNumericValue(rawScoreClean) {
					log.Error().Str("rawScore", rawScore).Msg("Failed to convert score to big.Float")
					return fmt.Errorf("Invalid Score: %s", rawScoreClean)
				} else {
					_, ok := score.SetString(rawScoreClean)
					if !ok {
						log.Error().Str("rawScore", rawScore).Msg("Failed to convert score to big.Float")
						return fmt.Errorf("Invalid Score: %s", rawScoreClean)
					}
				}
			case "block_height":
				cleanedValue := strings.Trim(attr.Value, "\"")
				blockHeight, err = strconv.Atoi(cleanedValue)
				if err != nil {
					return fmt.Errorf("failed to convert block_height to int: %w", err)
				}
			}
		}

		// Generate the placeholders for this row
		newStmt := fmt.Sprintf("($%d, $%d, $%d, $%d)", placeholderCounter, placeholderCounter+1,
			placeholderCounter+2, placeholderCounter+3)
		insertStatements = append(insertStatements, newStmt)
		scoreValue := score.Text('f', -1)
		values = append(values, actorType, topicID, blockHeight, scoreValue)
		placeholderCounter += 4 // Increase counter for next row
	}

	if len(insertStatements) > 0 {
		sqlStatement := fmt.Sprintf(`
			INSERT INTO %s (actor_type, topic_id, block_height, score) 
			VALUES %s ON CONFLICT (topic_id, actor_type, block_height)
			DO UPDATE SET score=EXCLUDED.score`, TB_TOPIC_INITIAL_EMA_SCORES,
			strings.Join(insertStatements, ","))
		log.Trace().Str("Event - Initial EMA Scores SQL Statement", sqlStatement).Interface("Values", values).Msg("Executing batch insert for initial ema scores")
		_, err := dbPool.Exec(context.Background(), sqlStatement, values...)
		if err != nil {
			return fmt.Errorf("initial ema scores insert failed: %v", err)
		}
	} else {
		log.Info().Msg("No initial ema scores data to insert")
	}

	return nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func isInvalidNumericValue(value string) bool {
	return strings.Contains(strings.ToLower(value), "infinity") || strings.Contains(strings.ToLower(value), "nan")
}

// Helper function to strip denomination and parse amount
func parseAmount(amountStr string) (*big.Float, error) {
	// Strip 'uallo' from the end if present
	cleanedValue := strings.TrimSuffix(strings.Trim(amountStr, "\""), "uallo")
	amount := new(big.Float)
	_, ok := amount.SetString(cleanedValue)
	if !ok {
		return nil, fmt.Errorf("failed to parse amount: %s", amountStr)
	}
	return amount, nil
}
