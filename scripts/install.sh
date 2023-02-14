#!/usr/bin/env bash

# This script installs the Stevedore tool into your system
# by downloading the a binary distribution and 
# running the install script https://github.com/gostevedore/stevedore/scripts/install.sh

set -eo pipefail

CURL_CMD="curl"
WGET_CMD="wget"
GITHUB_API_URL="https://api.github.com"
SOURCE_RELEASE_DEST_BASE_PATH=/opt/stevedore
BINARY_PATH=/usr/local/bin/stevedore

fail() {
    echo "$0": "$@" >&2
    exit 1
}

require_command() {
    command -v "$1" > /dev/null 2>&1 || fail "To install Stevedore '$1' is required"
}

must() {
    cmd=("$@")
    res=$("${cmd[@]}")
    if [[ "$?" == "1" ]]; then
        fail "${res}"
    fi

    if [ -n "${res}" ]; then
        echo "${res}"
    fi
}

create_dir() {
    require_command mkdir

    if [ "$#" -ne "1" ]; then
        fail "Invalid parameters: create_folder <dir>"
    fi

    local dir="${1}"
    
    must mkdir -p "${dir}"
}

create_tmp_dir() {
    require_command mktemp

    dir=$(must mktemp -d --suffix _stevedore)
    
    echo "${dir}"
}

cleanup() {
    if [ "$#" -ne "1" ]; then
        fail "Invalid parameters: cleanup <dir>"
    fi

    local dir="${1}"

    rm -rf "${1}"
}

curl_get_http_cmd() {
    if [ "$#" -ne "1" ]; then
        fail "Invalid parameters: curl_get_http <url>"
    fi

    local url=${1}

    require_command "${CURL_CMD}"

    echo "${CURL_CMD} -sL ${url}"
    return
}

curl_download_file_cmd() {

    if [ "$#" -ne "2" ]; then
        fail "Invalid parameters: curl_download_file <url> <dest>"
    fi

    local url=${1}
    local dest=${2}

    require_command "${CURL_CMD}"
    require_command dirname

    create_dir "$(must dirname "${dest}")"

    echo "${CURL_CMD} -sOL ${url} --output-dir ${dest}"
    return
}

wget_get_http_cmd(){
    if [ "$#" -ne "1" ]; then
        fail "Invalid parameters: wget_get_http <url>"
    fi

    local url=${1}

    require_command "${WGET_CMD}"

    echo "${WGET_CMD} -q --output-document - ${url}"
    return
}

wget_download_file_cmd() {
    if [ "$#" -ne "2" ]; then
        fail "Invalid parameters: curl_download_file <url> <dir>"
    fi

    local url=${1}
    local dir=${2}

    require_command "${CURL_CMD}"
    require_command dirname

    create_dir "${dir}"

    echo "${WGET_CMD} -q --show-progress --directory-prefix ${dir} ${url}"
    return
}

http_client_factory() {
    if command -v ${CURL_CMD} > /dev/null 2>&1; then
        echo "curl_get_http_cmd" "curl_download_file_cmd"
        return
    elif command -v ${WGET_CMD} > /dev/null 2>&1; then
        echo "wget_get_http_cmd" "wget_download_file_cmd"
        return
    else
        err "You require either curl or wget to install Stevedore"
    fi 
}

get_kernel_name() {
    require_command uname
    uname -s
}

get_arch() {
    require_command uname
    uname -m
}

fetch_release() {
    if [ "$#" -ne "1" ]; then
        fail "Invalid parameters: fetch_release <cmd>"
    fi

    local cmd=("$@")

    ${cmd} | grep -Po '"tag_name": "\K.*?(?=")'
    return
}

generate_artefact_name() {
    if [ "$#" -ne "1" ]; then
        fail "Invalid parameters: generate_artefact_name <release>"
    fi

    local arch
    local kernel
    local release="${1}"

    arch="$(must get_arch)"
    kernel="$(must get_kernel_name)"

    echo "stevedore_$(echo "${release}" | sed 's/^v//' )_${kernel}-${arch}.tar.gz"
    return
}

fetch_artefact_url() {
    if [ "$#" -ne "1" ]; then
        fail "Invalid parameters: fetch_artefact_url <release>"
    fi

    local arch
    local kernel
    local release="${1}"


    arch="$(must get_arch)"
    kernel="$(must get_kernel_name)"

    # https://github.com/gostevedore/stevedore/releases/download/v0.10.3/stevedore_0.10.3_Linux-x86_64.tar.gz
    artefact=$(must generate_artefact_name "${release}")
    local url="https://github.com/gostevedore/stevedore/releases/download/${release}/${artefact}"

    echo "${url}"
    return
}

extract_artefact() {
    require_command tar

    if [ "$#" -ne "2" ]; then
        fail "Invalid parameters: extract_artefact <source> <dest>"
    fi

    local source="${1}"
    local dest="${2}"

    if [ ! -f "${source}" ]; then
        echo "Source file ${source} does not exist"
        return 1
    fi

    create_dir "${dest}"

    must tar -zxf "${source}" -C "${dest}"
    return
}

install_artefact() {
    require_command ln
    require_command dirname

    if [ "$#" -ne "2" ]; then
        fail "Invalid parameters: install_artefact <source> <dest>"
    fi

    local source="${1}"
    local dest="${2}"

    create_dir "$(must dirname "${dest}")"

    must ln -sf "${source}" "${dest}"
    return
}

read -r get_http_cmd download_file_cmd < <(must http_client_factory)

download_dir=$(must create_tmp_dir)
trap 'cleanup "${download_dir}"' EXIT

release=$(fetch_release "$($get_http_cmd ${GITHUB_API_URL}/repos/gostevedore/stevedore/releases/latest)")
if [ -z "${release}" ]; then
    fail "Release can not be achieved"
fi

artefact=$(must generate_artefact_name "${release}")
if [ -z "${artefact}" ]; then
    fail "Artefact name can not be achieved"
fi

source_release_dest_path="${SOURCE_RELEASE_DEST_BASE_PATH}/${release}"
if [ -d "${source_release_dest_path}" ]; then
    fail "${source_release_dest_path} already exists. Remove the directory and reinstall Stevedore"
fi

echo " Installing Stevedore ${release}..."
echo "  artefact: ${artefact}"

eval "$($download_file_cmd "$(fetch_artefact_url "${release}")" "${download_dir}")"

extract_artefact "${download_dir}/${artefact}" "${source_release_dest_path}"
install_artefact "${source_release_dest_path}/$(basename "${BINARY_PATH}")" "${BINARY_PATH}"

must "${BINARY_PATH}" version
echo