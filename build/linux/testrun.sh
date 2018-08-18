#!/usr/bin/env bash

echo "df -h output:"
df -h | grep -v tmp | grep -v loop | grep -v udev
echo
echo "aether-report output:"
./aether-report -h
echo
echo "df blocks output:"
df | grep -v tmp | grep -v loop | grep -v udev
echo
echo "aether-report Blocks:"
./aether-report
echo
echo "df -i output:"
df -i | grep -v tmp | grep -v loop | grep -v udev
echo
echo "aether-report Inodes:"
./aether-report -i
echo
echo "aether-report -i and -h join error:"
./aether-report -i -h
echo
echo "aether-report JSON:"
./aether-report -o json
