#!/bin/bash

DEBIAN_RELEASES=$(debian-distro-info --supported)
UBUNTU_RELEASES=$(ubuntu-distro-info --supported-esm)

mkdir -p cli-repo/deb
cd cli-repo/deb

for release in ${DEBIAN_RELEASES[@]} ${UBUNTU_RELEASES[@]}; do
  echo "Removing deb package of $release"
  reprepro -A i386 remove $release opslevel
  reprepro -A amd64 remove $release opslevel
done

for release in ${DEBIAN_RELEASES[@]} ${UBUNTU_RELEASES[@]}; do
  echo "Adding deb package to $release"
  reprepro includedeb $release ../../src/dist/*linux-64bit.deb
  reprepro includedeb $release ../../src/dist/*linux-32bit.deb
done

git add .
git commit -m "Update deb packages"
git push origin main