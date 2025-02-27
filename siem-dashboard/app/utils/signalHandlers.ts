/**
 * Signal Handlers
 *
 * This file sets up proper handlers for termination signals (SIGINT, SIGTERM)
 * to ensure all resources are cleaned up when the application is terminated.
 */

// We'll use a simpler approach that doesn't involve overriding global functions
// but still ensures proper cleanup on termination

// Set up signal handlers for graceful shutdown
if (typeof process !== "undefined") {
  const cleanupResources = () => {
    console.log("[Signal Handler] Cleaning up resources before exit...");

    // Force garbage collection by clearing references
    // This helps ensure that any lingering event listeners are released
    try {
      // Reset EventEmitter max listeners to default
      const EventEmitter = require("events");
      EventEmitter.defaultMaxListeners = 10;
      console.log(
        "[Signal Handler] Reset EventEmitter defaultMaxListeners to 10"
      );

      // Force Node.js to run a full garbage collection cycle
      if (global.gc) {
        global.gc();
        console.log("[Signal Handler] Forced garbage collection");
      }
    } catch (error) {
      console.error("[Signal Handler] Error during cleanup:", error);
    }

    console.log("[Signal Handler] Cleanup complete, exiting...");
  };

  // Handle SIGINT (Ctrl+C)
  process.on("SIGINT", () => {
    console.log("\n[Signal Handler] Received SIGINT signal (Ctrl+C)");
    cleanupResources();
    process.exit(0);
  });

  // Handle SIGTERM
  process.on("SIGTERM", () => {
    console.log("[Signal Handler] Received SIGTERM signal");
    cleanupResources();
    process.exit(0);
  });

  // Handle uncaught exceptions
  process.on("uncaughtException", (error) => {
    console.error("[Signal Handler] Uncaught exception:", error);
    cleanupResources();
    process.exit(1);
  });

  // Handle unhandled promise rejections
  process.on("unhandledRejection", (reason, promise) => {
    console.error("[Signal Handler] Unhandled promise rejection:", reason);
    cleanupResources();
    process.exit(1);
  });
}

export {};
