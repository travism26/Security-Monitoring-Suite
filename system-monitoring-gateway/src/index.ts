import { app } from './app';
import { kafkaWrapper } from './kafka/kafka-wrapper';

const start = async () => {
  try {
    console.log('Starting server...');

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
