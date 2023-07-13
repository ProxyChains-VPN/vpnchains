#!/bin/env sh

# Variables.
export VPNCHAINS_LIB_NAME=libvpnchains_inject.so
export VPNCHAINS_EXECUTABLE_NAME=vpnchains
export VPNCHAINS_OUTPUT_DIR=build

LIBBSON_INCLUDE_DIR=/usr/include/libbson-1.0

# $1 - exit code
exit_if_unsuccessful() {
    if [ ! $1 -eq 0 ] ; then
        echo "Build unsuccessful! Exiting with code 1..."
        exit 1
    fi
}

check_if_libbson_installed() {
    echo "Checking requirements"

    if [ ! -d "$LIBBSON_INCLUDE_DIR" ]; then
        echo "$LIBBSON_INCLUDE_DIR is not exist; respectively, libbson is not installed globally."
        echo "One is supplied in mongo-c-driver package (package name may be varied)"
        echo "Exiting with code 1..."
        exit 1
    fi

    echo "Requirements are ok!"
}

compile_and_mv_lib() {
    echo "Compiling lib"
    make lib -B

    exit_if_unsuccessful $?

    echo "Moving $VPNCHAINS_LIB_NAME to /usr/lib; need sudo"
    sudo mv $VPNCHAINS_OUTPUT_DIR/$VPNCHAINS_LIB_NAME /usr/lib/$VPNCHAINS_LIB_NAME
}

compile_app() {
    echo "Compiling app"
    make app -B

    exit_if_unsuccessful $?

    echo "Executable compiled and is located at $VPNCHAINS_OUTPUT_DIR/$VPNCHAINS_EXECUTABLE_NAME"
    echo "Done!"
}


# The script itself.

echo "Running vpnchains install script..."

check_if_libbson_installed

if [ -f "/usr/lib/$VPNCHAINS_LIB_NAME" ]; then
    read -r -p "$VPNCHAINS_LIB_NAME already exist! Recompile and replace? [Y/n] " response
    case "$response" in
        [nN][oO])
            echo "Skipping moving $VPNCHAINS_LIB_NAME to /usr/lib; be sure you have a correct library installed"
            ;;
        *)
            compile_and_mv_lib

    esac
else
    compile_and_mv_lib
fi

compile_app

# FOR REMOVAL ON RELEASE
echo "Compiling test"
make test -B
exit_if_unsuccessful $?
echo "Truly done"
