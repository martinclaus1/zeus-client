#!/usr/bin/env bash

package=zeus-client.go

platforms=("windows/amd64" "darwin/amd64"  "darwin/arm64" "linux/amd64" "linux/arm64")

rm -rf bin
for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	output_name='bin/'$GOOS'/zeus-client-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi

	env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
done