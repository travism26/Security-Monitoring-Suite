import express, { Request, Response } from 'express';
import { body } from 'express-validator';
import { validateRequest } from '../middlewares/validate-request';
import { metricsRegistry } from './metrics';
import { Counter } from 'prom-client';
import { SystemMetrics } from '../payload/system-metrics';

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
      console.log('Received metrics:', util.inspect(data, false, null, true));

      // Update Prometheus counter for incoming metrics
      const counter = metricsRegistry.getSingleMetric(
        'system_metrics_received_total'
      ) as Counter<string>;
      if (counter) {
        counter.inc();
      }

      // TODO: Add Kafka producer logic here

      res.status(202).json({
        status: 'accepted',
        timestamp: new Date().toISOString(),
      });
    } catch (error) {
      console.error('Error processing metrics:', error);
      res.status(500).json({
        errors: [{ message: 'Error processing metrics' }],
      });
    }
  }
);

export { router as systemMetricsRouter };
