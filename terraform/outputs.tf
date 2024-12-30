output "connection_string" {
  value = {
    for k in toset(local.schema_versioning_methods) : k => neon_project.schema_thesis[k].connection_uri
  }
}
