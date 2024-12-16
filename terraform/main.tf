provider "neon" {}

resource "neon_project" "schema_thesis" {
  for_each = toset(local.schema_versioning_methods)

  name                      = "schema_thesis_${each.key}"
  pg_version                = "16"
  history_retention_seconds = 86400
  region_id                 = "aws-eu-central-1"
}

resource "neon_role" "thesis_role" {
  for_each = toset(local.schema_versioning_methods)

  project_id = neon_project.schema_thesis[each.key].id
  branch_id  = neon_project.schema_thesis[each.key].default_branch_id
  name       = "thesis_role"
}

resource "neon_database" "movies" {
  for_each = toset(local.schema_versioning_methods)

  project_id = neon_project.schema_thesis[each.key].id
  branch_id  = neon_project.schema_thesis[each.key].default_branch_id
  owner_name = neon_role.thesis_role[each.key].name
  name       = "movies"
}
