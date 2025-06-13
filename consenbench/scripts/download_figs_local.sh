rm -r logs
mkdir logs
scp -i ~/.ssh/merge_key -r -J mergejump pasindut@torture-pasindut:~/toture-testing-consensus/logs/*  logs/