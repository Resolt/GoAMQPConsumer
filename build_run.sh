#!/bin/sh

export AMQP_USER="guest"
export AMQP_PASS="guest"
export AMQP_HOST="localhost"
export AMQP_PORT="5672"
export AMQP_VHOST=""
export AMQP_EXCHANGE="ggce"
export AMQP_QUEUE="ggce"
export TAG="LELELELE"

go build -o gac && ./gac