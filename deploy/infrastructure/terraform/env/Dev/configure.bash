#!/bin/bash

set -e

echo "Setting up SSH for Multipass VM"
NAME=lilo-finance-manager
IP=$( multipass list | awk "/${NAME}/"'{print $3}')
SSH_KEY="lilo_rsa"
USER="ubuntu"


ssh_setup() {
    echo "Setting up SSH server for $NAME at IP $IP"
    ssh_generate_key
    ssh_authorize_key
    ssh_check
}

ssh_generate_key() {
    SSH_PUB_KEY=".ssh/${SSH_KEY}.pub"
    if [ -f $SSH_PUB_KEY ]; then
        echo "SSH key already exists"
        return
    fi
    echo "Generating SSH key"
    mkdir -p .ssh
    ssh-keygen -t rsa -b 4096 -C "lxd" -f .ssh/$SSH_KEY -N "" > /dev/null
    echo "SSH key generated"
}

ssh_authorize_key() {
    echo "Authorizing SSH key"
    local SSH_FOLDER="/home/$USER/.ssh"
    local pubkey=$(cat $SSH_PUB_KEY)
    multipass exec $NAME -- sh -c "echo $pubkey >> $SSH_FOLDER/authorized_keys"
}

ssh_check() {
    echo "Checking SSH connection"
    ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i .ssh/$SSH_KEY $USER@$IP exit
    echo "SSH connection established"
}

add_ip_to_files() {
    add_ip_to_inventory
    add_ip_to_istio
}

add_ip_to_inventory() {
    echo "Adding IP to inventory"
    sed -i "/    dev:/,/      ansible_host:/ s/ansible_host: .*/ansible_host: $IP/" ../../../ansible/hosts.yaml
    echo "IP added to inventory"
}

add_ip_to_istio() {
    echo "Adding IP to istio"
    sed -i "/    hosts:/,/      - \"/ s/- \".*/- \"$IP\"/" ../../../../application/istio/ingress.yaml
    echo "IP added to istio"
}

ssh_setup
add_ip_to_files
