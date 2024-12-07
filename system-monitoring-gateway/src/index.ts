import { app } from './app';
import { kafkaWrapper } from './kafka/kafka-wrapper';
import { SystemMetricsPublisher } from './events/system-metrics-publisher';
import { Topics } from './kafka/topics';

const start = async () => {
  console.log('Starting server...');
  let systemMetricsPublisher: SystemMetricsPublisher;
  if (process.env.KAFKA_BROKER) {
    // Initialize Kafka and add the system metrics producer
    const clientId = process.env.KAFKA_CLIENT_ID || 'system-monitoring-gateway';
    await kafkaWrapper.initialize([process.env.KAFKA_BROKER], clientId);
    await kafkaWrapper.addProducer(
      Topics.SystemMetrics,
      new SystemMetricsPublisher(kafkaWrapper.getClient())
    );
  }
  try {
    const shutdown = async () => {
      process.exit(0);
    };

    process.on('SIGTERM', shutdown);
    process.on('SIGINT', shutdown);

    app.listen(3000, () => {
      console.log('Listening on port 3000');
    });
  } catch (err) {
    console.error(err);
  }
};

start();
