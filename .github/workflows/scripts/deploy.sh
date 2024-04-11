releaseDate=$(date +%m-%d-%y-%H:%M:%S)
lkSiteReleaseHash=$(git rev-parse --short $GITHUB_SHA)
referenceLink="krmckone/lk-site@$lkSiteReleaseHash"
git --version
git config --global user.name "Kaleb's GitHub Actions Bot from lk-site"
git config --global user.email "20476319+krmckone@users.noreply.github.com"
git checkout -b "lk-site-deploy-${{ github.ref_name }}-$lkSiteReleaseHash"
ls ~/site
cp -r ~/site/* ./
siteModified=$(git status --porcelain)
echo $siteModified
if [[ -z $siteModified ]]; then
  echo "No changes to the static site. Exiting without asset deployment."
  exit 0
else
  echo "Deploying changes to the static site."
fi
git add .
git commit -m "New Release $referenceLink $releaseDate"
git push origin HEAD
gh pr create --title "Automatic PR $releaseDate" --body "* This pull request was automatically created from ${{ github.server_url }}/krmckone/lk-site/actions/runs/${{ github.run_id }}. * This release is based on $referenceLink"
gh pr merge --auto --merge