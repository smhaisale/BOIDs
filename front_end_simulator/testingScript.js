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

var bitmaps = [];
var scene = new Phoria.Scene();
var sphereList = [];
var droneList = [];
var droneMap = {};

var pause = true;
var animate = false;
var debug = false;

var drone1 = new Drone();

function getRandomArbitrary(min, max) {
  return Math.floor(Math.random() * (max - min)) + min;
}

function createSphere(id, size, newX, newY, newZ, r, g, b) {

    var rValue = Math.floor(r);//getRandomArbitrary(0,255);
    var gValue = Math.floor(g);//getRandomArbitrary(0,255);
    var bValue = Math.floor(b);//getRandomArbitrary(0,255);

    /**
   var blueLightObj = Phoria.Entity.create({
      id: id,
      points: [{x:newX, y:newY, z:newZ}],
      style: {
         color: [rValue,gValue,bValue],
         drawmode: "point",
         shademode: "plain",
         linewidth: 10,
         linescale: 10
      }
   });
   var blueLight = Phoria.PointLight.create({
      position: {x:newX, y:newY, z:newZ},
      color: [0,0,0]
   });
   blueLightObj.children.push(blueLight);

   return blueLightObj;
     **/

    var s = Phoria.Util.generateSphere(size, 24, 48);

    var offsetPoints = [];

    for(var pointNumber = 0; pointNumber < s.points.length; pointNumber++) {
        offsetPoints.push({
            x: s.points[pointNumber].x + newX,
            y: s.points[pointNumber].y + newY,
            z: s.points[pointNumber].z + newZ
        });
    }

    return Phoria.Entity.create({
        id: id,
        points: offsetPoints,
        edges: s.edges,
        polygons: s.polygons,
        style: {
            color: [rValue,gValue,bValue],
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
}

//Used in frontend to refresh drone positions
function flipAnimate() {
    animate = !animate;
    if (animate) {
        // rotate the camera around the scene
        scene.onCamera(function(position, lookAt, up) {
            var rotMatrix = mat4.create();
            mat4.rotateY(rotMatrix, rotMatrix, Math.sin(Date.now()/10000)*Phoria.RADIANS*360);
            vec4.transformMat4(position, position, rotMatrix);
        });
    } else {
        scene.onCameraHandlers = null;
    }
    console.log("Paused: " + pause);
}

//Used in frontend to refresh drone positions
function flipPause() {
    pause = !pause;
    console.log("Paused: " + pause);
}

//Used in frontend to refresh drone positions
function flipDebug() {
    debug = !debug;
    if (debug) {
        for (var id in droneMap) {
            console.log(droneMap[id]);
            Phoria.Entity.debug(droneMap[id].sphere, {
                showId: true,
                showPosition: true
            });
        }
    } else {
        for (var id in droneMap) {
            Phoria.Entity.debug(droneMap[id].sphere, {
                showId: false,
                showPosition: false
            });
        }
    }
    console.log("Debug: " + debug);
}

function addDroneToEnvironment(address) {
    var addDroneUrl = 'http://' + document.location.hostname + ':18842/addDrone?messageType=type&data=' + address;

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

function formPolygon() {
    var formPolygonUrl = 'http://' + document.location.hostname + ':18842/formPolygon';
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: formPolygonUrl,
        success: function (data) {
            console.log("Sent form polygon request");
        }
    })
}

function formPrism() {
    var formPolygonUrl = 'http://' + document.location.hostname + ':18842/formShape?shape=prism';
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: formPolygonUrl,
        success: function (data) {
            console.log("Sent form prism request");
        }
    })
}

function formPyramid() {
    var formPolygonUrl = 'http://' + document.location.hostname + ':18842/formShape?shape=pyramid';
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: formPolygonUrl,
        success: function (data) {
            console.log("Sent form pyramid request");
        }
    })
}

function formBipyramid() {
    var formPolygonUrl = 'http://' + document.location.hostname + ':18842/formShape?shape=bipyramid';
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: formPolygonUrl,
        success: function (data) {
            console.log("Sent form bipyramid request");
        }
    })
}

function randomPositions() {
    var randomPositionsUrl = 'http://' + document.location.hostname + ':18842/randomPositions';
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: randomPositionsUrl,
        success: function (data) {
            console.log("Sent random positions request");
        }
    })
}

function setDronePosition(idPosition) {
    console.log("setting the drone position");
    var id = idPosition.substring(0,idPosition.indexOf(':'));
    var position = idPosition.substring(idPosition.indexOf(':') + 1);
    var positionData =  position.split(",");
    var droneToChange = droneMap[id];
    /**
    console.log(position);
    console.log(positionData);
    console.log(positionData[0]);
    console.log(droneToChange);

    droneToChange.newX = parseFloat(positionData[0]);
    droneToChange.newY = parseFloat(positionData[1]);
    droneToChange.newZ = parseFloat(positionData[2]);
    droneToChange.dX = (droneToChange.newX - droneToChange.currentX) * droneToChange.speed;
    droneToChange.dY = (droneToChange.newY - droneToChange.currentY) * droneToChange.speed;
    droneToChange.dZ = (droneToChange.newZ - droneToChange.currentZ) * droneToChange.speed;
    **/
    var droneAddress = droneToChange.address;
    var formPolygonUrl = droneAddress + '/moveToPosition?X=' + positionData[0] + '&Y=' + positionData[1] + '&Z=' + positionData[2];
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: formPolygonUrl,
        success: function (data) {
            console.log("Sent new drone position request");
        }
    })
}

function removeDroneFromEnvironment(address) {
    var killDroneUrl = 'http://' + document.location.hostname + ':18842/killDrone?messageType=type&data=' + address;
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: killDroneUrl,
        success: function (data) {
            console.log("added drone to the list");
        }
    })
}

function loadDroneData() {
    updateDronePositions();
}

function updateDronePositions() {

    if (!pause) {
        $.ajax({
            type: 'GET',
            dataType: 'json',
            url: 'http://' + document.location.hostname + ':18842/getAllDrones',
            success: function (data) {
                console.log(data);
                var mapSize = 0, key;
                for (key in droneMap) {
                    if (droneMap.hasOwnProperty(key)) mapSize++;
                }

                for (var i = 0; i < data.length; i++) {
                    var object = data[i];
                    if (mapSize <= i) {
                        var drone = new Drone(object.ID, object.DroneObject.pos.X, object.DroneObject.pos.Y, object.DroneObject.pos.Z, object.DroneObject.color.X, object.DroneObject.color.Y, object.DroneObject.color.Z, object.DroneObject.size);
                        var droneAddress = object.Address;
                        drone.address = droneAddress.substring(droneAddress.indexOf(':')+1);
                        console.log(drone.address);
                        console.log(drone);
                        var sphere = createSphere(object.ID, drone.size, drone.currentX, drone.currentY, drone.currentZ, drone.r, drone.g, drone.b);
                        //sphereList.push(sphere);
                        scene.graph.push(sphere);

                        drone.sphere = sphere;
                        droneMap[object.ID] = drone;
                    } else {
                        var currentDrone = droneMap[object.ID];
                        if (currentDrone.X != object.DroneObject.pos.X || currentDrone.Y != object.DroneObject.pos.Y || currentDrone.Z != object.DroneObject.pos.Z) {
                            currentDrone.setCoordinate(object.DroneObject.pos.X, object.DroneObject.pos.Y, object.DroneObject.pos.Z);
                        }
                    }
                }
                setTimeout(updateDronePositions, 16);
            }
        });
    } else {
        setTimeout(updateDronePositions, 16);
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
    var width = window.innerWidth, height = window.innerHeight;
    canvas.width = width; canvas.height = height;

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

    var light = Phoria.DistantLight.create({
        color: [1.0,1.0,1.0],
        direction: {x:1, y:-1, z:0}
    });
    var light2 = Phoria.DistantLight.create({
        color: [1.0,1.0,1.0],
        direction: {x:-1, y:-1, z:0}
    });
    scene.graph.push(light);
    scene.graph.push(light2);

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

    var obj = { add:function(){ console.log("clicked") }};
    var drone = {
        start : false,
        debug : false,
        animate : false,
        address : '',
        id : '',
        position : '',
        formPolygon: function() { formPolygon()},
        randomPositions: function() { randomPositions()},
        setDronePosition: function() { setDronePosition()},
        formPrism: function() { formPrism()},
        formPyramid: function() { formPyramid()},
        formBipyramid: function() { formBipyramid()}
    };

    // add GUI controls
    var gui = new dat.GUI();

    var f = gui.addFolder('Perspective');
    f.add(scene.perspective, "fov").min(5).max(175);
    f.add(scene.perspective, "near").min(1).max(100);
    f.add(scene.perspective, "far").min(1).max(200);

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

    f = gui.addFolder('Light');
    f.add(light.direction, "x").min(-25).max(25).step(0.1);
    f.add(light.direction, "y").min(-25).max(25).step(0.1);
    f.add(light.direction, "z").min(-25).max(25).step(0.1);
    f.add(light.color, "0").min(0).max(1).step(0.1).name("red");
    f.add(light.color, "1").min(0).max(1).step(0.1).name("green");
    f.add(light.color, "2").min(0).max(1).step(0.1).name("blue");
    f.add(light, "intensity").min(0).max(1).step(0.1);

    f = gui.addFolder('Drone Controls')
    f.add(drone, 'animate').name('Animate').onFinishChange(function(){flipAnimate()});
    f.add(drone, 'start').name('Running').onFinishChange(function(){flipPause()});
    f.add(drone, 'debug').name('Show Debug Info').onFinishChange(function(){flipDebug()});
    f.add(drone, 'address').name('Add Drone').onFinishChange(function(){addDroneToEnvironment(drone.address)});
   // f.add(drone, 'id').name('Kill Drone').onFinishChange(function(){addDroneToEnvironment(drone.id)});
    f.add(drone, 'position').name('Set drone position').onFinishChange(function(){setDronePosition(drone.position)});
    f.add(drone, 'formPolygon').name('Form Polygon');
    f.add(drone, 'randomPositions').name('Random Positions');
    f.add(drone, 'formPrism').name('Form Prism');
    f.add(drone, 'formPyramid').name('Form Pyramid');
    f.add(drone, 'formBipyramid').name('Form Bipyramid');
    f.open();

    // start animation
    requestAnimFrame(fnAnimate);
}