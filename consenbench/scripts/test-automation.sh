rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

echo "Running tests for all protocols..."

for protocol in dedis_paxos dedis_raft sadl_racs racs efficient quepaxa hotstuff_2 hotstuff_3 etcd zoo_keeper; do
  for attack in noop leader_1 leader_2 leader_3 leader_4 leader_5 leader_6 leader_7 leader_8 leader_9; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" "${attack}" eth1
  done
done
