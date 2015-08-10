#!/bin/bash
URL="http://localhost:8000"

curl -X POST --url ${URL} -H "X-Gitlab-Event: Push Hook" -d @push.json
curl -X POST --url ${URL} -H "X-Gitlab-Event: Tag Push Hook" -d @tag.json
curl -X POST --url ${URL} -H "X-Gitlab-Event: Issue Hook" -d @issues.json
curl -X POST --url ${URL} -H "X-Gitlab-Event: Note Hook" -d @commit_comment.json
curl -X POST --url ${URL} -H "X-Gitlab-Event: Note Hook" -d @merge_comment.json
curl -X POST --url ${URL} -H "X-Gitlab-Event: Note Hook" -d @issue_comment.json
curl -X POST --url ${URL} -H "X-Gitlab-Event: Note Hook" -d @snippet_comment.json
curl -X POST --url ${URL} -H "X-Gitlab-Event: Merge Request Hook" -d @merge.json
