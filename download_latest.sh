#!/bin/bash

# Assign the GitHub API URL to a variable
github_api_url="https://api.github.com/repos/pobyzaarif/git-webhook-receiver/releases/latest"

# Step 1: Fetch the JSON response from the GitHub API using curl and store it in a variable
response=$(curl -s "$github_api_url")

# Function to extract value from JSON response based on field name
extract_field_value() {
    local field_name=$1
    local value=$(echo "$response" | grep -o "\"$field_name\": \"[^\"]*\"" | sed -n "s/\"$field_name\": \"\(.*\)\"/\1/p")
    echo "$value"
}

# Extract tag_name and published_at values from the JSON response
tag_name=$(extract_field_value "tag_name")
published_at=$(extract_field_value "published_at")

# Display tag_name and published_at
echo "Latest release:"
echo "Tag Name: $tag_name"
echo "Published At: $published_at"

# Step 2: Extract asset names and create a prompt for the user to choose
assets=$(echo "$response" | grep -o "\"name\": \"[^\"]*\"" | sed -n "s/\"name\": \"git-webhook-receiver-\(.*\).tar.gz\"/\1/p")
IFS=$'\n' read -rd '' -a assets_array <<< "$assets"

# Display the prompt for the user to choose an asset
echo "Choose an asset to download:"

# Loop through the assets and display a numbered list
for i in "${!assets_array[@]}"; do
    echo "$((i + 1)). ${assets_array[$i]}"
done

# Read user input for the chosen asset
read -p "Enter the number of the asset you want to download: " choice

# Validate user input and download the selected asset
if [[ $choice =~ ^[0-9]+$ && $choice -ge 1 && $choice -le ${#assets_array[@]} ]]; then
    selected_asset_index=$((choice - 1))
    selected_asset="${assets_array[selected_asset_index]}"

    # Get the download URL for the selected asset using awk
    download_url=$(echo "$response" | awk -v idx="${selected_asset_index}" -F'"' '/"browser_download_url"/ { if (idx-- == 0) print $4 }')

    # Construct the output filename without special characters
    output_filename="git-webhook-receiver-${selected_asset}.tar.gz"

    # Display the download URL for debugging
    echo "Download URL: $download_url"

    # Download the selected asset
    echo "Downloading $output_filename..."
    curl -L "$download_url" -o "$output_filename"
    echo "Download complete!"
else
    echo "Invalid choice. Please enter a valid number corresponding to the asset you want to download."
fi
