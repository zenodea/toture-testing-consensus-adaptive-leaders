rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

/bin/bash build.sh

echo "Running tests for all protocols..."
echo
echo
echo

for protocol in dedis_paxos dedis_raft rabia sadl_racs racs efficient quepaxa mahi mysticeti hotstuff_2 hotstuff_3 tusk bullshark etcd zoo_keeper cft_dag; do
  for attack in noop LeaderPartition OneQuorumNodePartition MajorityHighDelay MinorityCrash MinorityStraggler; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" "${attack}" eth1
  done
done
