DROP TABLE IF EXISTS `ton_messages`;

CREATE EXTERNAL TABLE IF NOT EXISTS `ton_messages` (
    `type` string,
    `hash` string,
    `hash_hex` string,
    `src_address` struct<
        `hex`:string,
        `base64`:string
    >,
    `src_tx_lt` bigint,
    `src_tx_hash` string,
    `src_tx_hash_hex` string,
    `src_workchain` bigint,
    `src_shard` bigint,
    `src_block_seq_no` bigint,
    `src_state` string,
    `dst_address` struct<
        `hex`:string,
        `base64`:string
    >,
    `dst_tx_lt` bigint,
    `dst_tx_hash` string,
    `dst_tx_hash_hex` string,
    `dst_workchain` bigint,
    `dst_shard` bigint,
    `dst_block_seq_no` bigint,
    `dst_state` string,
    `bounce` boolean,
    `bounced` boolean,
    `amount` string,
    `ihr_disable` boolean,
    `ihr_fee` string,
    `fwd_fee` string,
    `body` string,
    `body_hash` string,
    `operation_id` bigint,
    `transfer_comment` string,
    `state_init_code` string,
    `state_init_data` string,
    `src_contract` string,
    `dst_contract` string,
    `operation_name` string,
    `data` string,
    `error` string,
    `created_at` string,
    `created_lt` bigint,
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
LOCATION 's3://nimbus-ton-indexer/messages/ton-messages/'
TBLPROPERTIES ('classification' = 'json');

MSCK REPAIR TABLE ton_messages;

SELECT * FROM ton_messages WHERE datekey = '2024-06-24' LIMIT 10;