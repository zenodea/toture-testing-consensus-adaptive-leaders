rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"


for protocol in mahi mysticeti hotstuff_2 tusk bullshark cft_dag; do
  for attack in noop; do
    for size in 32 512; do
      for load in 10000 20000 50000 70000 100000 120000 150000 180000 200000 250000 300000; do
        /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" "${attack}" eth1 "$load" "$size"
      done
    done
  done
done