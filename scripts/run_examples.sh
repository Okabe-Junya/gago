#!/bin/bash

examples=$(find ./examples -type f -name "*.go")

for example in $examples; do
    echo "Running example: $example"
    go run $example
    if [ $? -ne 0 ]; then
        echo "Example failed: $example"
        exit 1
    fi
done

echo "All examples ran successfully."
exit 0
