var requestAnimFrame = window.requestAnimationFrame || window.webkitRequestAnimationFrame ||
                       window.mozRequestAnimationFrame || window.msRequestAnimationFrame || 
                       function(c) {window.setTimeout(c, 15)};
/**
   Phoria
   pho·ri·a (fôr-)
   n. The relative directions of The Eyes during binocular fixation on a given object
*/

// bind to window onload event
window.addEventListener('load', onloadHandler, false);

var testData = [{"ID":"drone1","pos":{"X":1,"Y":1,"Z":1},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}},
                {"ID":"drone2","pos":{"X":2,"Y":2,"Z":2},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}},
                {"ID":"drone3","pos":{"X":3,"Y":3,"Z":3},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}}];

var testData2 = [{"ID":"drone1","pos":{"X":-1,"Y":-1,"Z":-1},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}},
                {"ID":"drone2","pos":{"X":-2,"Y":-2,"Z":-2},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}},
                {"ID":"drone3","pos":{"X":-3,"Y":-3,"Z":-3},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}}];
var bitmaps = [];
var scene = new Phoria.Scene();
var sphereList = [];
var droneList = [];

var currentX = 0;
var currentY = 2;
var currentZ = 0;
var newX = 0;
var newY = 2;
var newZ = 0;
var dX = 0;
var dY = 0;
var dZ = 0;

var speed = 0.01;
var pause = true;

var drone1 = new Drone();

function createSphere(size, x, y, z) {

   var s = Phoria.Util.generateSphere(size, 24, 48);

   var offsetPoints = [];

   for(var pointNumber = 0; pointNumber < s.points.length; pointNumber++) {
       offsetPoints.push({
           x: s.points[pointNumber].x + x,
           y: s.points[pointNumber].y + y,
           z: s.points[pointNumber].z + z
       });
   }

   return Phoria.Entity.create({
       points: offsetPoints,
       edges: s.edges,
       polygons: s.polygons,
       style: {
           diffuse: 1,
           specular: 128
       }
   });
}

function makeSphereWithValue() {
      var input = document.getElementById('values');
      var data = input.value.split(",");
      //console.log("making sphere at " +  data[0] + "," + data[1] + "," + data[2]);
      for (var i = 0; i < droneList.length; i++) {
         droneList[i].newX = data[0];
         droneList[i].newY = data[1];
         droneList[i].newZ = data[2];
         droneList[i].dX = (droneList[i].newX - droneList[i].currentX) * droneList[i].speed;
         droneList[i].dY = (droneList[i].newY - droneList[i].currentY) * droneList[i].speed;
         droneList[i].dZ = (droneList[i].newZ - droneList[i].currentZ) * droneList[i].speed;
         pause = false;

      }


      //var sphere = createSphere(parseFloat(data[0]), parseInt(data[1]), parseInt(data[2]), parseInt(data[3]));
      //sphereList.push(sphereList);
      //scene.graph.push(sphere);
}

function formSquare() {

   // insert some math stuff to form shapes.
   /**
   var x = 0;
   var y = 0;
   var z = 0;
   for (var i = 0; i < droneList.length; i++) {
      if (i === 0) {
         x = -2;
         y = 2;
         z = 0;
      } else if (i === 1) {
         x = 2;
         y = 2;
         z = 0;
      } else if (i === 2) {
         x = -2;
         y = -2;
         z = 0;
      } else if (i === 3) {
         x = 2;
         y = -2;
         z = 0;
      }
      droneList[i].setCoordinate(x,y,z);
      pause = false;
   }
   **/
   updateDronePositions();
}

function jsonCallback(json) {
   console.log(json);
}

function loadDroneData() {
   // insert some call to fetch the initial data
   for (var i = 0; i < testData.length; i++) {
      var obj = testData[i];
      var drone = new Drone(obj.ID,obj.pos.X,obj.pos.Y,obj.pos.Z);
      var sphere = createSphere(drone.size, drone.currentX, drone.currentY, drone.currentZ);
      sphereList.push(sphere);
      scene.graph.push(sphere);

      drone.sphere = sphere;
      droneList.push(drone);
   }
}

function updateDronePositions() {
   for (var i = 0; i < testData2.length; i++) {
      var currentDrone = droneList[i];
      var objDrone = testData2[i];
      if (currentDrone.X != objDrone.pos.X || currentDrone.Y != objDrone.pos.Y || currentDrone.Z != objDrone.pos.Z  ) {
         currentDrone.setCoordinate(objDrone.pos.X,objDrone.pos.Y,objDrone.pos.Z);
      }
   }
   pause = false;
}

function onloadHandler()
{
   /**
   $.ajax({
     url: 'http://localhost:18842/drones',
     dataType: "jsonp"
   });
   
   console.log("onloadHandler");

   var script = document.createElement('script');
   script.src = 'http://localhost:18842/drones?callback=hooray';
   document.body.appendChild(script);
   **/

   loadDroneData();
       // get conference list

   // get the images loading
   var loader = new Phoria.Preloader();
   for (var i=0; i<6; i++)
   {
      bitmaps.push(new Image());
      loader.addImage(bitmaps[i], 'images/texture'+i+'.png');
   }
   loader.onLoadCallback(init);
}
function init()
{
   console.log("init()");

   // get the canvas DOM element and the 2D drawing context
   var canvas = document.getElementById('canvas');
   
   // create the scene and setup camera, perspective and viewport
   scene.camera.position = {x:0.0, y:5.0, z:-15.0};
   scene.perspective.aspect = canvas.width / canvas.height;
   scene.viewport.width = canvas.width;
   scene.viewport.height = canvas.height;
   
   // create a canvas renderer
   var renderer = new Phoria.CanvasRenderer(canvas);
   
   // add a grid to help visualise camera position etc.
   var plane = Phoria.Util.generateTesselatedPlane(8,8,0,20);
   scene.graph.push(Phoria.Entity.create({
      points: plane.points,
      edges: plane.edges,
      polygons: plane.polygons,
      style: {
         drawmode: "wireframe",
         shademode: "plain",
         linewidth: 0.5,
         objectsortmode: "back"
      }
   }));


   //scene.graph.push(cube);
   scene.graph.push(Phoria.DistantLight.create({
      direction: {x:0, y:-0.5, z:1}
   }));

   // added sphere
   //for (var i = 0; i < 4; i++) {

      /**
   for (var i = 0; i < 4; i++) {
      var drone = new Drone(0.5,i,0,0);
      var sphere = createSphere(drone.size, drone.currentX, drone.currentY, drone.currentZ);
      sphereList.push(sphere);
      scene.graph.push(sphere);

      drone.sphere = sphere;
      droneList.push(drone);

   }
   **/
   

   //}

   var animateX = 0.0;
   var animateY = 0.0;
   var fnAnimate = function() {
      if (!pause)
      {
         for (var i = 0; i < droneList.length; i++) {
            var drone = droneList[i];
            var sphere = drone.sphere;
            sphere.translateX(drone.dX);
            sphere.translateY(drone.dY);
            sphere.translateZ(drone.dZ);

            drone.currentX += drone.dX;
            drone.currentY += drone.dY;
            drone.currentZ += drone.dZ;

            //console.log(currentX + ", " + currentY + ", " + currentZ);
            //console.log(newX + ", " + newY + ", " + newZ);
            //console.log((currentX - newX) + ", " + (currentY - newY) + ", " + (currentZ - newZ));

            if ((Math.abs(drone.currentX - drone.newX) < 0.001) && (Math.abs(drone.currentY - drone.newY) < 0.001) && (Math.abs(drone.currentZ - drone.newZ) < 0.001)) {
               drone.dX = 0;
               drone.dY = 0;
               drone.dZ = 0;
            } else {
               //pause = false;
            }
         }
         //sphere.translateY(0.01);
         // rotate local matrix of the cube
         // cube.rotateY(0.5*Phoria.RADIANS);
         /**
         for (var i = 0; i < sphereList.length; i++) {
            var sphereL = sphereList[i];
            sphereL.translateY(0.01);
         }
         **/
         //childCube.identity().translateY(Math.sin(Date.now() / 1000) + 3);
         
         // execute the model view 3D pipeline and render the scene
      }
               scene.modelView();
         renderer.render(scene);
      requestAnimFrame(fnAnimate);
   };
   
   // keep track of heading to generate position
   var heading = 0.0;
   var lookAt = vec3.fromValues(0,-5,15);

   /**
    * @param forward {vec3}   Forward movement offset
    * @param heading {float}  Heading in Phoria.RADIANS
    * @param lookAt {vec3}    Lookat projection offset from updated position
    */
   var fnPositionLookAt = function positionLookAt(forward, heading, lookAt) {
      // recalculate camera position based on heading and forward offset
      var pos = vec3.fromValues(
         scene.camera.position.x,
         scene.camera.position.y,
         scene.camera.position.z);
      var ca = Math.cos(heading), sa = Math.sin(heading);
      var rx = forward[0]*ca - forward[2]*sa,
          rz = forward[0]*sa + forward[2]*ca;
      forward[0] = rx;
      forward[2] = rz;
      vec3.add(pos, pos, forward);
      scene.camera.position.x = pos[0];
      scene.camera.position.y = pos[1];
      scene.camera.position.z = pos[2];

      // calcuate rotation based on heading - apply to lookAt offset vector
      rx = lookAt[0]*ca - lookAt[2]*sa,
      rz = lookAt[0]*sa + lookAt[2]*ca;
      vec3.add(pos, pos, vec3.fromValues(rx, lookAt[1], rz));

      // set new camera look at
      scene.camera.lookat.x = pos[0];
      scene.camera.lookat.y = pos[1];
      scene.camera.lookat.z = pos[2];
   }
   
   // key binding
   document.addEventListener('keydown', function(e) {
      switch (e.keyCode)
      {
          /**
         case 32: // spacebar 
            var sphere = createSphere(0.5,newX,0,0);
            newX-=1;
            scene.graph.push(sphere);
            console.log("making sphere of size 0.5 at " +  newX + ",0,0");
            break;
            **/
         case 27: // ESC
            console.log("pausing");
            pause = !pause;
            break;
           /**
         case 87: // W
            // move forward along current heading
            fnPositionLookAt(vec3.fromValues(0,0,1), heading, lookAt);
            break;
         case 83: // S
            // move back along current heading
            fnPositionLookAt(vec3.fromValues(0,0,-1), heading, lookAt);
            break;
         case 65: // A
            // strafe left from current heading
            fnPositionLookAt(vec3.fromValues(-1,0,0), heading, lookAt);
            break;
         case 68: // D
            // strafe right from current heading
            fnPositionLookAt(vec3.fromValues(1,0,0), heading, lookAt);
            break;
            
         case 37: // LEFT
            // turn left
            heading += Phoria.RADIANS*4;
            // recalculate lookAt
            // given current camera position, project a lookAt vector along current heading for N units
            fnPositionLookAt(vec3.fromValues(0,0,0), heading, lookAt);
            break;
         case 39: // RIGHT
            // turn right
            heading -= Phoria.RADIANS*4;
            // recalculate lookAt
            // given current camera position, project a lookAt vector along current heading for N units
            fnPositionLookAt(vec3.fromValues(0,0,0), heading, lookAt);
            break;
         case 38: // UP
            lookAt[1]++;
            fnPositionLookAt(vec3.fromValues(0,0,0), heading, lookAt);
            break;
         case 40: // DOWN
            lookAt[1]--;
            fnPositionLookAt(vec3.fromValues(0,0,0), heading, lookAt);
            break;
            **/
      }
   }, false);

   /*
   KEY:
   {
      SHIFT:16, CTRL:17, ESC:27, RIGHT:39, UP:38, LEFT:37, DOWN:40, SPACE:32,
      A:65, E:69, G:71, L:76, P:80, R:82, S:83, Z:90
   },
   */

   // add GUI controls
   var gui = new dat.GUI();
   var f = gui.addFolder('Perspective');
   f.add(scene.perspective, "fov").min(5).max(175);
   f.add(scene.perspective, "near").min(1).max(100);
   f.add(scene.perspective, "far").min(1).max(1000);
   f = gui.addFolder('Camera LookAt');
   f.add(scene.camera.lookat, "x").min(-100).max(100);
   f.add(scene.camera.lookat, "y").min(-100).max(100);
   f.add(scene.camera.lookat, "z").min(-100).max(100);
   f = gui.addFolder('Camera Position');
   f.add(scene.camera.position, "x").min(-100).max(100);
   f.add(scene.camera.position, "y").min(-100).max(100);
   f.add(scene.camera.position, "z").min(-100).max(100);
   f = gui.addFolder('Camera Up');
   f.add(scene.camera.up, "x").min(-10).max(10).step(0.1);
   f.add(scene.camera.up, "y").min(-10).max(10).step(0.1);
   f.add(scene.camera.up, "z").min(-10).max(10).step(0.1);

   // start animation
   requestAnimFrame(fnAnimate);
}