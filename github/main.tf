variable "GITHUB_TOKEN" {
  type = string
}

terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "6.2.3"
    }
  }
}

provider "github" {
  token = var.GITHUB_TOKEN
  owner = "gidoichi"
}

resource "github_repository" "this" {
  name                        = "yaml-path"
  allow_auto_merge            = true
  allow_merge_commit          = false
  allow_rebase_merge          = false
  delete_branch_on_merge      = true
  description                 = "Get \"xpath\" for a given line at column of a YAML file"
  has_issues                  = true
  squash_merge_commit_message = "BLANK"
  squash_merge_commit_title   = "PR_TITLE"
}

resource "github_branch_protection" "default" {
  repository_id = github_repository.this.node_id
  pattern       = "main"
  required_status_checks {
    strict = true
    contexts = [
      "go-test",
      "no-diff",
      "pull-request",
      "terraform-plan",
    ]
  }
}
