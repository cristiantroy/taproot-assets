// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: transfers.sql

package sqlite

import (
	"context"
	"database/sql"
	"time"
)

const applySpendDelta = `-- name: ApplySpendDelta :one
WITH old_script_key_id AS (
    SELECT script_key_id
    FROM script_keys
    WHERE tweaked_script_key = $5
)
UPDATE assets
SET amount = $1, script_key_id = $2, 
    split_commitment_root_hash = $3,
    split_commitment_root_value = $4
WHERE script_key_id in (SELECT script_key_id FROM old_script_key_id)
RETURNING asset_id
`

type ApplySpendDeltaParams struct {
	NewAmount                int64
	NewScriptKeyID           int32
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	OldScriptKey             []byte
}

func (q *Queries) ApplySpendDelta(ctx context.Context, arg ApplySpendDeltaParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, applySpendDelta,
		arg.NewAmount,
		arg.NewScriptKeyID,
		arg.SplitCommitmentRootHash,
		arg.SplitCommitmentRootValue,
		arg.OldScriptKey,
	)
	var asset_id int32
	err := row.Scan(&asset_id)
	return asset_id, err
}

const deleteAssetWitnesses = `-- name: DeleteAssetWitnesses :exec
DELETE FROM asset_witnesses
WHERE asset_id = ?
`

func (q *Queries) DeleteAssetWitnesses(ctx context.Context, assetID int32) error {
	_, err := q.db.ExecContext(ctx, deleteAssetWitnesses, assetID)
	return err
}

const deleteSpendProofs = `-- name: DeleteSpendProofs :exec
DELETE FROM transfer_proofs
WHERE transfer_id = ?
`

func (q *Queries) DeleteSpendProofs(ctx context.Context, transferID int32) error {
	_, err := q.db.ExecContext(ctx, deleteSpendProofs, transferID)
	return err
}

const fetchAssetDeltas = `-- name: FetchAssetDeltas :many
SELECT  
    deltas.old_script_key, deltas.new_amt, 
    script_keys.tweaked_script_key AS new_script_key_bytes,
    script_keys.tweak AS script_key_tweak,
    deltas.new_script_key AS new_script_key_id, 
    internal_keys.raw_key AS new_raw_script_key_bytes,
    internal_keys.key_family AS new_script_key_family, 
    internal_keys.key_index AS new_script_key_index,
    deltas.serialized_witnesses, split_commitment_root_hash, 
    split_commitment_root_value
FROM asset_deltas deltas
JOIN script_keys
    ON deltas.new_script_key = script_keys.script_key_id
JOIN internal_keys 
    ON script_keys.internal_key_id = internal_keys.key_id
WHERE transfer_id = ?
`

type FetchAssetDeltasRow struct {
	OldScriptKey             []byte
	NewAmt                   int64
	NewScriptKeyBytes        []byte
	ScriptKeyTweak           []byte
	NewScriptKeyID           int32
	NewRawScriptKeyBytes     []byte
	NewScriptKeyFamily       int32
	NewScriptKeyIndex        int32
	SerializedWitnesses      []byte
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
}

func (q *Queries) FetchAssetDeltas(ctx context.Context, transferID int32) ([]FetchAssetDeltasRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAssetDeltas, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchAssetDeltasRow
	for rows.Next() {
		var i FetchAssetDeltasRow
		if err := rows.Scan(
			&i.OldScriptKey,
			&i.NewAmt,
			&i.NewScriptKeyBytes,
			&i.ScriptKeyTweak,
			&i.NewScriptKeyID,
			&i.NewRawScriptKeyBytes,
			&i.NewScriptKeyFamily,
			&i.NewScriptKeyIndex,
			&i.SerializedWitnesses,
			&i.SplitCommitmentRootHash,
			&i.SplitCommitmentRootValue,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchAssetDeltasWithProofs = `-- name: FetchAssetDeltasWithProofs :many
SELECT  
    deltas.old_script_key, deltas.new_amt, 
    script_keys.tweaked_script_key AS new_script_key_bytes,
    script_keys.tweak AS script_key_tweak,
    deltas.new_script_key AS new_script_key_id, 
    internal_keys.raw_key AS new_raw_script_key_bytes,
    internal_keys.key_family AS new_script_key_family, 
    internal_keys.key_index AS new_script_key_index,
    deltas.serialized_witnesses, deltas.split_commitment_root_hash, 
    deltas.split_commitment_root_value, transfer_proofs.sender_proof,
    transfer_proofs.receiver_proof
FROM asset_deltas deltas
JOIN script_keys
    ON deltas.new_script_key = script_keys.script_key_id
JOIN internal_keys 
    ON script_keys.internal_key_id = internal_keys.key_id
JOIN transfer_proofs
    ON deltas.proof_id = transfer_proofs.proof_id
WHERE deltas.transfer_id = ?
`

type FetchAssetDeltasWithProofsRow struct {
	OldScriptKey             []byte
	NewAmt                   int64
	NewScriptKeyBytes        []byte
	ScriptKeyTweak           []byte
	NewScriptKeyID           int32
	NewRawScriptKeyBytes     []byte
	NewScriptKeyFamily       int32
	NewScriptKeyIndex        int32
	SerializedWitnesses      []byte
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	SenderProof              []byte
	ReceiverProof            []byte
}

func (q *Queries) FetchAssetDeltasWithProofs(ctx context.Context, transferID int32) ([]FetchAssetDeltasWithProofsRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAssetDeltasWithProofs, transferID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchAssetDeltasWithProofsRow
	for rows.Next() {
		var i FetchAssetDeltasWithProofsRow
		if err := rows.Scan(
			&i.OldScriptKey,
			&i.NewAmt,
			&i.NewScriptKeyBytes,
			&i.ScriptKeyTweak,
			&i.NewScriptKeyID,
			&i.NewRawScriptKeyBytes,
			&i.NewScriptKeyFamily,
			&i.NewScriptKeyIndex,
			&i.SerializedWitnesses,
			&i.SplitCommitmentRootHash,
			&i.SplitCommitmentRootValue,
			&i.SenderProof,
			&i.ReceiverProof,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchSpendProofs = `-- name: FetchSpendProofs :one
SELECT sender_proof, receiver_proof
FROM transfer_proofs
WHERE transfer_id = ?
`

type FetchSpendProofsRow struct {
	SenderProof   []byte
	ReceiverProof []byte
}

func (q *Queries) FetchSpendProofs(ctx context.Context, transferID int32) (FetchSpendProofsRow, error) {
	row := q.db.QueryRowContext(ctx, fetchSpendProofs, transferID)
	var i FetchSpendProofsRow
	err := row.Scan(&i.SenderProof, &i.ReceiverProof)
	return i, err
}

const insertAssetDelta = `-- name: InsertAssetDelta :exec
INSERT INTO asset_deltas (
    old_script_key, new_amt, new_script_key, serialized_witnesses, transfer_id,
    proof_id, split_commitment_root_hash, split_commitment_root_value
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
)
`

type InsertAssetDeltaParams struct {
	OldScriptKey             []byte
	NewAmt                   int64
	NewScriptKey             int32
	SerializedWitnesses      []byte
	TransferID               int32
	ProofID                  int32
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
}

func (q *Queries) InsertAssetDelta(ctx context.Context, arg InsertAssetDeltaParams) error {
	_, err := q.db.ExecContext(ctx, insertAssetDelta,
		arg.OldScriptKey,
		arg.NewAmt,
		arg.NewScriptKey,
		arg.SerializedWitnesses,
		arg.TransferID,
		arg.ProofID,
		arg.SplitCommitmentRootHash,
		arg.SplitCommitmentRootValue,
	)
	return err
}

const insertAssetTransfer = `-- name: InsertAssetTransfer :one
INSERT INTO asset_transfers (
    old_anchor_point, new_internal_key, new_anchor_utxo, transfer_time_unix
) VALUES (
    ?, ?, ?, ?
) RETURNING id
`

type InsertAssetTransferParams struct {
	OldAnchorPoint   []byte
	NewInternalKey   int32
	NewAnchorUtxo    int32
	TransferTimeUnix time.Time
}

func (q *Queries) InsertAssetTransfer(ctx context.Context, arg InsertAssetTransferParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertAssetTransfer,
		arg.OldAnchorPoint,
		arg.NewInternalKey,
		arg.NewAnchorUtxo,
		arg.TransferTimeUnix,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const insertSpendProofs = `-- name: InsertSpendProofs :one
INSERT INTO transfer_proofs (
   transfer_id, sender_proof, receiver_proof 
) VALUES (
    ?, ?, ?
) RETURNING proof_id
`

type InsertSpendProofsParams struct {
	TransferID    int32
	SenderProof   []byte
	ReceiverProof []byte
}

func (q *Queries) InsertSpendProofs(ctx context.Context, arg InsertSpendProofsParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertSpendProofs, arg.TransferID, arg.SenderProof, arg.ReceiverProof)
	var proof_id int32
	err := row.Scan(&proof_id)
	return proof_id, err
}

const queryAssetTransfers = `-- name: QueryAssetTransfers :many
SELECT 
    asset_transfers.old_anchor_point, utxos.outpoint AS new_anchor_point,
    utxos.taro_root, utxos.tapscript_sibling, 
    utxos.utxo_id AS new_anchor_utxo_id, txns.raw_tx AS anchor_tx_bytes, 
    txns.txid AS anchor_txid, txns.txn_id AS anchor_tx_primary_key, 
    txns.chain_fees, transfer_time_unix, keys.raw_key AS internal_key_bytes,
    keys.key_family AS internal_key_fam, keys.key_index AS internal_key_index,
    id AS transfer_id, transfer_time_unix
FROM asset_transfers
JOIN internal_keys keys
    ON asset_transfers.new_internal_key = keys.key_id
JOIN managed_utxos utxos
    ON asset_transfers.new_anchor_utxo = utxos.utxo_id
JOIN chain_txns txns
    ON utxos.utxo_id = txns.txn_id
WHERE (
    -- We'll use this clause to filter out for only transfers that are
    -- unconfirmed. But only if the unconf_only field is set.
    -- TODO(roasbeef): just do the confirmed bit, 
    (($1 == 0 OR $1 IS NULL)
        OR
    (($1 == 1) == (length(hex(txns.block_hash)) == 0)))

    AND
    
    -- Here we have another optional query clause to select a given transfer
    -- based on the new_anchor_point, but only if it's specified.
    (length(hex($2)) == 0 OR 
        utxos.outpoint = $2)
)
`

type QueryAssetTransfersParams struct {
	UnconfOnly     interface{}
	NewAnchorPoint interface{}
}

type QueryAssetTransfersRow struct {
	OldAnchorPoint     []byte
	NewAnchorPoint     []byte
	TaroRoot           []byte
	TapscriptSibling   []byte
	NewAnchorUtxoID    int32
	AnchorTxBytes      []byte
	AnchorTxid         []byte
	AnchorTxPrimaryKey int32
	ChainFees          int64
	TransferTimeUnix   time.Time
	InternalKeyBytes   []byte
	InternalKeyFam     int32
	InternalKeyIndex   int32
	TransferID         int32
	TransferTimeUnix_2 time.Time
}

func (q *Queries) QueryAssetTransfers(ctx context.Context, arg QueryAssetTransfersParams) ([]QueryAssetTransfersRow, error) {
	rows, err := q.db.QueryContext(ctx, queryAssetTransfers, arg.UnconfOnly, arg.NewAnchorPoint)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QueryAssetTransfersRow
	for rows.Next() {
		var i QueryAssetTransfersRow
		if err := rows.Scan(
			&i.OldAnchorPoint,
			&i.NewAnchorPoint,
			&i.TaroRoot,
			&i.TapscriptSibling,
			&i.NewAnchorUtxoID,
			&i.AnchorTxBytes,
			&i.AnchorTxid,
			&i.AnchorTxPrimaryKey,
			&i.ChainFees,
			&i.TransferTimeUnix,
			&i.InternalKeyBytes,
			&i.InternalKeyFam,
			&i.InternalKeyIndex,
			&i.TransferID,
			&i.TransferTimeUnix_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const reanchorAssets = `-- name: ReanchorAssets :exec
WITH assets_to_update AS (
    SELECT asset_id
    FROM assets
    JOIN managed_utxos utxos
        ON assets.anchor_utxo_id = utxos.utxo_id
    WHERE utxos.outpoint = $2
)
UPDATE assets
SET anchor_utxo_id = $1
WHERE asset_id IN (SELECT asset_id FROM assets_to_update)
`

type ReanchorAssetsParams struct {
	NewOutpointUtxoID sql.NullInt32
	OldOutpoint       []byte
}

func (q *Queries) ReanchorAssets(ctx context.Context, arg ReanchorAssetsParams) error {
	_, err := q.db.ExecContext(ctx, reanchorAssets, arg.NewOutpointUtxoID, arg.OldOutpoint)
	return err
}
