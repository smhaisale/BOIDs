package boids

type IDrone interface {

    formShape()

}

type Drone struct {
    position Position
    params DroneType
    speed Speed
}

