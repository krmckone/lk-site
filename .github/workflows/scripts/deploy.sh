#!/usr/bin/env bash
set -Eeuo pipefail
IFS=$'\n\t'

RELEASE_DATE=$(date +%m-%d-%y-%H:%M:%S)
SHORT_SHA=$(git rev-parse --short "$GITHUB_SHA")
DEPLOY_BRANCH="deploy-$GITHUB_REF_NAME-$SHORT_SHA"
REFERENCE_LINK="krmckone/lk-site@$SHORT_SHA"

SITE_DIR="$GITHUB_WORKSPACE/site"
TARGET_REPO_DIR="$GITHUB_WORKSPACE/krm-site"

cd "$TARGET_REPO_DIR"

git config user.name "lk-site GitHub Actions Bot"
git config user.email "20476319+krmckone@users.noreply.github.com"

git checkout -b "$DEPLOY_BRANCH"

if [[ ! -d "$SITE_DIR" ]] || [[ -z "$(ls -A "$SITE_DIR")" ]]; then
  echo "Error: site build output is empty"
  exit 1
fi


rsync -a --delete \
  --exclude='.git/' \
  --exclude='.github/' \
  --exclude='.gitignore' \
  --exclude='CNAME' \
  --exclude='README.md' \
  --exclude='game_of_life/' \
  "$GITHUB_WORKSPACE/site/" \
  "$GITHUB_WORKSPACE/krm-site/"

if [[ -z $(git status --porcelain) ]]; then
  echo "No changes to the static site. Exiting."
  exit 0
fi

git add .
git commit -m "New Release $REFERENCE_LINK $RELEASE_DATE"
git push origin "$DEPLOY_BRANCH" --set-upstream

if ! gh pr view "$DEPLOY_BRANCH" &>/dev/null; then
    PR_BODY="* This pull request was automatically created by $GITHUB_SERVER_URL/krmckone/lk-site/actions/runs/$GITHUB_RUN_ID
    * Release based on $REFERENCE_LINK"
    gh pr create \
        --title "Automatic Release $RELEASE_DATE" \
        --body "$PR_BODY" \
        --base main \
        --head "$DEPLOY_BRANCH"
fi

PR_NUMBER=$(gh pr list --head "$DEPLOY_BRANCH" --base main --state open --json number --jq '.[0].number')

if [[ -n "$PR_NUMBER" ]]; then
    gh pr merge "$PR_NUMBER" --merge --delete-branch || echo "PR already merged or merge failed"
fi