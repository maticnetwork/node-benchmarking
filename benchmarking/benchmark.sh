#!/usr/bin/env sh

for i in {1..50}; do
  ./monitoring-tools fire -txs 10000 -clients 1 -seed $i -delay $i &
done
