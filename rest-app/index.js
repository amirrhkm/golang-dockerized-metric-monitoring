/**
 * Starts the express server and handles incoming requests.
 * @module index
 */

require("./meter");
const { metrics } = require("@opentelemetry/api");

/**
 * Retrieves the OpenTelemetry Meter for the express server.
 * @type {Meter}
 */
const meter = metrics.getMeter("express-server");

/**
 * Creates a counter metric to track the number of requests per name.
 * @type {Counter}
 */
let counter = meter.createCounter("name-req-count", {
  description: "The number of requests per name the server got",
});

const express = require("express");
const app = express();

/**
 * Handles GET requests to the /user/:name route.
 * @param {Request} req - The request object.
 * @param {Response} res - The response object.
 */
app.get("/user/:name", (req, res) => {
  const data = {
    route: "/user/:name",
    name: req.params.name,
  };
  counter.add(1, data);
  console.log({ data });
  res.send("Hello " + req.params.name);
});

/**
 * Starts the express server on port 8008.
 */
app.listen(8008, () => {
  console.log("Server is up and running");
});
