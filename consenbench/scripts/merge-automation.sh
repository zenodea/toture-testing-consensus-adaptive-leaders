rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"

for protocol in mahi mysticeti hotstuff_2 tusk bullshark; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop eth1
  size=32
  for load in 600000 500000 400000 300000 200000 100000 80000 50000 30000 20000 10000 5000; do
    /bin/bash  consenbench/scripts/dry_run.sh  "${protocol}" noop eth1 "$load" "$size"
  done
  size=512
  for load in 200000 180000 150000 120000 100000 80000 50000 40000 30000 20000 10000 5000; do
    /bin/bash  consenbench/scripts/dry_run.sh  "${protocol}" noop eth1 "$load" "$size"
  done
done