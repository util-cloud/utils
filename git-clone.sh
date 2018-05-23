#!/bin/bash

REPOS="a b c"

for repo in $REPOS
do
  git clone https://ipaddr/$repo.git
  echo "git clone https://ipaddr/$repo.git"
done
