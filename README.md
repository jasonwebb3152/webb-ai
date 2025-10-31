# WebbAi

## Setup

The setup process is currently compatible with Linux and WSL on Windows. Although, one could fairly easily extend the steps and scripts for MacOS.

1. (Windows only) If on Windows then install WSL2.

2. Install docker onto your machine.

3. Install Kubectl cli onto your machine

4. To setup local development for Webb-AI run the `setup.sh` script in the root folder. This script works for Linux or WSL, installs Kind (Kubernetes in Docker), Golang (but won't upgrade golang to necessary version), and all golang tools/dependencies.
