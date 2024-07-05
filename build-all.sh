#!/bin/bash

for dir in */ ; do
	if [ -f "$dir/pom.xml" ]; then
		(cd "$dir" && mvn compile jib:dockerBuild)
	fi

done
