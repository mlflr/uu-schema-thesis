output "connection_string" {
  sensitive = true
  value = {
    for k in toset(local.schema_versioning_methods) : k =>
    "postgresql://${neon_role.thesis_role[k].name}:${neon_role.thesis_role[k].password}@${neon_project.schema_thesis[k].database_host}/${neon_database.movies[k].name}?sslmode=require"
  }
}
