package calculator

import (
	"fmt"
	"layouts-orders-bot/internal/entity"
	"math"
	"strconv"
	"strings"
)

type Calculator struct{}

func New() *Calculator {
	return &Calculator{}
}

type priceRange struct {
	min float64
	max float64
}

func (p *priceRange) multiply(min, max float64) {
	p.min *= min
	p.max *= max
}

func (p *priceRange) add(min, max float64) {
	p.min += min
	p.max += max
}

func (c *Calculator) Calculate(props map[*entity.State]string) (int, int) {
	price := priceRange{
		min: 4500,
		max: 4500,
	}

	c.applyLayoutType(&price, props)
	c.applyPurpose(&price, props)
	c.applyScale(&price, props)
	c.applySize(&price, props)
	c.applyMaterial(&price, props)
	c.applyDetails(&price, props)
	c.apply3DPrint(&price, props)
	c.applyLandscape(&price, props)
	c.applyDrawings(&price, props)
	c.applyDeadline(&price, props)

	return roundDown(price.min), roundUp(price.max)
}

func (c *Calculator) applyLayoutType(price *priceRange, props map[*entity.State]string) {
	switch props[entity.StateLayoutType] {
	case entity.LayoutTypeArchitectural:
	case entity.LayoutTypeInterior:
		price.multiply(1.15, 1.25)
	case entity.LayoutTypeIndustrial:
	case entity.LayoutTypeLandscape:
	default:
		panic("unknown layout type")
	}
}

func (c *Calculator) applyPurpose(price *priceRange, props map[*entity.State]string) {
	switch props[entity.StatePurpose] {
	case entity.PurposeEducational:
	case entity.PurposeSchool:
		price.multiply(0.9, 0.9)
	case entity.PurposeExhibition:
		price.multiply(1.2, 1.4)
	case entity.PurposeGift:
		price.multiply(1.05, 1.05)
	default:
		panic("unknown purpose")
	}
}

func (c *Calculator) applyScale(_ *priceRange, props map[*entity.State]string) {
	scale := props[entity.StateScale]
	if _, err := entity.StateScale.SpecialValidation(scale); err != nil {
		panic(fmt.Sprintf("invalid scale %q", scale))
	}
}

func (c *Calculator) applySize(price *priceRange, props map[*entity.State]string) {
	size := props[entity.StateSize]
	if _, err := entity.StateSize.SpecialValidation(size); err != nil {
		panic(fmt.Sprintf("invalid size %q", size))
	}
	k := estimateSize(size)
	price.multiply(k, k)
}

func estimateSize(size string) float64 {
	num1 := strings.Split(size, "x")[0]
	num2 := strings.Split(size, "x")[1]

	N, _ := strconv.Atoi(num1)
	M, _ := strconv.Atoi(num2)

	return math.Sqrt(float64(N*M) / (15 * 15))
}

func (c *Calculator) applyMaterial(price *priceRange, props map[*entity.State]string) {
	switch props[entity.StateMaterial] {
	case entity.MaterialCardboard:
	case entity.MaterialFoamBoard:
		price.multiply(2, 3)
	case entity.MaterialPlastic:
		price.multiply(4, 5)
	case entity.MaterialCombined:
		price.multiply(2, 3)
	default:
		panic("unknown material")
	}
}

func (c *Calculator) applyDetails(price *priceRange, props map[*entity.State]string) {
	switch props[entity.StateDetails] {
	case entity.DetailsBasic:
		price.multiply(0.8, 1)
	case entity.DetailsMedium:
	case entity.DetailsHigh:
		price.multiply(1.2, 1.3)
	default:
		panic("unknown details")
	}
}

func (c *Calculator) apply3DPrint(price *priceRange, props map[*entity.State]string) {
	switch props[entity.State3DPrint] {
	case entity.Print3DYes:
		price.add(1000, 7000)
	case entity.Print3DNo:
	default:
		panic("unknown 3d print option")
	}
}

func (c *Calculator) applyLandscape(price *priceRange, props map[*entity.State]string) {
	switch props[entity.StateLandscape] {
	case entity.LandscapeYes:
		price.add(500, 3000)
	case entity.LandscapeNo:
	default:
		panic("unknown landscape option")
	}
}
func (c *Calculator) applyDrawings(price *priceRange, props map[*entity.State]string) {
	switch props[entity.StateDrawings] {
	case entity.DrawingsNone, entity.DrawingsPartial:
		price.add(2500, 2500)
	case entity.DrawingsComplete:
	default:
		panic("unknown drawings option")
	}
}

func (c *Calculator) applyDeadline(price *priceRange, props map[*entity.State]string) {
	switch props[entity.StateDeadline] {
	case entity.Deadline2Days:
		price.multiply(1.7, 2)
	case entity.Deadline7Days:
		price.multiply(1.2, 1.3)
	case entity.Deadline14Days:
	case entity.DeadlineMonth:
	case entity.DeadlineLonger:
	default:
		panic("unknown deadline")
	}
}

func roundDown(cost float64) int {
	return int(math.Floor(cost/1000)) * 1000
}

func roundUp(cost float64) int {
	return int(math.Ceil(cost/1000)) * 1000
}
