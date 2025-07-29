rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"

for protocol in cft_dag dedis_paxos sadl_racs quepaxa efficient; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop ens5
  size=18
  for load in 10000 100000 200000 250000 300000 400000 500000 600000 800000; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" noop ens5 "$load" "$size"
  done
done