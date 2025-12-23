#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

RELEASE_DATE=$(date +%m-%d-%y-%H:%M:%S)
REFERENCE_LINK="krmckone/lk-site@$(git rev-parse --short "$GITHUB_SHA")"

configure_git() {
  git config --global --type bool push.autoSetupRemote true
  git config --global user.name "lk-site GitHub Actions Bot"
  git config --global user.email "20476319+krmckone@users.noreply.github.com"
  git config user.name
  git config user.email
}

create_deployment_branch() {
  local release_hash
  release_hash=$(git rev-parse --short "$GITHUB_SHA")
  local safe_ref
  safe_ref=$(echo "$GITHUB_REF_NAME" | tr '/' '-' | cut -c1-30)

  git checkout -b "lk-site-deploy-$safe_ref-$release_hash"
}

copy_site_files() {
  pwd
  echo "Copying site artifacts to krm-site"
  ls -lsa "$HOME/site/"
  ls ./
  rsync -a --delete "$HOME/site/" ./
  ls "$HOME/site/"
  ls ./
}

commit_changes() {
  if [[ -z $(git status --porcelain) ]]; then
    echo "No changes to the static site. Exiting."
    exit 0
  fi

  git add .
  git commit -m "New Release $REFERENCE_LINK $RELEASE_DATE"
  git push origin HEAD
}

create_and_merge_pr() {
  local body
  body=$(cat <<EOF
* This pull request was automatically created by $GITHUB_SERVER_URL/krmckone/lk-site/actions/runs/$GITHUB_RUN_ID
* This release is based on $REFERENCE_LINK
EOF
)

  gh pr create \
    --title "Automatic Release $RELEASE_DATE" \
    --body "$body" \
    || echo "PR already exists, continuing"

  local retries=0
  local max_retries=3

  until gh pr merge --auto --merge; do
    (( retries++ ))
    if (( retries >= max_retries )); then
      echo "Failed to enable auto-merge after $retries attempts"
      exit 1
    fi
    sleep 3
  done
}

git --version
configure_git
create_deployment_branch
copy_site_files
commit_changes
create_and_merge_pr
