sudo rm -r logs
mkdir logs
scp -r -i ~/.ssh/merge_key -J mergejump pasindut@torture-pasindut:~/toture-testing-consensus/logs/  logs/