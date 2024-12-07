import express, { Request, Response } from 'express';
import { body } from 'express-validator';
import { validateRequest } from '../middlewares/validate-request';
import { metricsRegistry } from './metrics';
import { Counter } from 'prom-client';
import { SystemMetrics } from '../payload/system-metrics';
import { kafkaWrapper } from '../kafka/kafka-wrapper';

const router = express.Router();

// Validation middleware
const validateMetrics = [
  body('data').notEmpty().withMessage('Data is required'),
  body('data.metrics').notEmpty().withMessage('Metrics data is required'),
  body('timestamp').isISO8601().withMessage('Invalid timestamp format'),
];

router.post(
  '/api/v1/system-metrics/ingest',
  validateMetrics,
  validateRequest,
  async (req: Request<{}, {}, SystemMetrics>, res: Response) => {
    try {
      const { data, timestamp } = req.body;
      const util = require('util');

      // Create a copy of data without processes
      const filteredData = { ...data };
      delete filteredData.metrics.processes;

      console.log(
        'Received metrics:',
        util.inspect(filteredData, false, null, true)
      );

      // Update Prometheus counter for incoming metrics
      const counter = metricsRegistry.getSingleMetric(
        'system_metrics_received_total'
      ) as Counter<string>;
      if (counter) {
        counter.inc();
      }

      // Attempt to publish to Kafka
      if (!kafkaWrapper.isInitialized()) {
        return res.status(503).json({
          errors: [
            {
              message: 'Metrics service temporarily unavailable',
              details: 'Kafka connection not established',
            },
          ],
        });
      }

      try {
        const kafkaProducer = kafkaWrapper.getProducer('system-metrics');
        await kafkaProducer.publish({
          ...data,
        });

        // Only send success response if Kafka publish succeeds
        return res.status(202).json({
          status: 'accepted',
          timestamp: new Date().toISOString(),
        });
      } catch (kafkaError) {
        console.error('Error producing metrics to Kafka:', kafkaError);
        // Return early with Kafka-specific error
        return res.status(202).json({
          errors: [
            {
              message: 'Metrics service temporarily unavailable',
              details: 'Kafka connection not established',
            },
          ],
        });
      }
    } catch (error) {
      console.error('Error processing metrics:', error);
      // Handle any other unexpected errors
      return res.status(500).json({
        errors: [{ message: 'Internal server error while processing metrics' }],
      });
    }
  }
);

export { router as systemMetricsRouter };
