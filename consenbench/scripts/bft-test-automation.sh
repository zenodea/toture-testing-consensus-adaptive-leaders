rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Num_Replicas,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"

for protocol in tusk bullshark hotstuff_2; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop eth1
  size=32
  for load in 10000 20000 50000 100000 150000 200000 300000 400000 500000 600000; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" noop eth1 "$load" "$size"
  done
done
