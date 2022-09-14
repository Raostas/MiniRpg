package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Weapon struct {
	damageRange                           []int
	price                                 int
	strengthReq, agilityReq, intellectReq int
}

type Consumable struct {
	hitpointsEffect int
	price           int
}

type Gopher struct {
	name       string
	hitpoints  int
	weaponName string

	strength, agility, intellect int
	coins                        int
}

var weapons = map[string]Weapon{
	"bare-handed":               {damageRange: []int{1, 1}},
	"knife":                     {damageRange: []int{2, 3}, price: 10},
	"sword":                     {damageRange: []int{3, 5}, price: 35, strengthReq: 2},
	"ninjaku":                   {damageRange: []int{1, 7}, price: 25, agilityReq: 2},
	"wand":                      {damageRange: []int{3, 3}, price: 30, intellectReq: 2},
	"gophermourne":              {damageRange: []int{6, 7}, price: 65, strengthReq: 5},
	"warglaives_of_gopherinoth": {damageRange: []int{2, 9}, price: 55, agilityReq: 5},
	"codeseeker":                {damageRange: []int{4, 4}, price: 60, intellectReq: 5},
}

var consumables = map[string]Consumable{
	"small_health_potion":  {hitpointsEffect: 5, price: 5},
	"medium_health_potion": {hitpointsEffect: 10, price: 9},
	"big_health_potion":    {hitpointsEffect: 20, price: 18},
}

var gopher1 = &Gopher{
	name:       "Pufik",
	hitpoints:  30,
	weaponName: "sword",

	coins: 20,
}

var gopher2 = &Gopher{
	name:       "Stolik",
	hitpoints:  30,
	weaponName: "bare-handed",

	coins: 20,
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	r := bufio.NewReader(os.Stdin)
	HandleInput(r)
	winner := getWinner()
	fmt.Printf("Game over! %s is the winner!\n", winner)
}

func HandleInput(r *bufio.Reader) {
	Action := ""
	turn := 0
	for Action != "exit" || isGameOver() {
		var currentGopher, otherGopher *Gopher
		if turn%2 == 0 {
			currentGopher = gopher1
			otherGopher = gopher2
		} else {
			currentGopher = gopher2
			otherGopher = gopher1
		}

		fmt.Printf("%s's turn\n", currentGopher.name)
		fmt.Printf(">")
		line, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}

		tokens := strings.Split(strings.TrimSpace(line), " ")
		Action = tokens[0]
		args := tokens[1:]
		switch Action {
		case "attack":
			err = attack(currentGopher, otherGopher)
		case "buy":
			if args[0] == "weapon" {
				err = buyWeapon(currentGopher, args[1])
			} else if args[0] == "consumable" {
				err = buyConsumable(currentGopher, args[1])
			} else {
				err = errors.New("invalid command")
			}

		case "work":
			err = work(currentGopher)
		case "train":
			err = train(currentGopher, args[0])
		}

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println()
			fmt.Println("--------")
			printStats(gopher1)
			fmt.Println("--------")
			printStats(gopher2)
			fmt.Println("--------")
			fmt.Println()
			turn++

		}
	}

}

func attack(attacker *Gopher, defender *Gopher) error {
	attackerWeapon, ok := weapons[attacker.name]
	if !ok {
		return fmt.Errorf("invalid state.The weapon of the attacker doesn't exist")
	}
	damage := intInRange(attackerWeapon.damageRange[0], attackerWeapon.damageRange[1])
	defender.hitpoints -= damage
	return nil

}

func buyWeapon(gopher *Gopher, weaponName string) error {
	weapon, ok := weapons[weaponName]
	if !ok {
		return fmt.Errorf("the weapon %s doesn't exist", weaponName)
	}

	if gopher.coins < weapon.price {
		return fmt.Errorf("insufficient coins. Need %d, got %d", weapon.price, gopher.coins)
	}

	if gopher.strength < weapon.strengthReq ||
		gopher.agility < weapon.agilityReq ||
		gopher.intellect < weapon.intellectReq {
		return errors.New("insufficient attributes")
	}

	gopher.weaponName = weaponName
	gopher.coins -= weapon.price
	return nil
}

func buyConsumable(gopher *Gopher, consumableName string) error {
	consumable, ok := consumables[consumableName]
	if !ok {
		return fmt.Errorf("the consumable %s doesn't exist", consumableName)
	}

	if gopher.coins < consumable.price {
		return fmt.Errorf("insufficient coins. Need %d, got %d",
			consumable.price, gopher.coins)
	}
	gopher.coins -= consumable.price
	gopher.hitpoints = int(math.Min(30, float64(gopher.hitpoints+consumable.hitpointsEffect)))
	return nil

}

func work(gopher *Gopher) error {

	salary := intInRange(5, 15)
	gopher.coins += salary
	return nil
}

func train(gopher *Gopher, skill string) error {
	//TODO: Implement
	if skill != "strength" && skill != "agility" && skill != "intellect" {
		return errors.New("bad attribut chosen ")
	}

	if gopher.coins < 5 {
		return fmt.Errorf(
			"insufficient coins.. Need 5, you have %d", gopher.coins)
	}

	switch skill {
	case "strength":
		gopher.strength += 2
	case "agility":
		gopher.agility += 2
	case "intellect":
		gopher.intellect += 2
	default:
		return fmt.Errorf(
			"insufficient coins.. Need 5, you have %d", gopher.coins)
	}

	gopher.coins -= 5
	return nil
}

func isGameOver() bool {
	return gopher1.hitpoints > 0 || gopher2.hitpoints < 0
}

func getWinner() string {
	if gopher1.hitpoints < 0 {
		return gopher2.name
	} else if gopher2.hitpoints < 0 {
		return gopher1.name
	}
	panic("no winner yet")
}

func printStats(gopher *Gopher) {

	fmt.Printf(
		"name : %s\n hitpoints:\nweapon:%s\n"+
			" %d\n strength%d\nagility%d\nintellect%d\ncoins%d\n",
		gopher.name, gopher.hitpoints, gopher.weaponName,
		gopher.strength, gopher.agility, gopher.intellect, gopher.coins)
}

func intInRange(start, end int) int {
	return rand.Intn(end-start+1) + start
}
