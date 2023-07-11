#!/bin/env sh

# $1 - exit code
exit_if_unsuccessful() {
    if [ ! $1 -eq 0 ] ; then
        echo "Build unsuccessful! Exiting with code 1..."
        exit 1
    fi
}

export VPNCHAINS_LIB_NAME=libvpnchains_inject.so
export VPNCHAINS_EXECUTABLE_NAME=vpnchains
export VPNCHAINS_OUTPUT_DIR=build

LIBBSON_INCLUDE_DIR=/usr/include/libbson-1.0
LIBBSON_LIB_PATH=/usr/lib/libbson-1.0.so

echo "Running vpnchains install script..."
echo "Checking requirements"

if [ ! -d "$LIBBSON_INCLUDE_DIR" ]; then
    echo "$LIBBSON_INCLUDE_DIR is not exist; respectively, libbson is not installed globally."
    echo "One is supplied in mongo-c-driver package (package name may be varied)"
    echo "Exiting with code 1..."
    exit 1
fi

if [ ! -f "$LIBBSON_LIB_PATH" ]; then
    echo "$LIBBSON_LIB_PATH is not exist; respectively, libbson is not installed globally."
    echo "One is supplied in mongo-c-driver package (package name may be varied)"
    echo "Exiting with code 1..."
    exit 1
fi

echo "Requirements are ok!"

if [ ! -f "/usr/lib/$VPNCHAINS_LIB_NAME" ]; then
    read -r -p "$VPNCHAINS_LIB_NAME already exist! Recompile and replace? [Y/n] " response
    case "$response" in
        [nN][oO])
            echo "Skipping moving $VPNCHAINS_LIB_NAME to /usr/lib; be sure you have a correct library installed"
            ;;
        *)
            echo "Compiling lib"
            make lib -B

            exit_if_unsuccessful $?

            echo "Moving $VPNCHAINS_LIB_NAME to /usr/lib; need sudo"
            sudo cp $VPNCHAINS_OUTPUT_DIR/$VPNCHAINS_LIB_NAME /usr/lib/$VPNCHAINS_LIB_NAME
            ;;

    esac
else
    echo "Moving $VPNCHAINS_LIB_NAME to /usr/lib; need sudo"
    sudo mv $VPNCHAINS_OUTPUT_DIR/$VPNCHAINS_LIB_NAME /usr/lib/$VPNCHAINS_LIB_NAME
fi

echo "Compiling app"
make app -B

exit_if_unsuccessful $?

echo "Executable compiled and is located at $VPNCHAINS_OUTPUT_DIR/$VPNCHAINS_EXECUTABLE_NAME"
echo "Done!"

# TMP
echo "Compiling test"
make test -B
exit_if_unsuccessful $?
echo "Truly done"
