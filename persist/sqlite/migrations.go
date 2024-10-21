package sqlite

import (
	"fmt"

	"go.sia.tech/core/types"
	"go.uber.org/zap"
)

// migrateVersion4 splits the height and ID of the last indexed tip into two
// separate columns for easier querying.
func migrateVersion4(tx *txn, _ *zap.Logger) error {
	var dbVersion int
	var indexMode int
	var elementNumLeaves uint64
	var index types.ChainIndex
	err := tx.QueryRow(`SELECT db_version, index_mode, element_num_leaves, last_indexed_tip FROM global_settings`).Scan(&dbVersion, &indexMode, &elementNumLeaves, decode(&index))
	if err != nil {
		return fmt.Errorf("failed to get last indexed tip: %w", err)
	} else if _, err := tx.Exec(`DROP TABLE global_settings`); err != nil {
		return fmt.Errorf("failed to drop global_settings: %w", err)
	}

	_, err = tx.Exec(`CREATE TABLE global_settings (
	id INTEGER PRIMARY KEY NOT NULL DEFAULT 0 CHECK (id = 0), -- enforce a single row
	db_version INTEGER NOT NULL, -- used for migrations
	index_mode INTEGER, -- the mode of the data store
	last_indexed_height INTEGER NOT NULL, -- the height of the last chain index that was processed
	last_indexed_id BLOB NOT NULL, -- the block ID of the last chain index that was processed
	element_num_leaves INTEGER NOT NULL -- the number of leaves in the state tree
);`)
	if err != nil {
		return fmt.Errorf("failed to create global_settings: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO global_settings (id, db_version, index_mode, last_indexed_height, last_indexed_id, element_num_leaves) VALUES (0, ?, ?, ?, ?, ?)`, dbVersion, indexMode, index.Height, encode(index.ID), elementNumLeaves)
	return err
}

// migrateVersion3 adds additional indices to event_addresses and wallet_addresses
// to improve query performance.
func migrateVersion3(tx *txn, _ *zap.Logger) error {
	_, err := tx.Exec(`CREATE INDEX event_addresses_event_id_address_id_idx ON event_addresses (event_id, address_id);
CREATE INDEX wallet_addresses_wallet_id_address_id_idx ON wallet_addresses (wallet_id, address_id);`)
	return err
}

// migrateVersion2 recreates indices and speeds up event queries
func migrateVersion2(tx *txn, _ *zap.Logger) error {
	_, err := tx.Exec(`DROP INDEX IF EXISTS chain_indices_height;
DROP INDEX IF EXISTS siacoin_elements_address_id;
DROP INDEX IF EXISTS siacoin_elements_maturity_height_matured;
DROP INDEX IF EXISTS siacoin_elements_chain_index_id;
DROP INDEX IF EXISTS siacoin_elements_spent_index_id;
DROP INDEX IF EXISTS siacoin_elements_address_id_spent_index_id;
DROP INDEX IF EXISTS siafund_elements_address_id;
DROP INDEX IF EXISTS siafund_elements_chain_index_id;
DROP INDEX IF EXISTS siafund_elements_spent_index_id;
DROP INDEX IF EXISTS siafund_elements_address_id_spent_index_id;
DROP INDEX IF EXISTS events_chain_index_id;
DROP INDEX IF EXISTS event_addresses_event_id_idx;
DROP INDEX IF EXISTS event_addresses_address_id_idx;
DROP INDEX IF EXISTS wallet_addresses_wallet_id;
DROP INDEX IF EXISTS wallet_addresses_address_id;
DROP INDEX IF EXISTS syncer_bans_expiration_index;

CREATE INDEX IF NOT EXISTS chain_indices_height_idx ON chain_indices (block_id, height);
CREATE INDEX IF NOT EXISTS siacoin_elements_address_id_idx ON siacoin_elements (address_id);
CREATE INDEX IF NOT EXISTS siacoin_elements_maturity_height_matured_idx ON siacoin_elements (maturity_height, matured);
CREATE INDEX IF NOT EXISTS siacoin_elements_chain_index_id_idx ON siacoin_elements (chain_index_id);
CREATE INDEX IF NOT EXISTS siacoin_elements_spent_index_id_idx ON siacoin_elements (spent_index_id);
CREATE INDEX IF NOT EXISTS siacoin_elements_address_id_spent_index_id_idx ON siacoin_elements(address_id, spent_index_id);
CREATE INDEX IF NOT EXISTS siafund_elements_address_id_idx ON siafund_elements (address_id);
CREATE INDEX IF NOT EXISTS siafund_elements_chain_index_id_idx ON siafund_elements (chain_index_id);
CREATE INDEX IF NOT EXISTS siafund_elements_spent_index_id_idx ON siafund_elements (spent_index_id);
CREATE INDEX IF NOT EXISTS siafund_elements_address_id_spent_index_id_idx ON siafund_elements(address_id, spent_index_id);
CREATE INDEX IF NOT EXISTS events_chain_index_id_idx ON events (chain_index_id);
CREATE INDEX IF NOT EXISTS events_maturity_height_id_idx ON events (maturity_height DESC, id DESC);
CREATE INDEX IF NOT EXISTS event_addresses_event_id_idx ON event_addresses (event_id);
CREATE INDEX IF NOT EXISTS event_addresses_address_id_idx ON event_addresses (address_id);
CREATE INDEX IF NOT EXISTS wallet_addresses_wallet_id_idx ON wallet_addresses (wallet_id);
CREATE INDEX IF NOT EXISTS wallet_addresses_address_id_idx ON wallet_addresses (address_id);
CREATE INDEX IF NOT EXISTS syncer_bans_expiration_index_idx ON syncer_bans (expiration);`)
	return err
}

// migrations is a list of functions that are run to migrate the database from
// one version to the next. Migrations are used to update existing databases to
// match the schema in init.sql.
var migrations = []func(tx *txn, log *zap.Logger) error{
	migrateVersion2,
	migrateVersion3,
	migrateVersion4,
}
