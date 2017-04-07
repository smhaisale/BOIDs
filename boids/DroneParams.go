package boids

type Position struct {
    X, Y, Z float64                 // 3D position of the drone
}

type Speed struct {
    vX, vY, vZ float64              // Dimensional velocities
    speed float64                   // The actual velocity
}

type Dimensions struct {
    dX, dY, dZ float64              // Any set of X, Y and Z dimensions
}

type DroneType struct {
    typeId string                   // Identifier for a type
    typeDescription string          // Description
    size Dimensions                 // Size of this drone type
    maxRange Dimensions             // Communication with another drone fails outside this range
    maxSpeed Speed                  // Maximum speed for this drone type
}