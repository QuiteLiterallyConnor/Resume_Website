(function() {
    const updateInfo = () => { camera.updateProjectionMatrix(); updateModelInfo(); };

    const tweenMesh = (prop, obj, vals) => new TWEEN.Tween(obj[prop]).to(vals, 500).easing(TWEEN.Easing.Cubic.InOut).start();

    function moveModel(mesh, y) {
        if (mesh) tweenMesh('position', mesh, { y });
    }

    function scaleModel(mesh, s) {
        if (mesh) tweenMesh('scale', mesh, { x: s, y: s, z: s }).onComplete(updateModelInfo);
    }

    function moveCamera(y) {
        tweenMesh('position', camera, { y }).onUpdate(updateCameraInfo);
    }

    function clickedExpand() {
        const settings = {
            nose: { pos: [55, 10, -10], scale: [0.4, 0.1, 0.1], cam: 35, model: 'nose' },
            fuselage: { pos: [50, 0, -25], scale: [0.1, 0.2, 0.1], cam: -5, model: 'fuselage' },
            base: { pos: [30, -5, -45], scale: [0.1, 0.1, 0.3], cam: -50, model: 'base' }
        };

        const { pos, scale, cam } = settings["nose"];
        moveModel(meshes.nose, pos[0]);
        moveModel(meshes.fuselage, pos[1]);
        moveModel(meshes.base, pos[2]);
        // moveCamera(cam);
        // updateInfo();
    };

    document.getElementById('expandButton').addEventListener('click', function() {
        clickedExpand(); // Call with appropriate settings
    });


})();
