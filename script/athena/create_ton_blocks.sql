DROP TABLE IF EXISTS `ton_blocks`;

CREATE EXTERNAL TABLE IF NOT EXISTS `ton_blocks` (
    `workchain` bigint,
    `shard` bigint,
    `seq_no` bigint,
    `file_hash` string,
    `file_hash_hex` string,
    `root_hash` string,
    `root_hash_hex` string,
    `master` string,
    `shards` string,
    `scanned_at` string,
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
LOCATION 's3://nimbus-ton-indexer/blocks/ton-blocks/'
TBLPROPERTIES ('classification' = 'json');

MSCK REPAIR TABLE ton_blocks;

SELECT * FROM ton_blocks WHERE datekey = '2024-06-24' LIMIT 10;