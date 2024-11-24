import express from 'express';
import 'express-async-errors'; // to handle async errors in express
import { json } from 'body-parser';
import cookieSession from 'cookie-session';
import { NotFoundError } from './errors/not-found-error';
import { errorHandler } from './middlewares/error-handlers';
// Import the metrics router
import { metricsRouter } from './routes/metrics';
import { systemMetricsRouter } from './routes/system-metrics';

const app = express();
app.set('trust proxy', true);
app.use(json());
app.use(
  cookieSession({ signed: false, secure: false }) //process.env.NODE_ENV !== 'test'
);

// Conditionally use the metrics router
if (process.env.NODE_ENV !== 'test') {
  app.use(metricsRouter);
}
app.use(systemMetricsRouter);
app.all('*', async (req, res) => {
  throw new NotFoundError();
});
app.use(errorHandler);

export { app };
