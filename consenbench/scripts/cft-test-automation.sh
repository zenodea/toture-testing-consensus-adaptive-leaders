rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Num_Replicas,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"

for protocol in dedis_paxos dedis_raft sadl_racs racs quepaxa; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop eth1

  size=18
  for load in 1000 10000 50000 100000 150000 200000 250000 300000 350000 400000; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" noop eth1 "$load" "$size"
  done

  size=32
  for load in 1000 10000 50000 100000 150000 200000 250000 300000 350000 400000; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" noop eth1 "$load" "$size"
  done

  size=512
  for load in 1000 5000 10000 20000 30000 50000; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" noop eth1 "$load" "$size"
  done

done

for protocol in rabia efficient efficient-exec; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop eth1
  size=18
  for load in 1000 10000 50000 100000 150000 200000 300000 400000 500000 600000 700000 800000 900000 1000000; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" noop eth1 "$load" "$size"
  done
done
