output "cluster_name" {
  value = aws_ecs_cluster.main.name
}

output "ecs_tasks_security_group_id" {
  value = aws_security_group.ecs_tasks.id
}
