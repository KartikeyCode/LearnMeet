extends Node2D

@onready var welcome_sign_label: Label = $welcomeSign/Sign/popup/base/NinePatchRect/MarginContainer/VBoxContainer/Label
@onready var left_sign_label: Label = $leftSign/Sign/popup/base/NinePatchRect/MarginContainer/VBoxContainer/Label


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	left_sign_label.text = "North: Playground \n East: StudyHall"
							  

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(_delta: float) -> void:
	pass



func _on_left_transition_body_entered(body: Node2D) -> void:
	if body is Player:
		get_tree().change_scene_to_file("res://scenes/left.tscn")
		Global.set_spawn_point(183,325)
		print("Entering Left")
