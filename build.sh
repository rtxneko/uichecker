#!/bin/bash

# Define the target architectures
architectures=("amd64" "386" "arm" "arm64" "loong64" "s390x" "mips64" "mips64" "mips64le" "mips64p32" "mips64p32le" "mipsle" "riscv" "riscv64" "s390")

# Define the target operating systems (e.g., linux, windows, darwin)
operating_systems=("linux" "windows" "darwin" "android" "aix" "dragonfly" "freebsd" "solaris")

# Specify the output directory
output_dir="bin"

# Name of the Go program to build
program_name="xuiScanner"
rm $output_dir/*
# Iterate over each architecture
for arch in "${architectures[@]}"; do
  # Iterate over each operating system
  for os in "${operating_systems[@]}"; do
    # Set environment variables for cross-compilation
    export GOARCH="$arch"
    export GOOS="$os"

    # Build the Go program and output binary to the specified directory
    go build -o "$output_dir/$program_name-$os-$arch" -ldflags="-s -w" . && tar -czvf $output_dir/$program_name-$os-$arch.tar.gz $output_dir/$program_name-$os-$arch && rm $output_dir/$program_name-$os-$arch
    echo build $program_name for $os/$arch complete

    # Clear environment variables
    unset GOARCH
    unset GOOS
  done
done

echo "Cross-build complete! Binaries are located in the $output_dir directory."
