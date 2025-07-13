rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"

for protocol in mahi mysticeti hotstuff_2 tusk bullshark; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop eth1
  for attack in noop; do
    for size in 32; do
      for load in 500 5000; do
        /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" "${attack}" eth1 "$load" "$size"
      done
    done
  done
done