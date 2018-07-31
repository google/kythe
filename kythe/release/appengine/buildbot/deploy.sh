#!/bin/bash -e
#
# Copyright 2018 The Kythe Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Script to deploy the kythe-repo/buildbot AppEngine module.

cd "$(dirname "$0")"

# Cleanup secrets on exit
trap "rm -rf '$PWD/secrets'* '$PWD/master/master.cfg'" EXIT ERR INT

gsutil cp gs://kythe-buildbot/secrets.tar.enc secrets.tar.enc
gcloud kms decrypt --location=global --keyring=Buildbot --key=secrets \
  --plaintext-file=secrets.tar --ciphertext-file=secrets.tar.enc
tar xf secrets.tar

SECRETS_SED=
for secret in secrets/*; do
  name="$(basename "$secret")"
  SECRETS_SED="$SECRETS_SED
s/@secret{$name}/$(cat "$secret" | tr -d '\n' | sed 's#/#\\/#g')/g"
done

sed "$SECRETS_SED" master/master.cfg.template >master/master.cfg
buildbot checkconfig master

VERSION=v1
docker build -t gcr.io/kythe_repo/buildbot.$VERSION .
docker push gcr.io/kythe_repo/buildbot.$VERSION
gcloud app deploy --image-url=gcr.io/kythe_repo/buildbot.$VERSION --stop-previous-version --version $VERSION
