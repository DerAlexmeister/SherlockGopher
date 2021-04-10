from kafka.consumer import KafkaConsumer

from kafka.admin import KafkaAdminClient
from kafka.admin import NewTopic

admin_client = KafkaAdminClient(
    bootstrap_servers="0.0.0.0:9092",
    client_id='test',
    api_version=(2,7,0)
)

topic_list = []
topic_list.append(NewTopic(name="testurl1", num_partitions=1, replication_factor=1))
admin_client.create_topics(new_topics=topic_list, validate_only=False)

consumer = KafkaConsumer(
    bootstrap_servers="0.0.0.0:9092",
    group_id='test'
)
topics = consumer.topics()
print(topics, type(topics))