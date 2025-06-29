import re
import csv
import sys

def extract_to_csv(input_filename, output_filename="logs/extracted_data.csv"):
    pattern = re.compile(r'^\.*Running (\S+)\s+(\S+)\.*$')

    results = []
    current_algo = None
    current_attack = None
    current_content_lines = []

    with open(input_filename, 'r') as file:
        for line in file:
            line = line.strip()
            match = pattern.match(line)
            if match:

                if current_algo and current_attack:
                    content = '\n'.join(current_content_lines).strip()
                    results.append((current_algo, current_attack, content))
                    current_content_lines = []

                current_algo, current_attack = match.groups()
            else:
                if current_algo and current_attack:
                    current_content_lines.append(line)

        # Append the final block if it exists
        if current_algo and current_attack and current_content_lines:
            content = '\n'.join(current_content_lines).strip()
            results.append((current_algo, current_attack, content))

    # Write to CSV
    with open(output_filename, 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(['Algo', 'Attack', 'Content'])
        for row in results:
            writer.writerow(row)

    print(f"Data successfully written to {output_filename}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python extract_log_to_csv.py <log_file>")
    else:
        extract_to_csv(sys.argv[1])
