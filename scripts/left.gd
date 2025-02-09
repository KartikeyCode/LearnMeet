extends Node2D


func _on_world_transition_body_entered(body: Node2D) -> void:
	if body is Player:
		get_tree().change_scene_to_file("res://scenes/world.tscn")
		Global.set_spawn_point(183,16)
		print("Entering World")
