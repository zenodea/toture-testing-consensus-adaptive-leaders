rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

echo "Running tests for all protocols..."

for protocol in dedis_paxos; do
  for attack in noop LeaderPartition OneQuorumNodePartition MajorityHighDelay MinorityCrash MinorityStraggler; do
    /bin/bash  consenbench/scripts/dry_run.sh   "${protocol}" "${attack}" eth1
  done
done
