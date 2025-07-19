rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"

#for protocol in mahi mysticeti hotstuff_2 tusk bullshark; do
for protocol in mahi; do
  /bin/bash  consenbench/scripts/dry_setup.sh  "${protocol}" noop ens5
  for attack in noop; do
    for size in 512; do
      for load in 100000; do
        /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" "${attack}" ens5 "$load" "$size"
      done
    done
  done
done