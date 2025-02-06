import axios from "axios";
import { Request, Response, NextFunction } from "express";

// This is a helper function that will create an axios client with
// some default options that are customized to our ingress-nginx
// setup. This will allow us to make requests to different services within
// the k8s cluster.
export default ({ req }: { req: Request }) => {
  if (typeof window === "undefined") {
    // we are on the server
    // requests should be made to http://SERVICENAME.NAMESPACE.svc.cluster.local
    // http://ingress-nginx-controller.ingress-nginx.svc.cluster.local
    return axios.create({
      baseURL:
        "http://ingress-nginx-controller.ingress-nginx.svc.cluster.local",
      headers: req.headers,
      withCredentials: true,
    });
  } else {
    // we are on the browser
    // requests can be made with a base url of ''
    return axios.create({
      baseURL: "/",
      withCredentials: true,
    });
  }
};

// Usage:
// import buildClient from "../api/buildClient";
// const { data } = await buildClient({ req }).get("/api/tickets");
