#!/bin/bash

echo "Compiling Goify..."
go build -o goify ./app

if [ $? -ne 0 ]; then
    echo "Compilation failed. Please make sure Go is installed and set up correctly."
    read -p "Press any key to exit..." -n1 -s
    exit 1
fi

echo "Enter the installation directory (default is /opt/goify):"
read -p "Install Path: " INSTALL_DIR
INSTALL_DIR=${INSTALL_DIR:-/opt/goify}

echo "Creating installation folder..."
mkdir -p "$INSTALL_DIR"

echo "Moving goify to $INSTALL_DIR..."
mv goify "$INSTALL_DIR"

echo "Adding $INSTALL_DIR to PATH..."
echo "export PATH=\$PATH:$INSTALL_DIR" >> ~/.bashrc
source ~/.bashrc

if [ $? -ne 0 ]; then
    echo "Failed to update PATH. You may need to do it manually."
    read -p "Press any key to exit..." -n1 -s
    exit 1
fi

echo "Installation complete. Goify is ready to use from the terminal!"
read -p "Press any key to exit..." -n1 -s