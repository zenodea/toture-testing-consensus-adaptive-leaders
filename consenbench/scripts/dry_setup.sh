algo=$1
attack=$2
device=$3

./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 0 --controller_operation_type bootstrap  --consensus_algorithm ${algo}  --attack_duration 60 --attack ${attack} --device ${device}
./consenbench/bin/bench --is_controller --id 1 --debug_on --debug_level 0 --controller_operation_type copy       --consensus_algorithm ${algo}  --attack_duration 60 --attack ${attack} --device ${device}