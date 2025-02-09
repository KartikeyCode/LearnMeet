extends StaticBody2D

@onready var interactable: Area2D = $Interactable
@onready var ui = $Sign
@onready var label = $Sign/popup/base/NinePatchRect/MarginContainer/VBoxContainer/Label

func _ready() -> void:
	ui.visible=false
	interactable.interact = _on_interact
	
	
func _on_interact():
	_toggle_ui()

func _toggle_ui():
	if ui:
		ui.visible = not ui.visible
	
