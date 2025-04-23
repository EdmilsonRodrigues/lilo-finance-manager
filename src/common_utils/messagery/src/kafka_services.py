import json

from kafka import KafkaConsumer, KafkaProducer
from kafka.errors import KafkaError


class ProducerService:
    def __init__(
        self, brokers: list[str], producer: KafkaProducer | None = None
    ) -> None:
        self.brokers = brokers
        self.producer_service = producer or KafkaProducer(
            bootstrap_servers=self.brokers,
            value_serializer=self._serialize_json,
        )

    def _serialize_json(self, value):
        return json.dumps(value).encode('utf-8')

    def send(self, topic: str, message: dict) -> None:
        try:
            self.producer_service.send(topic, value=message)
        except KafkaError as exc:
            raise exc


class ConsumerService:
    def __init__(
        self,
        brokers: list[str],
        topics: list[str],
        group_id: str,
        consumer: KafkaConsumer | None = None,
    ) -> None:
        self.brokers = brokers
        self.topics = topics
        self.group_id = group_id
        self.consumer = consumer or KafkaConsumer(
            *self.topics,
            bootstrap_servers=self.brokers,
            group_id=self.group_id,
            value_deserializer=self._deserialize_json,
            auto_offset_reset='earliest',
            enable_auto_commit=True,
            max_poll_records=1,
        )

    def _deserialize_json(self, value):
        try:
            return json.loads(value.decode('utf-8'))
        except json.JSONDecodeError:
            return value.decode('utf-8')

    def subscribe(self, topics: list[str]) -> None:
        self.consumer.subscribe(topics=topics)

    def poll(self, timeout_ms: int = 100) -> dict:
        return {
            topic.topic: [record.value for record in records]
            for topic, records in self.consumer.poll(
                timeout_ms=timeout_ms
            ).items()
        }

    def commit(self) -> None:
        self.consumer.commit()
