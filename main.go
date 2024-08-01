package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Returns true if ray intersects with rectangle, ratio of distance over direction vector, and pointer to intersection point
//
// rayOrigin: starting point of the ray
//
// rayDir: direction of the ray (not normalized)
//
// rect: rectangle to check collision
func RayVSRect(rayOrigin, rayDir rl.Vector2, rect rl.Rectangle) (bool, float32, *rl.Vector2) {
	// tNear and tFar is the distance ratio vector from origin of the ray to the near/far intersection
	// point with rectangle over the ray direction (ray direction vector assumed to be not normalized)
	tNear := rl.Vector2{
		X: (rect.X - rayOrigin.X) / rayDir.X,
		Y: (rect.Y - rayOrigin.Y) / rayDir.Y,
	}
	tFar := rl.Vector2{
		X: (rect.X + rect.Width - rayOrigin.X) / rayDir.X,
		Y: (rect.Y + rect.Height - rayOrigin.Y) / rayDir.Y,
	}
	// swap if necessary, so they are in the correct order
	if tNear.X > tFar.X {
		tNear.X, tFar.X = tFar.X, tNear.X
	}
	if tNear.Y > tFar.Y {
		tNear.Y, tFar.Y = tFar.Y, tNear.Y
	}

	// if collision
	if (tNear.X > tFar.Y) || (tNear.Y > tFar.X) {
		return false, -1, nil
	}

	tHitNear := float32(math.Max(float64(tNear.X), float64(tNear.Y)))
	tHitFar := math.Min(float64(tFar.X), float64(tFar.Y))

	// if tHitFar is negative, it means the ray is pointing away from the rectangle
	if tHitFar < 0 {
		return false, -1, nil
	}

	hitPoint := rl.Vector2Add(rayOrigin, rl.Vector2Scale(rayDir, tHitNear))
	// hit, t, *hitPoint
	return true, tHitNear, &hitPoint
}

// Return true if rectangle/actor will collide with the solid, and point in next frame rectangle/actor should be placed
//
// returned position is the position of the actor after the collision
//
// solid: solid rectangle to check collision with
//
// P: previous rectangle/actor C: current rectangle/actor
func SweptAABB(solid, P, C rl.Rectangle) (bool, *rl.Vector2) {
	PC := rl.Vector2{X: P.X + P.Width/2, Y: P.Y + P.Height/2} // center of the previous actor rectangle
	CC := rl.Vector2{X: C.X + C.Width/2, Y: C.Y + C.Height/2} // center of the current actor rectangle
	dir := rl.NewVector2(CC.X-PC.X, CC.Y-PC.Y)                // direction from PC to CC

	// extend the solid rectangle in each side by half of the actor rectangle
	// this is to check the collision with ray from the center of the previous(P) rectangle
	// to the center of the current(C) rectangle
	extendedSolid := rl.Rectangle{
		X:      solid.X - P.Width/2,
		Y:      solid.Y - P.Height/2,
		Width:  solid.Width + P.Width,
		Height: solid.Height + P.Height,
	}
	hit, t, hitPoint := RayVSRect(PC, dir, extendedSolid)
	// t > 1 means collision will not happen in this frame
	if hit && t <= 1 {
		hitPointPos := &rl.Vector2{
			X: hitPoint.X - P.Width/2,
			Y: hitPoint.Y - P.Height/2,
		}
		return true, hitPointPos
	}

	return false, nil
}
