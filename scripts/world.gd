extends Node2D

@onready var welcome_sign_label: Label = $welcomeSign/Sign/popup/base/NinePatchRect/MarginContainer/VBoxContainer/Label
@onready var left_sign_label: Label = $leftSign/Sign/popup/base/NinePatchRect/MarginContainer/VBoxContainer/Label

const packets := preload("res://scripts/packets.gd")

@onready var _log := $CanvasLayer/Log

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	WS.connected_to_server.connect(_on_ws_connected_to_server)
	WS.connection_closed.connect(_on_ws_connection_closed)
	WS.packet_received.connect(_on_ws_packet_received)
	
	_log.info("Connecting to server")
	WS.connect_to_url("ws://localhost:8080/ws")
	
	welcome_sign_label.text = "Hello! \n Welcome to LearnMeet"
	left_sign_label.text = "North: Playground \n East: StudyHall"
	var packet := packets.Packet.new()
	packet.set_sender_id(69)
	var chat_msg := packet.new_chat()
	chat_msg.set_msg("Hello!")
	push_warning(packet)			
	var data := packet.to_bytes()				  

func _on_ws_connected_to_server() -> void:
	var packet := packets.Packet.new()
	var chat_msg := packet.new_chat()
	chat_msg.set_msg("Hello!")
	var err := WS.send(packet)
	if err:
		_log.error("ERROR SENDING PACKET")
	else:
		_log.success("Packet Sent Successfully")
	
func _on_ws_connection_closed() -> void:
	_log.warn("Connection Closed")
	
func _on_ws_packet_received(packet: packets.Packet) -> void:
	_log.info("Received packet from the server: %s" % packet)

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(_delta: float) -> void:
	pass



func _on_left_transition_body_entered(body: Node2D) -> void:
	if body is Player:
		get_tree().change_scene_to_file("res://scenes/left.tscn")
		Global.set_spawn_point(183,325)
		print("Entering Left")
