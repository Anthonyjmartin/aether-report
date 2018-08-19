#!/usr/bin/env bash

echo "df -h output:"
df -h | grep -v tmp | grep -v loop | grep -v udev
echo
echo "aether-report disk output:"
./aether-report disk -h
echo
echo "df blocks output:"
df | grep -v tmp | grep -v loop | grep -v udev
echo
echo "aether-report disk Blocks:"
./aether-report disk
echo
echo "df -i output:"
df -i | grep -v tmp | grep -v loop | grep -v udev
echo
echo "aether-report disk Inodes:"
./aether-report disk -i
echo
echo "aether-report disk JSON:"
./aether-report disk -o json
echo
echo "aether-report disk -i and -h join error:"
./aether-report disk -i -h
echo
echo "aether-report missing subcommand error:"
./aether-report
