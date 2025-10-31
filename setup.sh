set -e

##########
# Set script's key variables
##########

REQUIRED_GO_VERSION="1.25.3"
GO_AMD_FILE="go${REQUIRED_GO_VERSION}.linux-amd64.tar.gz"
GO_AMD_URL="https://go.dev/dl/${GO_AMD_FILE}"
GO_ARM_FILE="go${REQUIRED_GO_VERSION}.linux-arm64.tar.gz"
GO_ARM_URL="https://go.dev/dl/${GO_ARM_FILE}"
KIND_VERSION="v0.30.0"
KIND_AMD_URL="https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-amd64"
KIND_ARM_URL="https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-arm64"

ARCH=$(uname -m)

##########
# Update tools such as make
##########

sudo apt update
sudo apt install build-essential

##########
# Install golang version, fail if wrong version installed already
##########

if command -v go &> /dev/null; then
  INSTALLED_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
  if [[ "$INSTALLED_VERSION" < "$REQUIRED_GO_VERSION" ]]; then
    echo "ERROR: golang version $INSTALLED_VERSION is too old. Please upgrade to >= $REQUIRED_GO_VERSION"
    exit 1
  else
    echo "Go version $INSTALLED_VERSION is sufficient."
  fi
else
  echo "Go not found. Installing $REQUIRED_GO_VERSION..."
  if [ "$ARCH" = "x86_64" ]; then
    wget -P /tmp "$GO_AMD_URL"
    sudo tar -C /usr/local -xzf "/tmp/${GO_AMD_FILE}"

  elif [ "$ARCH" = "aarch64" ]; then
    wget -P /tmp "$GO_ARM_URL"
    sudo tar -C /usr/local -xzf "/tmp/${GO_ARM_FILE}"

  else
    echo "ERROR: Script does not support your architecture: $ARCH"
    exit 1
  fi
fi

##########
# Install Kind for local development cluster
##########

echo "Installing Kind for local kubernetes development"

sudo rm -rf /usr/local/bin/kind
sudo rm -rf /tmp/kind

# For AMD
if [ "$ARCH" = "x86_64" ]; then
  curl -L "$KIND_AMD_URL" -o /tmp/kind

# For ARM
elif [ "$ARCH" = "aarch64" ]; then
  curl -L "$KIND_ARM_URL" -o /tmp/kind

else
  echo "ERROR: Script does not support your architecture: $ARCH"
    exit 1
fi

chmod +x /tmp/kind
sudo mv /tmp/kind /usr/local/bin/kind

##########
# Install necessary golang binaries and depencies for project
##########

# First add go to path
GO_BIN=$(command -v go)
GOROOT=$(go env GOROOT)
GOPATH=$(go env GOPATH)
export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

echo "Running 'go install tool'"
go install tool
exit 0
