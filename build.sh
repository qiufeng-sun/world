#!/bin/bash

#find dest dir
p=$(dirname $0)
echo "dir:$p"

#go to dest dir && build
out="world.run"
(cd $p && cd ../.. && . gvp && go env && cd - && pwd && go build -o $out && echo "build ok! out: $out")
