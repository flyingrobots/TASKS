#!/bin/bash

# Check if a string argument was provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <string>"
    echo "Example: $0 'MyLogo'"
    exit 1
fi

# The string to generate logos for
STRING="$1"

# Array of fonts to use
FONTS=(
    '3d'
    'block'
    'chrome'
    'console'
    'grid'
    'huge'
    'pallet'
    'shade'
    'simple'
    'simple3d'
    'simpleBlock'
    'slick'
    'tiny'
)

# Generate logo for each font
for font in "${FONTS[@]}"; do
    echo "Generating logo with font: $font"
    npx oh-my-logo "$STRING" --filled --block-font "$font"
    echo "----------------------------------------"
done

echo "All logos generated successfully!"