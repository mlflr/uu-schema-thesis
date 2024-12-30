output "connection_string" {
  sensitive = true
  value = {
    for k in toset(local.schema_versioning_methods) : k => neon_project.schema_thesis[k].connection_uri
  }
}
