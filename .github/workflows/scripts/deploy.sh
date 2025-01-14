#!/bin/bash

# Function to configure Git
configure_git() {
  git config --global user.name "Kaleb's GitHub Actions Bot from lk-site"
  git config --global user.email "20476319+krmckone@users.noreply.github.com" # TODO: Parameterize
}

# Function to create a new branch for deployment
create_deployment_branch() {
  local release_hash=$(git rev-parse --short $GITHUB_SHA)
  git checkout -b "lk-site-deploy-${{ github.ref_name }}-$release_hash"
}

# Function to copy site files
copy_site_files() {
  cp -r ~/site/* ./
}

# Function to check for changes and commit them
commit_changes() {
  local site_modified=$(git status --porcelain)
  if [[ -z $site_modified ]]; then
    echo "No changes to the static site. Exiting without asset deployment."
    exit 0
  else
    echo "Deploying changes to the static site."
    git add .
    git commit -m "New Release $reference_link $release_date"
    git push origin HEAD
  fi
}

# Function to create and merge a pull request
create_and_merge_pr() {
  local release_date=$(date +%m-%d-%y-%H:%M:%S)
  local reference_link="krmckone/lk-site@$(git rev-parse --short $GITHUB_SHA)"

  echo -e "* This pull request was automatically created and merged by ${{ github.server_url }}/krmckone/lk-site/actions/runs/${{ github.run_id }}.\n* This release is based on $reference_link" > body
  gh pr create --title "Automatic Release $release_date" --body-file body

  local max_retries=4
  local retries=0
  until gh pr merge --auto --merge; do
    sleep 5
    [[ $retries -eq $max_retries ]] && echo 'Unable to merge PR in krm-site!' && exit 1
    ((retries++))
    echo "GH CLI returned non-zero after $retries attempts, going to retry auto merge"
  done
  echo "PR merged successfully in krm-site"
}

# Main script execution
git --version
git pull origin main
configure_git
create_deployment_branch
copy_site_files
commit_changes
create_and_merge_pr