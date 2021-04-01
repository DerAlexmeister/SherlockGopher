from kafka.consumer import KafkaConsumer

from kafka.admin import KafkaAdminClient
from kafka.admin import NewTopic

consumer = KafkaConsumer(
    bootstrap_servers="0.0.0.0:9092",
    group_id='test'
)
topics = consumer.topics()
print(topics, type(topics))