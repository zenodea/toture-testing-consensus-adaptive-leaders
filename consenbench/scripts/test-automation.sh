rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Load,Size,Protocol,Attack,Throughput,Latency(ms),Percentile_99(ms),CPU,MEM,NET_IN,NET_OUT"


for protocol in mahi mysticeti hotstuff_2 tusk bullshark; do
  for attack in noop; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" "${attack}" eth1 60000 256
  done
done