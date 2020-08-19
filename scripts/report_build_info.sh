#!/bin/bash

# WARNING: DO NOT EDIT, THIS FILE IS PROBABLY A COPY
#
# The original version of this file is located in the https://github.com/istio/common-files repo.
# If you're looking at this file in a different repo and want to make a change, please go to the
# common-files repo, make the change there and check it in. Then come back to this repo and run
# "make update-common".

# Copyright Istio Authors
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

# adapt to beego

if BUILD_GIT_REVISION=$(git rev-parse HEAD 2> /dev/null); then
  if [[ -n "$(git status --porcelain 2>/dev/null)" ]]; then
    BUILD_GIT_REVISION=${BUILD_GIT_REVISION}"-dirty"
  fi
else
  BUILD_GIT_REVISION=unknown
fi

# Check for local changes
if git diff-index --quiet HEAD --; then
  tree_status="Clean"
else
  tree_status="Modified"
fi

# security wanted VERSION='unknown'
VERSION="${BUILD_GIT_REVISION}"
if [[ -n ${BEEGO_VERSION} ]]; then
  VERSION="${BEEGO_VERSION}"
fi

GIT_DESCRIBE_TAG=$(git describe --tags)

echo "github.com/astaxie/beego.BuildVersion=${VERSION}"
echo "github.com/astaxie/beego.BuildGitRevision=${BUILD_GIT_REVISION}"
echo "github.com/astaxie/beego.BuildStatus=${tree_status}"
echo "github.com/astaxie/beego.BuildTag=${GIT_DESCRIBE_TAG}"