#! /bin/bash -e
function createSudoersFile {
    execPath=$(which cloud-provider-kind)
    whoami=$(whoami)
    sudo rm -rf /private/etc/sudoers.d/cloud-provider-kind
sudo tee /private/etc/sudoers.d/cloud-provider-kind << EOF
$whoami ALL=(ALL) NOPASSWD: $execPath
EOF
}

function installCloudProviderKind {
    tmpDir=$(mktemp -d)
    pushd $tmpDir
        git clone https://github.com/momentumai-team/cloud-provider-kind.git
        pushd cloud-provider-kind
            make build
            chmod +x bin/cloud-provider-kind
            mv bin/cloud-provider-kind $GOPATH/bin
        popd
    popd
}

function installCloudProviderKindTray {
    go build -v -o "$GOPATH/bin/cloud-provider-kind-tray" main.go
}
