#!/bin/env sh

export VPNCHAINS_LIB_NAME=vpnchains_inject.so
export VPNCHAINS_EXECUTABLE_NAME=vpnchains
export VPNCHAINS_OUTPUT_DIR=build

LIBBSON_INCLUDE_DIR=/usr/include/libbson-1.0

echo "Running vpnchains install script..."
echo "Checking requirements"

if [ ! -d "$LIBBSON_INCLUDE_DIR" ]; then
    echo "$LIBBSON_INCLUDE_DIR is not exist; respectively, libbson is not installed globally."
    echo "One is supplied in mongo-c-driver package (package name may be varied)"
    exit 1
fi

echo "Compiling lib"
make lib -B

if [ -f "/usr/lib/$VPNCHAINS_LIB_NAME" ]; then
    read -r -p "$VPNCHAINS_LIB_NAME already exist! Replace? [y/N] " response
    case "$response" in
        [yY][eE][sS]|[yY])
            echo "Moving $VPNCHAINS_LIB_NAME to /usr/lib; need sudo"
            sudo mv $VPNCHAINS_OUTPUT_DIR/$VPNCHAINS_LIB_NAME /usr/lib/$VPNCHAINS_LIB_NAME
            ;;
        *)
            echo "Skipping moving $VPNCHAINS_LIB_NAME to /usr/lib; be sure you have a correct library installed"
            ;;
    esac
fi

echo "Compiling app"
make app -B
echo "Executable compiled and is located at $VPNCHAINS_OUTPUT_DIR/$VPNCHAINS_EXECUTABLE_NAME"
echo "Done!"
