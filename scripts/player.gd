extends CharacterBody2D

class_name Player

const speed = 100
var current_direction = "none"

func _ready():
	position.x = Global.player_spawn_position_x
	position.y = Global.player_spawn_position_y
	$AnimatedSprite2D.play("front_idle")
	
	


func _physics_process(delta):
	player_movement(delta)

func player_movement(_delta):
	if Input.is_action_pressed("ui_right"):
		current_direction = "right"
		play_anim(1)
		velocity.x = speed
		velocity.y = 0
	elif Input.is_action_pressed("ui_left"):
		current_direction = "left"
		play_anim(1)
		velocity.x = -speed
		velocity.y = 0
	elif Input.is_action_pressed("ui_up"):
		current_direction = "up"
		play_anim(1)
		velocity.x = 0
		velocity.y = -speed
	elif Input.is_action_pressed("ui_down"):
		current_direction = "down"
		play_anim(1)
		velocity.x = 0
		velocity.y = speed
	else:
		play_anim(0)
		velocity.x = 0 
		velocity.y = 0
		
	move_and_slide()

func play_anim(movement):
	var dir = current_direction
	var anim = $AnimatedSprite2D
	if dir == "right":
		anim.flip_h = false
		if movement == 1:
			anim.play("side_walk")
		elif movement == 0:
			anim.play("side_idle")
	if dir == "left":
		anim.flip_h = true
		if movement == 1:
			anim.play("side_walk")
		elif movement == 0:
			anim.play("side_idle")
	if dir == "down":
		anim.flip_h = true
		if movement == 1:
			anim.play("front_walk")
		elif movement == 0:
			anim.play("front_idle")
	if dir == "up":
		anim.flip_h = true
		if movement == 1:
			anim.play("back_walk")
		elif movement == 0:
			anim.play("back_idle")
