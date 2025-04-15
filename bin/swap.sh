#!/bin/bash

# Swap Setup Script
# Usage: sudo ./setup_swap.sh [swap_size_in_GB(default:4)] [add_to_fstab(0|1)(default:1)]

# Check if running as root
if [ "$(id -u)" -ne 0 ]; then
    echo "This script requires root privileges. Please run with sudo."
    exit 1
fi

# Set swap size, default 4GB
SWAP_SIZE=${1:-4}

# Set whether to add to fstab, default is yes (1)
ADD_TO_FSTAB=${2:-1}

echo "=== Starting setup of ${SWAP_SIZE}GB swap space ==="

# Check if swap already exists
echo "Checking current swap status..."
CURRENT_SWAP=$(swapon --show)
if [ -n "$CURRENT_SWAP" ]; then
    echo "Existing swap detected:"
    echo "$CURRENT_SWAP"
    read -p "Continue setting up new swap? (y/n): " CONTINUE
    if [[ "$CONTINUE" != "y" && "$CONTINUE" != "Y" ]]; then
        echo "Operation cancelled"
        exit 0
    fi
    echo "Continuing with new swap setup..."
else
    echo "No swap detected, proceeding with setup..."
fi

# Set swap file path
SWAP_FILE="/swapfile"

# Check if swap file already exists
if [ -f "$SWAP_FILE" ]; then
    echo "Existing swap file found, turning off and removing..."
    swapoff "$SWAP_FILE" 2>/dev/null
    rm -f "$SWAP_FILE"
fi

# Create swap file
echo "Creating ${SWAP_SIZE}GB swap file..."
if command -v fallocate &>/dev/null; then
    fallocate -l "${SWAP_SIZE}G" "$SWAP_FILE"
    if [ $? -ne 0 ]; then
        echo "fallocate command failed, attempting to use dd..."
        dd if=/dev/zero of="$SWAP_FILE" bs=1G count="$SWAP_SIZE"
    fi
else
    echo "System doesn't support fallocate, using dd..."
    dd if=/dev/zero of="$SWAP_FILE" bs=1G count="$SWAP_SIZE"
fi

# Set permissions
echo "Setting file permissions..."
chmod 600 "$SWAP_FILE"

# Set up swap
echo "Setting up file as swap space..."
mkswap "$SWAP_FILE"

# Enable swap
echo "Enabling swap..."
swapon "$SWAP_FILE"

# Check if swap is enabled
echo "Verifying swap status..."
SWAP_ENABLED=$(swapon --show | grep "$SWAP_FILE")
if [ -n "$SWAP_ENABLED" ]; then
    echo "Swap successfully enabled:"
    swapon --show
    free -h | grep -E 'Mem|Swap'
else
    echo "Swap enabling failed, please check system logs"
    exit 1
fi

# Add to fstab if requested
if [ "$ADD_TO_FSTAB" -eq 1 ]; then
    echo "Adding to /etc/fstab to ensure automatic activation at boot..."
    # Check if entry already exists
    if ! grep -q "$SWAP_FILE" /etc/fstab; then
        echo "$SWAP_FILE none swap sw 0 0" >>/etc/fstab
        echo "Successfully added to fstab"
    else
        echo "Swap entry already exists in fstab, no addition needed"
    fi
else
    echo "Skipping fstab entry as requested. Swap will not automatically activate on reboot."
fi

# Set swappiness
SWAPPINESS=10
echo "Setting vm.swappiness=${SWAPPINESS}..."
sysctl vm.swappiness="$SWAPPINESS"
if ! grep -q "vm.swappiness" /etc/sysctl.conf; then
    echo "vm.swappiness=$SWAPPINESS" >>/etc/sysctl.conf
    echo "Successfully added swappiness to sysctl.conf"
else
    # Update existing swappiness setting
    sed -i "s/vm.swappiness=.*/vm.swappiness=$SWAPPINESS/" /etc/sysctl.conf
    echo "Updated swappiness setting in sysctl.conf"
fi

echo "=== ${SWAP_SIZE}GB swap setup completed ==="
if [ "$ADD_TO_FSTAB" -eq 1 ]; then
    echo "Swap will automatically activate on system reboot"
else
    echo "Note: Swap will NOT automatically activate on system reboot"
    echo "To manually enable after reboot: sudo swapon ${SWAP_FILE}"
fi
echo "You can check memory and swap usage with 'free -h' command"
