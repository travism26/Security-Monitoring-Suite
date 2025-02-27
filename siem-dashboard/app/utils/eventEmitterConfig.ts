/**
 * Event Emitter Configuration
 *
 * This file configures the Node.js EventEmitter to increase the default max listeners limit.
 * This helps prevent the "MaxListenersExceededWarning: Possible EventEmitter memory leak detected"
 * warning that can occur when multiple components with polling intervals create HTTP requests.
 */

import { EventEmitter } from "events";

// Increase the default max listeners from 10 to 20
// This should be sufficient for our dashboard components that use polling
EventEmitter.defaultMaxListeners = 20;

console.log("[EventEmitter Config] Increased defaultMaxListeners to 20");

// Export a function to set a custom max listeners value if needed
export function setMaxListeners(count: number): void {
  EventEmitter.defaultMaxListeners = count;
  console.log(`[EventEmitter Config] Updated defaultMaxListeners to ${count}`);
}

// Export a function to get the current max listeners value
export function getMaxListeners(): number {
  return EventEmitter.defaultMaxListeners;
}
