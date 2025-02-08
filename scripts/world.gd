extends Node2D


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	pass # Replace with function body.


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass


func _on_child_entered_tree(body):
	if body is Player:
		print(body)


func _on_left_transition_body_entered(body: Node2D) -> void:
	pass # Replace with function body.
