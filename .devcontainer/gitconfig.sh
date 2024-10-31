#!/usr/bin/env zsh

rm -f git.env
echo "GIT_AUTHOR_EMAIL=\"$(git config get --global user.email)\"" >> .devcontainer/git.env
echo "GIT_COMMITTER_EMAIL=\"$(git config get --global user.email)\"" >> .devcontainer/git.env
echo "GIT_AUTHOR_NAME=\"$(git config get --global user.name)\"" >> .devcontainer/git.env
echo "GIT_COMMITTER_NAME=\"$(git config get --global user.name)\"" >> .devcontainer/git.env
