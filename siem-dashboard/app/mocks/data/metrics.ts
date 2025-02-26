import {
  SystemMetrics,
  CPUMetrics,
  MemoryMetrics,
  DiskMetrics,
  NetworkMetrics,
  DiskPartition,
} from "./types";
import { generateResourceMetric, generateTrafficValue } from "./utils";

function generateCPUMetrics(): CPUMetrics {
  const usage = Math.floor(Math.random() * 100);
  return {
    usage,
    temperature: 40 + Math.floor(Math.random() * 40), // 40-80°C
    processes: 100 + Math.floor(Math.random() * 900), // 100-1000 processes
    loadAverage: [Math.random() * 4, Math.random() * 3, Math.random() * 2],
  };
}

function generateMemoryMetrics(): MemoryMetrics {
  const total = 32 * 1024 * 1024 * 1024; // 32GB in bytes
  const used = Math.floor(Math.random() * total * 0.8); // Up to 80% usage
  const cached = Math.floor(Math.random() * (total - used) * 0.5);
  const free = total - used - cached;

  const swapTotal = 8 * 1024 * 1024 * 1024; // 8GB swap
  const swapUsed = Math.floor(Math.random() * swapTotal * 0.3); // Up to 30% swap usage

  return {
    total,
    used,
    free,
    cached,
    swapTotal,
    swapUsed,
  };
}

function generateDiskPartitions(): DiskPartition[] {
  const partitions = [
    { name: "/", size: 256 },
    { name: "/home", size: 512 },
    { name: "/var", size: 128 },
    { name: "/tmp", size: 64 },
  ];

  return partitions.map(({ name, size }) => {
    const total = size * 1024 * 1024 * 1024; // Convert GB to bytes
    const used = Math.floor(Math.random() * total * 0.9); // Up to 90% usage
    const free = total - used;

    return {
      mount: name,
      total,
      used,
      free,
    };
  });
}

function generateDiskMetrics(): DiskMetrics {
  const partitions = generateDiskPartitions();
  const total = partitions.reduce((acc, curr) => acc + curr.total, 0);
  const used = partitions.reduce((acc, curr) => acc + curr.used, 0);
  const free = total - used;

  return {
    total,
    used,
    free,
    readRate: Math.floor(Math.random() * 500 * 1024 * 1024), // 0-500 MB/s
    writeRate: Math.floor(Math.random() * 300 * 1024 * 1024), // 0-300 MB/s
    partitions,
  };
}

function generateNetworkMetrics(): NetworkMetrics {
  const bytesIn = generateTrafficValue(5000000, 2000000); // ~5MB/s ± 2MB/s
  const bytesOut = generateTrafficValue(3000000, 1000000); // ~3MB/s ± 1MB/s
  const packetsIn = Math.floor(bytesIn / 1000); // Rough estimate of packets
  const packetsOut = Math.floor(bytesOut / 1000);

  return {
    bytesIn,
    bytesOut,
    packetsIn,
    packetsOut,
    errors: Math.floor(Math.random() * 100),
    dropped: Math.floor(Math.random() * 50),
  };
}

export function generateMockSystemMetrics(): SystemMetrics {
  return {
    cpu: generateCPUMetrics(),
    memory: generateMemoryMetrics(),
    disk: generateDiskMetrics(),
    network: generateNetworkMetrics(),
    timestamp: new Date().toISOString(),
  };
}

// Generate a series of metrics over time
export function generateMetricsTimeSeries(
  minutes: number = 60,
  intervalSeconds: number = 60
): SystemMetrics[] {
  const now = new Date();
  const samples = (minutes * 60) / intervalSeconds;

  return Array.from({ length: samples }, (_, i) => ({
    ...generateMockSystemMetrics(),
    timestamp: new Date(
      now.getTime() - (samples - 1 - i) * intervalSeconds * 1000
    ).toISOString(),
  }));
}

// Generate metrics with a specific trend (e.g., increasing CPU usage)
export function generateTrendingMetrics(
  minutes: number = 60,
  intervalSeconds: number = 60,
  trend: {
    cpu?: boolean;
    memory?: boolean;
    disk?: boolean;
    network?: boolean;
  } = { cpu: true }
): SystemMetrics[] {
  const series = generateMetricsTimeSeries(minutes, intervalSeconds);
  const samples = series.length;

  return series.map((metrics, i) => {
    const trendFactor = i / (samples - 1); // 0 to 1

    if (trend.cpu) {
      metrics.cpu.usage = Math.min(
        95,
        30 + Math.floor(60 * trendFactor) + Math.floor(Math.random() * 10)
      );
    }

    if (trend.memory) {
      const memory = metrics.memory;
      const trendUsage = Math.min(
        0.9,
        0.3 + 0.5 * trendFactor + Math.random() * 0.1
      );
      memory.used = Math.floor(memory.total * trendUsage);
      memory.free = memory.total - memory.used - memory.cached;
    }

    if (trend.disk) {
      const disk = metrics.disk;
      const trendUsage = Math.min(
        0.95,
        0.4 + 0.4 * trendFactor + Math.random() * 0.1
      );
      disk.used = Math.floor(disk.total * trendUsage);
      disk.free = disk.total - disk.used;
    }

    if (trend.network) {
      const network = metrics.network;
      const baseFactor = 1 + trendFactor * 2; // 1x to 3x increase
      network.bytesIn *= baseFactor;
      network.bytesOut *= baseFactor;
      network.packetsIn = Math.floor(network.bytesIn / 1000);
      network.packetsOut = Math.floor(network.bytesOut / 1000);
    }

    return metrics;
  });
}
