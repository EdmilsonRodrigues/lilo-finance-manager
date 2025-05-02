from unittest.mock import MagicMock

import pytest
from kafka.errors import KafkaError

from kafka_services import ConsumerService, ProducerService


@pytest.fixture
def mock_kafka_producer():
    return MagicMock()


@pytest.fixture
def mock_kafka_consumer():
    return MagicMock()


def test_producer_initialization(mock_kafka_producer):
    producer_service = ProducerService(
        brokers=['localhost:9092'], producer=mock_kafka_producer
    )
    assert producer_service.producer_service is mock_kafka_producer
    assert producer_service.brokers == ['localhost:9092']


def test_producer_send_message(mock_kafka_producer):
    producer_service = ProducerService(
        brokers=['localhost:9092'], producer=mock_kafka_producer
    )
    topic = 'test_topic'
    message = {'key': 'value'}

    producer_service.send(topic, message)
    mock_kafka_producer.send.assert_called_once_with(topic, value=message)


def test_producer_send_failure(mock_kafka_producer):
    mock_kafka_producer.send.side_effect = KafkaError('Connection error')
    producer_service = ProducerService(
        brokers=['localhost:9092'], producer=mock_kafka_producer
    )
    topic = 'test_topic'
    message = {'key': 'value'}

    with pytest.raises(KafkaError):
        producer_service.send(topic, message)


def test_consumer_initialization(mock_kafka_consumer):
    topics = ['test_topic']
    group_id = 'my_group'
    consumer_service = ConsumerService(
        brokers=['localhost:9092'],
        topics=topics,
        group_id=group_id,
        consumer=mock_kafka_consumer,
    )

    assert consumer_service.consumer is mock_kafka_consumer
    assert consumer_service.brokers == ['localhost:9092']
    assert consumer_service.topics == topics
    assert consumer_service.group_id == group_id


def test_consumer_subscribe(mock_kafka_consumer):
    consumer_service = ConsumerService(
        brokers=['localhost:9092'],
        topics=[],
        group_id='my_group',
        consumer=mock_kafka_consumer,
    )
    topics_to_subscribe = ['test_topic_1', 'test_topic_2']
    consumer_service.subscribe(topics_to_subscribe)
    mock_kafka_consumer.subscribe.assert_called_once_with(
        topics=topics_to_subscribe
    )


def test_consumer_poll(mock_kafka_consumer):
    mock_message = MagicMock(value={'data': 'test'})
    mock_kafka_consumer.poll.return_value = {
        MagicMock(topic='test_partition'): [MagicMock(value=mock_message)]
    }
    consumer_service = ConsumerService(
        brokers=['localhost:9092'],
        topics=[],
        group_id='my_group',
        consumer=mock_kafka_consumer,
    )
    messages = consumer_service.poll()
    assert messages == {'test_partition': [mock_message]}
