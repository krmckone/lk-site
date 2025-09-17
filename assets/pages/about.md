### About me

My name is Kaleb McKone. I'm a software engineer with a current focus in web development. I currently work at TPM startup Vividy. In the past I've worked at at Thomson Reuters, the University of Iowa, Target HQ, an ed tech startup, and a freight tech SaaS startup.

My current hobby/goal is slowly reaching into computer graphics with a focus on real time rendering. Part of this site's purpose is to document my learning journey as I try to break into the computer graphics space.

I studied computer science and mathematics at the University of Iowa in Iowa City, Iowa. I graduated in 2019. Some of my favorite topics that I studied at undergrad include discrete mathematics, linear and abstract algebra, numerical analysis, algorithms, and programming language design. During my Junior and Senior years I worked in Iowa's math tutoring lab part time.

### How does this site work?

This site is a project in progress, both in the content and infrastructure. Today, it is a static content site hosted by GitHub pages at [krm-site](https://github.com/krmckone/krm-site). The domain krmckone.com is managed by cloudflare. The source of the site is a combination of markdown and HTML templates. Individual pages are written in markdown and a static site generator consumes those files and inserts the content into HTML templates to form the site itself that you are viewing now. For example, this about page is authored in markdown but is converted to HTML using a template. The styling is primarily implemented by [mcss](https://mikemai.net/mcss/) under the Verdana style.

All of the implementation details of how this works are visible in the mentioned repositories.

#### Static Site Generator

The static site generator is [lk-site](https://github.com/krmckone/lk-site). This is a custom tool implemented in Go that converts a combination of markdown content and raw HTML snippets into full HTML pages. It takes as input markdown files and HTML template files and produces the HTML content that your browser is presenting to you. This is in continuous development and still could use a handful of features implemented for making more content on this site easier. I try to track the big ones using GitHub issues for the lk-site repository.

The lk-site repository utilizes GitHub Actions for CI/CD workflows. For all pull requests, unit tests are run and merging is blocked until they pass. For pull requests that merge against main, a deploy is triggered that runs the tests again and builds the site itself. After building the static assets, they are uploaded to GitHub's artifact repository. A shell script then checks out krm-site and creates a new branch, adds any modifications, then creates a new pull request and automatically merges once any checks pass. I will generally ask GitHub to merge the pull request against lk-site automatically by using `gh pr merge --auto`, which effectively kicks off the deploy.

Following the automatic merge in krm-site against main, the final deploy step is automatically managed by GitHub pages. Generally, changes that are merged against lk-site are visible at krmckone.com after GitHub Pages finishes with its deploy after 1-2 minutes.
