releaseDate=$(date +%m-%d-%y-%H:%M:%S)
lkSiteReleaseHash=$(git rev-parse --short $GITHUB_SHA)
referenceLink="krmckone/lk-site@$lkSiteReleaseHash"
git --version
git config --global user.name "Kaleb's GitHub Actions Bot from lk-site"
git config --global user.email "20476319+krmckone@users.noreply.github.com" # TODO: Parameterize
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
echo -e "* This pull request was automatically created and merged by ${{ github.server_url }}/krmckone/lk-site/actions/runs/${{ github.run_id }}.\n* This release is based on $referenceLink" > body
gh pr create --title "Automatic Release $releaseDate" --body-file body
maxRetries=4
retries=0
echo "Creating and merging PR in krm-site"
until gh pr merge --auto --merge
do
  sleep 5
  [[ $retries -eq $maxRetries ]] && echo 'Unable to merge PR in krm-site!' && exit 1
  ((retries++))
  echo "GH CLI returned non-zero after $retries attempts, going to retry auto merge"
done