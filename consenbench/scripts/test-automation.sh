rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

echo "Running tests for all protocols..."

/bin/bash  consenbench/scripts/dry_run.sh dedis_paxos   noop eth1
/bin/bash  consenbench/scripts/dry_run.sh dedis_raft    noop eth1
/bin/bash  consenbench/scripts/dry_run.sh sadl_racs     noop eth1
/bin/bash  consenbench/scripts/dry_run.sh racs          noop eth1
/bin/bash  consenbench/scripts/dry_run.sh efficient     noop eth1
/bin/bash  consenbench/scripts/dry_run.sh quepaxa       noop eth1
/bin/bash  consenbench/scripts/dry_run.sh hotstuff_2    noop eth1
/bin/bash  consenbench/scripts/dry_run.sh hotstuff_3    noop eth1
/bin/bash  consenbench/scripts/dry_run.sh etcd          noop eth1
/bin/bash  consenbench/scripts/dry_run.sh zoo_keeper    noop eth1