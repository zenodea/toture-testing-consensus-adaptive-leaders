rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

echo "Running tests for all protocols..."

for protocol in dedis_paxos dedis_raft rabia sadl_racs racs efficient quepaxa mahi mysticeti hotstuff_2 hotstuff_3 tusk bullshark etcd zoo_keeper; do
  for attack in noop LeaderPartition OneQuorumNodePartition MajorityHighDelay MinorityCrash MinorityStraggler; do
    echo "..........Running ${protocol} ${attack}........"
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" "${attack}" eth1
  done
done
