#!/bin/bash

ROOT_DIR=~/Documents/toture-testing-consensus/logs
OUTPUT_BASE=~/Documents/toture-testing-consensus/logs/pdf-to-jpg-output

ROOT_DIR="${ROOT_DIR/#\~/$HOME}"
OUTPUT_BASE="${OUTPUT_BASE/#\~/$HOME}"

mkdir -p "$OUTPUT_BASE"

find "$ROOT_DIR" -type f -iname "*.pdf" | while IFS= read -r pdf_file; do
    rel_path="${pdf_file#$ROOT_DIR/}"
    rel_dir=$(dirname "$rel_path")

    output_dir="$OUTPUT_BASE/$rel_dir"
    mkdir -p "$output_dir"

    base_name=$(basename "$pdf_file" .pdf)

    pdftoppm -jpeg "$pdf_file" "$output_dir/${base_name}"
done