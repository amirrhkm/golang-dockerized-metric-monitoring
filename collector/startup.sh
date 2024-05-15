#!/bin/bash

# Start the agent
otelcol-contrib --config=/etc/config.agent.yaml &

# Start the collector
otelcol-contrib --config=/etc/config.collector.yaml

# Wait for both processes to finish
wait
