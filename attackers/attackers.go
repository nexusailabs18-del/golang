package main

import "fmt"

// Attacker is an interface.
// Anything that has Attack() satisfies Attacker.
type Attacker interface {
	Attack() int
}

// Player represents the player.
type Player struct {
	Name   string
	Weapon string
	Power  int
}

// Enemy represents a normal enemy.
type Enemy struct {
	Name  string
	Power int
}

// Boss represents a powerful enemy.
type Boss struct {
	Name  string
	Power int
}

// Player's version of Attack()
func (p Player) Attack() int {
	damage := p.Power + 10
	fmt.Println(p.Name, "attacks with", p.Weapon)
	fmt.Println("Damage:", damage)
	return damage
}

// Enemy's version of Attack()
func (e Enemy) Attack() int {
	damage := e.Power
	fmt.Println(e.Name, "attacks!")
	fmt.Println("Damage:", damage)
	return damage
}

// Boss's version of Attack()
func (b Boss) Attack() int {
	damage := b.Power * 2
	fmt.Println(b.Name, "uses a DEVASTATING ATTACK!")
	fmt.Println("Damage:", damage)
	return damage
}

// fight doesn't care what the attacker actually is.
// It only cares that it satisfies Attacker.
func fight(attacker Attacker) {
	fmt.Println("--------------------")
	attacker.Attack()
	fmt.Println("--------------------")
}

func main() {

	player := Player{
		Name:   "Hero",
		Weapon: "Sword",
		Power:  20,
	}

	goblin := Enemy{
		Name:  "Goblin",
		Power: 8,
	}

	boss := Boss{
		Name:  "Dragon",
		Power: 30,
	}

	fight(player)
	fight(goblin)
	fight(boss)
}
