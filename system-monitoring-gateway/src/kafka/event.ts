import { Topics } from "./topics";

export interface Event<T = any> {
  topic: Topics;
  data: T;
}
