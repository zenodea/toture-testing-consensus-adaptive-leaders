rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Num_Replicas,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"

for protocol in cft_dag dedis_paxos quepaxa sadl_racs; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop ens5
  size=18
  for load in 10000 50000 100000 150000 200000 250000 300000 350000 400000 450000 500000 600000; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" noop ens5 "$load" "$size"
  done
done