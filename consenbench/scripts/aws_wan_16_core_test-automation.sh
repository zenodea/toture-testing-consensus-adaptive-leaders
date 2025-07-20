rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"

for protocol in mahi mysticeti hotstuff_2 tusk bullshark; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop ens5
  for size in 512; do
    for load in 20000; do
      /bin/bash  consenbench/scripts/dry_run.sh  "${protocol}" noop ens5 "$load" "$size"
    done
  done
done

for protocol in dedis_paxos rabia sadl_racs efficient quepaxa cft_dag; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop ens5
  for size in 18; do
    for load in 100000; do
      /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" noop ens5 "$load" "$size"
    done
  done
done