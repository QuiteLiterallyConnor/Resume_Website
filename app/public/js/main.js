let scene, camera, renderer;
const meshes = { base: null, fuselage: null, nose: null };

const currentPositions = { base: { x: 0, y: 0, z: 0 }, fuselage: { x: 0, y: 0, z: 0 } };
const rotationSpeed = (2 * Math.PI) / (10 * 120);

function init() {
    scene = new THREE.Scene();
    const aspect = (window.innerWidth - 200) / window.innerHeight;
    const d = 50;  // Decreased from 100 to 50 for a more zoomed-in view
    camera = new THREE.OrthographicCamera(-d * aspect, d * aspect, d, -d, 1, 1000);
    camera.position.set(0, 0, 50);
    camera.lookAt(new THREE.Vector3(0, 0, 0));
    renderer = new THREE.WebGLRenderer({ antialias: true });
    renderer.setSize(window.innerWidth - 200, window.innerHeight);
    document.getElementById('threeContainer').appendChild(renderer.domElement);
    scene.add(new THREE.AmbientLight(0x404040, 1));
    const directionalLight = new THREE.DirectionalLight(0xffffff, 1);
    directionalLight.position.set(1, 1, 1).normalize();
    scene.add(directionalLight);
    loadModels(scene);
    window.addEventListener('resize', onWindowResize, false);

    const buttons = ['base', 'fuselage', 'nose'];
    buttons.forEach(model => {
        const btn = document.getElementById(`move${model}Button`);
        if (btn) {
            btn.addEventListener('click', () => handleButtonClick(model));
        }
    });
}

function onWindowResize() {
    const aspect = (window.innerWidth - 200) / window.innerHeight;
    const d = 50;  // Decreased from 100 to 50 for a more zoomed-in view
    Object.assign(camera, { left: -d * aspect, right: d * aspect, top: d, bottom: -d });
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth - 200, window.innerHeight);
}

function reset() {
    Object.values(meshes).forEach(mesh => mesh && scene.remove(mesh));
    Object.keys(meshes).forEach(key => meshes[key] = null);
    camera.position.set(0, 0, 50);
    camera.lookAt(new THREE.Vector3(0, 0, 0));
    loadModels(scene);
}

function animate() {
    requestAnimationFrame(animate);
    TWEEN.update();
    Object.values(meshes).forEach(mesh => mesh && (mesh.rotation.y += rotationSpeed));
    renderer.render(scene, camera);
}

function startThreeJS() {
    document.getElementById('introVideo').style.display = 'none';
    document.querySelector('.sidebar').classList.remove('hidden');
    init();
    animate();
}

document.addEventListener('DOMContentLoaded', () => {
    const introVideo = document.getElementById('introVideo');
    const startButton = document.getElementById('startButton');

    startButton.addEventListener('click', () => {
        startButton.style.display = 'none';
        introVideo.classList.remove('hidden');

        // Ensure the video plays after user interaction
        introVideo.play().catch(error => {
            console.error("Error attempting to play the video:", error);
        });
    });

    introVideo.onended = startThreeJS;

    document.getElementById('introVideo').addEventListener('play', function() {
        this.style.display = 'block';
    });
});

function startThreeJS() {
    document.getElementById('introVideo').style.display = 'none';
    document.querySelector('.sidebar').classList.remove('hidden');
    init();
    animate();
}
