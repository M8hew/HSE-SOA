-- Create the Kafka engine table
CREATE TABLE IF NOT EXISTS kafka_table (
    id UInt32,
    event_type Int64,
    post_id Int64,
    author Int64
) ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka:29092',
    kafka_topic_list = 'engagement-events',
    kafka_group_name = 'clickhouse_group',
    kafka_format = 'JSONEachRow',
    kafka_num_consumers = 1;

-- Create the real table where data will be stored
CREATE TABLE event_table (
    id UInt32,
    event_type Int64,
    post_id Int64,
    author Int64
) ENGINE = MergeTree()
ORDER BY id;

-- Create the materialized view to populate the real table from the Kafka engine table
CREATE MATERIALIZED VIEW kafka_to_event_table
TO event_table
AS SELECT
    id,
    event_type,
    post_id,
    author
FROM kafka_table;

