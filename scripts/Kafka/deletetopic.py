from kafka.consumer import KafkaConsumer

from kafka.admin import KafkaAdminClient
from kafka.admin import NewTopic

admin_client = KafkaAdminClient(
    bootstrap_servers="0.0.0.0:9092",
    client_id='test',
    api_version=(2,7,0)
)

topic_list = []
topic_list.append("testtask")
admin_client.delete_topics(topic_list, timeout_ms=30)

consumer = KafkaConsumer(
    bootstrap_servers="0.0.0.0:9092",
    group_id='test'
)
topics = consumer.topics()
print(topics, type(topics))