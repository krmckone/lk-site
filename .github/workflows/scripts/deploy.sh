#!/bin/bash

RELEASE_DATE=$(date +%m-%d-%y-%H:%M:%S)
REFERENCE_LINK="krmckone/lk-site@$(git rev-parse --short $GITHUB_SHA)"

configure_git() {
  git config --global --type bool push.autoSetupRemote true
  git config --global user.name "Kaleb's GitHub Actions Bot from lk-site"
  git config --global user.email "20476319+krmckone@users.noreply.github.com" # TODO: Parameterize
}

create_deployment_branch() {
  local release_hash=$(git rev-parse --short $GITHUB_SHA)
  git checkout -b "lk-site-deploy-$GITHUB_REF_NAME-$release_hash"
}

copy_site_files() {
  cp -r ~/site/* ./
}

commit_changes() {
  git status --porcelain
  local site_modified=$(git status --porcelain)
  if [[ -z $site_modified ]]; then
    echo "No changes to the static site. Exiting without asset deployment."
    exit 0
  else
    echo "Deploying changes to the static site."
    git add .
    git commit -m "New Release $REFERENCE_LINK $RELEASE_DATE"
    git push origin
  fi
}

create_and_merge_pr() {
  echo -e "* This pull request was automatically created and merged by $GITHUB_SERVER_URL/krmckone/lk-site/actions/runs/$GITHUB_RUN_ID.\n* This release is based on $REFERENCE_LINK" > body
  gh pr create --title "Automatic Release $RELEASE_DATE" --body-file body

  local max_retries=2
  local retries=0
  until gh pr merge --auto --merge; do
    sleep 3
    [[ $retries -eq $max_retries ]] && echo 'Unable to merge PR in krm-site!' && exit 1
    ((retries++))
    echo "GH CLI returned non-zero after $retries attempts, going to retry auto merge"
  done
  echo "PR merged successfully in krm-site"
}

# Main script execution
git --version
configure_git
create_deployment_branch
copy_site_files
commit_changes
create_and_merge_pr