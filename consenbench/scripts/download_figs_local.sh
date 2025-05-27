sudo rm -r logs
mkdir logs
scp -i ~/.ssh/merge_key -J mergejump pasindut@torture-pasindut:~/toture-testing-consensus/logs/*.pdf  logs/