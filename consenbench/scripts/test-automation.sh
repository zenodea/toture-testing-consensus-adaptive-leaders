rm -r final-results
rm -r logs

mkdir final-results
mkdir logs

echo "Running tests for all protocols..."

/bin/bash  consenbench/scripts/dry_run.sh quepaxa noop eth1
/bin/bash  consenbench/scripts/dry_run.sh dedis_paxos noop eth1
/bin/bash  consenbench/scripts/dry_run.sh quepaxa leader_1 eth1