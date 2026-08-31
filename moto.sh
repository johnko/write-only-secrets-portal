#!/usr/bin/env bash
set -exuo pipefail

## a script to start a simulation of AWS secretsmanager service

# renovate: datasource=python-version depName=python packageName=python
PYTHON_VERSION=3.14.7

# renovate: datasource=pypi depName=pipx packageName=pipx
PIPX_VERSION=1.16.6

# renovate: datasource=pypi depName=poetry packageName=poetry
POETRY_VERSION=2.4.1

WORKDIR=$(dirname "$0")
pushd "$WORKDIR"
cd ./test/aws/
mise exec python@$PYTHON_VERSION -- pip install --user pipx==$PIPX_VERSION
mise exec python@$PYTHON_VERSION -- "$HOME/.local/bin/pipx" install poetry==$POETRY_VERSION
mise exec python@$PYTHON_VERSION -- "$HOME/.local/bin/poetry" --version
mise exec python@$PYTHON_VERSION -- "$HOME/.local/bin/poetry" install --no-root
mise exec python@$PYTHON_VERSION -- "$HOME/.local/bin/poetry" run moto_server -p3000
popd
