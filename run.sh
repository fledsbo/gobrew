#!/bin/bash

until ./gobrew; do
  echo "Failure, restarting"
  sleep 1
done

