extends Node

#Global Variables:
var current_scene = "world"
var transition_scene:bool 
var player_spawn_position_x = 37
var player_spawn_position_y = 249
var scene_name:String

#Global Functions:
func set_spawn_point(x,y):
	player_spawn_position_x = x
	player_spawn_position_y = y

	
