DROP TABLE IF EXISTS `ton_txs`;

CREATE EXTERNAL TABLE IF NOT EXISTS `ton_txs` (
    `address` struct<
        `hex`:string,
        `base64`:string
    >,
    `hash` string,
    `hash_hex` string,
    `created_lt` bigint,
    `account` string,
    `workchain` bigint,
    `shard` bigint,
    `block_seq_no` bigint,
    `prev_tx_hash` string,
    `prev_tx_hash_hex` string,
    `prev_tx_lt` bigint,
    `in_msg_hash` string,
    `in_msg_hash_hex` string,
    `in_msg` string,
    `in_amount` string,
    `out_msg` string,
    `out_msg_count` bigint,
    `out_amount` string,
    `total_fees` string,
    `description_boc` string,
    `description` string,
    `compute_phase_exit_code` bigint,
    `action_phase_result_code` bigint,
    `orig_status` string,
    `end_status` string,
    `created_at` string,
    `timestamp_ms` bigint
)
PARTITIONED BY (`datekey` string)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
WITH SERDEPROPERTIES (
  'ignore.malformed.json' = 'TRUE',
  'dots.in.keys' = 'FALSE',
  'case.insensitive' = 'TRUE',
  'mapping' = 'TRUE'
)
STORED AS INPUTFORMAT 'org.apache.hadoop.mapred.TextInputFormat' OUTPUTFORMAT 'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION 's3://nimbus-ton-indexer/txs/ton-txs/'
TBLPROPERTIES ('classification' = 'json');

MSCK REPAIR TABLE ton_txs;

SELECT * FROM ton_txs WHERE datekey = '2024-06-24' LIMIT 10;