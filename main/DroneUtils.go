package main

import (
	"math"
    "math/rand"
)

/**
  This class is intended to encapsulate all the physical attributes and actions of a single drone.
*/

// 3D position of the drone
type Position struct {
	X float64 `json: "x"`
	Y float64 `json: "y"`
	Z float64 `json: "z"`
}

// Dimensional velocities
type Speed struct {
	VX float64 `json: "vX"`
	VY float64 `json: "vY"`
	VZ float64 `json: "vZ"`
}

// Any set of X, Y and Z dimensions
type Dimensions struct {
	DX float64 `json: "dX"`
	DY float64 `json: "dY"`
	DZ float64 `json: "dZ"`
}

// Attributes for a certain type of drone : e.g. a DJI Phantom X2
type DroneType struct {
	TypeId          string     `json: "id"`          // Identifier for a type
	TypeDescription string     `json: "description"` // Description
	Size            Dimensions `json: "size"`        // Size of this drone type
	MaxRange        Dimensions `json: "maxRange"`    // Max communication range for this drone type
	MaxSpeed        Speed      `json: "maxSpeed"`    // Maximum speed for this drone type
}

// Drone object with all attributes for one running drone
type DroneObject struct {
	Pos   Position  `json:"pos"`
	Type  DroneType `json:"type"`
	Speed Speed     `json:"speed"`
}

// System representation of a drone including ID, URL, paxos role etc.
type Drone struct {
	ID			string		`json: "id"`
	Address		string		`json: "address"`
	DroneObject	DroneObject	`json: "droneObject"`
}

func (d Drone) moveTo(newPos Position, speed Speed) (time float64) {
	timeX := math.Abs((newPos.X - d.DroneObject.Pos.X) / speed.VX)
	timeY := math.Abs((newPos.Y - d.DroneObject.Pos.Y) / speed.VY)
	timeZ := math.Abs((newPos.Z - d.DroneObject.Pos.Z) / speed.VZ)
	return math.Max(timeX, math.Max(timeY, timeZ))
}

var samplePosition = Position{5, 4, 6}

var sampleSpeed = Speed{2, 1, -1}

var sampleDimension = Dimensions{1, 1, 1}

var sampleDroneType = DroneType{"type1", "Simple sample drone type", sampleDimension, Dimensions{10, 10, 10}, Speed{10, 10, 10}}

var sampleDroneObject = DroneObject{samplePosition, sampleDroneType, sampleSpeed}

var sampleDrone = Drone{"drone1", "localhost:1111", sampleDroneObject}

func max(list ...int) int {
	max := list[0]
	for _, i := range list[1:] {
		if i > max {
			max = i
		}
	}
	return max
}

func min(list ...int) int {
	min := list[0]
	for _, i := range list[1:] {
		if i < min {
			min = i
		}
	}
	return min
}

func getPolygonCoordinates(size int, total int) []Position {
    positions := make([]Position, total, 2*total)
    polygon := calculateCoordinates(size, 2, 5 + rand.Float64() * 10)
    for i := 0; i < total; i++ {
        if i < size {
            positions[i] = polygon[i]
        } else {
            positions[i] = Position{rand.Float64() * 20 - 10, 0, rand.Float64() * 20 - 10}
        }
    }
    return positions
}

func calculateCoordinatesForPrism(n int, usedDrones int, radius float64) []Position{
	posArray := make([]Position, n, 2*n)
        var angle float64
        angle = float64 (2) * math.Pi / float64(usedDrones/2)
        var x,y,z float64
        for i := 0; i < usedDrones/2; i++ {
            x = 0 + radius * math.Sin(float64(i) * angle)
            y = 5                 
            z = 0 + radius * math.Cos(float64(i) * angle)
            posArray[i] = Position{ x, y, z } 
            posArray[n/2 + i] = Position{ x, y + radius, z }   
        }         
        if (usedDrones%2 != 0) {
            posArray[usedDrones-1] = Position {float64(usedDrones - 1), 0, 0}
        }
        for i := usedDrones; i < n - usedDrones; i++ {
            posArray[i] = Position {float64(2*i), 0 , 0}
        } 
        return posArray 
}

func calculateCoordinatesForPyramid(n int, usedDrones int, radius float64) []Position{
        posArray := make([]Position, n, 2*n)
        var angle float64
        angle = float64 (2) * math.Pi / float64(usedDrones-1)
        var x,y,z float64
        for i := 0; i < usedDrones-1; i++ {
            x = 0 + radius * math.Sin(float64(i) * angle)
            y = 5
            z = 0 + radius * math.Cos(float64(i) * angle)
            posArray[i] = Position{ x, y, z }
        }
        posArray[usedDrones-1] = Position {0, 5 + (radius/math.Sqrt(2)), 0}
        for i := usedDrones; i < n - usedDrones; i++ {
            posArray[i] = Position {float64(2*i), 0 , 0}
        }
        return posArray
}

func calculateCoordinatesForBipyramid(n int, usedDrones int, radius float64) []Position{
        posArray := make([]Position, n, 2*n)
        var angle float64
        angle = float64 (2) * math.Pi / float64(usedDrones-2)
        var x,y,z float64
        for i := 0; i < usedDrones-2; i++ {
            x = 0 + radius * math.Sin(float64(i) * angle)
            y = 5
            z = 0 + radius * math.Cos(float64(i) * angle)
            posArray[i] = Position{ x, y, z }
        }
        posArray[usedDrones-2] = Position {0, 5 + (radius/math.Sqrt(2)), 0}
        posArray[usedDrones-1] = Position {0, (5 - (radius/math.Sqrt(2))), 0}
        for i := usedDrones; i < n - usedDrones; i++ {
            posArray[i] = Position {-(float64(2*i)), 0 , 0}
        }
        return posArray
}

func calculateCoordinates(n int, dimension int, radius float64) []Position{
        posArray := make([]Position, n, 2*n)
        if dimension == 2 {
            var angle float64
            angle = float64 (2) * math.Pi / float64(n)
            //posArray := make([]Position, n, 2*n)
            //var posArray [n]Position
            var x,y,z float64
            for i := 0; i < n; i++ {
                x = 0 + radius * math.Sin(float64(i) * angle)
                y = 5
                z = 0 + radius * math.Cos(float64(i) * angle)
                posArray[i] = Position{ x, y, z }
            }
        } else {
            if n == 4 {
                posArray[0] = Position {0, radius/2, (radius/2)/math.Sqrt(2)}
                posArray[1] = Position {0, -(radius/2), (radius/2)/math.Sqrt(2)}
                posArray[2] = Position {radius/2, 0, -(radius/2)/math.Sqrt(2)}
                posArray[3] = Position {-(radius/2), 0, -(radius/2)/math.Sqrt(2)}
            } else if n == 5 {
                posArray[0] = Position {radius, 5, 0}
                posArray[1] = Position {0, 5, radius}
                posArray[2] = Position {-(radius), 5, 0}
                posArray[3] = Position {0, 5 , -(radius)}
                posArray[4] = Position {0,  5 + (radius/math.Sqrt(2)), 0}
            } else if n == 6 {
                posArray[0] = Position {radius, 5, 0}
                posArray[1] = Position {0, 5, radius}
                posArray[2] = Position {-radius, 5, 0}
                posArray[3] = Position {0, 5, -radius}
                posArray[4] = Position {0, 5 + (radius/math.Sqrt(2)), 0}
                posArray[5] = Position {0, (-5 - (radius/math.Sqrt(2))), 0}
            } else {
                var angle float64
                angle = float64 (2) * math.Pi / float64(n)
                //posArray := make([]Position, n, 2*n)
                //var posArray [n]Position
                var x,y,z float64
                for i := 0; i < n-2; i++ {
                    x = 0 + radius * math.Sin(float64(i) * angle)
                    y = 5
                    z = 0 + radius * math.Cos(float64(i) * angle)
                    posArray[i] = Position{ x, y, z }
                }
                posArray[n-2] = Position {0, 5 + (radius/math.Sqrt(2)), 0}
                posArray[n-1] = Position {0, (-5 - (radius/math.Sqrt(2))), 0} 
           }
        }    
        return posArray
}
