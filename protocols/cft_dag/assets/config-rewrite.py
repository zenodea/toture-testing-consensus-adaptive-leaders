import re
import sys

file_path = sys.argv[1]

with open(file_path, 'r') as file:
    lines = file.readlines()

updated_lines = []
for line in lines:
    match = re.match(r'(\s{2}network_address:\s*\d+\.\d+\.\d+\.\d+):\d+', line)
    if match:
        updated_line = f"{match.group(1)}:1500\n"
        updated_lines.append(updated_line)
    else:
        updated_lines.append(line)

with open(file_path, 'w') as file:
    file.writelines(updated_lines)

print("Updated network addresses to use port 1500.")
