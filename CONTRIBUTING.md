# Contributing

We'd love your help making Instahelper the very best Instagram Automation Tool!

If you'd like to add new features or help contribute to one, please [open an issue][open-issue]
describing your proposal &mdash; discussing API changes ahead of time makes
pull request review much smoother.


## Setup

[Fork][fork], then clone the repository:

```
go get github.com/socialplanner/instahelper
cd $GOPATH/src/github.com/instahelper
git remote add upstream https://github.com/socialplanner/instahelper
git fetch upstream
make assets
```

## Making Changes

To ensure your contribution gets added we request you to.
* Include documentation
* Lint all code
* Test all code beforehand
* Setup proper logging
* Add notable additions/improvements to the [changelog][changelog]

```
golint
git push origin cool_new_feature
```

Then use the GitHub UI to open a pull request.

At this point, you're waiting on us to review your changes. We *try* to respond
to issues and pull requests within a few days, and we may suggest some
improvements or alternatives. Once your changes are approved, one of the
project maintainers will merge them.

We're much more likely to approve your changes if you:

* Add tests for new functionality.
* Write a [good commit message][commit-message].
* Maintain backward compatibility.

[fork]: https://github.com/socialplanner/instahelper/fork
[open-issue]: https://github.com/socialplanner/instahelper/issues/new
[commit-message]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html
[changelog]: https://github.com/socialplanner/instahelper/blob/master/TERMINAL.md