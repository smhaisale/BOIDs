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

// var testData = [{"ID":"drone1","pos":{"X":0,"Y":0,"Z":0},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}},
//     {"ID":"drone2","pos":{"X":0,"Y":0,"Z":0},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}},
//     {"ID":"drone3","pos":{"X":0,"Y":0,"Z":0},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}}];
//
// var testData2 = [{"ID":"drone1","pos":{"X":-1,"Y":-1,"Z":-1},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}},
//     {"ID":"drone2","pos":{"X":-2,"Y":-2,"Z":-2},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}},
//     {"ID":"drone3","pos":{"X":-3,"Y":-3,"Z":-3},"type":{"TypeId":"type1","TypeDescription":"Simple sample drone type","Size":{"DX":1,"DY":1,"DZ":1},"MaxRange":{"DX":10,"DY":10,"DZ":10},"MaxSpeed":{"VX":10,"VY":10,"VZ":10}},"speed":{"VX":2,"VY":1,"VZ":-1}}];

var bitmaps = [];
var scene = new Phoria.Scene();
var sphereList = [];
var droneList = [];
var droneMap = {};

var pause = true;
var debug = false;

var drone1 = new Drone();

function createSphere(id, size, x, y, z) {

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
        id: id,
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

//Used in frontend to refresh drone positions
function flipPause() {
    pause = !pause;
    if (pause) {
        document.getElementById("flipPause").innerHTML = "Resume";
    } else {
        document.getElementById("flipPause").innerHTML = "Pause";
    }
    console.log("Paused: " + pause);
}

//Used in frontend to refresh drone positions
function flipDebug() {
    debug = !debug;
    if (debug) {
        for (var id in droneMap) {
            Phoria.Entity.debug(droneMap[id].sphere, {
                showId: true,
                showPosition: true
            });
        }
        /**
        for (var i = 0; i < droneList.length; i++) {
            Phoria.Entity.debug(droneList[i].sphere, {
                showId: true,
                showPosition: true
            });
        }
        **/
        document.getElementById("flipDebug").innerHTML = "Remove Debug Info";
    } else {
        for (var id in droneMap) {
            Phoria.Entity.debug(droneMap[id].sphere, {
                showId: false,
                showPosition: false
            });
        }
        /**
        for (var i = 0; i < droneList.length; i++) {
            Phoria.Entity.debug(droneList[i].sphere, {
                showId: false,
                showPosition: false
            });
        }
        **/
        document.getElementById("flipDebug").innerHTML = "Add Debug Info";
    }
    console.log("Debug: " + debug);
}

function addDroneToEnvironment() {
    var address = document.getElementById("droneAddress").value;
    var addDroneUrl = 'http://localhost:18842/addDrone?messageType=type&data=' + address;

    console.log(addDroneUrl);
    $.ajax({
    type: 'GET',
    dataType: 'json',
    url: addDroneUrl ,
    success: function (data) {
        console.log("added drone to the list");
    }
});

}

function removeDroneFromEnvironment() {
    var address = document.getElementById("droneId").value;
    // kill drone not implemented yet, call the function here when it's done
    //var addDroneUrl = 'http://localhost:18842/addDrone?messageType=type&data=' + address;
    /**
    $.ajax({
    type: 'GET',
    dataType: 'json',
    url: addDroneUrl ,
    success: function (data) {
        console.log("added drone to the list");
    }
    **/
}

function loadDroneData() {
    updateDronePositions();
}

function updateDronePositions() {

    if (!pause) {
        $.ajax({
            type: 'GET',
            dataType: 'json',
            url: 'http://localhost:18842/getAllDrones',
            success: function (data) {
                console.log(data);
                var mapSize = 0, key;
                for (key in droneMap) {
                    if (droneMap.hasOwnProperty(key)) mapSize++;
                }
                

                for (var i = 0; i < data.length; i++) {
                    var object = data[i];
                    if (mapSize <= i) {
                        var drone = new Drone(object.ID, object.DroneObject.pos.X, object.DroneObject.pos.Y, object.DroneObject.pos.Z);
                        var sphere = createSphere(object.ID, drone.size, drone.currentX, drone.currentY, drone.currentZ);
                        sphereList.push(sphere);
                        scene.graph.push(sphere);

                        drone.sphere = sphere;
                        droneMap[object.ID] = drone;

                    /**
                    if (droneList.length <= i) {
                        var drone = new Drone(object.ID, object.pos.X, object.pos.Y, object.pos.Z);
                        var sphere = createSphere(object.ID, drone.size, drone.currentX, drone.currentY, drone.currentZ);
                        sphereList.push(sphere);
                        scene.graph.push(sphere);

                        drone.sphere = sphere;
                        droneList.push(drone);
                        **/
                    } else {
                        var currentDrone = droneMap[object.ID];
                        if (currentDrone.X != object.DroneObject.pos.X || currentDrone.Y != object.DroneObject.pos.Y || currentDrone.Z != object.DroneObject.pos.Z) {
                            currentDrone.setCoordinate(object.DroneObject.pos.X, object.DroneObject.pos.Y, object.DroneObject.pos.Z);
                        }
                        /**
                        var currentDrone = droneList[i];
                        if (currentDrone.X != object.pos.X || currentDrone.Y != object.pos.Y || currentDrone.Z != object.pos.Z) {
                            currentDrone.setCoordinate(object.pos.X, object.pos.Y, object.pos.Z);
                        }
                        **/
                    }
                }
                setTimeout(updateDronePositions, 1000);
            }
        });
    } else {
        setTimeout(updateDronePositions, 1000);
    }

}

function onloadHandler()
{
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

    var Rnd = function(s) {return Math.random() * s};

    // get the canvas DOM element and the 2D drawing context
    var canvas = document.getElementById('canvas');

    // remove frame margin and scrollbars when maxing out size of canvas
    document.body.style.margin = "0px";
    document.body.style.overflow = "hidden";

    // get dimensions of window and resize the canvas to fit
    // var width = window.innerWidth, height = window.innerHeight - 200;
    // canvas.width = width; canvas.height = height;

    // create the scene and setup camera, perspective and viewport
    scene.camera.position = {x:0, y:25.0, z:-60.0};
    scene.camera.lookat = {x:0.0, y:10.0, z:0.0};
    scene.perspective.aspect = canvas.width / canvas.height;
    scene.viewport.width = canvas.width;
    scene.viewport.height = canvas.height;

    // create a canvas renderer
    var renderer = new Phoria.CanvasRenderer(canvas);
    // add a grid to help visualise camera position etc.
    var plane = Phoria.Util.generateTesselatedPlane(16,16,0,40);
    scene.graph.push(Phoria.Entity.create({
        points: plane.points,
        edges: plane.edges,
        polygons: plane.polygons,
        style: {
            //color: [160,255,160],
            drawmode: "wireframe",
            shademode: "plain",
            linewidth: 0.5,
            objectsortmode: "back"
        }
    }));

    var fnGenerateStarfield = function(num, scale) {
        scale = scale || 1;
        var s = scale / 2;
        var points = [];
        for (var i = 0; i < num; i++) {
            // TODO: points too close to the origin (i.e. camera view point) shoud be discared
            points.push({x:Math.random()*scale-s, y:Math.random()*scale-s, z:Math.random()*scale-s});
        }
        return Phoria.Entity.create({
            points: points,
            style: {
                color: [100+~~(Math.random()*55),100+~~(Math.random()*55),100+~~(Math.random()*55)],
                drawmode: "point",
                shademode: "plain",
                linewidth: 0.75,
                objectsortmode: "back"
            }
        });
    };
    scene.graph.push(fnGenerateStarfield(500,2000));

    // rotate the camera around the scene
    scene.onCamera(function(position, lookAt, up) {
        var rotMatrix = mat4.create();
        mat4.rotateY(rotMatrix, rotMatrix, Math.sin(Date.now()/10000)*Phoria.RADIANS*360);
        vec4.transformMat4(position, position, rotMatrix);
    });

    scene.graph.push(Phoria.DistantLight.create({
        direction: {x:0, y:-0.5, z:1}
    }));

    var fnAnimate = function() {
        if (!pause)
        {
            for (var id in droneMap) {
         //   for (var i = 0; i < droneList.length; i++) {
         //       var drone = droneList[i];
                var drone = droneMap[id];
                var sphere = drone.sphere;
                sphere.translateX(drone.dX);
                sphere.translateY(drone.dY);
                sphere.translateZ(drone.dZ);

                drone.currentX += drone.dX;
                drone.currentY += drone.dY;
                drone.currentZ += drone.dZ;

                if ((Math.abs(drone.currentX - drone.newX) < 0.001) && (Math.abs(drone.currentY - drone.newY) < 0.001) && (Math.abs(drone.currentZ - drone.newZ) < 0.001)) {
                    drone.dX = 0;
                    drone.dY = 0;
                    drone.dZ = 0;
                }
            }
        }
        scene.modelView();
        renderer.render(scene);
        requestAnimFrame(fnAnimate);
    };

    // key binding
    document.addEventListener('keydown', function(e) {
        if (e.keyCode == 27)
        {
            console.log("pausing");
            pause = !pause;
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
    //f = gui.addFolder('Camera LookAt');
    //f.add(scene.camera.lookat, "x").min(-100).max(100);
    //f.add(scene.camera.lookat, "y").min(-100).max(100);
    //f.add(scene.camera.lookat, "z").min(-100).max(100);
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