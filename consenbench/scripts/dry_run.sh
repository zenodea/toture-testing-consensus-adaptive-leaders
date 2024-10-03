algo=$1

git pull origin main
/bin/bash build.sh
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type bootstrap  --consensus_algorithm ${algo}  --attack_duration 60 --attack noop --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type copy       --consensus_algorithm ${algo}  --attack_duration 60 --attack noop --device enp1s0

./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack basic         --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack noop          --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack partition_1   --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack partition_2   --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack partition_3   --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack partition_4_1 --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack partition_4_2 --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack partition_5   --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack partition_6   --device enp1s0
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 5 --controller_operation_type run        --consensus_algorithm ${algo}  --attack_duration 60 --attack partition_7   --device enp1s0