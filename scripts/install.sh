#!/bin/sh

# This script installs the Stevedore tool into your system
# by downloading the a binary distribution and 
# running the install script https://github.com/${OWNER}/${REPO}/scripts/install.sh

set -e

GITHUB_API_URL="https://api.github.com"
OWNER=gostevedore
REPO=stevedore
SOURCE_RELEASE_DEST_BASE_PATH=/opt/stevedore
BINARY_PATH=/usr/local/bin/stevedore
FORCE="0"
INSTALL_SCRIPT="install.sh"

alternatives() {
    echo "
If you have some problems installing Stevedore using the install script, look for installation alternatives in the documentation.
https://gostevedore.github.io/docs/getting-started/install/
    "
}

fail() {
    echo "$0": "$@" >&2
    alternatives >&2
    exit 1
}

cleanup() {
    if [ "$#" -ne "1" ]; then
        fail "Invalid parameters: cleanup <dir>"
    fi

    dir="${1}"
    shift

    echo " Cleanup directory ${dir}..."

    if ! rm -rf "${dir}";
    then
        fail "cleanup: Error cleaning up the directory ${dir}."
    fi
}

require_command() {
    command -v "$1" > /dev/null 2>&1 || fail "To install Stevedore '$1' is required"
}

usage() {
    cmd="${1}"
    cat << EOF

 Usage: ${cmd} [-f] [-h|?] [-b <binary-path>] [-d <releases-directory>] [-v <version>]

  -b: Specifies the path to the Stevedore binary file [default: ${BINARY_PATH}]
  -d: Specifies the directory to install the Stevedore release packages [default: ${SOURCE_RELEASE_DEST_BASE_PATH}]
  -f: Forces the installation of the binary
  -h: Shows the usage command
  -v: Defines the Stevedore version to install [default: the latest]
  -x: Sets the execution mode to debug

EOF
    exit 2
}

parse_args() {

  while getopts "b:fh?d:xv:" arg; do
    case "$arg" in
      b) BINARY_PATH="$OPTARG" ;;
      d) SOURCE_RELEASE_DEST_BASE_PATH="$OPTARG" ;;
      f) FORCE="1" ;;
      h | \?) usage "$0" ;;
      v) VERSION="$OPTARG" ;;
      x) set -x ;;
    esac
  done
  shift $((OPTIND - 1))
}

create_dir() {
    if [ "$#" -ne "1" ]; then
        fail "create_dir: Invalid parameters: create_folder <dir>."
    fi

    dir="${1}"
    shift
    
    echo " Ensuring directory ${dir} exist..."
    if ! mkdir -p "${dir}";
    then
        fail "create_dir: Error creating ${dir}."
    fi
}

create_tmp_dir() {
    dir=$(mktemp -d)
    if [ ! -d "${dir}" ]; then
        fail "create_tmp_dir: Error creating ${dir}."
    fi
    echo "${dir}"
}

get_os() {
    if ! uname -s;
    then
        fail "get_os: Error getting the kernel name."
    fi   
}

check_os() {
    if [ "$#" -ne "1" ]; then
        fail "check_os: Invalid parameters: check_os <kernel_name>."
    fi

    os="${1}"
    shift

    case "${os}" in
        Darwin) return 0 ;;
        Linux) return 0 ;;
        *) fail "check_arch: ${arch} is not supported by the installation script."
    esac
}

get_arch() {
    arch=$(uname -m)
    if [ "$?" -ne "0" ]; then
        fail "get_arch: Error getting the architecture."
    fi

    case $arch in
        x86_64) arch="x86_64" ;;
        x86) arch="386" ;;
        i686) arch="386" ;;
        i386) arch="386" ;;
        aarch64) arch="arm64" ;;
        armv5*) arch="arm" ;;
        armv6*) arch="arm" ;;
        armv7*) arch="arm" ;;
    esac
    echo "${arch}"

    return 0
}

check_arch() {
    if [ "$#" -ne "1" ]; then
        fail "check_arch: Invalid parameters: check_arch <arch>."
    fi

    arch="${1}"
    shift

    case $arch in
        x86_64) return 0 ;;
        x86) return 0 ;;
        i686) return 0 ;;
        i386) return 0 ;;
        aarch64) return 0 ;;
        armv5*) return 0 ;;
        armv6*) return 0 ;;
        armv7*) return 0 ;;
        *) fail "check_arch: ${arch} is not supported by the installation script."
    esac
}

fetch_release() {
    url="${GITHUB_API_URL}/repos/${OWNER}/${REPO}/releases/latest"

    # if ! get_http_cmd "${url}" | grep -Po '"tag_name": "\K.*?(?=")';
    if ! get_http_cmd "${url}" | grep '"tag_name":' | cut -d '"' -f 4;
    then
        fail "fetch_release: Error fetching the release tag name."
    fi
    return 0
}

generate_artefact_name() {
    if [ "$#" -ne "3" ]; then
        fail "generate_artefact_name: Invalid parameters: generate_artefact_name <release> <arch> <kernel>."
    fi

    release="${1}"
    shift
    arch="${1}"
    shift
    kernel="${1}"
    shift

    echo "stevedore_$(echo "${release}" | sed 's/^v//' )_${kernel}_${arch}.tar.gz"
    return 0
}

fetch_artefact_url() {
    if [ "$#" -ne "2" ]; then
        fail "fetch_artefact_url: Invalid parameters: fetch_artefact_url <release> <artefact>."
    fi

    release="${1}"
    shift
    artefact="${1}"
    shift

    url="https://github.com/${OWNER}/${REPO}/releases/download/${release}/${artefact}"

    echo "${url}"
    return
}

get_http_cmd() {
    if [ "$#" -ne "1" ]; then
        fail "get_http_cmd: Invalid parameters: get_http_cmd <source-url>."
    fi

    url="${1}"
    shift

    if ! curl  -w '%{http_code}' -sqL "${url}";
    then
        fail "get_http_cmd: Received HTTP status $code from $url."
    fi

    return 0
}

download_file_cmd() {
    if [ "$#" -ne "2" ]; then
        fail "download_file_cmd: Invalid parameters: download_file_cmd <source-url> <dest>."
    fi

    url="${1}"
    shift
    dest="${1}"
    shift

    dest_dir=$(dirname "${dest}")
    if [ ! -d "${dest_dir}" ]; then
        fail "download_file_cmd: ${dest_dir} does not exist."
    fi

    code=$(curl -w '%{http_code}' -fsSL "${url}" --output "${dest}")
    if [ "$code" != "200" ]; then
        fail "download_file_cmd: received HTTP status $code."
    fi

    return 0
}

extract_artefact() {
    if [ "$#" -ne "2" ]; then
        fail "extract_artefact: Invalid parameters: extract_artefact <source> <dest>."
    fi

    source="${1}"
    shift
    dest="${1}"
    shift

    if [ ! -f "${source}" ]; then
        fail "extract_artefact: Source file ${source} does not exist."
    fi

    if [ ! -d "${dest}" ]; then
        fail "extract_artefact: ${dest} does not exist."
    fi

    echo " Extracting artefact to ${dest}... " 
    if ! tar -zxf "${source}" -C "${dest}";
    then
        fail "extract_artefact: Error extracting artefact from ${source} to ${dest}."
    fi

    return 0
}

install_artefact() {
    if [ "$#" -ne "2" ]; then
        fail "Invalid parameters: install_artefact <source> <dest>"
    fi

    source="${1}"
    shift
    dest="${1}"
    shift

    if [ ! -d "${dest_dir}" ]; then
        fail "install_artefact: ${dest_dir} does not exist."
    fi

    if ! ln -sf "${source}" "${dest}";
    then
        fail "install_artefact: Error creating the symbolic link from ${source} to ${dest}."
    fi

    return 0
}
    
require_command curl
require_command dirname
require_command ln
require_command mkdir
require_command mktemp
require_command rm
require_command tar
require_command uname

parse_args "$@"

download_dir=$(create_tmp_dir)
trap 'cleanup "${download_dir}"' EXIT

arch="$(get_arch)"
if [ -z "${arch}" ]; then
    fail "System's architecture can not be achieved"
fi
check_arch "${arch}"

os="$(get_os)"
if [ -z "${os}" ]; then
    fail "Kernel name can not be achieved"
fi
check_os "${os}"

if [ -z "${VERSION}" ]; then
    release=$(fetch_release)
    if [ -z "${release}" ]; then
        fail "Release can not be achieved"
    fi
else
    release=${VERSION}
fi

artefact=$(generate_artefact_name "${release}" "${arch}" "${os}")
if [ -z "${artefact}" ]; then
    fail "Artefact name can not be achieved"
fi

artefact_url=$(fetch_artefact_url "${release}" "${artefact}")
if [ -z "${artefact_url}" ]; then
    fail "Artefact URL name can not be achieved"
fi

source_release_dest_path="${SOURCE_RELEASE_DEST_BASE_PATH}/${release}"
if [ -d "${source_release_dest_path}" ] && [ ${FORCE} -eq "0" ] ; then
    fail "Error installing Stevedore.\n ${source_release_dest_path} already exists.\n\n Remove the directory and reinstall Stevedore. You can also force the installation by using the flag '-f'"
fi

create_dir "${source_release_dest_path}"
create_dir "$(dirname "${BINARY_PATH}")"

echo " Installing Stevedore ${release} using the artefact ${artefact}..."
echo " Downloading artefact from ${artefact_url}"

local_file="${download_dir}/${artefact}"
download_file_cmd "${artefact_url}" "${local_file}"
extract_artefact "${local_file}" "${source_release_dest_path}"
install_artefact "${source_release_dest_path}/$(basename "${BINARY_PATH}")" "${BINARY_PATH}"

echo " Installation completed successfully!"

echo
"${BINARY_PATH}" version
echo
