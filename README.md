# ghapp

Investigation of creating a GitHub App.

- The feature we want to use is [Commit Statuses](https://docs.github.com/en/rest/commits/statuses?apiVersion=2022-11-28)
  - [Create a commit status](https://docs.github.com/en/rest/commits/statuses?apiVersion=2022-11-28#create-a-commit-status)

### Step 1

I use [webhook.site](https://webhook.site) for testing.

This is the webhook payload I got when I installed the app on @bartekpacia

<details>
<summary>Payload</summary>

```json
{
  "action": "created",
  "installation": {
    "id": 52534914,
    "account": {
      "login": "bartekpacia",
      "id": 40357511,
      "node_id": "MDQ6VXNlcjQwMzU3NTEx",
      "avatar_url": "https://avatars.githubusercontent.com/u/40357511?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/bartekpacia",
      "html_url": "https://github.com/bartekpacia",
      "followers_url": "https://api.github.com/users/bartekpacia/followers",
      "following_url": "https://api.github.com/users/bartekpacia/following{/other_user}",
      "gists_url": "https://api.github.com/users/bartekpacia/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/bartekpacia/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/bartekpacia/subscriptions",
      "organizations_url": "https://api.github.com/users/bartekpacia/orgs",
      "repos_url": "https://api.github.com/users/bartekpacia/repos",
      "events_url": "https://api.github.com/users/bartekpacia/events{/privacy}",
      "received_events_url": "https://api.github.com/users/bartekpacia/received_events",
      "type": "User",
      "site_admin": false
    },
    "repository_selection": "all",
    "access_tokens_url": "https://api.github.com/app/installations/52534914/access_tokens",
    "repositories_url": "https://api.github.com/installation/repositories",
    "html_url": "https://github.com/settings/installations/52534914",
    "app_id": 938460,
    "app_slug": "silesiaci",
    "target_id": 40357511,
    "target_type": "User",
    "permissions": {
      "contents": "read",
      "metadata": "read",
      "pull_requests": "read",
      "statuses": "write"
    },
    "events": [
      "pull_request",
      "push"
    ],
    "created_at": "2024-07-06T02:08:35.000+02:00",
    "updated_at": "2024-07-06T02:08:36.000+02:00",
    "single_file_name": null,
    "has_multiple_single_files": false,
    "single_file_paths": [],
    "suspended_by": null,
    "suspended_at": null
  },
  "repositories": [
    {
      "id": 138645750,
      "node_id": "MDEwOlJlcG9zaXRvcnkxMzg2NDU3NTA=",
      "name": "spitfire",
      "full_name": "bartekpacia/spitfire",
      "private": false
    },
    {
      "id": 154893669,
      "node_id": "MDEwOlJlcG9zaXRvcnkxNTQ4OTM2Njk=",
      "name": "visiting-card-android",
      "full_name": "bartekpacia/visiting-card-android",
      "private": false
    },
  ],
  "requester": null,
  "sender": {
    "login": "bartekpacia",
    "id": 40357511,
    "node_id": "MDQ6VXNlcjQwMzU3NTEx",
    "avatar_url": "https://avatars.githubusercontent.com/u/40357511?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/bartekpacia",
    "html_url": "https://github.com/bartekpacia",
    "followers_url": "https://api.github.com/users/bartekpacia/followers",
    "following_url": "https://api.github.com/users/bartekpacia/following{/other_user}",
    "gists_url": "https://api.github.com/users/bartekpacia/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/bartekpacia/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/bartekpacia/subscriptions",
    "organizations_url": "https://api.github.com/users/bartekpacia/orgs",
    "repos_url": "https://api.github.com/users/bartekpacia/repos",
    "events_url": "https://api.github.com/users/bartekpacia/events{/privacy}",
    "received_events_url": "https://api.github.com/users/bartekpacia/received_events",
    "type": "User",
    "site_admin": false
  }
}
```

</details>

This schema for this type of event is described [here](https://docs.github.com/en/webhooks/webhook-events-and-payloads#installation).

### Step 2

Then I pushed a commit to bartekpacia/dumbpkg repository.

This resulted in the following webhook being sent:

<details>
<summary>Payload</summary>

```json
{
  "ref": "refs/heads/master",
  "before": "fa90a7d104ac220613fcd34bf59eeafc2d71bcad",
  "after": "fa3fd6baa2e1532793e51e0dd1239bd55e480bd8",
  "repository": {
    "id": 609253504,
    "node_id": "R_kgDOJFB4gA",
    "name": "dumbpkg",
    "full_name": "bartekpacia/dumbpkg",
    "private": false,
    "owner": {
      "name": "bartekpacia",
      "email": "barpac02@gmail.com",
      "login": "bartekpacia",
      "id": 40357511,
      "node_id": "MDQ6VXNlcjQwMzU3NTEx",
      "avatar_url": "https://avatars.githubusercontent.com/u/40357511?v=4",
      "gravatar_id": "",
      "url": "https://api.github.com/users/bartekpacia",
      "html_url": "https://github.com/bartekpacia",
      "followers_url": "https://api.github.com/users/bartekpacia/followers",
      "following_url": "https://api.github.com/users/bartekpacia/following{/other_user}",
      "gists_url": "https://api.github.com/users/bartekpacia/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/bartekpacia/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/bartekpacia/subscriptions",
      "organizations_url": "https://api.github.com/users/bartekpacia/orgs",
      "repos_url": "https://api.github.com/users/bartekpacia/repos",
      "events_url": "https://api.github.com/users/bartekpacia/events{/privacy}",
      "received_events_url": "https://api.github.com/users/bartekpacia/received_events",
      "type": "User",
      "site_admin": false
    },
    "html_url": "https://github.com/bartekpacia/dumbpkg",
    "description": null,
    "fork": false,
    "url": "https://github.com/bartekpacia/dumbpkg",
    "forks_url": "https://api.github.com/repos/bartekpacia/dumbpkg/forks",
    "keys_url": "https://api.github.com/repos/bartekpacia/dumbpkg/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/bartekpacia/dumbpkg/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/bartekpacia/dumbpkg/teams",
    "hooks_url": "https://api.github.com/repos/bartekpacia/dumbpkg/hooks",
    "issue_events_url": "https://api.github.com/repos/bartekpacia/dumbpkg/issues/events{/number}",
    "events_url": "https://api.github.com/repos/bartekpacia/dumbpkg/events",
    "assignees_url": "https://api.github.com/repos/bartekpacia/dumbpkg/assignees{/user}",
    "branches_url": "https://api.github.com/repos/bartekpacia/dumbpkg/branches{/branch}",
    "tags_url": "https://api.github.com/repos/bartekpacia/dumbpkg/tags",
    "blobs_url": "https://api.github.com/repos/bartekpacia/dumbpkg/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/bartekpacia/dumbpkg/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/bartekpacia/dumbpkg/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/bartekpacia/dumbpkg/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/bartekpacia/dumbpkg/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/bartekpacia/dumbpkg/languages",
    "stargazers_url": "https://api.github.com/repos/bartekpacia/dumbpkg/stargazers",
    "contributors_url": "https://api.github.com/repos/bartekpacia/dumbpkg/contributors",
    "subscribers_url": "https://api.github.com/repos/bartekpacia/dumbpkg/subscribers",
    "subscription_url": "https://api.github.com/repos/bartekpacia/dumbpkg/subscription",
    "commits_url": "https://api.github.com/repos/bartekpacia/dumbpkg/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/bartekpacia/dumbpkg/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/bartekpacia/dumbpkg/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/bartekpacia/dumbpkg/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/bartekpacia/dumbpkg/contents/{+path}",
    "compare_url": "https://api.github.com/repos/bartekpacia/dumbpkg/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/bartekpacia/dumbpkg/merges",
    "archive_url": "https://api.github.com/repos/bartekpacia/dumbpkg/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/bartekpacia/dumbpkg/downloads",
    "issues_url": "https://api.github.com/repos/bartekpacia/dumbpkg/issues{/number}",
    "pulls_url": "https://api.github.com/repos/bartekpacia/dumbpkg/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/bartekpacia/dumbpkg/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/bartekpacia/dumbpkg/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/bartekpacia/dumbpkg/labels{/name}",
    "releases_url": "https://api.github.com/repos/bartekpacia/dumbpkg/releases{/id}",
    "deployments_url": "https://api.github.com/repos/bartekpacia/dumbpkg/deployments",
    "created_at": 1677865384,
    "updated_at": "2024-04-11T01:24:51Z",
    "pushed_at": 1720225285,
    "git_url": "git://github.com/bartekpacia/dumbpkg.git",
    "ssh_url": "git@github.com:bartekpacia/dumbpkg.git",
    "clone_url": "https://github.com/bartekpacia/dumbpkg.git",
    "svn_url": "https://github.com/bartekpacia/dumbpkg",
    "homepage": null,
    "size": 11,
    "stargazers_count": 1,
    "watchers_count": 1,
    "language": "Dart",
    "has_issues": true,
    "has_projects": true,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": false,
    "has_discussions": false,
    "forks_count": 0,
    "mirror_url": null,
    "archived": false,
    "disabled": false,
    "open_issues_count": 0,
    "license": {
      "key": "mit",
      "name": "MIT License",
      "spdx_id": "MIT",
      "url": "https://api.github.com/licenses/mit",
      "node_id": "MDc6TGljZW5zZTEz"
    },
    "allow_forking": true,
    "is_template": false,
    "web_commit_signoff_required": false,
    "topics": [],
    "visibility": "public",
    "forks": 0,
    "open_issues": 0,
    "watchers": 1,
    "default_branch": "master",
    "stargazers": 1,
    "master_branch": "master"
  },
  "pusher": {
    "name": "bartekpacia",
    "email": "barpac02@gmail.com"
  },
  "sender": {
    "login": "bartekpacia",
    "id": 40357511,
    "node_id": "MDQ6VXNlcjQwMzU3NTEx",
    "avatar_url": "https://avatars.githubusercontent.com/u/40357511?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/bartekpacia",
    "html_url": "https://github.com/bartekpacia",
    "followers_url": "https://api.github.com/users/bartekpacia/followers",
    "following_url": "https://api.github.com/users/bartekpacia/following{/other_user}",
    "gists_url": "https://api.github.com/users/bartekpacia/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/bartekpacia/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/bartekpacia/subscriptions",
    "organizations_url": "https://api.github.com/users/bartekpacia/orgs",
    "repos_url": "https://api.github.com/users/bartekpacia/repos",
    "events_url": "https://api.github.com/users/bartekpacia/events{/privacy}",
    "received_events_url": "https://api.github.com/users/bartekpacia/received_events",
    "type": "User",
    "site_admin": false
  },
  "installation": {
    "id": 52534914,
    "node_id": "MDIzOkludGVncmF0aW9uSW5zdGFsbGF0aW9uNTI1MzQ5MTQ="
  },
  "created": false,
  "deleted": false,
  "forced": false,
  "base_ref": null,
  "compare": "https://github.com/bartekpacia/dumbpkg/compare/fa90a7d104ac...fa3fd6baa2e1",
  "commits": [
    {
      "id": "fa3fd6baa2e1532793e51e0dd1239bd55e480bd8",
      "tree_id": "e5da987b75df9d818e5058d8d14d669ba52773b0",
      "distinct": true,
      "message": "empty",
      "timestamp": "2024-07-06T02:21:22+02:00",
      "url": "https://github.com/bartekpacia/dumbpkg/commit/fa3fd6baa2e1532793e51e0dd1239bd55e480bd8",
      "author": {
        "name": "Bartek Pacia",
        "email": "barpac02@gmail.com",
        "username": "bartekpacia"
      },
      "committer": {
        "name": "Bartek Pacia",
        "email": "barpac02@gmail.com",
        "username": "bartekpacia"
      },
      "added": [],
      "removed": [],
      "modified": []
    }
  ],
  "head_commit": {
    "id": "fa3fd6baa2e1532793e51e0dd1239bd55e480bd8",
    "tree_id": "e5da987b75df9d818e5058d8d14d669ba52773b0",
    "distinct": true,
    "message": "empty",
    "timestamp": "2024-07-06T02:21:22+02:00",
    "url": "https://github.com/bartekpacia/dumbpkg/commit/fa3fd6baa2e1532793e51e0dd1239bd55e480bd8",
    "author": {
      "name": "Bartek Pacia",
      "email": "barpac02@gmail.com",
      "username": "bartekpacia"
    },
    "committer": {
      "name": "Bartek Pacia",
      "email": "barpac02@gmail.com",
      "username": "bartekpacia"
    },
    "added": [],
    "removed": [],
    "modified": []
  }
}
```

</details>
