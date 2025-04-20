#!/bin/bash

echo "Downloading playwright with depenedencies:"
go get -u github.com/playwright-community/playwright-go

echo "Compiling Goify..."
go build -o goify ./app

if [ $? -ne 0 ]; then
    echo "Compilation failed. Please make sure Go is installed and set up correctly."
    read -p "Press any key to exit..." -n1 -s
    exit 1
fi

echo "Enter the installation directory (default is ~/Goify):"
read -p "Install Path: " INSTALL_DIR
INSTALL_DIR=${INSTALL_DIR:-~/Goify}

echo "Creating installation folder..."
mkdir -p "$INSTALL_DIR"
if [ $? -ne 0 ]; then
    echo "Failed to create directory '$INSTALL_DIR'. Please check your permissions."
    read -p "Press any key to exit..." -n1 -s
    exit 1
fi

echo "Moving goify to $INSTALL_DIR..."
mv goify "$INSTALL_DIR"
if [ $? -ne 0 ]; then
    echo "Failed to move goify to '$INSTALL_DIR'. Please check your permissions."
    read -p "Press any key to exit..." -n1 -s
    exit 1
fi

echo "Adding $INSTALL_DIR to PATH..."
echo -e "\nexport PATH=\$PATH:$INSTALL_DIR" >> ~/.bashrc

if [ $? -ne 0 ]; then
    echo "Failed to update PATH. You may need to do it manually."
    read -p "Press any key to exit..." -n1 -s
    exit 1
fi

echo "Installation complete. Goify is ready to use from the terminal!"
echo "If the goify command isn't yet recognized please run 'source ~/.bashrc'"
read -p "Press any key to exit..." -n1 -s
source ~/.bashrc
