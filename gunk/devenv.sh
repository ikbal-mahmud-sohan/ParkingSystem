#!/usr/bin/env bash

readonly gunk_dir=$(git rev-parse --show-toplevel)/gunk

PATH=$gunk_dir/bin:$gunk_dir/node_modules/.bin:$PATH

export PATH
