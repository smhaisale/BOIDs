function Drone(id,x,y,z) {
   this.ID;
   this.size = 0.5;
   this.currentX = x;
   this.currentY = y;
   this.currentZ = z;
   this.newX = 0;
   this.newY = 0;
   this.newZ = 0;
   this.dX = 0;
   this.dY = 0;
   this.dZ = 0;
   this.speed = 0.01;
   this.sphere;   
}

Drone.prototype.setCoordinate = function(x,y,z) {
   this.newX = x;
   this.newY = y;
   this.newZ = z;
   this.dX = (this.newX - this.currentX) * this.speed;
   this.dY = (this.newY - this.currentY) * this.speed;
   this.dZ = (this.newZ - this.currentZ) * this.speed;
}