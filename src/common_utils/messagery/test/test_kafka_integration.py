import pytest
from testcontainers.kafka import KafkaContainer

from kafka_services import ConsumerService, ProducerService


@pytest.fixture(scope='session')
def kafka_container():
    with KafkaContainer() as kafka:
        yield kafka


def test_services_integration(kafka_container):
    brokers = [kafka_container.get_bootstrap_server()]
    producer_service = ProducerService(brokers=brokers)
    topic = 'integration_test_topic_producer'
    message = {'key': 'value'}

    producer_service.send(topic, message)

    import time

    time.sleep(1)

    consumer_service = ConsumerService(
        brokers=brokers,
        topics=[topic],
        group_id='integration_test_group_producer',
    )
    consumer_service.subscribe([topic])
    records = consumer_service.poll(timeout_ms=5000)

    assert records == {topic: [message]}
    consumer_service.commit()
