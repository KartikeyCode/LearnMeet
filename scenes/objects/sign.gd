extends StaticBody2D

@onready var interactable: Area2D = $Interactable

func _ready() -> void:
	interactable.interact = _on_interact
	
func _on_interact():
	push_error("Button was clicked")
