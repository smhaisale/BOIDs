package main

import (
	"math"
	"strconv"
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

func toString(position Position) string {
	return "(" + strconv.FormatFloat(position.X, 'f', -1, 64) + "," + strconv.FormatFloat(position.Y, 'f', -1, 64) + "," + strconv.FormatFloat(position.Z, 'f', -1, 64) + ")"
}

func add(pos1 Position, pos2 Position) Position {
	return Position{pos1.X + pos2.X, pos1.Y + pos2.Y, pos1.Z + pos2.Z}
}

func sub(pos1 Position, pos2 Position) Position {
	return Position{pos1.X - pos2.X, pos1.Y - pos2.Y, pos1.Z - pos2.Z}
}

func dotMul(pos1 Position, pos2 Position) float64 {
	return pos1.X * pos2.X + pos1.Y * pos2.Y + pos1.Z * pos2.Z
}

func scalaMul(x float64, pos Position) Position {
	return Position{x * pos.X, x * pos.Y, x * pos.Z}
}

func norm(pos Position) float64{
	return math.Sqrt(dotMul(pos, pos))
}

func dist3D_Segment_to_Segment(path1 PathLock, path2 PathLock) float64 {
	SMALL_NUM := 0.00000001
	u := sub(path1.To, path1.From)
	v := sub(path2.To, path2.From)
	w := sub(path1.From, path2.From)
	a := dotMul(u, u)
	b := dotMul(u, v)
	c := dotMul(v, v)
	d := dotMul(u, w)
	e := dotMul(v, w)
	D := a * c - b * b
	sc, sN, sD := D, D, D
	tc, tN, tD := D, D, D

	if D < SMALL_NUM {
		sN = 0.0
		sD = 1.0
		tN = e
		tD = c
	} else {
		sN = b * e - c * d
		tN = a * e - b * d
		if sN < 0.0 {
			sN = 0.0
			tN = e
			tD = c
		} else if sN > sD {
			sN = sD
			tN = e + b
			tD = c
		}
	}
	if tN < 0.0 {
		tN = 0.0
		if -d < 0.0 {
			sN = 0.0
		} else if -d > a {
			sN = sD
		} else {
			sN = -d
			sD = a
		}
	} else if tN > tD {
		tN = tD
		if -d + b < 0.0 {
			sN = 0
		} else if -d + b > a {
			sN = sD
		} else {
			sN = -d +  b
			sD = a
		}
	}
	if math.Abs(sN) < SMALL_NUM {
		sc = 0.0
	} else {
		sc = sN / sD
	}
	if math.Abs(tN) < SMALL_NUM {
		tc = 0.0
	} else {
		tc = tN / tD
	}

	// get the difference of the two closest points
	dP := sub(add(w, scalaMul(sc, u)), scalaMul(tc, v))

	return norm(dP);   // return the closest distance

}

