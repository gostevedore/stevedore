#!/bin/bash

set -e

echo -e "\e[1;31mRun Jekyll serve to watch for changes"
bundle exec jekyll serve --no-watch --livereload --drafts --host 0.0.0.0
